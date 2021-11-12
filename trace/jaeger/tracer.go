package jaeger

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	jaegertransport "github.com/uber/jaeger-client-go/transport"
	jaegerzipkin "github.com/uber/jaeger-client-go/transport/zipkin"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"io"
	"strings"
)

var (
	logger       = log.NewLogger("msx.trace.jaeger")
	jaegerLogger = log.NewLogger("jaeger")
)

type tracer struct {
	closer io.Closer
}

func (t *tracer) LogContext(span trace.Span) map[string]interface{} {
	return nil
}

func (t *tracer) Configure(ctx context.Context, tracingConfig *trace.TracingConfig) error {
	jaegerConfig := tracingConfig.ToJaegerConfig()
	jaegerSampler := jaeger.NewConstSampler(true)
	jaegerLogger := &loggerAdapter{logger: jaegerLogger}
	var reporter jaegerconfig.Option
	if tracingConfig.Reporter.Enabled {
		jaegerTransport, err := t.newTransport(ctx, tracingConfig.Reporter)
		if err != nil {
			return err
		}
		reporter = jaegerconfig.Reporter(jaeger.NewRemoteReporter(jaegerTransport))
	} else {
		reporter = jaegerconfig.Reporter(jaeger.NewNullReporter())
	}
	jaegerMetricsFactory := prometheus.New()

	closer, err := jaegerConfig.InitGlobalTracer(
		jaegerConfig.ServiceName,
		jaegerconfig.Sampler(jaegerSampler),
		reporter,
		jaegerconfig.Logger(jaegerLogger),
		jaegerconfig.Metrics(jaegerMetricsFactory),
		jaegerconfig.ZipkinSharedRPCSpan(true),
	)
	if err != nil {
		return err
	}

	t.closer = closer
	return nil
}

func (t *tracer) newTransport(ctx context.Context, reporterConfig trace.TracingReporterConfig) (jaeger.Transport, error) {
	switch reporterConfig.Name {
	case "zipkin":
		logger.WithContext(ctx).Infof("Sending traces to zipkin: %q", reporterConfig.Url)
		return jaegerzipkin.NewHTTPTransport(reporterConfig.Url, jaegerzipkin.HTTPBatchSize(1))
	case "jaeger":
		logger.WithContext(ctx).Infof("Sending traces to jaeger: %q", reporterConfig.Address())
		return jaeger.NewUDPTransport(reporterConfig.Address(), 0)
	case "jaeger-http":
		logger.WithContext(ctx).Infof("Sending traces to jaeger: %q", reporterConfig.Address())
		return jaegertransport.NewHTTPTransport(reporterConfig.Url, jaegertransport.HTTPBatchSize(1)), nil
	}

	return nil, errors.New("Unknown transport: " + reporterConfig.Name)
}

func (t *tracer) StartSpan(operationName string, options ...trace.StartSpanOption) trace.Span {
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
		spanContext := related.Ref.(SpanContext)
		switch related.Type {
		case trace.RefChildOf:
			openTracingOptions = append(openTracingOptions, opentracing.ChildOf(spanContext.OpenTracingSpanContext))
		case trace.RefFollowsFrom:
			openTracingOptions = append(openTracingOptions, opentracing.FollowsFrom(spanContext.OpenTracingSpanContext))
		}
	}

	openTracingSpan := opentracing.StartSpan(operationName, openTracingOptions...)

	return &Span{
		OpenTracingSpan: openTracingSpan,
		Tracer:          t,
		Error:           nil,
	}
}

func (t *tracer) Extract(carrier trace.TextMapCarrier) (trace.SpanContext, error) {
	var traceIdValue, spanIdValue, parentSpanIdValue, sampledValue string
	err := carrier.ForeachKey(func(key, value string) error {
		name := strings.ToLower(key)
		switch name {
		case trace.HeaderTraceId:
			traceIdValue = strings.Trim(value, "\"")
		case trace.HeaderSpanId:
			spanIdValue = strings.Trim(value, "\"")
		case trace.HeaderParentSpanId:
			parentSpanIdValue = strings.Trim(value, "\"")
		case trace.HeaderSampled:
			sampledValue = strings.Trim(value, "\"")
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to extract trace context")
	}


	if traceIdValue == "" {
		return nil, errors.Wrapf(trace.ErrNoTrace, "Missing %s header", trace.HeaderTraceId)
	}
	traceId, err := jaeger.TraceIDFromString(traceIdValue)
	if err != nil {
		return nil, errors.Wrapf(trace.ErrInvalidTrace, "Invalid %s header", trace.HeaderTraceId)
	}

	if spanIdValue == "" {
		return nil, errors.Wrapf(trace.ErrNoTrace, "Missing %s header", trace.HeaderSpanId)
	}
	spanId, err := jaeger.SpanIDFromString(spanIdValue)
	if err != nil {
		return nil, errors.Wrapf(trace.ErrInvalidTrace, "Invalid %s header", trace.HeaderSpanId)
	}

	sampled := sampledValue == "1"

	parentSpanId, _ := jaeger.SpanIDFromString("0")
	if parentSpanIdValue != "" {
		parentSpanId, err = jaeger.SpanIDFromString(parentSpanIdValue)
	}
	if err != nil {
		return nil, errors.Wrapf(trace.ErrInvalidTrace, "Invalid %s header", trace.HeaderParentSpanId)
	}

	openTracingSpanContext := jaeger.NewSpanContext(traceId, spanId, parentSpanId, sampled, nil)
	return SpanContext{
		OpenTracingSpanContext: openTracingSpanContext,
		SpanID:                 trace.SpanId(spanId),
		TraceID:                trace.TraceId(traceId),
	}, nil
}

func (t *tracer) Inject(spanContext trace.SpanContext, carrier trace.TextMapCarrier) error {
	jaegerSpanContext := spanContext.(SpanContext).OpenTracingSpanContext.(jaeger.SpanContext)
	carrier.Set(trace.HeaderTraceId, fmt.Sprintf("%016x", jaegerSpanContext.TraceID().Low))
	carrier.Set(trace.HeaderSpanId, fmt.Sprintf("%016x", uint64(jaegerSpanContext.SpanID())))
	if jaegerSpanContext.IsSampled() {
		carrier.Set(trace.HeaderSampled, "1")
	} else {
		carrier.Set(trace.HeaderSampled, "0")
	}
	if jaegerSpanContext.ParentID() != 0 {
		carrier.Set(trace.HeaderParentSpanId, fmt.Sprintf("%016x", uint64(jaegerSpanContext.ParentID())))
	}
	return nil
}

func (t *tracer) Shutdown(_ context.Context) error {
	if t.closer != nil {
		return t.closer.Close()
	}
	return nil
}

type loggerAdapter struct {
	logger *log.Logger
}

func (j *loggerAdapter) Error(msg string) {
	j.logger.Error(msg)
}

func (j *loggerAdapter) Infof(msg string, args ...interface{}) {
	j.logger.Infof(msg, args...)
}

func RegisterTracer(ctx context.Context) error {
	trace.RegisterTracer("jaeger", new(tracer))
	return nil
}
