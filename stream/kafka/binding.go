package kafka

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"fmt"
)

const (
	configRootKafkaBindings = "spring.cloud.stream.kafka.bindings"
)

type BindingProducerConfig struct {
	Sync bool `config:"default=true"`
}

type BindingConfiguration struct {
	Producer            BindingProducerConfig
	StreamBindingConfig *stream.BindingConfiguration `config:"-"`
}

func NewBindingConfiguration(ctx context.Context, key string, streamBindingConfig *stream.BindingConfiguration) (*BindingConfiguration, error) {
	cfg := config.FromContext(ctx)
	return NewBindingConfigurationFromConfig(cfg, key, streamBindingConfig)
}

// Deprecated: NewBindingConfigurationFromConfig should be replaced by NewBindingConfiguration.
func NewBindingConfigurationFromConfig(cfg *config.Config, key string, streamBindingConfig *stream.BindingConfiguration) (*BindingConfiguration, error) {
	prefix := fmt.Sprintf("%s.%s", configRootKafkaBindings, key)

	bindingConfig := &BindingConfiguration{}
	if err := cfg.Populate(bindingConfig, prefix); err != nil {
		return nil, err
	}

	bindingConfig.StreamBindingConfig = streamBindingConfig

	return bindingConfig, nil
}
