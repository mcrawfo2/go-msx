package tokensource

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/hashicorp/vault/api"
	"io/ioutil"
)

const (
	configRootKubernetes = "spring.cloud.vault.tokensource.kubernetes"
)

type KubernetesSource struct {
}

type KubernetesConfig struct {
	Path    string `config:"default=/auth/kubernetes/login"`
	JWTPath string `config:"default=/run/secrets/kubernetes.io/serviceaccount/token"`
	Role    string
}

func NewKubernetesConfig(cfg *config.Config) (*KubernetesConfig, error) {
	kubernetesConfig := KubernetesConfig{}
	if err := cfg.Populate(&kubernetesConfig, configRootKubernetes); err != nil {
		return nil, err
	}
	return &kubernetesConfig, nil
}

func (k *KubernetesSource) GetToken(client *api.Client, cfg *config.Config) (token string, err error) {
	kubernetesConfig, err := NewKubernetesConfig(cfg)
	if err != nil {
		return "", err
	}

	jwt, err := ioutil.ReadFile(kubernetesConfig.JWTPath)
	if err != nil {
		return "", err
	}
	data := make(map[string]interface{})
	data["jwt"] = string(jwt)
	data["role"] = kubernetesConfig.Role
	login, err := client.Logical().Write(kubernetesConfig.Path, data)
	if err != nil {
		return "", err
	}
	return login.Auth.ClientToken, nil
}

func (k *KubernetesSource) StartRenewer(client *api.Client) {
	r, err := initRenewer(client)
	if err != nil {
		logger.Error("Error initializing token renewer: ", err)
	}
	logger.Info("Starting token renewal.")
	startRenewer(r)
}
