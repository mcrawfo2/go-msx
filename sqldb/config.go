// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
)

const configRootSpringDatasourceConfig = "spring.datasource"

type Config struct {
	Driver         string `config:"default=${sql.driver}"`
	DataSourceName string `config:"default=${sql.data-source-name}"`
	Enabled        bool   `config:"default=false"`
	Disconnected   bool   `config:"default=${cli.flag.disconnected:false}"`
}

func NewSqlConfigFromConfig(cfg *config.Config) (*Config, error) {
	var sqlConfig Config
	if err := cfg.Populate(&sqlConfig, configRootSpringDatasourceConfig); err != nil {
		return nil, err
	}
	return &sqlConfig, nil
}

func NewSqlConfig(ctx context.Context) (*Config, error) {
	return NewSqlConfigFromConfig(config.FromContext(ctx))
}
