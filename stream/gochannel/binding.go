package gochannel

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"fmt"
)

const (
	configRootGoChannelBindings = "spring.cloud.stream.gochannel.bindings"
)

type BindingConfiguration struct {
	Producer struct {
		OutputChannelBuffer            int64 `config:"default=16"`
		Persistent                     bool  `config:"default=false"`
		BlockPublishUntilSubscriberAck bool  `config:"default=false"`
	}
	StreamBindingConfig *stream.BindingConfiguration `config:"-"`
}

func NewBindingConfigurationFromConfig(cfg *config.Config, key string, streamBindingConfig *stream.BindingConfiguration) (*BindingConfiguration, error) {
	prefix := fmt.Sprintf("%s.%s", configRootGoChannelBindings, key)

	bindingConfig := &BindingConfiguration{}
	if err := cfg.Populate(bindingConfig, prefix); err != nil {
		return nil, err
	}

	bindingConfig.StreamBindingConfig = streamBindingConfig

	return bindingConfig, nil
}
