package config

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
)

const (
	configRootVaultToken = "spring.cloud.vault.token"
)

type ConfigSource struct {
}

func (c *ConfigSource) GetToken(cfg *config.Config) (token string, err error) {
	return cfg.String(configRootVaultToken)
}


