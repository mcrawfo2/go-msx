package skel

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"path"
)

func generateScreamingSnakeName(s1 string, s2 string) string {
	return strcase.ToScreamingSnake(fmt.Sprintf("%s_%s", s1, s2))
}

func generateEnumConstants(schema Schema, file *File) {
	var constCode []Code
	idName := generateScreamingSnakeName(schema.TypeName(), "INVALID")
	ident := Id(strcase.ToScreamingSnake(idName)).
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
			ident = Id(idName).
				Op("=").
				Iota()
		} else {
			ident = Id(idName)
		}
		constCode = append(constCode, ident)
	}

	file.Const().Defs(constCode...)
}

func generateStringFunc(schema Schema, file *File) {
	lowerCamel := strcase.ToLowerCamel(schema.TypeName())
	ident := fmt.Sprintf("%sNames[v]", lowerCamel)
	file.Func().Parens(Id("v").Id(schema.TypeName())).Id("String").Params().Id("string").
		Block(
			Return(Id(ident)),
		)
}

func generateIDFunc(schema Schema, file *File) {
	var blockCode []Code
	lowerCamel := strcase.ToLowerCamel(schema.TypeName())
	camelCase := strcase.ToCamel(schema.TypeName())
	funcName := fmt.Sprintf("New%s", camelCase)
	ident := fmt.Sprintf("%sIds[val]", lowerCamel)
	invalidName := generateScreamingSnakeName(schema.TypeName(), "INVALID")
	errInvalid := fmt.Sprintf("ErrInvalid%s", camelCase)

	if schema.schemaRef.Value.Nullable {
		blockCode = append(blockCode, If(
			Id("result").Op(",").Id("ok").Op(":=").Id(ident).Op(";").Id("ok").Block(
				Return(Id("result").Op(",").Id("nil")))),
			If(Id("val").Op("==").Lit("null").Op("||").Id("val").Op("==").Lit("").Block(
				Return(Id(generateScreamingSnakeName(schema.TypeName(), "NULL")).Op(",").Id("nil"))),
			),
			Return(Id(invalidName).Op(",").Id(errInvalid)),
		)
	} else {
		blockCode = append(blockCode, If(
			Id("result").Op(",").Id("ok").Op(":=").Id(ident).Op(";").Id("ok").Block(
				Return(Id("result").Op(",").Id("nil")))),
			Return(Id(invalidName).Op(",").Id(errInvalid)),
		)
	}
	file.Func().Id(funcName).Parens(Id("val").String()).Parens(Id(schema.TypeName()).Op(",").Error()).
		Block(blockCode...)
}

func generateVariables(schema Schema, file *File) {
	var EnumNamesCode []Code
	lowerCamel := strcase.ToLowerCamel(schema.TypeName())
	ident := fmt.Sprintf("%sNames", lowerCamel)

	for _, v := range schema.Enum() {
		litValue := "null"
		if v != nil {
			litValue = v.(string)
		}
		screamingSnake := generateScreamingSnakeName(schema.TypeName(), litValue)
		EnumNamesCode = append(EnumNamesCode, Id(screamingSnake).Op(":").Lit(litValue))
	}
	file.Var().Id(ident).Op("=").Index(Op("...")).Id("string").Values(EnumNamesCode...).Line()

	var EnumIdsCode []Code
	ident = fmt.Sprintf("%sIds", lowerCamel)
	for _, v := range schema.Enum() {
		litValue := "null"
		if v != nil {
			litValue = v.(string)
		}
		screamingSnake := generateScreamingSnakeName(schema.TypeName(), litValue)
		EnumIdsCode = append(EnumIdsCode, Lit(litValue).Op(":").Id(screamingSnake))
	}
	file.Var().Id(ident).Op("=").Map(String()).Id(schema.TypeName()).Values(EnumIdsCode...).Line()
}

func generateMarshalJSONFunc(schema Schema, file *File) {
	file.Func().Parens(
		Id("v").Op("*").Id(schema.TypeName())).
		Id("MarshalJSON").Params().Parens(Index().Byte().Op(",").Error()).
		Block(
			Return(Index().Byte().Parens(Id("v.String()")).Op(",").Id("nil")),
		)
}

func generateUnmarshalJSONFunc(schema Schema, file *File) {
	lowerCamel := strcase.ToLowerCamel(schema.TypeName())
	ident := fmt.Sprintf("%sIds[strVal]", lowerCamel)
	var blockCode []Code

	if schema.schemaRef.Value.Nullable {
		blockCode = append(blockCode, Id("strVal").Op(":=").Id("strings").Dot("ReplaceAll").
			Parens(String().Parens(Id("val")).Op(",").Lit("\"").Op(",").Lit("")),
			If(
				Id("idVal").Op(",").Id("ok").Op(":=").Id(ident).Op(";").Id("ok").Block(
					Op("*").Id("v").Op("=").Id("idVal"),
					Return(Id("nil")),
				)),
			If(Id("strVal").Op("==").Lit("null").Block(
				Op("*").Id("v").Op("=").Id(generateScreamingSnakeName(schema.TypeName(), "NULL")),
				Return(Id("nil")),
			)),
			Return(Id("ErrInvalid"+strcase.ToCamel(schema.TypeName()))))
	} else {
		blockCode = append(blockCode, Id("strVal").Op(":=").Id("strings").Dot("ReplaceAll").
			Parens(String().Parens(Id("val")).Op(",").Lit("\"").Op(",").Lit("")),
			If(
				Id("idVal").Op(",").Id("ok").Op(":=").Id(ident).Op(";").Id("ok").Block(
					Op("*").Id("v").Op("=").Id("idVal"),
					Return(Id("nil")),
				)),
			Return(Id("ErrInvalid"+strcase.ToCamel(schema.TypeName()))))
	}
	file.Func().Parens(Id("v").Op("*").Id(schema.TypeName())).Id("UnmarshalJSON").
		Params(Id("val").Index().Byte()).Error().
		Block(blockCode...)
}

func generateEnums(schema Schema) error {
	f := NewFile("api")

	f.Id("import").Parens(Lit("errors").Line().Lit("strings"))

	f.Type().Id(schema.TypeName()).Int().Line()

	f.Var().Id("ErrInvalid" + strcase.ToCamel(schema.TypeName())).
		Op("=").Id("errors").Op(".").Id("New").
		Parens(
			Lit("invalid " + schema.TypeName() + " value")).
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
