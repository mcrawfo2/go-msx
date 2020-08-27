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

func (c *ConfigSource) GetToken(client *api.Client, cfg *config.Config) (token string, err error) {
	return cfg.StringOr(configRootVaultToken, "replace_with_token_value")
}

func (c *ConfigSource) StartRenewer(client *api.Client) {
	r, err := initRenewer(client)
	if err != nil {
		logger.Error("Error initializing token renewer: ", err)
	}
	logger.Info("Starting token renewal.")
	startRenewer(r)
}
