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

func NewSqlConfig(ctx context.Context) (*Config, error) {
	var cfg Config
	if err := config.FromContext(ctx).Populate(&cfg, configRootSpringDatasourceConfig); err != nil {
		return nil, err
	}
	return &cfg, nil
}
