package vault

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
)

const (
	configKeyVaultToken = "spring.cloud.vault.token"
)

type ConfigSource struct {
	cfg *config.Config
}

func (c ConfigSource) GetToken(_ context.Context) (token string, err error) {
	return c.cfg.StringOr(configKeyVaultToken, "replace_with_token_value")
}

func (c ConfigSource) Renewable() bool {
	return false
}

func NewConfigSource(cfg *config.Config) ConfigSource {
	return ConfigSource{cfg:cfg}
}
