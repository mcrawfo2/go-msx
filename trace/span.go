package trace

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/opentracing/opentracing-go"
	tracelog "github.com/opentracing/opentracing-go/log"
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

func Operation(ctx context.Context, operationName string, operation func(context.Context) error) (err error) {
	ctx, span := NewSpan(ctx, operationName)
	defer span.Finish()

	err = operation(ctx)

	if err != nil {
		span.LogFields(Error(err))
	}

	return err
}

var Error = tracelog.Error
var Int = tracelog.Int
var String = tracelog.String

func HttpCode(code int) tracelog.Field {
	return Int(FieldHttpCode, code)
}

func Status(status string) tracelog.Field {
	return String(FieldStatus, status)
}

type contextTraceKey int

const (
	contextKeyUntracedContext contextTraceKey = iota
)

func ContextWithUntracedContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyUntracedContext, ctx)
}

func UntracedContextFromContext(ctx context.Context) context.Context {
	return ctx.Value(contextKeyUntracedContext).(context.Context)
}
