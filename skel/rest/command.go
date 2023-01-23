package rest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"path"
	"regexp"
	"strings"
)

const (
	ActionList     = "list"
	ActionRetrieve = "retrieve"
	ActionCreate   = "create"
	ActionUpdate   = "update"
	ActionDelete   = "delete"

	ComponentGlobals    = "pkg"
	ComponentPayloads   = "payloads"
	ComponentController = "controller"
	ComponentConverter  = "converter"
	ComponentService    = "service"
	ComponentRepository = "repository"
	ComponentModel      = "model"

	TenantNone      = ""
	TenantSingle    = "single"
	TenantHierarchy = "hierarchy"

	StyleV2 = "v2"
	StyleV8 = "v8"

	MatchesAny = "*"
)

var (
	ActionOptions = []string{
		ActionList,
		ActionRetrieve,
		ActionCreate,
		ActionUpdate,
		ActionDelete,
	}

	ComponentOptions = []string{
		ComponentGlobals,
		ComponentPayloads,
		ComponentController,
		ComponentConverter,
		ComponentService,
		ComponentRepository,
		ComponentModel,
	}

	TenantOptions = []string{
		TenantNone,
		TenantSingle,
		TenantHierarchy,
	}

	StyleOptions = []string{
		StyleV2,
		StyleV8,
	}
)

type GeneratorConfig struct {
	Style      string
	Domain     string
	Spec       string
	Folder     string
	Tenant     string
	Actions    []string
	Components []string
}

func (c GeneratorConfig) Validate() error {
	return types.ErrorMap{
		"domain": validation.Validate(&c.Domain, validation.Required,
			validation.NewStringRule(
				regexp.MustCompile(`^[A-Za-z][A-Za-z ]+$`).MatchString,
				"Domain must contain only letters or spaces, and start with a letter")),
		"folder": validation.Validate(&c.Folder, validation.Required,
			validation.NewStringRule(
				regexp.MustCompile(`^^[A-Za-z]+(/[A-Za-z]+)*/?$`).MatchString,
				"Folder must contain 1 or more components of only letters, separated by slashes")),
		"style":      validation.Validate(&c.Style, validation.Required, validation.In(types.Slice[string](StyleOptions).AnySlice()...)),
		"tenant":     validation.Validate(&c.Tenant, validation.In(types.Slice[string](TenantOptions).AnySlice()...)),
		"actions":    validation.Validate(c.Actions, validation.Each(validation.In(types.Slice[string](ActionOptions).AnySlice()...))),
		"components": validation.Validate(c.Components, validation.Each(validation.In(types.Slice[string](ComponentOptions).AnySlice()...))),
	}
}

func (c GeneratorConfig) PackageName() string {
	_, packageName := path.Split(c.Folder)
	return packageName
}

func (c GeneratorConfig) PackagePath() string {
	return path.Join(skel.Config().AppPackageUrl(), c.Folder)
}

var generatorConfig GeneratorConfig

func init() {
	cmd := skel.AddTarget("generate-domain", "Create domain", GenerateDomain)
	cmd.PreRunE = generateDomainOptions
	cmd.Args = cobra.MinimumNArgs(1)
	cmd.Flags().StringVar(&generatorConfig.Folder, "folder", "internal/lowerplural", "Output folder")
	cmd.Flags().StringVar(&generatorConfig.Style, "style", StyleV8, "API style.  One of: v2, v8")
	cmd.Flags().StringVar(&generatorConfig.Spec, "spec", "", "OpenApi specification")
	cmd.Flags().StringVar(&generatorConfig.Tenant, "tenant", TenantNone, "Tenant access control.  One of: single, hierarchy")
	cmd.Flags().StringArrayVar(&generatorConfig.Actions, "actions", []string{}, "Generate specific actions.  Any of: list, retrieve, create, update, delete")
	cmd.Flags().StringArrayVar(&generatorConfig.Components, "components", []string{}, "Generate specific components.  Any of: pkg, api, controller, converter, service, repository, model")
}

func generateDomainOptions(_ *cobra.Command, args []string) error {
	generatorConfig.Domain = strings.TrimSpace(strings.Join(args, " "))

	if len(generatorConfig.Actions) == 0 {
		generatorConfig.Actions = ActionOptions
	}

	if len(generatorConfig.Components) == 0 {
		generatorConfig.Components = ComponentOptions
	}

	generatorConfig.Folder = skel.NewInflector(generatorConfig.Domain).Inflect(generatorConfig.Folder)

	err := validate.Validate(generatorConfig)
	if err != nil {
		return err
	}

	return nil
}

type ComponentGenerator interface {
	Generate() error
	Render() string
	Filename() string
	Variables() map[string]string
	Conditions() map[string]bool
}

type ScopedComponentGenerator struct {
	Style     string
	Tenant    string
	Component string
	Factory   func(spec Spec) ComponentGenerator
}

var componentGenerators = []ScopedComponentGenerator{
	{
		Style:     MatchesAny,
		Tenant:    MatchesAny,
		Component: ComponentGlobals,
		Factory:   NewDomainGlobalsGenerator,
	},
	{
		Style:     MatchesAny,
		Tenant:    MatchesAny,
		Component: ComponentPayloads,
		Factory:   NewDomainPayloadsGenerator,
	},
	{
		Style:     StyleV8,
		Tenant:    MatchesAny,
		Component: ComponentController,
		Factory:   NewDomainControllerGeneratorV8,
	},
	{
		Style:     StyleV8,
		Tenant:    TenantNone,
		Component: ComponentService,
		Factory:   NewDomainServiceGeneratorV8,
	},
	{
		Style:     StyleV8,
		Tenant:    TenantNone,
		Component: ComponentConverter,
		Factory:   NewDomainConverterGeneratorV8,
	},
	{
		Style:     MatchesAny,
		Tenant:    MatchesAny,
		Component: ComponentModel,
		Factory:   NewDomainModelGenerator,
	},
	{
		Style:     MatchesAny,
		Tenant:    MatchesAny,
		Component: ComponentRepository,
		Factory:   NewDomainRepositoryGenerator,
	},
}

func findGenerator(style, tenant, component string) *ScopedComponentGenerator {
	for _, componentGenerator := range componentGenerators {
		if componentGenerator.Style != MatchesAny && componentGenerator.Style != style {
			continue
		}
		if componentGenerator.Tenant != MatchesAny && componentGenerator.Tenant != tenant {
			continue
		}
		if componentGenerator.Component != MatchesAny && componentGenerator.Component != component {
			continue
		}
		return &componentGenerator
	}

	return nil
}

// GenerateDomain is the CLI entry point for generating generic REST domains
func GenerateDomain(_ []string) (err error) {
	inflector := skel.NewInflector(generatorConfig.Domain)
	tagName := inflector.Inflect(skel.InflectionUpperCamelSingular)
	spec, err := loadSpecification(tagName, inflector)
	if err != nil {
		return err
	}

	var generators []ScopedComponentGenerator
	for _, component := range generatorConfig.Components {
		generator := findGenerator(generatorConfig.Style, generatorConfig.Tenant, component)
		if generator != nil {
			generators = append(generators, *generator)
		} else {
			err = errors.Errorf("Failed to identify generator for style=%s tenant=%s component=%s",
				generatorConfig.Style, generatorConfig.Tenant, component)
			// return err
			logger.WithError(err).Error("Could not generate component")
		}
	}

	for _, scopedGenerator := range generators {
		generator := scopedGenerator.Factory(spec)

		if err = generator.Generate(); err != nil {
			return err
		}

		template := skel.Template{
			Name:       scopedGenerator.Component,
			DestFile:   generator.Filename(),
			SourceData: []byte(generator.Render()),
			Format:     skel.FileFormatGo,
			Operation:  skel.OpAdd,
		}

		var opts = skel.NewRenderOptions()
		opts.AddVariables(generator.Variables())
		opts.AddConditions(generator.Conditions())
		if err = template.Render(opts); err != nil {
			return err
		}
	}

	return nil
}
