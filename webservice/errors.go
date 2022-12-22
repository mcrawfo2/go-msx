// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package webservice

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type StatusCodeProvider interface {
	StatusCode() int
}

type statusCodeProviderImpl struct {
	body       interface{}
	statusCode int
}

func (s statusCodeProviderImpl) StatusCode() int {
	return s.statusCode
}

func (s statusCodeProviderImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.body)
}

func NewStatusCodeProvider(body interface{}, status int) StatusCodeProvider {
	return statusCodeProviderImpl{body: body, statusCode: status}
}

type StatusError struct {
	cause  error
	status int
}

func NewStatusError(cause error, status int) error {
	return NewStatusCodeError(cause, status)
}

func NewStatusCodeError(cause error, status int) StatusCodeError {
	if cause == nil {
		cause = errors.New(fmt.Sprintf("Unknown status error: %d", status))
	}
	return &StatusError{
		cause:  cause,
		status: status,
	}
}

func (e *StatusError) Error() string {
	return e.cause.Error()
}

func (e *StatusError) StatusCode() int {
	return e.status
}

func (e *StatusError) Cause() error {
	return e.cause
}

func (e *StatusError) Unwrap() error {
	// leveraged by errors.As()
	return e.cause
}

func NewBadRequestError(cause error) error {
	return NewStatusError(cause, http.StatusBadRequest)
}

func NewUnauthorizedError(cause error) error {
	return NewStatusError(cause, http.StatusUnauthorized)
}

func NewForbiddenError(cause error) error {
	return NewStatusError(cause, http.StatusForbidden)
}

func NewNotFoundError(cause error) error {
	return NewStatusError(cause, http.StatusNotFound)
}

func NewInternalError(cause error) error {
	return NewStatusError(cause, http.StatusInternalServerError)
}

func NewConflictError(cause error) error {
	return NewStatusError(cause, http.StatusConflict)
}

func NewBadRequestStatusError(cause error) StatusCodeError {
	return NewStatusCodeError(cause, http.StatusBadRequest)
}

func NewUnauthorizedStatusError(cause error) StatusCodeError {
	return NewStatusCodeError(cause, http.StatusUnauthorized)
}

func NewForbiddenStatusError(cause error) StatusCodeError {
	return NewStatusCodeError(cause, http.StatusForbidden)
}

func NewNotFoundStatusError(cause error) StatusCodeError {
	return NewStatusCodeError(cause, http.StatusNotFound)
}

func NewInternalStatusError(cause error) StatusCodeError {
	return NewStatusCodeError(cause, http.StatusInternalServerError)
}

func NewConflictStatusError(cause error) StatusCodeError {
	return NewStatusCodeError(cause, http.StatusConflict)
}

type ErrorRaw interface {
	SetError(code int, err error, path string)
}

type ErrorApplier interface {
	ApplyError(err error)
}

type ErrorCoder interface {
	Code() string
}

type ErrorV8 struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ErrorV8) ApplyError(err error) {
	if errorCoder, ok := err.(ErrorCoder); ok {
		e.Code = errorCoder.Code()
	} else {
		e.Code = "UNKNOWN"
	}

	e.Message = err.Error()
}

type StatusCodeError interface {
	StatusCodeProvider
	error
}

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