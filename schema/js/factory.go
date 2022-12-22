// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package js

import (
	"github.com/swaggest/jsonschema-go"
)

func NewType(simpleType jsonschema.SimpleType) *jsonschema.Type {
	val := simpleType.Type()
	return &val
}

func NewSchemaPtr(t jsonschema.SimpleType) *jsonschema.Schema {
	return new(jsonschema.Schema).WithType(*NewType(t))
}

func StringSchema() *jsonschema.Schema {
	return NewSchemaPtr(jsonschema.String)
}

func BooleanSchema() *jsonschema.Schema {
	return NewSchemaPtr(jsonschema.Boolean)
}

func NumberSchema() *jsonschema.Schema {
	return NewSchemaPtr(jsonschema.Number)
}

func IntegerSchema() *jsonschema.Schema {
	return NewSchemaPtr(jsonschema.Integer)
}

func ArraySchema(items jsonschema.Schema) *jsonschema.Schema {
	return NewSchemaPtr(jsonschema.Array).
		WithItems(*(&jsonschema.Items{}).
			WithSchemaOrBool(items.
				ToSchemaOrBool()))
}

func ObjectSchema() *jsonschema.Schema {
	return NewSchemaPtr(jsonschema.Object)
}

func AnySchema() *jsonschema.Schema {
	return &jsonschema.Schema{}
}

func MapSchema(additionalProperties *jsonschema.Schema) *jsonschema.Schema {
	return ObjectSchema().
		WithAdditionalProperties(additionalProperties.
			ToSchemaOrBool())
}
