// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package redis

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"fmt"
	"time"
)

const (
	configRootRedisBindings = "spring.cloud.stream.redis.bindings"
)

type BindingProducerConfig struct{}

type BindingConsumerConfig struct {
	NackResendSleep time.Duration `config:"default=0s"`
	MaxIdleTime     time.Duration `config:"default=60s"`
	ClientId        string        `config:"default=${spring.application.name}"`
	ClientIdSuffix  string        `config:"default=${spring.application.instance}"`
}

type BindingConfiguration struct {
	Producer            BindingProducerConfig
	Consumer            BindingConsumerConfig
	StreamBindingConfig *stream.BindingConfiguration `config:"-"`
}

func NewBindingConfiguration(ctx context.Context, key string, streamBindingConfig *stream.BindingConfiguration) (*BindingConfiguration, error) {
	cfg := config.FromContext(ctx)

	prefix := fmt.Sprintf("%s.%s", configRootRedisBindings, key)

	bindingConfig := &BindingConfiguration{}
	if err := cfg.Populate(bindingConfig, prefix); err != nil {
		return nil, err
	}

	bindingConfig.StreamBindingConfig = streamBindingConfig

	return bindingConfig, nil
}
