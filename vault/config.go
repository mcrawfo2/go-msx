// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package vault

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"fmt"
	"github.com/hashicorp/vault/api"
)

const (
	configRootVaultConnection = "spring.cloud.vault"
)

type ConnectionTokenSourceConfig struct {
	Source string `config:"default=config"`
}

type ConnectionSslConfig struct {
	Cacert     string `config:"default="`
	ClientCert string `config:"default="`
	ClientKey  string `config:"default="`
	Insecure   bool   `config:"default=true"`
}

type ConnectionIssuerConfig struct {
	Mount string `config:"default=/pki"` // Mount sets targeted mount point for Vault PKI
}

type ConnectionKvConfig struct {
	Mount string `config:"default=/secret"`
}

type ConnectionKv2Config struct {
	Mount string `config:"default=/v2secret"`
}

type ConnectionKubernetesConfig struct {
	LoginPath string `config:"default=/auth/kubernetes/login"`
}

type ConnectionAppRoleConfig struct {
	LoginPath string `config:"default=/auth/approle/login"`
}

type ConnectionConfig struct {
	Enabled      bool   `config:"default=true"`
	Disconnected bool   `config:"default=${cli.flag.disconnected:false}"`
	Host         string `config:"default=localhost"`
	Port         int    `config:"default=8200"`
	Scheme       string `config:"default=http"`
	TokenSource  ConnectionTokenSourceConfig
	Ssl          ConnectionSslConfig
	Issuer       ConnectionIssuerConfig
	Kubernetes   ConnectionKubernetesConfig
	AppRole      ConnectionAppRoleConfig
	KV           ConnectionKvConfig
	KV2          ConnectionKv2Config
}

func (c ConnectionConfig) Address() string {
	return fmt.Sprintf("%s://%s:%d", c.Scheme, c.Host, c.Port)
}

func (c ConnectionConfig) ClientConfig() (*api.Config, error) {
	clientConfig := api.DefaultConfig()
	clientConfig.Address = c.Address()
	if c.Scheme == "https" {
		t := api.TLSConfig{
			CACert:     c.Ssl.Cacert,
			ClientCert: c.Ssl.ClientCert,
			ClientKey:  c.Ssl.ClientKey,
			Insecure:   c.Ssl.Insecure,
		}
		err := clientConfig.ConfigureTLS(&t)
		if err != nil {
			return nil, err
		}
	}
	return clientConfig, nil
}

func newConnectionConfig(cfg *config.Config) (*ConnectionConfig, error) {
	connectionConfig := &ConnectionConfig{}
	if err := cfg.Populate(connectionConfig, configRootVaultConnection); err != nil {
		return nil, err
	}
	return connectionConfig, nil
}
