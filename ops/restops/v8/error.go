// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package v8

type ErrorCoder interface {
	Code() string
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e Error) Example() interface{} {
	return Error{
		Code:    "BIZ001",
		Message: "Entity in busy state",
	}
}

func (e *Error) ApplyError(err error) {
	if errorCoder, ok := err.(ErrorCoder); ok {
		e.Code = errorCoder.Code()
	} else {
		e.Code = "UNKNOWN"
	}

	e.Message = err.Error()
}
