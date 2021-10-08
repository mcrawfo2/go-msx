package trace

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/pkg/errors"
)

type Tracer interface {
	Configure(ctx context.Context, tracingConfig *TracingConfig) error
	StartSpan(operationName string, options ...StartSpanOption) Span
	Extract(carrier TextMapCarrier) (SpanContext, error)
	Inject(spanContext SpanContext, carrier TextMapCarrier) error
	Shutdown(ctx context.Context) error
}

var tracers = make(map[string]Tracer)
var tracer Tracer = newNoopTracer()

func RegisterTracer(name string, t Tracer) {
	if t != nil {
		tracers[name] = t
	}
}

func SetTracer(t Tracer) {
	tracer = t
}

func ConfigureTracer(ctx context.Context) error {
	cfg := config.FromContext(ctx)

	tracingConfig, err := NewTracingConfig(cfg)
	if err != nil {
		return err
	}

	t, ok := tracers[tracingConfig.Collector]
	if !ok {
		return errors.Errorf("Unknown tracer: %q", tracingConfig.Collector)
	} else {
		tracer = t
	}

	return tracer.Configure(ctx, tracingConfig)
}

func ShutdownTracer(ctx context.Context) error {
	return tracer.Shutdown(ctx)
}
