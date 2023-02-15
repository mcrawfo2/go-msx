// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package rest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/text"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/mcrawfo2/go-jsonschema/pkg/codegen"
	"path"
)

type DomainConverterGenerator struct {
	Domain     string
	Folder     string
	Style      string
	Actions    types.ComparableSlice[string]
	Converters []DomainConversion
	Spec       Spec
	*text.GoFile
}

const (
	ConverterFilterQueryInputs = "FromFilterQueryInputs"
	ConverterCreateRequest     = "FromCreateRequest"
	ConverterUpdateRequest     = "FromUpdateRequest"
	ConverterListResponse      = "ToListResponse"
	ConverterResponse          = "ToResponse"
	ConverterSortOptions       = "SortOptions"
)

type DomainConversion struct {
	Converter string
	Actions   types.ComparableSlice[string]
}

func (g DomainConverterGenerator) createServiceSnippet() error {
	return g.AddNewText(
		"Service",
		"implementation",
		`
			// lowerCamelSingularConverter translates between API requests/responses and models.
			type lowerCamelSingularConverter struct{}
			`,
		[]codegen.Import{})
}

func (g DomainConverterGenerator) createConverterFilterQuerySnippet() error {
	return g.AddNewText(
		"Converters/List Query Filter",
		"filter",
		`
			// FromFilterQueryInputs maps query fields to a database filter
			func (c lowerCamelSingularConverter) FromFilterQueryInputs(freq lowerCamelSingularFilterQueryInputs) lowerCamelSingularFilters {
				return lowerCamelSingularFilters{}
			}
			`,
		[]codegen.Import{})
}

func (g DomainConverterGenerator) createConverterCreateRequestSnippet() error {
	return g.AddNewText(
		"Converters/Create",
		"create",
		`
			// FromUpperCamelSingularCreateRequest creates a new UpperCamelSingular model from the values in the UpperCamelSingularCreateRequest payload
			func (c lowerCamelSingularConverter) FromUpperCamelSingularCreateRequest(request UpperCamelSingularCreateRequest) UpperCamelSingular {
				return UpperCamelSingular{
					UpperCamelSingularId: uuid.New(),
					Data:  request.Data,
				}
			}
			`,
		[]codegen.Import{
			text.ImportUuid,
		})
}

func (g DomainConverterGenerator) createConverterUpdateRequestSnippet() error {
	return g.AddNewText(
		"Converters/Update",
		"update",
		`
			// FromUpperCamelSingularUpdateRequest updates an existing UpperCamelSingular model from the values in the UpperCamelSingularUpdateRequest payload
			func (c lowerCamelSingularConverter) FromUpperCamelSingularUpdateRequest(lowerCamelSingular UpperCamelSingular, request UpperCamelSingularUpdateRequest) UpperCamelSingular {
				return UpperCamelSingular{
					UpperCamelSingularId: lowerCamelSingular.UpperCamelSingularId,
					Data:  request.Data,
				}
			}
			`,
		[]codegen.Import{
			text.ImportUuid,
		})
}

func (g DomainConverterGenerator) createConverterListResponseSnippet() error {
	return g.AddNewText(
		"Converters/List Response",
		"listResponse",
		`
			// ToUpperCamelSingularListResponse maps a series of UpperCamelSingular models onto a series of UpperCamelSingularResponse payloads
			func (c lowerCamelSingularConverter) ToUpperCamelSingularListResponse(lowerCamelPlural []UpperCamelSingular) (responses []UpperCamelSingularResponse) {
				for _, lowerCamelSingular := range lowerCamelPlural {
					responses = append(responses, c.ToUpperCamelSingularResponse(lowerCamelSingular))
				}
				return
			}
			`,
		[]codegen.Import{})
}

func (g DomainConverterGenerator) createConverterResponseSnippet() error {
	return g.AddNewText(
		"Converters/Response",
		"retrieveResponse",
		`
			// ToUpperCamelSingularResponse maps a single UpperCamelSingular model onto a UpperCamelSingularResponse payload
			func (c lowerCamelSingularConverter) ToUpperCamelSingularResponse(lowerCamelSingular UpperCamelSingular) UpperCamelSingularResponse {
				return UpperCamelSingularResponse{
					Data:  lowerCamelSingular.Data,
					UpperCamelSingularId: db.ToApiUuid(lowerCamelSingular.UpperCamelSingularId),
				}
			}
			`,
		[]codegen.Import{
			text.ImportPrepared,
		})
}

func (g DomainConverterGenerator) createConverterSortByOptionsSnippet() error {
	return g.AddNewText(
		"Sort",
		"sortByOptions",
		`
			var lowerCamelSingularSortByOptions = paging.SortByOptions{
				DefaultProperty: "lowerCamelSingularId",
				Mapping: types.StringPairSlice{
					{
						Left:  "lowerCamelSingularId",
						Right: "lower_snake_singular_id",
					},
				},
			}
			`,
		[]codegen.Import{
			text.ImportPaging,
			text.ImportTypes,
		})
}

func (g DomainConverterGenerator) Generate() error {
	errs := types.ErrorList{
		g.createServiceSnippet(),
	}

	for _, converter := range g.Converters {
		if !g.Actions.ContainsAny(converter.Actions...) {
			continue
		}

		var err error
		switch converter.Converter {
		case ConverterFilterQueryInputs:
			err = g.createConverterFilterQuerySnippet()
		case ConverterCreateRequest:
			err = g.createConverterCreateRequestSnippet()
		case ConverterUpdateRequest:
			err = g.createConverterUpdateRequestSnippet()
		case ConverterListResponse:
			err = g.createConverterListResponseSnippet()
		case ConverterResponse:
			err = g.createConverterResponseSnippet()
		case ConverterSortOptions:
			err = g.createConverterSortByOptionsSnippet()
		}

		errs = append(errs, err)
	}

	return errs.Filter()
}

func (g DomainConverterGenerator) Filename() string {
	target := path.Join(g.Folder, fmt.Sprintf("converter_lowersingular_%s.go", g.Style))
	return g.GoFile.Inflector.Inflect(target)
}

func NewDomainConverterGenerator(spec Spec) ComponentGenerator {
	return DomainConverterGenerator{
		Domain:  generatorConfig.Domain,
		Folder:  generatorConfig.Folder,
		Style:   generatorConfig.Style,
		Actions: generatorConfig.Actions,
		Converters: []DomainConversion{
			{
				Converter: ConverterFilterQueryInputs,
				Actions:   []string{ActionList},
			},
			{
				Converter: ConverterCreateRequest,
				Actions:   []string{ActionCreate},
			},
			{
				Converter: ConverterUpdateRequest,
				Actions:   []string{ActionUpdate},
			},
			{
				Converter: ConverterListResponse,
				Actions:   []string{ActionList},
			},
			{
				Converter: ConverterResponse,
				Actions:   []string{ActionList, ActionRetrieve, ActionCreate, ActionUpdate},
			},
			{
				Converter: ConverterSortOptions,
				Actions:   []string{ActionList},
			},
		},
		GoFile: &text.GoFile{
			File: &text.File[text.GoSnippet]{
				Comment:   "API Converter for " + generatorConfig.Domain,
				Inflector: skel.NewInflector(generatorConfig.Domain),
				Sections: text.NewGoSections(
					"Service",
					&text.Section[text.GoSnippet]{
						Name: "Converters",
						Sections: text.NewGoSections(
							"List Query Filter",
							"Create",
							"Update",
							"List Response",
							"Response"),
					},
					"Sort",
				),
			},
			Package: generatorConfig.PackageName(),
		},
	}
}
