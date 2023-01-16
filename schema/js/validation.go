// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package js

import (
	"bytes"
	"crypto/md5"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/hex"
	jsv "github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/swaggest/jsonschema-go"
	"sync"
)

type ValidationSchema struct {
	schema *jsv.Schema
}

func NewValidationSchema(schema *jsv.Schema) ValidationSchema {
	for schema != nil && schema.Ref != nil {
		schema = schema.Ref
	}
	return ValidationSchema{schema: schema}
}

var anyType = NewValidationSchema(&jsv.Schema{
	Types: []string{
		string(jsonschema.Null),
		string(jsonschema.Array),
		string(jsonschema.Object),
		string(jsonschema.String),
		string(jsonschema.Integer),
		string(jsonschema.Number),
		string(jsonschema.Boolean),
	},
})

func (s ValidationSchema) Types() []string {
	var st = make(types.StringSet)
	st.AddAll(anyType.schema.Types...)

	if s.schema == nil {
		return st.Values()
	}

	if s.schema.Types != nil {
		schemaTypes := s.schema.Types
		intersectionTypes := st.Intersect(types.NewStringSet(schemaTypes...))
		st = types.NewStringSet(intersectionTypes...)
	}

	for _, allOfOne := range s.AllOf() {
		allOfOneTypes := allOfOne.Types()
		intersectionTypes := st.Intersect(types.NewStringSet(allOfOneTypes...))
		st = types.NewStringSet(intersectionTypes...)
	}

	return st.Values()
}

func (s ValidationSchema) Items() *jsv.Schema {
	return s.schema.Items2020
}

func (s ValidationSchema) AllOf() []ValidationSchema {
	var results []ValidationSchema
	for _, allOfOne := range s.schema.AllOf {
		results = append(results, NewValidationSchema(allOfOne))
	}
	return results
}

func (s ValidationSchema) IsPropertyAllowed(key string) bool {
	// Explicitly declared properties
	if _, ok := s.schema.Properties[key]; ok {
		return true
	}

	if s.schema.AdditionalProperties == nil {
		return false
	}

	if allowed, ok := s.schema.AdditionalProperties.(bool); ok {
		return allowed
	}

	if _, ok := s.schema.AdditionalProperties.(*jsv.Schema); ok {
		return true
	}

	return false

}

func (s ValidationSchema) Property(key string) ValidationSchema {
	// Explicitly declared properties
	if ps, ok := s.schema.Properties[key]; ok {
		return NewValidationSchema(ps)
	}

	if s.schema.AdditionalProperties == nil {
		return NewValidationSchema(new(jsv.Schema))
	}

	if allowed, ok := s.schema.AdditionalProperties.(bool); ok {
		if allowed {
			return anyType
		}
	}

	if ps, ok := s.schema.AdditionalProperties.(*jsv.Schema); ok {
		return NewValidationSchema(ps)
	}

	return ValidationSchema{}
}

func (s ValidationSchema) Validate(value interface{}) error {
	if s.schema == nil {
		// for testing
		return nil
	}

	return s.schema.Validate(value)
}

var validationCompilerMtx sync.Mutex
var validationCompiler = jsv.NewCompiler()

func NewValidationSchemaFromJsonSchema(schema *jsonschema.Schema) (vs ValidationSchema, err error) {
	validationCompilerMtx.Lock()
	defer validationCompilerMtx.Unlock()

	// Convert the json schema to json
	schemaBytes, err := schema.JSONSchemaBytes()
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

	vs = NewValidationSchema(s)
	return
}
