package sql

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"fmt"
	"time"
)

const (
	configRootSqlBindings = "spring.cloud.stream.sql.bindings"
)

type BindingProducerConfig struct{}

type BindingConsumerConfig struct {
	PollInterval   time.Duration `config:"default=1s"`
	ResendInterval time.Duration `config:"default=1s"`
	RetryInterval  time.Duration `config:"default=1s"`
}

type BindingConfiguration struct {
	Producer            BindingProducerConfig
	Consumer            BindingConsumerConfig
	StreamBindingConfig *stream.BindingConfiguration `config:"-"`
}

func NewBindingConfigurationFromConfig(cfg *config.Config, key string, streamBindingConfig *stream.BindingConfiguration) (*BindingConfiguration, error) {
	prefix := fmt.Sprintf("%s.%s", configRootSqlBindings, key)

	bindingConfig := &BindingConfiguration{}
	if err := cfg.Populate(bindingConfig, prefix); err != nil {
		return nil, err
	}

	bindingConfig.StreamBindingConfig = streamBindingConfig

	return bindingConfig, nil
}
