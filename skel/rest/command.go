// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package rest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	ComponentMigration  = "migration"

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
		ComponentMigration,
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

func (c GeneratorConfig) Apply(opts skel.RenderOptions) skel.RenderOptions {
	opts.AddConditions(map[string]bool{
		"TENANT_DOMAIN": types.ComparableSlice[string]{TenantSingle, TenantHierarchy}.Contains(c.Tenant),
	})
	return opts
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

type ActionFunc func() error

type ComponentGenerator interface {
	Generate() error
	Render() string
	Filename() string
	FileFormat() skel.FileFormat
}

type RenderOptionsTransformer interface {
	Apply(options skel.RenderOptions) skel.RenderOptions
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
		Style:     MatchesAny,
		Tenant:    MatchesAny,
		Component: ComponentController,
		Factory:   NewDomainControllerGenerator,
	},
	{
		Style:     MatchesAny,
		Tenant:    TenantNone,
		Component: ComponentService,
		Factory:   NewDomainServiceGenerator,
	},
	{
		Style:     MatchesAny,
		Tenant:    TenantNone,
		Component: ComponentConverter,
		Factory:   NewDomainConverterGenerator,
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
	{
		Style:     MatchesAny,
		Tenant:    MatchesAny,
		Component: ComponentMigration,
		Factory:   NewDomainMigrationGenerator,
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
		filename := generator.Filename()
		logger.Infof("📎 %s", cases.Title(language.English).String(scopedGenerator.Component))

		logger.Infof("  📖 Generating %s", path.Base(filename))

		if err = generator.Generate(); err != nil {
			return err
		}

		template := skel.Template{
			Name:       cases.Title(language.English).String(scopedGenerator.Component),
			DestFile:   filename,
			SourceData: []byte(generator.Render()),
			Format:     generator.FileFormat(),
			Operation:  skel.OpAdd,
		}

		var opts = skel.NewRenderOptions()
		opts = generatorConfig.Apply(opts)
		if optionsSource, ok := generator.(RenderOptionsTransformer); ok {
			opts = optionsSource.Apply(opts)
		}

		if err = template.Render(opts); err != nil {
			return err
		}
	}

	var actions = []ActionFunc{
		func() error {
			return skel.InitializePackageFromFile(
				path.Join(skel.Config().TargetDirectory(), "cmd", "app", "main.go"),
				path.Join(skel.Config().AppPackageUrl(), generatorConfig.Folder))
		},
	}

	for _, action := range actions {
		if err = action(); err != nil {
			return err
		}
	}

	return
}
