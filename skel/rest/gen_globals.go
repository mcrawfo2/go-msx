package rest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/mcrawfo2/go-jsonschema/pkg/codegen"
	"github.com/mcrawfo2/jennifer/jen"
	"path"
)

type DomainGlobalsGenerator struct {
	Domain string
	Folder string
	*File
}

func (g DomainGlobalsGenerator) createLoggerSnippet() error {
	return g.AddNewStatement(
		"Logger",
		"logger",
		jen.Var().Id("logger").Op("=").Qual(PkgLog, "NewPackageLogger").Call())
}

func (g DomainGlobalsGenerator) createContextKeyTypeSnippet() error {
	return g.AddNewDecl(
		"Context",
		"contextKeyNamed",
		&codegen.TypeDecl{
			Type: codegen.PrimitiveType{
				Type: "string",
			},
			Name: "contextKeyNamed",
		},
		nil)
}

func (g DomainGlobalsGenerator) Generate() error {
	return types.ErrorList{
		g.createLoggerSnippet(),
		g.createContextKeyTypeSnippet(),
	}.Filter()
}

func (g DomainGlobalsGenerator) Filename() string {
	target := path.Join(g.Folder, "pkg.go")
	return g.File.Inflector.Inflect(target)
}

func (g DomainGlobalsGenerator) Variables() map[string]string {
	return nil
}

func (g DomainGlobalsGenerator) Conditions() map[string]bool {
	return nil
}

func NewDomainGlobalsGenerator(_ Spec) ComponentGenerator {
	return DomainGlobalsGenerator{
		Domain: generatorConfig.Domain,
		Folder: generatorConfig.Folder,
		File: &File{
			Comment:   "Globals for " + generatorConfig.Domain,
			Package:   generatorConfig.PackageName(),
			Inflector: skel.NewInflector(generatorConfig.Domain),
			Sections: NewSections(
				"Logger",
				"Context",
			),
		},
	}
}
