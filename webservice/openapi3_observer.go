package webservice

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
)

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

type LogResponseObserver struct {
	Context context.Context
}

func (l LogResponseObserver) Success(code int) {
	// Recorded by the tracingFilter
}

func (l LogResponseObserver) Error(code int, responseError error) {
	// TODO: Recorded by the tracingFilter?
	logger.
		WithContext(l.Context).
		WithField("status", code).
		WithError(responseError).
		Error("Request failed")
}
