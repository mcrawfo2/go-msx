// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package rest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel/text"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/mcrawfo2/go-jsonschema/pkg/codegen"
	"github.com/mcrawfo2/jennifer/jen"
	"path"
)

type DomainGlobalsGenerator struct {
	Domain string
	Folder string
	*text.GoFile
}

func (g DomainGlobalsGenerator) createLoggerSnippet() error {
	return g.AddNewStatement(
		"Logger",
		"logger",
		jen.Var().Id("logger").Op("=").Qual(text.PkgLog, "NewPackageLogger").Call())
}

func (g DomainGlobalsGenerator) createContextKeyTypeSnippet() error {
	return g.AddNewGenerator(
		"Context",
		"contextKeyNamed",
		text.Decls{&codegen.TypeDecl{
			Type: codegen.PrimitiveType{
				Type: "string",
			},
			Name: "contextKeyNamed",
		}},
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
	return g.GoFile.Inflector.Inflect(target)
}

func NewDomainGlobalsGenerator(_ Spec) ComponentGenerator {
	return DomainGlobalsGenerator{
		Domain: generatorConfig.Domain,
		Folder: generatorConfig.Folder,
		GoFile: &text.GoFile{
			File: &text.File[text.GoSnippet]{
				Comment:   "Globals for " + generatorConfig.Domain,
				Inflector: text.NewInflector(generatorConfig.Domain),
				Sections: text.NewGoSections(
					"Logger",
					"Context",
				),
			},
			Package: generatorConfig.PackageName(),
		},
	}
}
