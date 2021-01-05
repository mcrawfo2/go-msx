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
