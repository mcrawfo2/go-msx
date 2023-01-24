// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package rest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/text"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/mcrawfo2/go-jsonschema/pkg/codegen"
	"path"
)

type DomainModel struct {
	Model   string
	Actions []string
}

const (
	DomainModelFilter   = "Filter"
	DomainModelInstance = "Instance"
)

type DomainModelGenerator struct {
	Domain  string
	Folder  string
	Actions types.ComparableSlice[string]
	Models  []DomainModel
	Spec    Spec
	*text.GoFile
}

func (g DomainModelGenerator) createFilterSnippet() error {
	return g.AddNewText(
		"Filter",
		"filters",
		`
			type petFilters struct{}
			
			func (f petFilters) Where() sqldb.WhereOption {
				return nil
			}
			`,
		[]codegen.Import{
			text.ImportSqldb,
		})
}

func (g DomainModelGenerator) createInstanceSnippet() error {
	return g.AddNewGenerator(
		"Instance",
		"instance",
		text.Decls{&codegen.TypeDecl{
			Name: "UpperCamelSingular",
			Type: &codegen.StructType{
				Fields: []codegen.StructField{
					{
						Name: "UpperCamelSingularId",
						Type: codegen.NamedType{
							Package: &codegen.Package{
								QualifiedName: text.PkgUuid,
							},
							Decl: &codegen.TypeDecl{
								Name: "UUID",
							},
						},
						Tags: codegen.StructTags{{
							Name:  "db",
							Value: "lower_snake_singular_id",
						}},
					},
					{
						Name: "Data",
						Type: codegen.PrimitiveType{Type: "string"},
						Tags: codegen.StructTags{{
							Name:  "db",
							Value: "data",
						}},
					},
				},
			},
		}},
		[]codegen.Import{
			text.ImportUuid,
		})
}

func (g DomainModelGenerator) Generate() error {
	errs := types.ErrorList{}

	for _, model := range g.Models {
		if !g.Actions.ContainsAny(model.Actions...) {
			continue
		}

		var err error
		switch model.Model {
		case DomainModelFilter:
			err = g.createFilterSnippet()
		case DomainModelInstance:
			err = g.createInstanceSnippet()
		}

		errs = append(errs, err)
	}

	return errs.Filter()
}

func (g DomainModelGenerator) Filename() string {
	target := path.Join(g.Folder, "model_lowersingular.go")
	return g.GoFile.Inflector.Inflect(target)
}

func NewDomainModelGenerator(spec Spec) ComponentGenerator {
	return DomainModelGenerator{
		Domain:  generatorConfig.Domain,
		Folder:  generatorConfig.Folder,
		Actions: generatorConfig.Actions,
		Models: []DomainModel{
			{
				Model:   DomainModelFilter,
				Actions: []string{ActionList},
			},
			{
				Model:   DomainModelInstance,
				Actions: []string{ActionList, ActionRetrieve, ActionCreate, ActionUpdate, ActionDelete},
			},
		},
		Spec: spec,
		GoFile: &text.GoFile{
			File: &text.File[text.GoSnippet]{
				Comment:   "Model for " + generatorConfig.Domain,
				Inflector: skel.NewInflector(generatorConfig.Domain),
				Sections: text.NewGoSections(
					"Filter",
					"Instance",
				),
			},
			Package: generatorConfig.PackageName(),
		},
	}
}
