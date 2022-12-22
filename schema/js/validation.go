// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package js

import (
	"bytes"
	"crypto/md5"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	jsv "github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/swaggest/jsonschema-go"
	"path"
	"strings"
	"sync"
)

func NewValidationFailure(p string) *ValidationFailure {
	if len(p) > 0 && p[0] != '/' {
		p = "/" + p
	}
	return &ValidationFailure{
		Path:     p,
		Children: make(map[string]*ValidationFailure),
	}
}

var ErrValidationFailed = errors.New("Validation failure")

type ValidationErrors types.Pojo

func (v ValidationErrors) Error() string {
	return ErrValidationFailed.Error()
}

func (v ValidationErrors) ToPojo() types.Pojo {
	return types.Pojo(v)
}

func (v ValidationErrors) LogFields() map[string]interface{} {
	return v
}

func (v ValidationErrors) MarshalJSON() ([]byte, error) {
	return json.Marshal(v)
}

type ValidationFailure struct {
	Path     string
	Failures []string
	Children map[string]*ValidationFailure
}

func (e *ValidationFailure) Error() string {
	return ErrValidationFailed.Error()
}

func (e *ValidationFailure) ToPojo() types.Pojo {
	if len(e.Failures) == 0 && len(e.Children) == 0 {
		return nil
	}

	var result = make(types.Pojo)
	if len(e.Failures) > 0 {
		result[".failures"] = e.Failures
	}
	for k, v := range e.Children {
		result[k] = v.ToPojo()
	}
	return result
}

func (e *ValidationFailure) LogFields() map[string]interface{} {
	return map[string]interface{}{
		"validation": e,
	}
}

func (e *ValidationFailure) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.ToPojo())
}

func (e *ValidationFailure) Apply(err error) *ValidationFailure {
	switch typedErr := err.(type) {
	case *jsv.ValidationError:
		return e.applyJsvValidationError(typedErr)

	case types.ErrorList:
		for _, v := range typedErr {
			e.Failures = append(e.Failures, v.Error())
		}

	case types.ErrorMap:
		for k, v := range typedErr {
			child := new(ValidationFailure)
			e.Children[k] = child.Apply(v)
		}

	default:
		e.Failures = append(e.Failures, err.Error())
	}

	return e
}

func (e *ValidationFailure) applyJsvValidationError(err *jsv.ValidationError) *ValidationFailure {
	if err.InstanceLocation != e.Path {
		// Walk to next child
		suffix := strings.TrimPrefix(err.InstanceLocation, e.Path)
		suffixParts := strings.SplitN(suffix, "/", 3)
		child, ok := e.Children[suffixParts[1]]
		if !ok {
			// Create new child
			child = NewValidationFailure(path.Join(e.Path, suffixParts[1]))
			e.Children[suffixParts[1]] = child
		}

		// Not an actual error return
		_ = child.Apply(err)
		return e
	}

	if len(err.Causes) == 0 {
		e.Failures = append(e.Failures, err.Message)
		return e
	}

	for _, cause := range err.Causes {
		// Not an actual error return
		_ = e.Apply(cause)
	}

	return e
}

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
