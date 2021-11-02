package trace

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"fmt"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
)

const (
	configRootTracing = "trace"
)

var (
	logger = log.NewLogger("msx.trace")
)

type TracingConfig struct {
	Enabled        bool   `config:"default=true"`
	ServiceName    string `config:"default=${info.app.name}"`
	ServiceVersion string `config:"default=${info.build.version}"`
	Collector      string `config:"default=jaeger"`
	Reporter       TracingReporterConfig
}

func (c TracingConfig) ToJaegerConfig() *jaegerconfig.Configuration {
	cfg := jaegerconfig.Configuration{
		ServiceName: c.ServiceName,
		Disabled:    !c.Enabled,
	}
	return &cfg
}

type TracingReporterConfig struct {
	Enabled bool   `config:"default=false"`
	Name    string `config:"default=jaeger"`
	Host    string `config:"default=localhost"`
	Port    int    `config:"default=6831"`
	Url     string `config:"default=http://localhost:9411/api/v1/spans"`
}

func (c TracingReporterConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func NewTracingConfig(cfg *config.Config) (*TracingConfig, error) {
	var tracingConfig TracingConfig
	if err := cfg.Populate(&tracingConfig, configRootTracing); err != nil {
		return nil, err
	}

	return &tracingConfig, nil
}
