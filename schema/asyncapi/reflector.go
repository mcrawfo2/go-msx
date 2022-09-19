// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"bytes"
	"crypto/md5"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/js"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	jsv "github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/swaggest/jsonschema-go"
	"reflect"
	"sync"
)

type asyncApiReflector struct {
	jsonschema.Reflector
	Spec *Spec
}

func (r asyncApiReflector) SpecEns() *Spec {
	if r.Spec == nil {
		r.Spec = &Spec{}
	}
	return r.Spec
}

var reflectorOptions = []func(*jsonschema.ReflectContext){
	jsonschema.DefinitionsPrefix("#/components/schemas/"),
	jsonschema.RootRef,
	js.DefNameInterceptor(),
	js.TypeTitleInterceptor(),
	js.StructRequiredInterceptor(),
	js.EnvelopNullability(),
	js.ExampleInterceptor(),
}

var documentationReflector = asyncApiReflector{
	Reflector: jsonschema.Reflector{
		DefaultOptions: reflectorOptions,
	},
	Spec: &Spec{},
}

func RegistrySpec() *Spec {
	return documentationReflector.SpecEns()
}

func LookupSchema(schemaName string) (*jsonschema.Schema, bool) {
	result, ok := documentationReflector.SpecEns().ComponentsEns().Schemas[schemaName]
	if !ok {
		return nil, false
	}
	return &result, ok
}

func CollectDefinition(name string, schema jsonschema.Schema) {
	if _, ok := LookupSchema(name); ok {
		// Don't override
		return
	}

	documentationReflector.SpecEns().ComponentsEns().WithSchemasItem(name, schema)
}

func Reflect(value interface{}) (jsonschema.Schema, error) {
	s, err := documentationReflector.Reflect(
		value,
		jsonschema.CollectDefinitions(CollectDefinition),
	)

	if err != nil {
		return jsonschema.Schema{}, err
	}

	return s, nil
}

func addJsonSchema(portField *ops.PortField, sf reflect.StructField) {
	jsonSchema, err := Reflect(reflect.Zero(sf.Type).Interface())
	if err != nil {
		logger.WithError(err).Errorf("Failed to reflect field %q", portField.Name)
	}

	if e := portField.Enum(); e != nil {
		jsonSchema.Enum = e
	}

	if err = js.PopulateFieldsFromTags(&jsonSchema, sf.Tag); err != nil {
		logger.WithError(err).Errorf("Failed to populate tags onto %q", portField.Name)
	}

	portFieldWithJsonSchema(portField, &jsonSchema)
}

var validationCompilerMtx sync.Mutex
var validationCompiler = jsv.NewCompiler()
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
	validationCompilerMtx.Lock()
	defer validationCompilerMtx.Unlock()

	// Reflect a standalone jsonschema
	jsob, err := reflectForJsonValidationSchema(field, sf)
	if err != nil {
		return
	}

	// Convert the json schema to json
	schemaBytes, err := jsob.JSONSchemaBytes()
	if err != nil {
		return
	}

	// Compile the standalone document
	sum := md5.Sum(schemaBytes)
	hash := hex.EncodeToString(sum[:])
	schemaUrl := "mem:///" + hash + ".json"

	if err = validationCompiler.AddResource(schemaUrl, bytes.NewReader(schemaBytes)); err != nil {
		return
	}

	s, err := validationCompiler.Compile(schemaUrl)
	if err != nil {
		return
	}

	// Save the compiled schema document to the port field
	portFieldWithJsonValidatorSchema(field, s)

	return nil
}

func PostProcessPortField(portField *ops.PortField, sf reflect.StructField) {
	addJsonSchema(portField, sf)
	if err := addJsonValidationSchema(portField, sf); err != nil {
		logger.Errorf("Failed to compile json schema for field %s", portField.Name)
	}
}

type portFieldKey int

const (
	portFieldKeyJsonSchema portFieldKey = iota
	portFieldKeyJsonValidatorSchema
)

func portFieldWithJsonSchema(portField *ops.PortField, schema *jsonschema.Schema) {
	portField.WithBaggageItem(portFieldKeyJsonSchema, schema)
}

func jsonSchemaFromPortField(portField *ops.PortField) *jsonschema.Schema {
	value, _ := portField.Baggage[portFieldKeyJsonSchema].(*jsonschema.Schema)
	return value
}

func portFieldWithJsonValidatorSchema(portField *ops.PortField, schema *jsv.Schema) {
	portField.WithBaggageItem(portFieldKeyJsonValidatorSchema, schema)
}

func jsonValidatorSchemaFromPortField(portField *ops.PortField) *jsv.Schema {
	value, _ := portField.Baggage[portFieldKeyJsonValidatorSchema].(*jsv.Schema)
	return value
}

func GetJsonValidationSchema(field *ops.PortField) (schema js.ValidationSchema, err error) {
	jsvSchema := jsonValidatorSchemaFromPortField(field)
	if jsvSchema == nil {
		panic(fmt.Sprintf("JSON Validation Schema not found for field %s", field.Name))
	}
	logger.Infof("JSON Validation Schema for field %s: %s", field.Name, jsvSchema.String())
	return js.NewValidationSchema(jsvSchema), nil
}

func init() {
	typeMappings := []struct {
		from interface{}
		to   interface{}
	}{
		{types.UUID{}, js.StringFormatUuid{}},
		{types.Time{}, js.StringFormatTime{}},
		{types.Duration(0), js.StringFormatDuration{}},
		{[]byte{}, ""},
		{[]rune{}, ""},
	}

	for _, mapping := range typeMappings {
		documentationReflector.Reflector.AddTypeMapping(mapping.from, mapping.to)
		validationReflector.AddTypeMapping(mapping.from, mapping.to)
	}

	// Attach schema and validation schema to fields
	streamops.PortReflectorPostProcessField = PostProcessPortField

	// Allow opaque retrieval of message validation schema
	streamops.RegisterPortFieldValidationSchemaFunc(GetJsonValidationSchema)
}
