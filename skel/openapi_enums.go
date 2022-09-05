// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package skel

import (
	"fmt"
	"path"

	"github.com/iancoleman/strcase"
	"github.com/mcrawfo2/jennifer/jen"
)

func generateScreamingSnakeName(s1 string, s2 string) string {
	return strcase.ToScreamingSnake(fmt.Sprintf("%s_%s", s1, s2))
}

func generateEnumConstants(schema Schema, file *jen.File) {
	var constCode []jen.Code
	idName := generateScreamingSnakeName(schema.TypeName(), "INVALID")
	ident := jen.Id(strcase.ToScreamingSnake(idName)).
		Id(schema.TypeName()).
		Op("=").
		Id("-1")
	constCode = append(constCode, ident)
	for i, v := range schema.Enum() {
		var (
			name string
			ok   bool
		)

		if name, ok = v.(string); !ok {
			name = "NULL"
		}

		idName = generateScreamingSnakeName(schema.TypeName(), name)
		if i == 0 {
			ident = jen.Id(idName).
				Op("=").
				Iota()
		} else {
			ident = jen.Id(idName)
		}
		constCode = append(constCode, ident)
	}

	file.Const().Defs(constCode...)
}

func generateStringFunc(schema Schema, file *jen.File) {
	lowerCamel := strcase.ToLowerCamel(schema.TypeName())
	ident := fmt.Sprintf("%sNames[v]", lowerCamel)
	file.Func().Parens(jen.Id("v").Id(schema.TypeName())).Id("String").Params().Id("string").
		Block(
			jen.Return(jen.Id(ident)),
		)
}

func generateIDFunc(schema Schema, file *jen.File) {
	var blockCode []jen.Code
	lowerCamel := strcase.ToLowerCamel(schema.TypeName())
	camelCase := strcase.ToCamel(schema.TypeName())
	funcName := fmt.Sprintf("New%s", camelCase)
	ident := fmt.Sprintf("%sIds[val]", lowerCamel)
	invalidName := generateScreamingSnakeName(schema.TypeName(), "INVALID")
	errInvalid := fmt.Sprintf("ErrInvalid%s", camelCase)

	if schema.schemaRef.Value.Nullable {
		blockCode = append(blockCode, jen.If(
			jen.Id("result").Op(",").Id("ok").Op(":=").Id(ident).Op(";").Id("ok").Block(
				jen.Return(jen.Id("result").Op(",").Id("nil")))),
			jen.If(jen.Id("val").Op("==").Lit("null").Op("||").Id("val").Op("==").Lit("").Block(
				jen.Return(jen.Id(generateScreamingSnakeName(schema.TypeName(), "NULL")).Op(",").Id("nil"))),
			),
			jen.Return(jen.Id(invalidName).Op(",").Id(errInvalid)),
		)
	} else {
		blockCode = append(blockCode, jen.If(
			jen.Id("result").Op(",").Id("ok").Op(":=").Id(ident).Op(";").Id("ok").Block(
				jen.Return(jen.Id("result").Op(",").Id("nil")))),
			jen.Return(jen.Id(invalidName).Op(",").Id(errInvalid)),
		)
	}
	file.Func().Id(funcName).Parens(jen.Id("val").String()).Parens(jen.Id(schema.TypeName()).Op(",").Error()).
		Block(blockCode...)
}

func generateVariables(schema Schema, file *jen.File) {
	var EnumNamesCode []jen.Code
	lowerCamel := strcase.ToLowerCamel(schema.TypeName())
	ident := fmt.Sprintf("%sNames", lowerCamel)

	for _, v := range schema.Enum() {
		litValue := "null"
		if v != nil {
			litValue = v.(string)
		}
		screamingSnake := generateScreamingSnakeName(schema.TypeName(), litValue)
		EnumNamesCode = append(EnumNamesCode, jen.Id(screamingSnake).Op(":").Lit(litValue))
	}
	file.Var().Id(ident).Op("=").Index(jen.Op("...")).Id("string").Values(EnumNamesCode...).Line()

	var EnumIdsCode []jen.Code
	ident = fmt.Sprintf("%sIds", lowerCamel)
	for _, v := range schema.Enum() {
		litValue := "null"
		if v != nil {
			litValue = v.(string)
		}
		screamingSnake := generateScreamingSnakeName(schema.TypeName(), litValue)
		EnumIdsCode = append(EnumIdsCode, jen.Lit(litValue).Op(":").Id(screamingSnake))
	}
	file.Var().Id(ident).Op("=").Map(jen.String()).Id(schema.TypeName()).Values(EnumIdsCode...).Line()
}

func generateMarshalJSONFunc(schema Schema, file *jen.File) {
	file.Func().Parens(
		jen.Id("v").Id(schema.TypeName())).
		Id("MarshalJSON").Params().Parens(jen.Index().Byte().Op(",").Error()).
		Block(
			jen.Return(jen.Index().Byte().Parens(jen.Id("v.String()")).Op(",").Id("nil")),
		)
}

func generateUnmarshalJSONFunc(schema Schema, file *jen.File) {
	lowerCamel := strcase.ToLowerCamel(schema.TypeName())
	ident := fmt.Sprintf("%sIds[strVal]", lowerCamel)
	var blockCode []jen.Code

	if schema.schemaRef.Value.Nullable {
		blockCode = append(blockCode, jen.Id("strVal").Op(":=").Id("strings").Dot("ReplaceAll").
			Parens(jen.String().Parens(jen.Id("val")).Op(",").Lit("\"").Op(",").Lit("")),
			jen.If(
				jen.Id("idVal").Op(",").Id("ok").Op(":=").Id(ident).Op(";").Id("ok").Block(
					jen.Op("*").Id("v").Op("=").Id("idVal"),
					jen.Return(jen.Id("nil")),
				)),
			jen.If(jen.Id("strVal").Op("==").Lit("null").Block(
				jen.Op("*").Id("v").Op("=").Id(generateScreamingSnakeName(schema.TypeName(), "NULL")),
				jen.Return(jen.Id("nil")),
			)),
			jen.Return(jen.Id("ErrInvalid"+strcase.ToCamel(schema.TypeName()))))
	} else {
		blockCode = append(blockCode, jen.Id("strVal").Op(":=").Id("strings").Dot("ReplaceAll").
			Parens(jen.String().Parens(jen.Id("val")).Op(",").Lit("\"").Op(",").Lit("")),
			jen.If(
				jen.Id("idVal").Op(",").Id("ok").Op(":=").Id(ident).Op(";").Id("ok").Block(
					jen.Op("*").Id("v").Op("=").Id("idVal"),
					jen.Return(jen.Id("nil")),
				)),
			jen.Return(jen.Id("ErrInvalid"+strcase.ToCamel(schema.TypeName()))))
	}
	file.Func().Parens(jen.Id("v").Op("*").Id(schema.TypeName())).Id("UnmarshalJSON").
		Params(jen.Id("val").Index().Byte()).Error().
		Block(blockCode...)
}

func generateEnums(schema Schema) error {
	f := jen.NewFile("api")

	f.Id("import").Parens(jen.Lit("errors").Line().Lit("strings"))

	f.Type().Id(schema.TypeName()).Int().Line()

	f.Var().Id("ErrInvalid" + strcase.ToCamel(schema.TypeName())).
		Op("=").Id("errors").Op(".").Id("New").
		Parens(
			jen.Lit("invalid " + schema.TypeName() + " value")).
		Line()

	generateEnumConstants(schema, f)
	generateVariables(schema, f)
	generateStringFunc(schema, f)
	generateIDFunc(schema, f)
	generateMarshalJSONFunc(schema, f)
	generateUnmarshalJSONFunc(schema, f)

	targetFileName := path.Join(
		skeletonConfig.TargetDirectory(),
		"pkg",
		"api",
		strcase.ToSnake(schema.TypeName())+".go",
	)

	return writeFile(targetFileName, f)
}
