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
