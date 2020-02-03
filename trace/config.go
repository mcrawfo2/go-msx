package trace

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	jaegerzipkin "github.com/uber/jaeger-client-go/transport/zipkin"
	"github.com/uber/jaeger-client-go/zipkin"
	"github.com/uber/jaeger-lib/metrics"
	"io"
)

const (
	configRootTracing    = "trace"
	configKeyInfoAppName = "info.app.name"
)

var (
	logger       = log.NewLogger("msx.trace")
	jaegerLogger = log.NewLogger("jaeger")
	jaegerCloser io.Closer
)

type TracingConfig struct {
	Enabled     bool   `config:"default=false"`
	ServiceName string `config:"default="`
	Reporter    TracingReporterConfig
}

func (c TracingConfig) ToJaegerConfig() *jaegerconfig.Configuration {
	cfg := jaegerconfig.Configuration{
		ServiceName: c.ServiceName,
		Disabled:    !c.Enabled,
	}
	return &cfg
}

type TracingReporterConfig struct {
	Name string `config:"default=jaeger"`
	Host string `config:"default=localhost"`
	Port int    `config:"default=6831"`
	Url  string `config:"default=http://localhost:9411/api/v1/spans"`
}

func (c TracingReporterConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type JaegerLoggerAdapter struct {
	logger *log.Logger
}

func (j *JaegerLoggerAdapter) Error(msg string) {
	j.logger.Error(msg)
}

func (j *JaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	j.logger.Infof(msg, args...)
}

func NewTracingConfig(cfg *config.Config) (*TracingConfig, error) {
	var tracingConfig TracingConfig
	if err := cfg.Populate(&tracingConfig, configRootTracing); err != nil {
		return nil, err
	}

	serviceName, err := cfg.String(configKeyInfoAppName)
	if err != nil {
		return nil, err
	}
	tracingConfig.ServiceName = serviceName

	return &tracingConfig, nil
}

func newTransport(ctx context.Context, reporterConfig TracingReporterConfig) (jaeger.Transport, error) {
	switch reporterConfig.Name {
	case "zipkin":
		logger.WithContext(ctx).Infof("Sending traces to zipkin: %q", reporterConfig.Url)
		return jaegerzipkin.NewHTTPTransport(reporterConfig.Url, jaegerzipkin.HTTPBatchSize(1))
	case "jaeger":
		logger.WithContext(ctx).Infof("Sending traces to jaeger: %q", reporterConfig.Address())
		return jaeger.NewUDPTransport(reporterConfig.Address(), 0)
	}

	return nil, errors.New("Unknown transport: " + reporterConfig.Name)
}

func ConfigureTracer(ctx context.Context) error {
	cfg := config.FromContext(ctx)

	tracingConfig, err := NewTracingConfig(cfg)
	if err != nil {
		return err
	}

	jaegerConfig := tracingConfig.ToJaegerConfig()
	jaegerSampler := jaeger.NewConstSampler(true)
	jaegerLogger := &JaegerLoggerAdapter{logger: jaegerLogger}
	jaegerTransport, err := newTransport(ctx, tracingConfig.Reporter)
	if err != nil {
		return err
	}
	jaegerRemoteReporter := jaeger.NewRemoteReporter(jaegerTransport)
	jaegerMetricsFactory := metrics.NullFactory
	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	sleuthPropagator := NewSleuthTextMapPropagator(zipkinPropagator)

	closer, err := jaegerConfig.InitGlobalTracer(
		jaegerConfig.ServiceName,
		jaegerconfig.Sampler(jaegerSampler),
		jaegerconfig.Reporter(jaegerRemoteReporter),
		jaegerconfig.Logger(jaegerLogger),
		jaegerconfig.Metrics(jaegerMetricsFactory),
		jaegerconfig.Injector(opentracing.HTTPHeaders, zipkinPropagator),
		jaegerconfig.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
		jaegerconfig.Injector(opentracing.TextMap, sleuthPropagator),
		jaegerconfig.Extractor(opentracing.TextMap, sleuthPropagator),
		jaegerconfig.ZipkinSharedRPCSpan(true),
	)
	if err != nil {
		return err
	}

	jaegerCloser = closer

	return nil
}

func ShutdownTracer(ctx context.Context) error {
	if jaegerCloser != nil {
		defer jaegerCloser.Close()
	}
	return nil
}
