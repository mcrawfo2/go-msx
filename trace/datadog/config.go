// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package datadog

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
)

const configRootDatadog = "datadog"

type Config struct {
	ServiceName string `config:"default=${app.info.name}"`
	ServiceVersion string `config:"default=${app.build.version}"`
	ServiceEnv string `config:"default="`
}

func NewConfig(ctx context.Context) (*Config, error) {
	var cfg Config
	if err := config.FromContext(ctx).Populate(&cfg, configRootDatadog); err != nil {
		return nil, err
	}
	return &cfg, nil
}
