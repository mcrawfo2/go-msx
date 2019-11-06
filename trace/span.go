package trace

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

func NewSpan(ctx context.Context, operationName string, options ...opentracing.StartSpanOption) (context.Context, opentracing.Span) {
	var span opentracing.Span
	if span = opentracing.SpanFromContext(ctx); span != nil {
		// Create a child span inside the existing parent
		options = append(options, opentracing.ChildOf(span.Context()))
		span = opentracing.StartSpan(operationName, options...)
	} else {
		// Create a new root span
		span = opentracing.StartSpan(operationName, options...)
	}
	ctx = opentracing.ContextWithSpan(ctx, span)

	if spanContext, ok := span.Context().(jaeger.SpanContext); ok {
		// Inject log fields
		ctx = log.ExtendContext(ctx, log.LogContext{
			"spanId":   spanContext.SpanID(),
			"traceId":  spanContext.TraceID(),
			"parentId": spanContext.ParentID(),
		})
	}

	return ctx, span
}

func SpanContextFromContext(ctx context.Context) *jaeger.SpanContext {
	var span opentracing.Span
	if span = opentracing.SpanFromContext(ctx); span == nil {
		return nil
	}

	if spanContext, ok := span.Context().(jaeger.SpanContext); ok {
		return &spanContext
	}

	return nil
}

func SpanFromContext(ctx context.Context) opentracing.Span {
	return opentracing.SpanFromContext(ctx)
}