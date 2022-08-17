// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package mock

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/opentracing/opentracing-go/mocktracer"
	"reflect"
	"time"
)

type MockTracer struct {
	MockTracer *mocktracer.MockTracer
}

func (n *MockTracer) Configure(ctx context.Context, tracingConfig *trace.TracingConfig) error {
	return nil
}

func (n *MockTracer) LogContext(span trace.Span) map[string]interface{} {
	return nil
}

func (n *MockTracer) StartSpan(operationName string, options ...trace.StartSpanOption) trace.Span {
	var startSpanConfig trace.StartSpanConfig
	for _, option := range options {
		option(&startSpanConfig)
	}

	var openTracingOptions []opentracing.StartSpanOption
	if !startSpanConfig.StartTime.IsZero() {
		openTracingOptions = append(openTracingOptions, opentracing.StartTime(startSpanConfig.StartTime))
	}
	if startSpanConfig.Tags != nil {
		for k, v := range startSpanConfig.Tags {
			openTracingOptions = append(openTracingOptions, opentracing.Tag{Key: k, Value: v})
		}
	}
	for _, related := range startSpanConfig.Related {
		spanContext := related.Ref.(MockSpanContext)
		switch related.Type {
		case trace.RefChildOf:
			openTracingOptions = append(openTracingOptions, opentracing.ChildOf(spanContext.MockSpanContext))
		case trace.RefFollowsFrom:
			openTracingOptions = append(openTracingOptions, opentracing.FollowsFrom(spanContext.MockSpanContext))
		}
	}

	openTracingSpan := n.MockTracer.StartSpan(operationName, openTracingOptions...)
	return &MockSpan{
		MockSpan: openTracingSpan.(*mocktracer.MockSpan),
	}
}

func (n *MockTracer) Extract(_ trace.TextMapCarrier) (trace.SpanContext, error) {
	return nil, trace.ErrNoTracer
}

func (n *MockTracer) Inject(spanContext trace.SpanContext, carrier trace.TextMapCarrier) error {
	otSpanContext := spanContext.(MockSpanContext).MockSpanContext
	return n.MockTracer.Inject(otSpanContext, opentracing.TextMap, carrier)
}

func (n *MockTracer) Shutdown(ctx context.Context) error {
	return nil
}

func NewMockTracer() *MockTracer {
	return &MockTracer{
		MockTracer: mocktracer.New(),
	}
}

type MockSpan struct {
	MockSpan *mocktracer.MockSpan
}

func (n *MockSpan) Finish(option ...trace.FinishSpanOption) {
	n.MockSpan.Finish()
}

func (n *MockSpan) SetTag(key string, value interface{}) {
	n.MockSpan.SetTag(key, value)
}

func (n *MockSpan) Context() trace.SpanContext {
	otSpanContext := n.MockSpan.Context().(mocktracer.MockSpanContext)
	return MockSpanContext{
		MockSpanContext: otSpanContext,
	}
}

func (n *MockSpan) SetError(err error) {
	n.MockSpan.SetTag(trace.FieldError, err)
}

func (n *MockSpan) LogFields(i ...log.Field) {
	n.MockSpan.LogFields(i...)
}

func (n *MockSpan) LogKV(i ...interface{}) {
	n.MockSpan.LogKV(i...)
}

func (n *MockSpan) Tags() map[string]interface{} {
	return n.MockSpan.Tags()
}

func (n *MockSpan) Logs() []LogRecord {
	var results []LogRecord
	for _, entry := range n.MockSpan.Logs() {
		result := LogRecord{
			Timestamp: entry.Timestamp,
		}
		for _, field := range entry.Fields {
			result.Fields = append(result.Fields, LogField(field))
		}
		results = append(results, result)
	}
	return results
}

type MockSpanContext struct {
	MockSpanContext mocktracer.MockSpanContext
}

func (n MockSpanContext) SpanId() trace.SpanId {
	return trace.SpanId(n.MockSpanContext.TraceID)
}

func (n MockSpanContext) TraceId() trace.TraceId {
	return trace.TraceId{Low: uint64(n.MockSpanContext.SpanID)}
}

func (n MockSpanContext) ForeachBaggageItem(fn func(k string, v string) bool) {
	n.MockSpanContext.ForeachBaggageItem(fn)
}

type LogRecord struct {
	Timestamp time.Time
	Fields    []LogField
}

type LogField struct {
	Key         string
	ValueKind   reflect.Kind
	ValueString string
}
