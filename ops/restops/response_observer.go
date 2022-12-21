// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
)

//go:generate mockery --name=ResponseObserver --inpackage --case=snake --testonly --with-expecter

type ResponseObserver interface {
	Success(code int)
	Error(code int, responseError error)
}

type CompositeResponseObserver []ResponseObserver

func (c CompositeResponseObserver) Success(code int) {
	for _, observer := range c {
		observer.Success(code)
	}
}

func (c CompositeResponseObserver) Error(code int, responseError error) {
	for _, observer := range c {
		observer.Error(code, responseError)
	}
}

type TracingResponseObserver struct {
	Context context.Context
}

func (t TracingResponseObserver) Success(code int) {
	span := trace.SpanFromContext(t.Context)
	if span != nil {
		span.LogFields(trace.Int(trace.FieldHttpCode, code))
	}
}

func (t TracingResponseObserver) Error(code int, responseError error) {
	span := trace.SpanFromContext(t.Context)
	if span != nil {
		span.LogFields(trace.Int(trace.FieldHttpCode, code))
		span.LogFields(trace.Error(responseError))
	}
}

type LoggingResponseObserver struct {
	Context context.Context
}

func (l LoggingResponseObserver) Success(code int) {
	// TODO: Recorded by the tracingFilter
}

func (l LoggingResponseObserver) Error(code int, responseError error) {
	// TODO: Recorded by the tracingFilter
}
