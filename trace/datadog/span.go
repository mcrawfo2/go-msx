// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package datadog

import (
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	tracelog "github.com/opentracing/opentracing-go/log"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	datadog "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Span struct {
	DataDogSpan ddtrace.Span
	Tracer      trace.Tracer
	Error       error
}

func (s *Span) SetError(err error) {
	s.Error = err
}

// LogFields is not supported by datadog
func (s *Span) LogFields(_ ...tracelog.Field) {}

// LogKV is not supported by datadog
func (s *Span) LogKV(_ ...interface{}) {}

func (s *Span) Finish(options ...trace.FinishSpanOption) {
	// Collect finish options
	var finishSpanConfig trace.FinishSpanConfig
	for _, option := range options {
		option(&finishSpanConfig)
	}

	// Convert to datadog
	var dataDogOptions []ddtrace.FinishOption
	if !finishSpanConfig.FinishTime.IsZero() {
		dataDogOptions = append(dataDogOptions, datadog.FinishTime(finishSpanConfig.FinishTime))
	}
	if nil != finishSpanConfig.Error {
		dataDogOptions = append(dataDogOptions, datadog.WithError(finishSpanConfig.Error))
	} else if nil != s.Error {
		dataDogOptions = append(dataDogOptions, datadog.WithError(s.Error))
	}

	// Finish the datadog span
	s.DataDogSpan.Finish(dataDogOptions...)
}

func (s *Span) SetTag(key string, value interface{}) {
	s.DataDogSpan.SetTag(key, value)
}

func (s *Span) Context() trace.SpanContext {
	return SpanContext{
		DataDogSpanContext: s.DataDogSpan.Context(),
	}
}

type SpanContext struct {
	DataDogSpanContext ddtrace.SpanContext
}

func (s SpanContext) SpanId() trace.SpanId {
	return trace.SpanId(s.DataDogSpanContext.SpanID())
}

func (s SpanContext) TraceId() trace.TraceId {
	return trace.TraceId{
		High: 0,
		Low:  s.DataDogSpanContext.TraceID(),
	}
}

func (s SpanContext) ForeachBaggageItem(fn func(k string, v string) bool) {
	s.DataDogSpanContext.ForeachBaggageItem(fn)
}
