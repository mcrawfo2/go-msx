// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package skel

import (
	"path"

	. "github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
)

func generateSchema(schema Schema) error {
	f := NewFile("api")

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
		Parens(Id("v").Op("*").Id(schema.TypeName())).
		Id("Validate").Params().Id("error").
		Block(Return(
			Qual(pkgTypes, "ErrorMap").Values(validation)))

	targetFileName := path.Join(
		skeletonConfig.TargetDirectory(),
		"pkg",
		"api",
		strcase.ToSnake(schema.TypeName())+".go")

	return writeFile(targetFileName, f)
}

func generateSchemaValidation(f *File, schema Schema) (Code, error) {
	properties, err := schema.Properties()
	if err != nil {
		return nil, err
	}

	f.ImportName(pkgValidation, "validation")

	var result = make(Dict)
	for _, p := range properties {
		validators, err := generateValidators(f, p.Schema)
		if err != nil {
			return nil, err
		}

		args := append([]Code{
			Op("&").Id("v").Dot(p.StructFieldName()),
		}, validators...)

		result[Lit(p.JsonName())] = Qual(pkgValidation, "Validate").Call(args...)
	}

	return result, nil
}

func generateSchemaProperties(f *File, schema Schema) ([]Code, map[string]string, error) {
	var (
		properties = make([]Code, 0)
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
			statement = Id(fieldName)

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
