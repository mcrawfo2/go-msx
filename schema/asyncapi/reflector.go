// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/js"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/swaggest/jsonschema-go"
	"reflect"
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

var Reflector = asyncApiReflector{
	Reflector: jsonschema.Reflector{
		DefaultOptions: []func(*jsonschema.ReflectContext){
			jsonschema.DefinitionsPrefix("#/components/schemas/"),
			jsonschema.RootRef,
			js.DefNameInterceptor(),
			js.TypeTitleInterceptor(),
			js.StructRequiredInterceptor(),
			js.EnvelopNullability(),
			js.ExampleInterceptor(),
		},
	},
	Spec: &Spec{},
}

func LookupSchema(schemaName string) (*jsonschema.Schema, bool) {
	result, ok := Reflector.SpecEns().ComponentsEns().Schemas[schemaName]
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

	Reflector.SpecEns().ComponentsEns().WithSchemasItem(name, schema)
}

func Reflect(value interface{}) (jsonschema.Schema, error) {
	s, err := Reflector.Reflect(
		value,
		jsonschema.CollectDefinitions(CollectDefinition),
	)

	if err != nil {
		return jsonschema.Schema{}, err
	}

	return s, nil
}

func PostProcessPortField(portField *ops.PortField, sf reflect.StructField) {
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

	portFieldWithSchema(portField, &jsonSchema)
}

type portFieldKey int

const (
	portFieldKeyJsonSchema portFieldKey = iota
)

func portFieldWithSchema(portField *ops.PortField, schema *jsonschema.Schema) {
	portField.WithBaggageItem(portFieldKeyJsonSchema, schema)
}

func schemaFromPortField(portField *ops.PortField) *jsonschema.Schema {
	value, _ := portField.Baggage[portFieldKeyJsonSchema].(*jsonschema.Schema)
	return value
}

func init() {
	Reflector.Reflector.AddTypeMapping(types.UUID{}, js.StringFormatUuid{})
	Reflector.Reflector.AddTypeMapping(types.Time{}, js.StringFormatTime{})
	Reflector.Reflector.AddTypeMapping(types.Duration(0), js.StringFormatDuration{})
	Reflector.Reflector.AddTypeMapping([]byte{}, "")
	Reflector.Reflector.AddTypeMapping([]rune{}, "")

	streamops.PortReflectorPostProcessField = PostProcessPortField
}
