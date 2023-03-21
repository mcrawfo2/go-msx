// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package rest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/text"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
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
	UnitTests  bool
}

func (c GeneratorConfig) Validate() error {
	return types.ErrorMap{
		"domain": validation.Validate(&c.Domain, validation.Required,
			validation.NewStringRule(
				regexp.MustCompile(`^[A-Za-z][A-Za-z ]+$`).MatchString,
				"Domain must contain only letters or spaces, and start with a letter")),
		"folder": validation.Validate(&c.Folder, validation.Required,
			validation.NewStringRule(
				regexp.MustCompile(`^^[A-Za-z0-9]+(/[A-Za-z0-9]+)*/?$`).MatchString,
				"Folder must contain 1 or more components of only letters and/or digits, separated by slashes")),
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
	cmd.Flags().BoolVar(&generatorConfig.UnitTests, "tests", true, "Generate unit tests")

	cmd = skel.AddTarget("disinflect", "Reverse inflections", Disinflect)
	cmd.Args = cobra.MinimumNArgs(2)
	cmd.Hidden = true
}

func generateDomainOptions(_ *cobra.Command, args []string) error {
	generatorConfig.Domain = strings.TrimSpace(strings.Join(args, " "))

	if len(generatorConfig.Actions) == 0 {
		generatorConfig.Actions = ActionOptions
	}

	if len(generatorConfig.Components) == 0 {
		generatorConfig.Components = ComponentOptions
	}

	generatorConfig.Folder = text.NewInflector(generatorConfig.Domain).Inflect(generatorConfig.Folder)

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
	FileFormat() text.FileFormat
}

type RenderOptionsTransformer interface {
	Apply(options skel.RenderOptions) skel.RenderOptions
}

type ScopedComponentGenerator struct {
	Component string
	UnitTest  bool
	Factory   func(spec Spec) ComponentGenerator
}

var componentGenerators = []ScopedComponentGenerator{
	{
		Component: ComponentGlobals,
		Factory:   NewDomainGlobalsGenerator,
	},
	{
		Component: ComponentGlobals,
		Factory:   NewDomainGlobalsUnitTestGenerator,
		UnitTest:  true,
	},
	{
		Component: ComponentPayloads,
		Factory:   NewDomainPayloadsGenerator,
	},
	{
		Component: ComponentController,
		Factory:   NewDomainControllerGenerator,
	},
	{
		Component: ComponentController,
		Factory:   NewDomainControllerUnitTestGenerator,
		UnitTest:  true,
	},
	{
		Component: ComponentService,
		Factory:   NewDomainServiceGenerator,
	},
	{
		Component: ComponentService,
		Factory:   NewDomainServiceUnitTestGenerator,
		UnitTest:  true,
	},
	{
		Component: ComponentConverter,
		Factory:   NewDomainConverterGenerator,
	},
	{
		Component: ComponentModel,
		Factory:   NewDomainModelGenerator,
	},
	{
		Component: ComponentModel,
		Factory:   NewDomainModelUnitTestGenerator,
		UnitTest:  true,
	},
	{
		Component: ComponentRepository,
		Factory:   NewDomainRepositoryGenerator,
	},
	{
		Component: ComponentRepository,
		Factory:   NewDomainRepositoryUnitTestGenerator,
		UnitTest:  true,
	},
	{
		Component: ComponentMigration,
		Factory:   NewDomainMigrationGenerator,
	},
}

func findGenerator(component string, ut bool) *ScopedComponentGenerator {
	for _, componentGenerator := range componentGenerators {
		if componentGenerator.UnitTest != ut {
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
	inflector := text.NewInflector(generatorConfig.Domain)
	tagName := inflector.Inflect(text.InflectionUpperCamelSingular)
	spec, err := loadSpecification(tagName, inflector)
	if err != nil {
		return err
	}

	var generators []ScopedComponentGenerator
	for _, component := range generatorConfig.Components {
		generator := findGenerator(component, false)
		if generator != nil {
			generators = append(generators, *generator)
			utGenerator := findGenerator(component, true)
			if utGenerator != nil && generatorConfig.UnitTests {
				generators = append(generators, *utGenerator)
			}
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
			// initialize the generated package from main
			return skel.InitializePackageFromFile(
				path.Join(skel.Config().TargetDirectory(), "cmd", "app", "main.go"),
				path.Join(skel.Config().AppPackageUrl(), generatorConfig.Folder))
		},
		func() error {
			// generate mocks
			return skel.GoGenerate(
				path.Join(skel.Config().TargetDirectory(), generatorConfig.Folder))
		},
	}

	for _, action := range actions {
		if err = action(); err != nil {
			return err
		}
	}

	return
}

func Disinflect(args []string) error {
	sourceFile := args[0]
	domain := strings.TrimSpace(strings.Join(args[1:], " "))

	inflector := text.NewInflector(domain)

	sourceBytes, err := os.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	inflections := types.StringPairSlice{
		{
			Left:  `\b` + inflector[text.InflectionLowerCamelPlural],
			Right: text.InflectionLowerCamelPlural,
		},
		{
			Left:  `_` + inflector[text.InflectionLowerCamelPlural],
			Right: `_` + text.InflectionLowerCamelPlural,
		},
		{
			Left:  inflector[text.InflectionUpperCamelPlural],
			Right: text.InflectionUpperCamelPlural,
		},
		{
			Left:  `\b` + inflector[text.InflectionLowerCamelSingular],
			Right: text.InflectionLowerCamelSingular,
		},
		{
			Left:  `_` + inflector[text.InflectionLowerCamelSingular],
			Right: `_` + text.InflectionLowerCamelSingular,
		},
		{
			Left:  inflector[text.InflectionUpperCamelSingular],
			Right: text.InflectionUpperCamelSingular,
		},
		{
			Left:  `\b` + inflector[text.InflectionLowerSnakeSingular] + `_`,
			Right: text.InflectionLowerSnakeSingular + `_`,
		},
		{
			Left:  `\b` + inflector[text.InflectionScreamingSnakePlural] + `_`,
			Right: text.InflectionScreamingSnakePlural + `_`,
		},
		{
			Left:  `\b` + inflector[text.InflectionScreamingSnakeSingular] + `_`,
			Right: text.InflectionScreamingSnakeSingular + `_`,
		},
	}
	for _, inflection := range inflections {
		re := regexp.MustCompile(inflection.Left)
		sourceBytes = re.ReplaceAll(sourceBytes, []byte(inflection.Right))
	}

	_, err = fmt.Println(string(sourceBytes))
	return err
}
