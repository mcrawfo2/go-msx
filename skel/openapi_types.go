// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package skel

import (
	"path"

	"github.com/iancoleman/strcase"
	"github.com/mcrawfo2/jennifer/jen"
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
	f := jen.NewFile("api")
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
