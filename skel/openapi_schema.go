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
	var properties []Code
	var imports = make(map[string]string)

	schemaProperties, err := schema.Properties()
	if err != nil {
		return nil, nil, err
	}

	for _, schemaProperty := range schemaProperties {
		// Name
		property := Id(schemaProperty.StructFieldName())

		// Type
		if !schemaProperty.Schema.Required() && !schemaProperty.Schema.IsReference() {
			property = property.Op("*")
		}

		err := generateTypeWithImport(f, schema.Namespace(skeletonConfig.AppPackageUrl()), property, schemaProperty.Schema)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Failed to generate property type")
		}

		// Tags
		property = property.Tag(map[string]string{
			"json": strcase.ToLowerCamel(schemaProperty.JsonName()),
		})

		properties = append(properties, property)
	}

	return properties, imports, nil
}
