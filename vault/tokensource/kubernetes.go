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

func (k *KubernetesSource) GetToken(client *api.Client, cfg *config.Config) (token string, err error) {
	k8sconfig := &KubernetesConfig{}
	if err := cfg.Populate(k8sconfig, configRootKubernetes); err != nil {
		return "", err
	}
	jwt, err := ioutil.ReadFile(k8sconfig.JWTPath)
	if err != nil {
		return "", err
	}
	data := make(map[string]interface{})
	data["jwt"] = string(jwt)
	data["role"] = k8sconfig.Role
	login, err := client.Logical().Write(k8sconfig.Path, data)
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
