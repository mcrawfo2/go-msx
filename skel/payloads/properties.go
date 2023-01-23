// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package payloads

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/swaggest/jsonschema-go"
)

const (
	ExtraPropertiesGoJsonSchema = "goJSONSchema"
	ExtraPropertiesComponents   = "components"
	ExtraPropertiesSchemas      = "schemas"
)

type Components map[string]any

func (c Components) Schemas() Schemas {
	if _, ok := c[ExtraPropertiesSchemas]; !ok {
		c[ExtraPropertiesSchemas] = make(map[string]jsonschema.SchemaOrBool)
	}
	return c[ExtraPropertiesSchemas].(map[string]jsonschema.SchemaOrBool)
}

func ComponentsForSchema(schema *jsonschema.Schema) Components {
	if _, ok := schema.ExtraProperties[ExtraPropertiesComponents]; !ok {
		schema.WithExtraPropertiesItem(
			ExtraPropertiesComponents,
			make(map[string]any))
	}
	return schema.ExtraProperties[ExtraPropertiesComponents].(map[string]any)
}

type Schemas map[string]jsonschema.SchemaOrBool

func (s Schemas) Each(fn func(key string, value jsonschema.SchemaOrBool)) {
	for k, v := range s {
		fn(k, v)
	}
}

type GoJsonSchema map[string]any

func (s GoJsonSchema) Tags() GoJsonSchemaTags {
	if _, ok := s["tags"]; !ok {
		s["tags"] = make(map[string]string)
	}
	return s["tags"].(map[string]string)
}

func (s GoJsonSchema) Type() types.Optional[string] {
	if t, ok := s["type"]; ok {
		return types.OptionalOf(t.(string))
	}
	return types.OptionalEmpty[string]()
}

func GoJsonSchemaForSchema(schema *jsonschema.Schema) GoJsonSchema {
	if _, ok := schema.ExtraProperties[ExtraPropertiesGoJsonSchema]; !ok {
		schema.WithExtraPropertiesItem(
			ExtraPropertiesGoJsonSchema,
			make(map[string]any))
	}
	return schema.ExtraProperties[ExtraPropertiesGoJsonSchema].(map[string]any)
}

type GoJsonSchemaTags map[string]string

func (t GoJsonSchemaTags) AddTag(key, value string) GoJsonSchemaTags {
	t[key] = value
	return t
}

func (t GoJsonSchemaTags) ClearTag(key string) GoJsonSchemaTags {
	t[key] = ""
	return t
}

func AddJsonSchemaObjectProperty(objectSchema *jsonschema.Schema, propertyName string, propertySchema *jsonschema.Schema, required bool) {
	if required {
		objectSchema.Required = append(objectSchema.Required, propertyName)
	}

	if len(propertySchema.Definitions) > 0 {
		for dk, dv := range propertySchema.Definitions {
			objectSchema.WithDefinitionsItem(dk, dv)
		}
		propertySchema.Definitions = nil
	}

	objectSchema.WithPropertiesItem(propertyName, propertySchema.ToSchemaOrBool())
}
