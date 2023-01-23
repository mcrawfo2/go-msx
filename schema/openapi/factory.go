// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package openapi

import (
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi3"
	"path"
)

func SchemaRefName(reference *openapi3.SchemaReference) string {
	return path.Base(reference.Ref)
}

func ParameterRefName(reference *openapi3.ParameterReference) string {
	return path.Base(reference.Ref)
}

func RequestBodyRefName(reference *openapi3.RequestBodyReference) string {
	return path.Base(reference.Ref)
}

func ResponseRefName(reference *openapi3.ResponseReference) string {
	return path.Base(reference.Ref)
}

func HeaderRefName(reference *openapi3.HeaderReference) string {
	return path.Base(reference.Ref)
}

func NewSchemaOrRef(s *openapi3.Schema) openapi3.SchemaOrRef {
	return openapi3.SchemaOrRef{
		Schema: s,
	}
}

func NewSchemaOrRefPtr(s *openapi3.Schema) *openapi3.SchemaOrRef {
	schemaOrRef := NewSchemaOrRef(s)
	return &schemaOrRef
}

func NewSchemaRef(name string) openapi3.SchemaOrRef {
	return openapi3.SchemaOrRef{
		SchemaReference: &openapi3.SchemaReference{
			Ref: "#/components/schemas/" + name,
		},
	}
}

func NewSchemaRefPtr(name string) *openapi3.SchemaOrRef {
	v := NewSchemaRef(name)
	return &v
}

func NewSchemaPtr(t openapi3.SchemaType) *openapi3.Schema {
	return (&openapi3.Schema{}).WithType(t)
}

func StringSchema() *openapi3.Schema {
	return NewSchemaPtr(openapi3.SchemaTypeString)
}

func BooleanSchema() *openapi3.Schema {
	return NewSchemaPtr(openapi3.SchemaTypeBoolean)
}

func NumberSchema() *openapi3.Schema {
	return NewSchemaPtr(openapi3.SchemaTypeNumber)
}

func IntegerSchema() *openapi3.Schema {
	return NewSchemaPtr(openapi3.SchemaTypeInteger)
}

func ArraySchema(items openapi3.SchemaOrRef) *openapi3.Schema {
	return NewSchemaPtr(openapi3.SchemaTypeArray).WithItems(items)
}

func ObjectSchema() *openapi3.Schema {
	return NewSchemaPtr(openapi3.SchemaTypeObject)
}

func AnySchema() *openapi3.Schema {
	return new(openapi3.Schema)
}

func MapSchema(additionalProperties *openapi3.Schema) *openapi3.Schema {
	ap := (&openapi3.SchemaAdditionalProperties{}).WithSchemaOrRef(NewSchemaOrRef(additionalProperties))
	return ObjectSchema().WithAdditionalProperties(*ap)
}

func NewType(simpleType jsonschema.SimpleType) *jsonschema.Type {
	val := simpleType.Type()
	return &val
}

func NewParameterRef(name string) openapi3.ParameterOrRef {
	return openapi3.ParameterOrRef{
		ParameterReference: &openapi3.ParameterReference{
			Ref: "#/components/parameters/" + name,
		},
	}
}

func NewParameterRefPtr(name string) *openapi3.ParameterOrRef {
	v := NewParameterRef(name)
	return &v
}

func NewHeaderRef(name string) openapi3.HeaderOrRef {
	return openapi3.HeaderOrRef{
		HeaderReference: &openapi3.HeaderReference{
			Ref: "#/components/headers/" + name,
		},
	}
}

func NewHeaderRefPtr(name string) *openapi3.HeaderOrRef {
	v := NewHeaderRef(name)
	return &v
}

func SwaggerTypeToOpenApiSchema(typ *string, format *string) (schema *openapi3.Schema) {
	if typ == nil {
		return
	}

	switch *typ {
	case "string":
		schema = StringSchema()
	case "integer":
		schema = IntegerSchema()
	case "number":
		schema = NumberSchema()
	case "object":
		schema = ObjectSchema()
	case "boolean":
		schema = BooleanSchema()
	case "array":
		schema = ArraySchema(NewSchemaOrRef(new(openapi3.Schema)))
	default:
		return
	}

	if format != nil {
		schema.WithFormat(*format)
	}

	return
}
