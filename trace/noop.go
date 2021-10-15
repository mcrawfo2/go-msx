package trace

import (
	"context"
	"github.com/opentracing/opentracing-go"
	tracelog "github.com/opentracing/opentracing-go/log"
)

type noopTracer struct {
	NoopTracer opentracing.NoopTracer
}

func (n *noopTracer) Configure(ctx context.Context, tracingConfig *TracingConfig) error {
	return nil
}

func (n *noopTracer) LogContext(span Span) map[string]interface{} {
	return nil
}

func (n *noopTracer) StartSpan(operationName string, options ...StartSpanOption) Span {
	otSpan := n.NoopTracer.StartSpan(operationName)
	return &noopSpan{
		NoopSpan: otSpan,
	}
}

func (n *noopTracer) Extract(_ TextMapCarrier) (SpanContext, error) {
	return nil, ErrNoTracer
}

func (n *noopTracer) Inject(spanContext SpanContext, carrier TextMapCarrier) error {
	otSpanContext := spanContext.(*noopSpanContext).NoopSpanContext
	return n.NoopTracer.Inject(otSpanContext, opentracing.TextMap, carrier)
}

func (n *noopTracer) Shutdown(ctx context.Context) error {
	return nil
}

func newNoopTracer() *noopTracer {
	return &noopTracer{
		NoopTracer: opentracing.NoopTracer{},
	}
}

type noopSpan struct {
	NoopSpan opentracing.Span
}

func (n noopSpan) Finish(option ...FinishSpanOption) {
	n.NoopSpan.Finish()
}

func (n noopSpan) SetTag(key string, value interface{}) {
	n.NoopSpan.SetTag(key, value)
}

func (n noopSpan) Context() SpanContext {
	otSpanContext := n.NoopSpan.Context()
	return &noopSpanContext{
		NoopSpanContext: otSpanContext,
	}
}

func (n noopSpan) SetError(err error) {
	n.NoopSpan.SetTag(FieldError, err)
}

func (n noopSpan) LogFields(i ...tracelog.Field) {}

func (n noopSpan) LogKV(i ...interface{}) {}

type noopSpanContext struct {
	NoopSpanContext opentracing.SpanContext
}

func (n *noopSpanContext) SpanId() SpanId {
	return SpanId(0)
}

func (n *noopSpanContext) TraceId() TraceId {
	return TraceId{0, 0}
}

func (n *noopSpanContext) ForeachBaggageItem(fn func(k string, v string) bool) {
}
