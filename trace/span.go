package trace

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
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
			log.FieldSpanId:   spanContext.SpanID().String(),
			log.FieldTraceId:  spanContext.TraceID().String(),
			log.FieldParentId: spanContext.ParentID().String(),
		})
	}

	return ctx, span
}

func SpanFromContext(ctx context.Context) opentracing.Span {
	return opentracing.SpanFromContext(ctx)
}

func SpanDecorator(operationName string, options ...opentracing.StartSpanOption) types.ActionFuncDecorator {
	return func(action types.ActionFunc) types.ActionFunc {
		return func(ctx context.Context) error {
			ctx, span := NewSpan(ctx, operationName, options...)
			defer span.Finish()

			err := action(ctx)

			if err != nil {
				span.LogFields(Error(err))
			}

			return err
		}
	}
}

// BackgroundOperation executes the action inside a new detached span in a separate goroutine
func BackgroundOperation(ctx context.Context, operationName string, action types.ActionFunc) {
	go func() {
		ForegroundOperation(ctx, operationName, action)
	}()
}

// ForegroundOperation executes the action inside a new detached span
func ForegroundOperation(ctx context.Context, operationName string, action types.ActionFunc) {
	_ = NewIsolatedOperation(operationName, action).Run(ctx)
}

// NewIsolatedOperation creates a new operation completely separated from the current trace
func NewIsolatedOperation(operationName string, action types.ActionFunc) types.Operation {
	return NewOperation(operationName, action).
		WithDecorator(log.ErrorLogDecorator(logger, operationName)).
		WithFilter(types.NewOrderedDecorator(1000, UntracedContextDecorator))
}

// Operation executes the action inside a new child span
func Operation(ctx context.Context, operationName string, action types.ActionFunc) (err error) {
	return NewOperation(operationName, action).Run(ctx)
}

func NewOperation(operationName string, action types.ActionFunc) types.Operation {
	return types.NewOperation(action).
		WithDecorator(SpanDecorator(operationName)).
		WithDecorator(log.RecoverLogDecorator(logger))
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
	untracedContext, ok := ctx.Value(contextKeyUntracedContext).(context.Context)
	if !ok {
		logger.WithContext(ctx).Error("Context does not have untraced context stored")
		return nil
	}
	return ContextWithUntracedContext(untracedContext)
}

func UntracedContextDecorator(action types.ActionFunc) types.ActionFunc {
	return func(ctx context.Context) error {
		ctx = UntracedContextFromContext(ctx)
		return action(ctx)
	}
}
