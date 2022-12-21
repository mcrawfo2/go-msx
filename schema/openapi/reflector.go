// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package openapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/restops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/js"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi3"
	"mime/multipart"
	"reflect"
)

var reflectorOptions = []func(*jsonschema.ReflectContext){
	jsonschema.DefinitionsPrefix("#/components/schemas/"),
	jsonschema.RootRef,
	js.DefNameInterceptor(),
	js.TypeTitleInterceptor(),
	OverrideSchema(),
	js.StructRequiredInterceptor(),
	js.EnvelopNullability(),
	js.ExampleInterceptor(),
}

var documentationReflector = openapi3.Reflector{
	Reflector: jsonschema.Reflector{
		DefaultOptions: reflectorOptions,
	},
	Spec: &openapi3.Spec{
		Openapi: "3.0.3",
	},
}

func RegistrySpec() *openapi3.Spec {
	return documentationReflector.SpecEns()
}

var overrideSchema = make(map[reflect.Type]jsonschema.Schema)
var overrideSchemaById = make(map[string]reflect.Type)

func OverrideSchema() func(*jsonschema.ReflectContext) {
	return jsonschema.InterceptType(func(value reflect.Value, schema *jsonschema.Schema) (bool, error) {
		valueType := value.Type()
		if valueSchemaOverride, ok := overrideSchema[valueType]; ok {
			*schema = valueSchemaOverride
		}

		// Continue with type interceptor chain
		return false, nil
	})
}

func RegisterTypeSchema(t reflect.Type, schema jsonschema.Schema) {
	overrideSchema[t] = schema
	if schema.ID != nil {
		overrideSchemaById[*schema.ID] = t
	}

	s := ConvertToOpenApiSchema(schema)
	documentationReflector.SpecEns().ComponentsEns().SchemasEns().WithMapOfSchemaOrRefValuesItem(*schema.ID, s)
}

type Patcher interface {
	PatchOpenApiSchema(schema *openapi3.Schema)
}

func ConvertToOpenApiSchema(schema jsonschema.Schema) openapi3.SchemaOrRef {
	s := openapi3.SchemaOrRef{}
	s.FromJSONSchema(schema.ToSchemaOrBool())
	return s
}

func PromoteNullableSchema(osr openapi3.SchemaOrRef) openapi3.SchemaOrRef {
	if osr.Schema == nil {
		return osr
	}

	// Convert enveloped schema
	if len(osr.Schema.AnyOf) == 2 {
		for i := 0; i < len(osr.Schema.AnyOf); i++ {
			one := osr.Schema.AnyOf[i]
			if one.Schema == nil {
				continue
			}
			if types.NewOptionalBool(one.Schema.Nullable).OrElse(false) {
				osr.Schema.Nullable = one.Schema.Nullable
				osr.Schema.AllOf = []openapi3.SchemaOrRef{
					osr.Schema.AnyOf[2-i-1],
				}
				osr.Schema.AnyOf = nil
				break
			}
		}
	}

	for key, property := range osr.Schema.Properties {
		osr.Schema.Properties[key] = PromoteNullableSchema(property)
	}

	if osr.Schema.Items != nil {
		*osr.Schema.Items = PromoteNullableSchema(*osr.Schema.Items)
	}

	if osr.Schema.AdditionalProperties != nil && osr.Schema.AdditionalProperties.SchemaOrRef != nil {
		*osr.Schema.AdditionalProperties.SchemaOrRef = PromoteNullableSchema(*osr.Schema.AdditionalProperties.SchemaOrRef)
	}

	return osr
}

func CollectDefinition(name string, schema jsonschema.Schema) {
	if _, ok := LookupSchema(name); ok {
		// Don't override
		return
	}

	s := ConvertToOpenApiSchema(schema)
	s = PromoteNullableSchema(s)

	if s.SchemaReference != nil {
		// Get type name
		schemaRefName := SchemaRefName(s.SchemaReference)

		// patch the schema if desired
		if directSchema, ok := LookupSchema(schemaRefName); ok {
			value := reflect.New(schema.ReflectType).Interface()
			if patcher, ok := value.(Patcher); ok {
				// filter the schema through the patcher
				patcher.PatchOpenApiSchema(directSchema)
			}
		}
	}

	documentationReflector.SpecEns().ComponentsEns().SchemasEns().WithMapOfSchemaOrRefValuesItem(name, s)
}

// OpenApi Schema

func Reflect(value interface{}) (*openapi3.SchemaOrRef, error) {
	schema, err := ReflectOpenApiJson(value)
	if err != nil {
		return nil, err
	}

	s := ConvertToOpenApiSchema(schema)
	s = PromoteNullableSchema(s)

	if s.SchemaReference != nil {
		// Get type name
		schemaRefName := SchemaRefName(s.SchemaReference)

		// patch the schema if desired
		if directSchema, ok := LookupSchema(schemaRefName); ok {
			if patcher, ok := value.(Patcher); ok {
				// filter the schema through the patcher
				patcher.PatchOpenApiSchema(directSchema)
			}
		}
	}

	return &s, nil
}

func addOpenApiSchema(portField *ops.PortField, sf reflect.StructField) error {
	value := reflect.Zero(sf.Type).Interface()

	openApiSchemaOrRef, err := Reflect(value)
	if err != nil {
		logger.WithError(err).Errorf("Failed to reflect field %q", portField.Name)
	}

	openApiSchema := openApiSchemaOrRef.Schema
	if openApiSchema != nil {
		if e := portField.Enum(); e != nil {
			openApiSchema.Enum = e
		}

		if err = PopulateFieldsFromTags(openApiSchema, portField.Tags()); err != nil {
			return err
		}
	}

	portFieldWithOpenApiSchema(portField, openApiSchemaOrRef)
	return nil
}

// Json Schema

func ReflectOpenApiJson(value interface{}) (jsonschema.Schema, error) {
	s, err := documentationReflector.Reflect(
		value,
		jsonschema.CollectDefinitions(CollectDefinition),
	)

	if err != nil {
		return jsonschema.Schema{}, err
	}

	return s, nil
}

func addJsonSchema(portField *ops.PortField, sf reflect.StructField) error {
	value := reflect.Zero(sf.Type).Interface()

	schema, err := ReflectOpenApiJson(value)
	if err != nil {
		return err
	}

	if schema.Ref == nil {
		if err = js.PopulateFieldsFromTags(&schema, portField.Tags()); err != nil {
			return err
		}
	}

	portFieldWithJsonSchema(portField, &schema)
	return nil
}

// Json Validator Schema

var validationReflector = jsonschema.Reflector{
	DefaultOptions: reflectorOptions[1:],
}

func reflectForJsonValidationSchema(field *ops.PortField, sf reflect.StructField) (schema jsonschema.Schema, err error) {
	// Compile json-validation specific json schema
	definitions := map[string]*jsonschema.Schema{}
	schema, err = validationReflector.Reflect(
		reflect.Zero(field.Type.Type).Interface(),
		jsonschema.CollectDefinitions(func(name string, schema jsonschema.Schema) {
			definitions[name] = &schema
		}),
	)

	if err != nil {
		err = errors.Wrap(err, "Failed to reflect json schema")
		return
	}

	// Attach any references
	for k, v := range definitions {
		s := v.ToSchemaOrBool()
		schema.WithDefinitionsItem(k, s)
	}

	// Add any enumerations
	if e := field.Enum(); e != nil {
		schema.Enum = e
	}

	// Add any tag fields
	if err = js.PopulateFieldsFromTags(&schema, sf.Tag); err != nil {
		logger.WithError(err).Errorf("Failed to populate tags onto %q", field.Name)
	}

	return
}

func addJsonValidationSchema(field *ops.PortField, sf reflect.StructField) (err error) {
	// Reflect a standalone jsonschema
	jsonSchema, err := reflectForJsonValidationSchema(field, sf)
	if err != nil {
		return
	}

	// Convert to js.ValidationSchema
	vs, err := js.NewValidationSchemaFromJsonSchema(&jsonSchema)
	if err != nil {
		return
	}

	// Save the compiled schema document to the port field
	portFieldWithJsonValidatorSchema(field, vs)

	return nil
}

// PortField Baggage

func PostProcessPortField(portField *ops.PortField, sf reflect.StructField) {
	if err := addJsonSchema(portField, sf); err != nil {
		logger.WithError(err).Errorf("Failed to compile json schema for field %s", portField.Name)
	}
	if err := addOpenApiSchema(portField, sf); err != nil {
		logger.WithError(err).Errorf("Failed to compile openapi schema for field %q", portField.Name)
	}
	if err := addJsonValidationSchema(portField, sf); err != nil {
		logger.WithError(err).Errorf("Failed to compile json validation schema for field %s", portField.Name)
	}
}

type portFieldKey int

const (
	portFieldKeyOpenApiSchema portFieldKey = iota
	portFieldKeyJsonSchema
	portFieldKeyJsonValidatorSchema
)

func portFieldWithOpenApiSchema(portField *ops.PortField, schema *openapi3.SchemaOrRef) {
	portField.WithBaggageItem(portFieldKeyOpenApiSchema, schema)
}

func openApiSchemaFromPortField(portField *ops.PortField) *openapi3.SchemaOrRef {
	value, _ := portField.Baggage[portFieldKeyOpenApiSchema].(*openapi3.SchemaOrRef)
	return value
}

func portFieldWithJsonSchema(portField *ops.PortField, schema *jsonschema.Schema) {
	portField.WithBaggageItem(portFieldKeyJsonSchema, schema)
}

func jsonSchemaFromPortField(portField *ops.PortField) *jsonschema.Schema {
	value, _ := portField.Baggage[portFieldKeyJsonSchema].(*jsonschema.Schema)
	return value
}

func portFieldWithJsonValidatorSchema(portField *ops.PortField, schema js.ValidationSchema) {
	portField.WithBaggageItem(portFieldKeyJsonValidatorSchema, schema)
}

func jsonValidatorSchemaFromPortField(portField *ops.PortField) js.ValidationSchema {
	value, _ := portField.Baggage[portFieldKeyJsonValidatorSchema].(js.ValidationSchema)
	return value
}

func GetJsonValidationSchema(field *ops.PortField) (schema js.ValidationSchema, err error) {
	return jsonValidatorSchemaFromPortField(field), nil
}

func init() {
	typeMappings := []struct {
		from interface{}
		to   interface{}
	}{
		{types.UUID{}, js.StringFormatUuid{}},
		{types.Time{}, js.StringFormatTime{}},
		{types.Duration(0), js.StringFormatDuration{}},
		{&multipart.FileHeader{}, js.StringFormatBinary{}},
		{[]byte{}, ""},
		{[]rune{}, ""},
	}
	for _, mapping := range typeMappings {
		documentationReflector.Reflector.AddTypeMapping(mapping.from, mapping.to)
		validationReflector.AddTypeMapping(mapping.from, mapping.to)
	}

	// Attach schema and validation schema to fields
	restops.PortReflectorPostProcessField = PostProcessPortField

	// Allow opaque retrieval of message validation schema
	restops.RegisterPortFieldValidationSchemaFunc(GetJsonValidationSchema)
}
