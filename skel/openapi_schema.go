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
	f.Type().Id(schema.TypeName()).Struct(properties...)

	targetFileName := path.Join(
		skeletonConfig.TargetDirectory(),
		"pkg",
		"api",
		strcase.ToSnake(schema.TypeName())+".go")

	return writeFile(targetFileName, f)
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
		if !schemaProperty.SchemaType.Required() && !schemaProperty.SchemaType.IsReference() {
			property = property.Op("*")
		}

		err := generateTypeWithImport(f, schema.Namespace(skeletonConfig.AppPackageUrl()), property, schemaProperty.SchemaType)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Failed to generate property type")
		}

		// Tags
		property = property.Tag(map[string]string{
			"json": strcase.ToLowerCamel(schemaProperty.JsonFieldName()),
		})

		properties = append(properties, property)
	}

	return properties, imports, nil
}
