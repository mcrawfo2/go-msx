package tracetest

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/thejerf/abtime"
	"time"
)

type SpanContext struct{
	Baggage       map[string]string
}

func (n SpanContext) ForeachBaggageItem(handler func(k, v string) bool) {
	for k, v := range n.Baggage {
		if !handler(k, v) {
			break
		}
	}
}

type Span struct {
	OperationName string
	StartTime     time.Time
	FinishTime    time.Time
	Ctx           *SpanContext
	Tags          map[string]interface{}
	Logs          []opentracing.LogRecord
	tracer        *Tracer
}

func (s Span) Context() opentracing.SpanContext {
	return s.Ctx
}

func (s *Span) SetBaggageItem(key, val string) opentracing.Span {
	s.Ctx.Baggage[key] = val
	return s
}

func (s Span) BaggageItem(key string) string {
	return s.Ctx.Baggage[key]
}

func (s *Span) SetTag(key string, value interface{}) opentracing.Span {
	s.Tags[key] = value
	return s
}

func (s *Span) LogFields(fields ...log.Field) {
	s.Logs = append(s.Logs, opentracing.LogRecord{
		Fields:    fields,
		Timestamp: s.tracer.Clock.Now(),
	})
}
func (s *Span) LogKV(keyVals ...interface{}) {
	s.Logs = append(s.Logs, opentracing.LogRecord{
		Timestamp: s.tracer.Clock.Now(),
		Fields: func() []log.Field {
			fields, _ := log.InterleavedKVToFields(keyVals)
			return fields
		}(),
	})
}

func (s *Span) Finish() {
	s.FinishWithOptions(opentracing.FinishOptions{})
}

func (s *Span) FinishWithOptions(opts opentracing.FinishOptions) {
	if opts.FinishTime.IsZero() {
		s.FinishTime = s.tracer.Clock.Now()
	} else {
		s.FinishTime = opts.FinishTime
	}

	s.Logs = append(s.Logs, opts.LogRecords...)
	for _, ld := range opts.BulkLogData {
		s.Logs = append(s.Logs, ld.ToLogRecord())
	}
}

func (s *Span) SetOperationName(operationName string) opentracing.Span {
	s.OperationName = operationName
	return s
}

func (s Span) Tracer() opentracing.Tracer {
	return s.tracer
}

func (s *Span) LogEvent(event string) {
	s.Log(opentracing.LogData{Event: event})
}

func (s *Span) LogEventWithPayload(event string, payload interface{}) {
	s.Log(opentracing.LogData{Event: event, Payload: payload})
}

func (s *Span) Log(data opentracing.LogData) {
	if data.Timestamp.IsZero() {
		data.Timestamp = s.tracer.Clock.Now()
	}
	s.Logs = append(s.Logs, data.ToLogRecord())
}

type Tracer struct {
	Clock abtime.AbstractTime
	Spans []*Span
}

// StartSpan belongs to the Tracer interface.
func (t *Tracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	result := new(Span)
	result.OperationName = operationName
	result.Ctx = new(SpanContext)
	result.Ctx.Baggage = make(map[string]string)

	startSpanOptions := new(opentracing.StartSpanOptions)

	for _, opt := range opts {
		opt.Apply(startSpanOptions)
	}

	if startSpanOptions.StartTime.IsZero() {
		result.StartTime = t.Clock.Now()
	} else {
		result.StartTime = startSpanOptions.StartTime
	}

	result.Tags = make(map[string]interface{})
	for k, v := range startSpanOptions.Tags {
		result.Tags[k] = v
	}

	result.tracer = t

	t.Spans = append(t.Spans, result)

	return result
}

func (t Tracer) Inject(sp opentracing.SpanContext, format interface{}, carrier interface{}) error {
	return nil
}

func (t Tracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	return nil, opentracing.ErrSpanContextNotFound
}

func (t *Tracer) Reset() {
	t.Spans = make([]*Span, 0, 32)
}

func RecordTracing() *Tracer {
	var result = &Tracer{
		Clock: abtime.NewRealTime(),
		Spans: make([]*Span, 0, 32),
	}

	opentracing.SetGlobalTracer(result)

	return result
}
