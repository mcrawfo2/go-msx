package datadog

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"github.com/pkg/errors"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	datadog "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"os"
	"regexp"
	"strconv"
)

var (
	logger        = log.NewLogger("msx.trace.datadog")
	datadogLogger = log.NewLogger("datadog.tracer")
)

type tracer struct {
	cfg *trace.TracingConfig
}

func convertId(id string) string {
	if len(id) < 16 {
		return ""
	}
	if len(id) > 16 {
		id = id[16:]
	}
	intValue, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatUint(intValue, 10)
}

func (t *tracer) LogContext(span trace.Span) map[string]interface{} {
	results := log.LogContext{
		"dd.trace_id": convertId(span.Context().SpanId().String()),
		"dd.span_id": convertId(span.Context().TraceId().String()),
		"dd.service": t.cfg.ServiceName,
		"dd.version": t.cfg.ServiceVersion,
	}

	return results
}

func (t *tracer) Extract(carrier trace.TextMapCarrier) (trace.SpanContext, error) {
	dataDogSpanContext, err := datadog.Extract(carrier)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to extract trace context")
	}

	return SpanContext{
		DataDogSpanContext: dataDogSpanContext,
	}, nil
}

func (t *tracer) Inject(spanContext trace.SpanContext, carrier trace.TextMapCarrier) error {
	ddSpanContext, _ := spanContext.(SpanContext)
	err := datadog.Inject(ddSpanContext.DataDogSpanContext, carrier)
	if err != nil {
		return errors.Wrap(err, "Failed to inject trace context")
	}
	return nil
}

func (t *tracer) Configure(ctx context.Context, tracingConfig *trace.TracingConfig) error {
	_ = os.Setenv("DD_PROPAGATION_STYLE_INJECT", "b3")
	_ = os.Setenv("DD_PROPAGATION_STYLE_EXTRACT", "b3")

	options := []datadog.StartOption{
		datadog.WithAgentAddr(tracingConfig.Reporter.Address()),
		datadog.WithService(tracingConfig.ServiceName),
		datadog.WithLogger(loggerAdapter{logger: datadogLogger}),
		datadog.WithServiceVersion(tracingConfig.ServiceVersion),
		datadog.WithSampler(datadog.NewAllSampler()),
	}

	if tracingConfig.Reporter.Enabled {
		options = append(options, []datadog.StartOption{
			datadog.WithAgentAddr(tracingConfig.Reporter.Address()),
			datadog.WithDebugMode(true),
		}...)
	}

	datadog.Start(options...)

	t.cfg = tracingConfig

	return nil
}

func (t *tracer) StartSpan(operationName string, options ...trace.StartSpanOption) trace.Span {
	var startSpanConfig trace.StartSpanConfig
	for _, option := range options {
		option(&startSpanConfig)
	}

	var dataDogOptions []ddtrace.StartSpanOption
	if !startSpanConfig.StartTime.IsZero() {
		dataDogOptions = append(dataDogOptions, datadog.StartTime(startSpanConfig.StartTime))
	}
	if startSpanConfig.Tags != nil {
		for k, v := range startSpanConfig.Tags {
			dataDogOptions = append(dataDogOptions, datadog.Tag(k, v))
		}
	}
	for _, related := range startSpanConfig.Related {
		spanContext := related.Ref.(SpanContext)
		switch related.Type {
		case trace.RefChildOf:
			dataDogOptions = append(dataDogOptions, datadog.ChildOf(spanContext.DataDogSpanContext))
		case trace.RefFollowsFrom:
			// Not supported by DataDog
			dataDogOptions = append(dataDogOptions, datadog.ChildOf(spanContext.DataDogSpanContext))
		}
	}

	dataDogSpan := datadog.StartSpan(operationName, dataDogOptions...)

	return &Span{
		DataDogSpan: dataDogSpan,
		Tracer:      t,
		Error:       nil,
	}
}

func (t *tracer) Shutdown(_ context.Context) error {
	datadog.Stop()
	return nil
}

type loggerAdapter struct {
	logger *log.Logger
	ctx    context.Context
}

var logMessageRegexp = regexp.MustCompile(`(.*) (DEBUG|WARN|INFO|ERROR): (.*)`)
func (l loggerAdapter) Log(msg string) {
	parts := logMessageRegexp.FindStringSubmatch(msg)
	if parts == nil {
		l.logger.WithContext(l.ctx).Info(msg)
	} else {
		level := log.LevelFromName(parts[2])
		msg = parts[3]
		l.logger.WithContext(l.ctx).Log(level, msg)
	}
}

func RegisterTracer(ctx context.Context) error {
	trace.RegisterTracer("datadog", new(tracer))
	return nil
}
