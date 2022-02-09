package webservice

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

func MapSchema(additionalProperties *openapi3.Schema) *openapi3.Schema {
	ap := (&openapi3.SchemaAdditionalProperties{}).WithSchemaOrRef(NewSchemaOrRef(additionalProperties))
	return ObjectSchema().WithAdditionalProperties(*ap)
}

func NewType(simpleType jsonschema.SimpleType) *jsonschema.Type {
	val := simpleType.Type()
	return &val
}
