// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package v8

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
)

type ErrorCoder interface {
	Code() string
}

type Error struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func (e Error) Example() interface{} {
	return Error{
		Code:    "BIZ001",
		Message: "Entity in busy state",
	}
}

type pojoer interface {
	ToPojo() types.Pojo
}

func (e *Error) ApplyError(err error) {
	var errorCoder ErrorCoder
	if errors.As(err, &errorCoder) {
		e.Code = errorCoder.Code()
	} else {
		e.Code = "UNKNOWN"
	}

	var pojoErr pojoer
	if errors.As(err, &pojoErr) {
		e.Details = pojoErr.ToPojo()
	}

	e.Message = err.Error()
}
