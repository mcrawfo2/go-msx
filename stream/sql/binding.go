// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sql

import (
	"context"
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

func NewBindingConfiguration(ctx context.Context, key string, streamBindingConfig *stream.BindingConfiguration) (*BindingConfiguration, error) {
	cfg := config.FromContext(ctx)
	return NewBindingConfigurationFromConfig(cfg, key, streamBindingConfig)
}

// Deprecated: NewBindingConfigurationFromConfig should be replaced by NewBindingConfiguration.
func NewBindingConfigurationFromConfig(cfg *config.Config, key string, streamBindingConfig *stream.BindingConfiguration) (*BindingConfiguration, error) {
	prefix := fmt.Sprintf("%s.%s", configRootSqlBindings, key)

	bindingConfig := &BindingConfiguration{}
	if err := cfg.Populate(bindingConfig, prefix); err != nil {
		return nil, err
	}

	bindingConfig.StreamBindingConfig = streamBindingConfig

	return bindingConfig, nil
}
