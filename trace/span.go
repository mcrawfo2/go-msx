package trace

import (
	"bytes"
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	tracelog "github.com/opentracing/opentracing-go/log"
	"time"
)

type TraceId struct {
	High, Low uint64
}

func (t TraceId) String() string {
	var byteBuffer bytes.Buffer
	if t.High != 0 {
		byteBuffer.WriteString(fmt.Sprintf("%016x", t.High))
	}
	byteBuffer.WriteString(fmt.Sprintf("%016x", t.Low))
	return byteBuffer.String()
}

type SpanId uint64

func (s SpanId) String() string {
	return fmt.Sprintf("%016x", uint64(s))
}

type Span interface {
	Finish(...FinishSpanOption)
	SetTag(key string, value interface{})
	Context() SpanContext
	SetError(err error)
	LogFields(...tracelog.Field)
	LogKV(...interface{})
}

type SpanContext interface {
	SpanId() SpanId
	TraceId() TraceId
	ForeachBaggageItem(fn func(k, v string) bool)
}

type SpanReference struct {
	Type string
	Ref  SpanContext
}

type StartSpanConfig struct {
	Related   []SpanReference
	StartTime time.Time
	Tags      map[string]interface{}
}

type StartSpanOption func(*StartSpanConfig)

func StartWithRelated(refType string, refSpanContext SpanContext) StartSpanOption {
	return func(options *StartSpanConfig) {
		options.Related = append(options.Related, SpanReference{
			Type: refType,
			Ref:  refSpanContext,
		})
	}
}

func StartWithStartTime(t time.Time) StartSpanOption {
	return func(options *StartSpanConfig) {
		options.StartTime = t
	}
}

func StartWithTag(field string, value interface{}) StartSpanOption {
	return func(options *StartSpanConfig) {
		if options.Tags == nil {
			options.Tags = make(map[string]interface{})
		}
		options.Tags[field] = value
	}
}

type FinishSpanConfig struct {
	FinishTime time.Time
	Error      error
}

type FinishSpanOption func(*FinishSpanConfig)

func FinishWithFinishTime(t time.Time) FinishSpanOption {
	return func(options *FinishSpanConfig) {
		options.FinishTime = t
	}
}

func FinishWithError(err error) FinishSpanOption {
	return func(options *FinishSpanConfig) {
		options.Error = err
	}
}

func NewSpan(ctx context.Context, operationName string, options ...StartSpanOption) (context.Context, Span) {
	var span Span
	var parentSpan Span
	if parentSpan = SpanFromContext(ctx); parentSpan != nil {
		// Create a child span inside the existing parent
		options = append(options, StartWithRelated(RefChildOf, parentSpan.Context()))
	}

	span = tracer.StartSpan(operationName, options...)
	ctx = ContextWithSpan(ctx, span)
	ctx = contextWithTraceLogContext(ctx, tracer, span, parentSpan)

	return ctx, span
}

func baseLogContext(span, parentSpan Span) log.LogContext {
	// Calculate log fields
	spanId := span.Context().SpanId().String()
	traceId := span.Context().TraceId().String()
	parentSpanId := SpanId(0)
	if parentSpan != nil {
		parentSpanId = parentSpan.Context().SpanId()
	}
	parentId := parentSpanId.String()

	return log.LogContext{
		log.FieldSpanId:   spanId,
		log.FieldTraceId:  traceId,
		log.FieldParentId: parentId,
	}
}

func contextWithTraceLogContext(ctx context.Context, tracer Tracer, span, parentSpan Span) context.Context {
	if span.Context().SpanId() == SpanId(0) {
		// NoopTracer is in effect
		return ctx
	}

	// Generate generic log context
	logContext := baseLogContext(span, parentSpan)

	// Add tracer-specific log context
	for k, v := range tracer.LogContext(span) {
		logContext[k] = v
	}

	// Inject log fields
	return log.ExtendContext(ctx, logContext)
}

func SpanDecorator(operationName string, options ...StartSpanOption) types.ActionFuncDecorator {
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
