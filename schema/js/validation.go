// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package js

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"github.com/pkg/errors"
	jsv "github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/swaggest/jsonschema-go"
	"path"
	"strings"
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

func (s ValidationSchema) AllOf() []ValidationSchema {
	var results []ValidationSchema
	for _, allOfOne := range s.schema.AllOf {
		results = append(results, NewValidationSchema(allOfOne))
	}
	return results
}

func (s ValidationSchema) Validate(value interface{}) error {
	if s.schema == nil {
		// for testing
		return nil
	}

	return s.schema.Validate(value)
}
