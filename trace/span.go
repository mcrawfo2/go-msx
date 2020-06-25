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

func BackgroundOperation(ctx context.Context, operationName string, operation func(context.Context) error) {
	newCtx := UntracedContextFromContext(ctx)
	go func() {
		defer func() {
			c := recover()
			if c != nil {
				logger.WithContext(newCtx).Errorf("Operation %q panicked: %+v", operationName, c)
			}
		}()

		err := Operation(newCtx, operationName, operation)
		if err != nil {
			logger.WithContext(newCtx).WithError(err).Errorf("Operation %q failed", operationName)
		}
	}()
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
	untracedContext, ok := ctx.Value(contextKeyUntracedContext).(context.Context)
	if !ok {
		logger.WithContext(ctx).Error("Context does not have untraced context stored")
		return nil
	}
	return ContextWithUntracedContext(untracedContext)
}
