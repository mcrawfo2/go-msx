// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/schema/js"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	jsv "github.com/santhosh-tekuri/jsonschema/v5"
	"path"
	"strings"
)

type PortFieldValidationSchemaFunc func(field *PortField) (schema js.ValidationSchema, err error)

var ErrValidationFailed = errors.New("Validation failure")

type ValidationFailure struct {
	Path     string
	Failures []string
	Children map[string]*ValidationFailure
}

func (e *ValidationFailure) Error() string {
	return ErrValidationFailed.Error()
}

func (e *ValidationFailure) Unwrap() error {
	return ErrValidationFailed
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
		return e.applyErrorMap(typedErr)

	case validation.Errors:
		return e.applyErrorMap(typedErr)

	default:
		e.Failures = append(e.Failures, err.Error())
	}

	return e
}

func (e *ValidationFailure) applyErrorMap(err map[string]error) *ValidationFailure {
	for k, v := range err {
		child := NewValidationFailure(path.Join(e.Path, k))
		e.Children[k] = child.Apply(v)
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

func NewValidationFailure(p string) *ValidationFailure {
	if len(p) > 0 && p[0] != '/' {
		p = "/" + p
	}
	return &ValidationFailure{
		Path:     p,
		Children: make(map[string]*ValidationFailure),
	}
}
