package webservice

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type StatusCodeProvider interface {
	StatusCode() int
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

func NewBadRequestError(cause error) error {
	return NewStatusError(cause, http.StatusBadRequest)
}

func NewUnauthorizedError(cause error) error {
	return NewStatusError(cause, http.StatusUnauthorized)
}

func NewNotFoundError(cause error) error {
	return NewStatusError(cause, http.StatusNotFound)
}

func NewInternalError(cause error) error {
	return NewStatusError(cause, http.StatusInternalServerError)
}
