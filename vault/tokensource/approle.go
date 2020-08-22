package tokensource

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/hashicorp/vault/api"
)

const ( configRootApprole = "spring.cloud.vault.tokensource.approle" )

type ApproleSource struct {
}

type ApproleConfig struct {
	Path string `config:"default=/auth/approle/login"`
	Role_Id string
	Secret_Id string
}

func (c *ApproleSource) GetToken(client *api.Client, cfg  *config.Config) (token string, err error) {
	approleconfig := &ApproleConfig{}
	if err := cfg.Populate(approleconfig, configRootApprole); err != nil {
		return "", err
	}
	data := make(map[string]interface{})
	data["role_id"] = approleconfig.Role_Id
	data["secret_id"] = approleconfig.Secret_Id
	login, err := client.Logical().Write(approleconfig.Path, data)
	if err != nil {
		return "", err
	}
	return login.Auth.ClientToken, nil
}

func (a *ApproleSource) StartRenewer(client *api.Client) {
	r, err := initRenewer(client)
	if err != nil {
		logger.Error("Error initializing token renewer: ", err)
	}
	logger.Info("Starting token renewal.")
	startRenewer(r)
}
