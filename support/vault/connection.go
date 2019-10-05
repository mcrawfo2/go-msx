package vault

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/support/config"
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"io"
)

var (
	logger           = log.NewLogger("msx.support.vault")
	ErrVaultDisabled = errors.New("Consul connection disabled")
)

type ConnectionConfig struct {
	Enabled bool   `properties:"enabled,default=true"`
	Host    string `properties:"host,default=localhost"`
	Port    int    `properties:"port,default=8200"`
	Scheme  string `properties:"scheme,default=http"`
	Token   string `properties:"token,default=replace_with_token_value"`
	Ssl     struct {
		Cacert     string `properties:"cacert,default="`
		ClientCert string `properties:"clientCert,default="`
		ClientKey  string `properties:"clientKey,default="`
		Insecure   bool   `properties:"insecure,default=true"`
	} `properties:"ssl" json:"ssl"`
}

func (c ConnectionConfig) Address() string {
	return fmt.Sprintf("%s://%s:%d", c.Scheme, c.Host, c.Port)
}

type Connection struct {
	config *ConnectionConfig
	client *api.Client
}

func (c *Connection) Host() string {
	return c.config.Host
}

func (c *Connection) ListSecrets(ctx context.Context, path string) (map[string]string, error) {
	var results = make(map[string]string)

	if secrets, err := c.read(ctx, path); err != nil {
		return nil, errors.Wrap(err, "Failed to list vault secrets")
	} else if secrets != nil {
		logger.Info("Retrieved %d secrets from vault", len(secrets.Data))
		for key, val := range secrets.Data {
			results[key] = val.(string)
		}
	} else {
		logger.Warning("No secrets retrieved from vault: ", path)
	}

	return results, nil
}

// Copied from vault/logical to allow custom context
func (c *Connection) read(ctx context.Context, path string) (*api.Secret, error) {
	r := c.client.NewRequest("GET", "/v1/"+path)

	resp, err := c.client.RawRequestWithContext(ctx, r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		secret, parseErr := api.ParseSecret(resp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
			return nil, nil
		default:
			return nil, err
		}
		if secret != nil && (len(secret.Warnings) > 0 || len(secret.Data) > 0) {
			return secret, nil
		}
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return api.ParseSecret(resp.Body)
}

func NewConnection(cfg *config.Config) (*Connection, error) {
	connectionConfig := &ConnectionConfig{}
	if err := cfg.Populate(connectionConfig, "spring.cloud.vault"); err != nil {
		return nil, err
	}

	if !connectionConfig.Enabled {
		return nil, ErrVaultDisabled
	}

	clientConfig := api.DefaultConfig()
	clientConfig.Address = connectionConfig.Address()
	if connectionConfig.Scheme == "https" {
		t := api.TLSConfig{
			CACert:     connectionConfig.Ssl.Cacert,
			ClientCert: connectionConfig.Ssl.ClientCert,
			ClientKey:  connectionConfig.Ssl.ClientKey,
			Insecure:   connectionConfig.Ssl.Insecure,
		}
		err := clientConfig.ConfigureTLS(&t)
		if err != nil {
			return nil, err
		}
	}

	client, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	client.SetToken(connectionConfig.Token)

	return &Connection{
		config: connectionConfig,
		client: client,
	}, nil
}
