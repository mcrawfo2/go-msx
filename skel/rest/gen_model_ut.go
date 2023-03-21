// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package rest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel/text"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/mcrawfo2/go-jsonschema/pkg/codegen"
	"path"
)

type DomainModelUnitTestGenerator struct {
	Domain  string
	Folder  string
	Actions types.ComparableSlice[string]
	Models  []DomainModel
	Spec    Spec
	*text.GoFile
}

func (g DomainModelUnitTestGenerator) createFilterWhereTestSnippet() error {
	return g.AddNewText(
		"Filter/Where",
		"where",
		`
			func Test_lowerCamelSingularFilters_Where(t *testing.T) {
				tests := []struct {
					name string
					want sqldb.WhereOption
				}{
					{
						name: "Success",
						want: nil,
					},
				}
				for _, tt := range tests {
					t.Run(tt.name, func(t *testing.T) {
						f := lowerCamelSingularFilters{}
						assert.Equalf(t, tt.want, f.Where(), "Where()")
					})
				}
			}
		`,
		[]codegen.Import{
			text.ImportTesting,
			text.ImportSqldb,
			text.ImportTestifyAssert,
		})
}

func (g DomainModelUnitTestGenerator) Generate() error {
	errs := types.ErrorList{}

	for _, model := range g.Models {
		if !g.Actions.ContainsAny(model.Actions...) {
			continue
		}

		var err error
		switch model.Model {
		case DomainModelFilter:
			err = g.createFilterWhereTestSnippet()
		}

		errs = append(errs, err)
	}

	return errs.Filter()
}

func (g DomainModelUnitTestGenerator) Filename() string {
	target := path.Join(g.Folder, "model_lowersingular_test.go")
	return g.GoFile.Inflector.Inflect(target)
}

func NewDomainModelUnitTestGenerator(spec Spec) ComponentGenerator {
	return DomainModelUnitTestGenerator{
		Domain:  generatorConfig.Domain,
		Folder:  generatorConfig.Folder,
		Actions: generatorConfig.Actions,
		Models: []DomainModel{
			{
				Model:   DomainModelFilter,
				Actions: []string{ActionList},
			},
		},
		Spec: spec,
		GoFile: &text.GoFile{
			File: &text.File[text.GoSnippet]{
				Comment:   "Model Unit Test for " + generatorConfig.Domain,
				Inflector: text.NewInflector(generatorConfig.Domain),
				Sections: text.Sections[text.GoSnippet]{{
					Name: "Filter",
					Sections: text.NewGoSections(
						"Where",
					),
				}},
			},
			Package: generatorConfig.PackageName(),
		},
	}
}
