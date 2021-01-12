package tokensource

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/hashicorp/vault/api"
)

const (
	configRootVaultToken = "spring.cloud.vault.token"
)

type ConfigSource struct {
}

func (c *ConfigSource) GetToken(_ *api.Client, cfg *config.Config) (token string, err error) {
	return cfg.StringOr(configRootVaultToken, "replace_with_token_value")
}

func (c *ConfigSource) StartRenewer(_ *api.Client) {
	logger.Warn("Config token renewal disabled")
}
