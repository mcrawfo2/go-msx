package skel

import (
	. "github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"path"
)

func stringFormat(schema Schema) string {
	switch schema.schemaRef.Value.Format {
	case "uuid":
		return "types.UUID"
	case "date-time":
		return "types.Time"
	default:
		return "string"
	}
}

func intFormat(schema Schema) string {
	switch schema.schemaRef.Value.Format {
	case "int32":
		return "int32"
	case "int64":
		return "int64"
	default:
		return "int"
	}
}

func doubleFormat(schema Schema) string {
	switch schema.schemaRef.Value.Format {
	case "float32":
		return "float32"
	default:
		return "float64"
	}
}

func arrayFormat(schema Schema) (Schema, error) {
	arraySchema, err := NewArrayType(schema.schemaRef, schema.required)
	if err != nil {
		return Schema{}, err
	}
	return arraySchema, nil
}

func generateTypeSchema(schema Schema) error {
	f := NewFile("api")
	ns := schema.Namespace(skeletonConfig.AppPackageUrl())
	statement := f.Type().Id(schema.TypeName())

	switch schema.schemaRef.Value.Type {
	case "string":
		statement.Id(stringFormat(schema))
	case "integer":
		statement.Id(intFormat(schema))
	case "number":
		statement.Id(doubleFormat(schema))
	case "boolean":
		statement.Id("bool")
	case "array":
		schema, err := arrayFormat(schema)
		if err != nil {
			return err
		}
		err = generateTypeWithImport(f, ns, statement, schema)
		if err != nil {
			return err
		}
	}

	targetFileName := path.Join(
		skeletonConfig.TargetDirectory(),
		"pkg",
		"api",
		strcase.ToSnake(schema.TypeName())+".go",
	)

	return writeFile(targetFileName, f)
}
