// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package rest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
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
	*File
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
			importSqldb,
		})
}

func (g DomainModelGenerator) createInstanceSnippet() error {
	return g.AddNewDecl(
		"Instance",
		"instance",
		&codegen.TypeDecl{
			Name: "UpperCamelSingular",
			Type: &codegen.StructType{
				Fields: []codegen.StructField{
					{
						Name: "UpperCamelSingularId",
						Type: codegen.NamedType{
							Package: &codegen.Package{
								QualifiedName: PkgUuid,
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
		},
		[]codegen.Import{
			importUuid,
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
	return g.File.Inflector.Inflect(target)
}

func (g DomainModelGenerator) Variables() map[string]string {
	return nil
}

func (g DomainModelGenerator) Conditions() map[string]bool {
	return nil
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
		File: &File{
			Comment:   "Model for " + generatorConfig.Domain,
			Package:   generatorConfig.PackageName(),
			Inflector: skel.NewInflector(generatorConfig.Domain),
			Sections: NewSections(
				"Filter",
				"Instance",
			),
		},
	}
}
