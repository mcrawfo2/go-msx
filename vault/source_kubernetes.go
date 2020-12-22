package vault

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"io/ioutil"
)

const (
	configRootTokenSourceKubernetes = "spring.cloud.vault.token-source.kubernetes"
)

type KubernetesConfig struct {
	JWTPath string `config:"default=/run/secrets/kubernetes.io/serviceaccount/token"`
	Role    string
}

type KubernetesSource struct {
	cfg    *KubernetesConfig
	conn   ConnectionApi
}

func (k *KubernetesSource) GetToken(ctx context.Context) (token string, err error) {
	jwt, err := ioutil.ReadFile(k.cfg.JWTPath)
	if err != nil {
		return "", err
	}
	token, err = k.conn.LoginWithKubernetes(ctx, string(jwt), k.cfg.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (k *KubernetesSource) Renewable() bool {
	return true
}

func NewKubernetesConfig(cfg *config.Config) (*KubernetesConfig, error) {
	kubernetesConfig := KubernetesConfig{}
	if err := cfg.Populate(&kubernetesConfig, configRootTokenSourceKubernetes); err != nil {
		return nil, err
	}
	return &kubernetesConfig, nil
}

func NewKubernetesSource(cfg *config.Config, conn ConnectionApi) (*KubernetesSource, error) {
	kubernetesConfig, err := NewKubernetesConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &KubernetesSource{cfg: kubernetesConfig, conn: conn}, nil
}
