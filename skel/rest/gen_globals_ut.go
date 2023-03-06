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
	"github.com/mcrawfo2/jennifer/jen"
	"path"
)

type DomainData struct {
	Data    string
	Actions []string
}

const (
	DomainDataId            = "Id"
	DomainDataInstance      = "Instance"
	DomainDataCreateRequest = "CreateRequest"
	DomainDataUpdateRequest = "UpdateRequest"
	DomainDataInstances     = "Instances"
	DomainDataFilter        = "Filter"
	DomainDataPagingReq     = "PagingReq"
	DomainDataPagingResp    = "PagingResp"
)

type DomainGlobalsUnitTestGenerator struct {
	Domain  string
	Folder  string
	Style   string
	Actions types.ComparableSlice[string]
	Data    []DomainData
	Spec    Spec
	*text.GoFile
}

func (g DomainGlobalsUnitTestGenerator) createUniqueIdSnippet() error {
	return g.AddNewStatement(
		"UniqueId",
		"id",
		jen.Var().Id("lowerCamelSingularTestUpperCamelSingularId").Op("=").Qual(text.PkgTypes, "MustNewUUID").Call())
}

func (g DomainGlobalsUnitTestGenerator) createModelTestDataStructSnippet() error {
	st := new(codegen.StructType)
	imports := map[string]struct{}{}

	ucs := g.Inflector.Inflect("UpperCamelSingular")
	ucp := g.Inflector.Inflect("UpperCamelPlural")
	lcs := g.Inflector.Inflect("lowerCamelSingular")

	for _, data := range g.Data {
		if !g.Actions.ContainsAny(data.Actions...) {
			continue
		}

		switch data.Data {
		case DomainDataId:
			st.AddField(codegen.StructField{
				Name: ucs + "Id",
				Type: codegen.CustomNameType{Type: "uuid.UUID"},
			})
			imports[text.PkgUuid] = struct{}{}
		case DomainDataInstance:
			st.AddField(codegen.StructField{
				Name: ucs,
				Type: codegen.CustomNameType{Type: ucs},
			})
		case DomainDataCreateRequest, DomainDataUpdateRequest:
			// nothing for the model
		case DomainDataInstances:
			st.AddField(codegen.StructField{
				Name: ucp,
				Type: codegen.CustomNameType{Type: "[]" + ucs},
			})
		case DomainDataFilter:
			st.AddField(codegen.StructField{
				Name: ucs + "Filters",
				Type: codegen.CustomNameType{Type: lcs + "Filters"},
			})
		case DomainDataPagingReq:
			st.AddField(codegen.StructField{
				Name: "PagingRequest",
				Type: codegen.CustomNameType{Type: "paging.Request"},
			})
			imports[text.PkgPaging] = struct{}{}
		case DomainDataPagingResp:
			st.AddField(codegen.StructField{
				Name: "PagingResponse",
				Type: codegen.CustomNameType{Type: "paging.Response"},
			})
			imports[text.PkgPaging] = struct{}{}
		}
	}

	namedStructType := &codegen.TypeDecl{
		Name: lcs + "TestModelData",
		Type: st,
	}

	return g.AddNewDecl(
		"Model/Structure",
		"structure",
		namedStructType,
		g.uniqueImports(imports))
}

func (g DomainGlobalsUnitTestGenerator) createModelTestDataConstructorSnippet() error {
	ucs := g.Inflector.Inflect("UpperCamelSingular")
	ucp := g.Inflector.Inflect("UpperCamelPlural")
	lcs := g.Inflector.Inflect("lowerCamelSingular")
	lss := g.Inflector.Inflect("lower_snake_singular")

	var setters = []jen.Code{
		jen.Id(lcs+"Id").Op(":=").Qual(text.PkgPrepared, "ToModelUuid").Call(
			jen.Id(fmt.Sprintf("%sTest%sId", lcs, ucs))),
		jen.Id("data").Op(":=").Lit("data"),
		jen.Line(),
		jen.Id("result").Op(":=").Id(lcs + "TestModelData").Values(),
		jen.Line(),
	}

	var imports = map[string]struct{}{
		text.PkgPrepared: {},
	}

	for _, data := range g.Data {
		if !g.Actions.ContainsAny(data.Actions...) {
			continue
		}

		var setter *jen.Statement

		switch data.Data {
		case DomainDataId:
			setter = jen.Id("result").Dot(ucs + "Id").Op("=").Id(lcs + "Id")
		case DomainDataInstance:
			setter = jen.Id("result").Dot(ucs).Op("=").Id(ucs).Values(jen.Dict{
				jen.Id(ucs + "Id"): jen.Id(lcs + "Id"),
				jen.Id("Data"):     jen.Id("data"),
			})
		case DomainDataCreateRequest, DomainDataUpdateRequest:
			// nothing for the model
		case DomainDataInstances:
			setter = jen.Id("result").Dot(ucp).Op("=").Op("[]").Id(ucs).Values(
				jen.Values(jen.Dict{
					jen.Id(ucs + "Id"): jen.Id(lcs + "Id"),
					jen.Id("Data"):     jen.Id("data"),
				}),
				jen.Values(jen.Dict{
					jen.Id(ucs + "Id"): jen.Id(lcs + "Id"),
					jen.Id("Data"):     jen.Id("data"),
				}),
			)
		case DomainDataFilter:
			setter = jen.Id("result").Dot(ucs + "Filters").Op("=").Id(lcs + "Filters").Values()
		case DomainDataPagingReq:
			setter = jen.Id("result").Dot("PagingRequest").Op("=").Qual(text.PkgPaging, "Request").Values(jen.Dict{
				jen.Id("Size"): jen.Lit(10),
				jen.Id("Sort"): jen.Op("[]").Qual(text.PkgPaging, "SortOrder").Values(jen.Values(jen.Dict{
					jen.Id("Property"):  jen.Lit(lss + "_id"),
					jen.Id("Direction"): jen.Qual(text.PkgPaging, "SortDirectionAsc"),
				})),
			})
		case DomainDataPagingResp:
			setter = jen.Id("result").Dot("PagingResponse").Op("=").Qual(text.PkgPaging, "Response").Values(jen.Dict{
				jen.Id("Size"):       jen.Id("uint").Call(jen.Lit(10)),
				jen.Id("TotalItems"): jen.Qual(text.PkgTypes, "PtrTo").Types(jen.Id("uint")).Call(jen.Lit(2)),
				jen.Id("Sort"): jen.Op("[]").Qual(text.PkgPaging, "SortOrder").Values(jen.Values(jen.Dict{
					jen.Id("Property"):  jen.Lit(lss + "_id"),
					jen.Id("Direction"): jen.Qual(text.PkgPaging, "SortDirectionAsc"),
				})),
			})
			imports[text.PkgTypes] = struct{}{}
		}

		if setter != nil {
			setters = append(setters, setter)
		}
	}

	setters = append(setters,
		jen.Line(),
		jen.Return(jen.Id("result")))

	ctor := jen.Func().Id(fmt.Sprintf("new%sTestModelData", ucs)).Params().Id(lcs + "TestModelData").Block(setters...)

	return g.AddNewStatement(
		"Model/Constructor",
		"constructor",
		ctor)
}

func (g DomainGlobalsUnitTestGenerator) createApiTestDataStructSnippet() error {
	st := new(codegen.StructType)
	imports := map[string]struct{}{}

	ucs := g.Inflector.Inflect("UpperCamelSingular")
	lcs := g.Inflector.Inflect("lowerCamelSingular")

	for _, data := range g.Data {
		if !g.Actions.ContainsAny(data.Actions...) {
			continue
		}

		switch data.Data {
		case DomainDataId:
			st.AddField(codegen.StructField{
				Name: ucs + "Id",
				Type: codegen.CustomNameType{Type: "types.UUID"},
			})
			imports[text.PkgTypes] = struct{}{}
		case DomainDataInstance:
			st.AddField(codegen.StructField{
				Name: ucs + "Response",
				Type: codegen.CustomNameType{Type: ucs + "Response"},
			})
		case DomainDataCreateRequest:
			st.AddField(codegen.StructField{
				Name: ucs + "CreateRequest",
				Type: codegen.CustomNameType{Type: ucs + "CreateRequest"},
			})
		case DomainDataUpdateRequest:
			st.AddField(codegen.StructField{
				Name: ucs + "UpdateRequest",
				Type: codegen.CustomNameType{Type: ucs + "UpdateRequest"},
			})
		case DomainDataInstances:
			st.AddField(codegen.StructField{
				Name: ucs + "Responses",
				Type: codegen.CustomNameType{Type: "[]" + ucs + "Response"},
			})
		case DomainDataFilter:
			// nothing for the api
		case DomainDataPagingReq:
			st.AddField(codegen.StructField{
				Name: "PagingRequest",
				Type: codegen.CustomNameType{Type: g.Style + ".PagingSortingInputs"},
			})
			imports[g.pkgStyle()] = struct{}{}
		case DomainDataPagingResp:
			st.AddField(codegen.StructField{
				Name: "PagingResponse",
				Type: codegen.CustomNameType{Type: g.Style + ".PagingResponse"},
			})
			imports[g.pkgStyle()] = struct{}{}
		}
	}

	namedStructType := &codegen.TypeDecl{
		Name: lcs + "TestApiData",
		Type: st,
	}

	var uniqueImports []codegen.Import
	for k := range imports {
		uniqueImports = append(
			uniqueImports,
			codegen.Import{
				QualifiedName: k,
			})
	}

	return g.AddNewDecl(
		"API/Structure",
		"structure",
		namedStructType,
		uniqueImports)
}

func (g DomainGlobalsUnitTestGenerator) createApiTestDataConstructorSnippet() error {
	ucs := g.Inflector.Inflect("UpperCamelSingular")
	lcs := g.Inflector.Inflect("lowerCamelSingular")

	var setters = []jen.Code{
		jen.Id(lcs + "Id").Op(":=").Id(fmt.Sprintf("%sTest%sId", lcs, ucs)),
		jen.Id("data").Op(":=").Lit("data"),
		jen.Line(),
		jen.Id("result").Op(":=").Id(lcs + "TestApiData").Values(),
		jen.Line(),
	}

	var imports = map[string]struct{}{}

	var pkgStyle = g.pkgStyle()

	for _, data := range g.Data {
		if !g.Actions.ContainsAny(data.Actions...) {
			continue
		}

		var setter *jen.Statement

		switch data.Data {
		case DomainDataId:
			setter = jen.Id("result").Dot(ucs + "Id").Op("=").Id(lcs + "Id")
		case DomainDataInstance:
			setter = jen.Id("result").Dot(ucs + "Response").Op("=").Id(ucs + "Response").Values(jen.Dict{
				jen.Id(ucs + "Id"): jen.Id(lcs + "Id"),
				jen.Id("Data"):     jen.Id("data"),
			})
		case DomainDataCreateRequest:
			setter = jen.Id("result").Dot(ucs + "CreateRequest").Op("=").
				Id(ucs + "CreateRequest").Values(jen.Dict{
				jen.Id("Data"): jen.Id("data"),
			})
		case DomainDataUpdateRequest:
			setter = jen.Id("result").Dot(ucs + "UpdateRequest").Op("=").
				Id(ucs + "UpdateRequest").Values(jen.Dict{
				jen.Id("Data"): jen.Id("data"),
			})
		case DomainDataInstances:
			setter = jen.Id("result").Dot(ucs+"Responses").Op("=").Op("[]").Id(ucs+"Response").Values(
				jen.Values(jen.Dict{
					jen.Id(ucs + "Id"): jen.Id(lcs + "Id"),
					jen.Id("Data"):     jen.Id("data"),
				}),
				jen.Values(jen.Dict{
					jen.Id(ucs + "Id"): jen.Id(lcs + "Id"),
					jen.Id("Data"):     jen.Id("data"),
				}),
			)
		case DomainDataFilter:
			// nothing for API
		case DomainDataPagingReq:
			setter = jen.Id("result").Dot("PagingRequest").Op("=").Qual(pkgStyle, "PagingSortingInputs").Values(jen.Dict{
				jen.Id("PagingInputs"): jen.Qual(pkgStyle, "PagingInputs").Values(jen.Dict{
					jen.Id("PageSize"): jen.Lit(10),
				}),
				jen.Id("SortingInputs"): jen.Qual(pkgStyle, "SortingInputs").Values(jen.Dict{
					jen.Id("SortBy"):    jen.Lit(lcs + "Id"),
					jen.Id("SortOrder"): jen.Lit("asc"),
				}),
			})
		case DomainDataPagingResp:
			var pagingResponseFields jen.Dict
			if g.Style == StyleV2 {
				pagingResponseFields = jen.Dict{
					jen.Id("Size"): jen.Lit(10),
					jen.Id("Pageable"): jen.Qual(pkgStyle, "PageableResponse").Values(jen.Dict{
						jen.Id("Size"): jen.Lit(10),
						jen.Id("Sort"): jen.Qual(pkgStyle, "SortResponse").Values(jen.Dict{
							jen.Id("Orders"): jen.Op("[]").Qual(pkgStyle, "SortOrderResponse").Values(
								jen.Values(jen.Dict{
									jen.Id("Property"):  jen.Lit(lcs + "Id"),
									jen.Id("Direction"): jen.Qual(pkgStyle, "SortDirectionAsc"),
								})),
						}),
					}),
				}
			} else {
				pagingResponseFields = jen.Dict{
					jen.Id("PageSize"):   jen.Lit(10),
					jen.Id("TotalItems"): jen.Qual(text.PkgTypes, "PtrTo").Types(jen.Id("int")).Call(jen.Lit(2)),
					jen.Id("SortBy"):     jen.Lit(lcs + "Id"),
					jen.Id("SortOrder"):  jen.Qual("strings", "ToLower").Call(jen.Qual(pkgStyle, "SortDirectionAsc")),
				}
			}

			setter = jen.Id("result").Dot("PagingResponse").Op("=").Qual(pkgStyle, "PagingResponse").Values(pagingResponseFields)
			imports[text.PkgTypes] = struct{}{}
		}

		if setter != nil {
			setters = append(setters, setter)
		}
	}

	setters = append(setters,
		jen.Line(),
		jen.Return(jen.Id("result")))

	ctor := jen.Func().Id(fmt.Sprintf("new%sTestApiData", ucs)).Params().Id(lcs + "TestApiData").Block(setters...)

	return g.AddNewStatement(
		"API/Constructor",
		"constructor",
		ctor)
}

func (g DomainGlobalsUnitTestGenerator) uniqueImports(imports map[string]struct{}) []codegen.Import {
	var uniqueImports []codegen.Import
	for k := range imports {
		uniqueImports = append(
			uniqueImports,
			codegen.Import{
				QualifiedName: k,
			})
	}
	return uniqueImports
}

func (g DomainGlobalsUnitTestGenerator) Generate() error {
	return types.ErrorList{
		g.createUniqueIdSnippet(),
		g.createModelTestDataStructSnippet(),
		g.createModelTestDataConstructorSnippet(),
		g.createApiTestDataStructSnippet(),
		g.createApiTestDataConstructorSnippet(),
	}.Filter()
}

func (g DomainGlobalsUnitTestGenerator) Filename() string {
	target := path.Join(g.Folder, "data_lowersingular_test.go")
	return g.GoFile.Inflector.Inflect(target)
}

func (g DomainGlobalsUnitTestGenerator) pkgStyle() string {
	switch g.Style {
	case StyleV2:
		return text.PkgRestopsV2
	default:
		return text.PkgRestopsV8
	}
}

func NewDomainGlobalsUnitTestGenerator(spec Spec) ComponentGenerator {
	return DomainGlobalsUnitTestGenerator{
		Domain:  generatorConfig.Domain,
		Folder:  generatorConfig.Folder,
		Style:   generatorConfig.Style,
		Actions: generatorConfig.Actions,
		Data: []DomainData{
			{
				Data:    DomainDataInstances,
				Actions: []string{ActionList},
			},
			{
				Data:    DomainDataFilter,
				Actions: []string{ActionList},
			},
			{
				Data:    DomainDataPagingReq,
				Actions: []string{ActionList},
			},
			{
				Data:    DomainDataPagingResp,
				Actions: []string{ActionList},
			},
			{
				Data:    DomainDataId,
				Actions: []string{ActionRetrieve, ActionUpdate, ActionDelete},
			},
			{
				Data:    DomainDataInstance,
				Actions: []string{ActionRetrieve, ActionCreate, ActionUpdate},
			},
			{
				Data:    DomainDataCreateRequest,
				Actions: []string{ActionCreate},
			},
			{
				Data:    DomainDataUpdateRequest,
				Actions: []string{ActionUpdate},
			},
		},
		Spec: spec,
		GoFile: &text.GoFile{
			File: &text.File[text.GoSnippet]{
				Comment:   "Unit Test Data for " + generatorConfig.Domain,
				Inflector: skel.NewInflector(generatorConfig.Domain),
				Sections: text.Sections[text.GoSnippet]{
					{
						Name: "UniqueId",
					},
					{
						Name: "Model",
						Sections: text.NewGoSections(
							"Structure",
							"Constructor",
						),
					},
					{
						Name: "API",
						Sections: text.NewGoSections(
							"Structure",
							"Constructor",
						),
					},
				},
			},
			Package: generatorConfig.PackageName(),
		},
	}
}
