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

type DomainConverterGeneratorV8 struct {
	Domain     string
	Folder     string
	Actions    types.ComparableSlice[string]
	Converters []DomainConversion
	Spec       Spec
	*File
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

func (g DomainConverterGeneratorV8) createServiceSnippet() error {
	return g.AddNewText(
		"Service",
		"implementation",
		`
			// lowerCamelSingularConverter translates between API requests/responses and models.
			type lowerCamelSingularConverter struct{}
			`,
		[]codegen.Import{})
}

func (g DomainConverterGeneratorV8) createConverterFilterQuerySnippet() error {
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

func (g DomainConverterGeneratorV8) createConverterCreateRequestSnippet() error {
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
			importUuid,
		})
}

func (g DomainConverterGeneratorV8) createConverterUpdateRequestSnippet() error {
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
			importUuid,
		})
}

func (g DomainConverterGeneratorV8) createConverterListResponseSnippet() error {
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

func (g DomainConverterGeneratorV8) createConverterResponseSnippet() error {
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
			importPrepared,
		})
}

func (g DomainConverterGeneratorV8) createConverterSortByOptionsSnippet() error {
	return g.AddNewText(
		"Sort",
		"sortByOptions",
		`
			var lowerCamelSingularSortByOptions = types.StringPairSlice{
				{
					Left:  "lowerCamelSingularId",
					Right: "lower_snake_singular_id",
				},
			}
			`,
		[]codegen.Import{
			importTypes,
		})
}

func (g DomainConverterGeneratorV8) Generate() error {
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

func (g DomainConverterGeneratorV8) Filename() string {
	target := path.Join(g.Folder, "converter_lowersingular_v8.go")
	return g.File.Inflector.Inflect(target)
}

func (g DomainConverterGeneratorV8) Variables() map[string]string {
	return nil
}

func (g DomainConverterGeneratorV8) Conditions() map[string]bool {
	return nil
}

func NewDomainConverterGeneratorV8(spec Spec) ComponentGenerator {
	return DomainConverterGeneratorV8{
		Domain:  generatorConfig.Domain,
		Folder:  generatorConfig.Folder,
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
		File: &File{
			Comment:   "V8 API Converter for " + generatorConfig.Domain,
			Package:   generatorConfig.PackageName(),
			Inflector: skel.NewInflector(generatorConfig.Domain),
			Sections: NewSections(
				"Service",
				&Section{
					Name: "Converters",
					Sections: NewSections(
						"List Query Filter",
						"Create",
						"Update",
						"List Response",
						"Response"),
				},
				"Sort",
			),
		},
	}
}
