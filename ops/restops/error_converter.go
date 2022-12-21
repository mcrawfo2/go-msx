// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"sync"
)

type StatusCodeProvider interface {
	StatusCode() int
}

// StatusCodeError is an error with an HTTP status code attached
type StatusCodeError interface {
	StatusCodeProvider
	error
}

func NewStatusCodeError(cause error, status int) StatusCodeError {
	if cause == nil {
		cause = errors.New(fmt.Sprintf("Unknown status error: %d", status))
	}
	return &statusError{
		cause:  cause,
		status: status,
	}
}

type statusError struct {
	cause  error
	status int
}

func (e *statusError) Error() string {
	return e.cause.Error()
}

func (e *statusError) StatusCode() int {
	return e.status
}

func (e *statusError) Cause() error {
	return e.cause
}

func (e *statusError) Unwrap() error {
	// leveraged by errors.As()
	return e.cause
}

// ErrorConverter transforms an error into a StatusCodeError
type ErrorConverter interface {
	Convert(error) StatusCodeError
}

type ErrorConverterFunc func(error) StatusCodeError

func (f ErrorConverterFunc) Convert(e error) StatusCodeError {
	return f(e)
}

// ErrorStatusCoder maps an error instance to an HTTP status code
type ErrorStatusCoder interface {
	Code(error) int
}

type ErrorStatusCoderFunc func(error) int

func (f ErrorStatusCoderFunc) Code(e error) int {
	return f(e)
}

type ErrorStatusCoderConverter struct {
	ErrorStatusCoder
}

func (f ErrorStatusCoderConverter) Convert(e error) StatusCodeError {
	code := f.Code(e)
	return NewStatusCodeError(e, code)
}

type Coder interface {
	Code() string
}

// CodedError is an error with a string code
type CodedError struct {
	code string
	error
}

func (c CodedError) Code() string {
	return c.code
}

func NewCodedError(code string, cause error) CodedError {
	return CodedError{
		code:  code,
		error: cause,
	}
}

var (
	defaultErrorStatusCodes    = map[error]int{}
	defaultErrorStatusCodesMtx sync.RWMutex
)

func SetMappedErrorStatusCode(err error, statusCode int) {
	defaultErrorStatusCodesMtx.Lock()
	defer defaultErrorStatusCodesMtx.Unlock()
	defaultErrorStatusCodes[err] = statusCode
}

func MappedErrorStatusCode(err error) types.Optional[int] {
	defaultErrorStatusCodesMtx.RLock()
	defer defaultErrorStatusCodesMtx.RUnlock()

	for k, v := range defaultErrorStatusCodes {
		if errors.Is(err, k) {
			return types.OptionalOf(v)
		}
	}
	return types.OptionalEmpty[int]()
}

func DefaultErrorStatusCoder(err error) int {
	optionalStatusCode := MappedErrorStatusCode(err)
	return optionalStatusCode.OrElse(http.StatusBadRequest)
}
