// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package skel

import (
	"path"

	"github.com/iancoleman/strcase"
	"github.com/mcrawfo2/jennifer/jen"
	"github.com/pkg/errors"
)

func generateSchema(schema Schema) error {
	f := jen.NewFile("api")

	properties, imports, err := generateSchemaProperties(f, schema)
	if err != nil {
		return errors.Wrap(err, "Failed to generate struct fields")
	}

	f.ImportNames(imports)
	f.Type().Id(schema.TypeName()).Struct(properties...).Line()

	validation, err := generateSchemaValidation(f, schema)
	if err != nil {
		return err
	}

	f.Func().
		Parens(jen.Id("v").Op("*").Id(schema.TypeName())).
		Id("Validate").Params().Id("error").
		Block(jen.Return(
			jen.Qual(pkgTypes, "ErrorMap").Values(validation)))

	targetFileName := path.Join(
		skeletonConfig.TargetDirectory(),
		"pkg",
		"api",
		strcase.ToSnake(schema.TypeName())+".go")

	return writeFile(targetFileName, f)
}

func generateSchemaValidation(f *jen.File, schema Schema) (jen.Code, error) {
	properties, err := schema.Properties()
	if err != nil {
		return nil, err
	}

	f.ImportName(pkgValidation, "validation")

	var result = make(jen.Dict)
	for _, p := range properties {
		validators, err := generateValidators(f, p.Schema)
		if err != nil {
			return nil, err
		}

		args := append([]jen.Code{
			jen.Op("&").Id("v").Dot(p.StructFieldName()),
		}, validators...)

		result[jen.Lit(p.JsonName())] = jen.Qual(pkgValidation, "Validate").Call(args...)
	}

	return result, nil
}

func generateSchemaProperties(f *jen.File, schema Schema) ([]jen.Code, map[string]string, error) {
	var (
		properties = make([]jen.Code, 0)
		imports    = make(map[string]string)
		ns         = schema.Namespace(skeletonConfig.AppPackageUrl())

		schemaProps []Property
		err         error
	)

	if schemaProps, err = schema.Properties(); err != nil {
		return nil, nil, err
	}

	for _, prop := range schemaProps {
		var (
			fieldName = prop.StructFieldName()
			statement = jen.Id(fieldName)

			jsonName  = strcase.ToLowerCamel(prop.JsonName())
			fieldTags = map[string]string{"json": jsonName}
		)

		// Type
		if !prop.Schema.Required() && !prop.Schema.IsReference() {
			statement = statement.Op("*")
		}

		if err = generateTypeWithImport(f, ns, statement, prop.Schema); err != nil {
			return nil, nil, errors.Wrap(err, "Failed to generate property type")
		}

		// Tags
		properties = append(properties, statement.Tag(fieldTags))
	}

	return properties, imports, nil
}
