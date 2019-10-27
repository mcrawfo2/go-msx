package stream

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"strings"
)

const (
	configRootSubscriberBindings = "spring.cloud.stream.bindings"
	configKeyAppName             = "spring.application.name"
)

type BindingConfiguration struct {
	Destination string `config:"default="`                 // Topic if different from binding key
	Group       string `config:"default="`                 // Consumer group id
	ContentType string `config:"default=application/json"` // Content-Type Header
	Binder      string `config:"default=kafka"`            // Stream Provider
	Retry       types.Retry
}

func NewBindingConfigurationFromConfig(cfg *config.Config, key string) (*BindingConfiguration, error) {
	prefix := fmt.Sprintf("%s.%s", configRootSubscriberBindings, key)

	bindingConfig := &BindingConfiguration{}
	if err := cfg.Populate(bindingConfig, prefix); err != nil {
		return nil, err
	}

	if bindingConfig.Destination == "" {
		// Default topic name to binding key
		bindingConfig.Destination = key
	}

	if bindingConfig.Group == "" {
		// Derive consumer group name automatically
		if appName, err := cfg.String(configKeyAppName); err != nil {
			return nil, err
		} else {
			bindingConfig.Group = strings.TrimSuffix(strings.ToUpper(appName), "SERVICE") + "_GP"
		}
	}

	return bindingConfig, nil
}
