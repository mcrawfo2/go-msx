package tokensource

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/hashicorp/vault/api"
)

const (
	configRootApprole = "spring.cloud.vault.tokensource.approle"
)

type AppRoleConfig struct {
	Path     string `config:"default=/auth/approle/login"`
	RoleId   string
	SecretId string
}

func NewAppRoleConfig(cfg *config.Config) (*AppRoleConfig, error) {
	appRoleConfig := &AppRoleConfig{}
	if err := cfg.Populate(appRoleConfig, configRootApprole); err != nil {
		return nil, err
	}

	return appRoleConfig, nil
}

type AppRoleSource struct {
}

func (c *AppRoleSource) GetToken(client *api.Client, cfg *config.Config) (token string, err error) {
	appRoleConfig, err := NewAppRoleConfig(cfg)
	if err != nil {
		return "", err
	}

	data := make(map[string]interface{})
	data["role_id"] = appRoleConfig.RoleId
	data["secret_id"] = appRoleConfig.SecretId
	login, err := client.Logical().Write(appRoleConfig.Path, data)
	if err != nil {
		return "", err
	}
	return login.Auth.ClientToken, nil
}

func (c *AppRoleSource) StartRenewer(client *api.Client) {
	r, err := initRenewer(client)
	if err != nil {
		logger.Error("Error initializing token renewer: ", err)
	}
	logger.Info("Starting token renewal.")
	startRenewer(r)
}
