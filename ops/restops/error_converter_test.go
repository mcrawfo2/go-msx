// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestErrorCoderConverter(t *testing.T) {
	f := ErrorStatusCoderConverter{
		ErrorStatusCoder: ErrorStatusCoderFunc(func(err error) int {
			return http.StatusNotFound
		}),
	}

	e := errors.New("some error")

	c := f.Convert(e)
	assert.Equal(t, "some error", c.Error())
	assert.Equal(t, 404, c.StatusCode())
}

func TestErrorConverterFunc(t *testing.T) {
	f := ErrorConverterFunc(func(err error) StatusCodeError {
		return NewStatusCodeError(err, http.StatusNotFound)
	})

	e := errors.New("some error")

	c := f.Convert(e)
	assert.Equal(t, "some error", c.Error())
	assert.Equal(t, 404, c.StatusCode())
}

func TestNewCodedError(t *testing.T) {
	e := errors.New("some error")
	codedError := NewCodedError("test", e)
	assert.Equal(t, "test", codedError.Code())
	assert.Equal(t, "some error", codedError.Error())
}

func Test_statusError(t *testing.T) {
	e := &statusError{
		cause:  errors.New("some error"),
		status: http.StatusOK,
	}
	assert.Equal(t, http.StatusOK, e.StatusCode())
	assert.Equal(t, "some error", e.Error())
	assert.Equal(t, e.cause, e.Cause())
	assert.Equal(t, e.cause, e.Unwrap())
}
