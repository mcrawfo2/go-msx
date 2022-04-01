// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package stream

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/retry"
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
	BindingId   string `config:"default=${spring.application.instance}"`
	LogMessages bool   `config:"default=true"`
	Retry       retry.RetryConfig
	Consumer    ConsumerConfiguration
}

type ConsumerConfiguration struct {
	AutoStartup            bool    `config:"default=${spring.cloud.stream.default.consumer.auto-startup:true}"`
	Concurrency            int     `config:"default=${spring.cloud.stream.default.consumer.concurrency:1}"`
	Partitioned            bool    `config:"default=${spring.cloud.stream.default.consumer.partitioned:false}"`
	HeaderMode             string  `config:"default=${spring.cloud.stream.default.consumer.header-mode:none}"`
	MaxAttempts            int     `config:"default=${spring.cloud.stream.default.consumer.max-attempts:3}"`
	BackOffInitialInterval int     `config:"default=${spring.cloud.stream.default.consumer.backoff-initial-interval:1000}"`
	BackOffMaxInterval     int     `config:"default=${spring.cloud.stream.default.consumer.backoff-max-interval:10000}"`
	BackOffMultiplier      float32 `config:"default=${spring.cloud.stream.default.consumer.backoff-multiplier:2.0}"`
	DefaultRetryable       bool    `config:"default=${spring.cloud.stream.default.consumer.default-retryable:true}"`
	InstanceIndex          int     `config:"default=${spring.cloud.stream.default.consumer.instance-index:-1}"`
	InstanceCount          int     `config:"default=${spring.cloud.stream.default.consumer.instance-count:-1}"`
}

func NewBindingConfiguration(ctx context.Context, key string) (*BindingConfiguration, error) {
	cfg := config.FromContext(ctx)
	if cfg == nil {
		return nil, config.ErrNotLoaded
	}
	return NewBindingConfigurationFromConfig(cfg, key)
}

// Deprecated: NewBindingConfigurationFromConfig should be replaced with NewBindingConfiguration.
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

	bindingConfig.Group = key + "-" + bindingConfig.Group

	return bindingConfig, nil
}
