package jaeger

import (
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"github.com/opentracing/opentracing-go"
	tracelog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
)

type Span struct {
	OpenTracingSpan opentracing.Span
	Tracer          trace.Tracer
	Error           error
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

	// Convert to opentracing
	var openTracingOptions opentracing.FinishOptions
	if !finishSpanConfig.FinishTime.IsZero() {
		openTracingOptions.FinishTime = finishSpanConfig.FinishTime
	}
	if nil != finishSpanConfig.Error {
		s.OpenTracingSpan.SetTag(trace.FieldError, finishSpanConfig.Error)
	} else if nil != s.Error {
		s.OpenTracingSpan.SetTag(trace.FieldError, s.Error)
	}

	// Finish the opentracing span
	s.OpenTracingSpan.FinishWithOptions(openTracingOptions)
}

func (s *Span) SetTag(key string, value interface{}) {
	s.OpenTracingSpan.SetTag(key, value)
}

func (s *Span) Context() trace.SpanContext {
	return SpanContext{
		OpenTracingSpanContext: s.OpenTracingSpan.Context(),
	}
}

type SpanContext struct {
	OpenTracingSpanContext opentracing.SpanContext
	SpanID                 trace.SpanId
	TraceID                trace.TraceId
}

func (s SpanContext) SpanId() trace.SpanId {
	jaegerSpanContext, _ := s.OpenTracingSpanContext.(jaeger.SpanContext)
	return trace.SpanId(jaegerSpanContext.SpanID())
}

func (s SpanContext) TraceId() trace.TraceId {
	jaegerSpanContext, _ := s.OpenTracingSpanContext.(jaeger.SpanContext)
	return trace.TraceId(jaegerSpanContext.TraceID())
}

func (s SpanContext) ForeachBaggageItem(fn func(k string, v string) bool) {
	s.OpenTracingSpanContext.ForeachBaggageItem(fn)
}
