// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"github.com/pkg/errors"
	"net/http"
	"testing"
)

func TestCompositeResponseObserver_Success(t *testing.T) {
	c := CompositeResponseObserver{
		NewMockResponseObserver(t),
		NewMockResponseObserver(t),
		NewMockResponseObserver(t),
		NewMockResponseObserver(t),
	}

	for _, mockResponseObserver := range c {
		mockResponseObserver.(*MockResponseObserver).EXPECT().Success(http.StatusOK).Return()
	}

	c.Success(http.StatusOK)
}

func TestCompositeResponseObserver_Error(t *testing.T) {
	err := errors.New("some error")

	c := CompositeResponseObserver{}

	for i := 0; i < 5; i++ {
		mockResponseObserver := NewMockResponseObserver(t)
		mockResponseObserver.EXPECT().Error(http.StatusBadRequest, err).Return()
		c = append(c, mockResponseObserver)
	}

	c.Error(http.StatusBadRequest, err)
}

func TestTracingResponseObserver_Success(t *testing.T) {
	ctx := context.Background()
	ctx, span := trace.NewSpan(ctx, "success")
	defer span.Finish()

	o := TracingResponseObserver{Context: ctx}
	o.Success(http.StatusOK)
}

func TestTracingResponseObserver_Error(t *testing.T) {
	ctx := context.Background()
	ctx, span := trace.NewSpan(ctx, "error")
	defer span.Finish()

	o := TracingResponseObserver{Context: ctx}
	o.Error(http.StatusBadRequest, errors.New("some error"))
}

func TestLoggingResponseObserver_Success(t *testing.T) {
	ctx := context.Background()
	o := LoggingResponseObserver{Context: ctx}
	o.Success(http.StatusOK)
}

func TestLoggingResponseObserver_Error(t *testing.T) {
	ctx := context.Background()

	o := LoggingResponseObserver{Context: ctx}
	o.Error(http.StatusBadRequest, errors.New("some error"))
}
