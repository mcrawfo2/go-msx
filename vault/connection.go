package vault

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"io"
)

const (
	configRootVaultConnection = "spring.cloud.vault"
)

var (
	ErrDisabled = errors.New("Vault connection disabled")
	logger      = log.NewLogger("msx.vault")
)

type ConnectionConfig struct {
	Enabled bool   `config:"default=true"`
	Host    string `config:"default=localhost"`
	Port    int    `config:"default=8200"`
	Scheme  string `config:"default=http"`
	Token   string `config:"default=replace_with_token_value"`
	Ssl     struct {
		Cacert     string `config:"default="`
		ClientCert string `config:"default="`
		ClientKey  string `config:"default="`
		Insecure   bool   `config:"default=true"`
	}
}

func (c ConnectionConfig) Address() string {
	return fmt.Sprintf("%s://%s:%d", c.Scheme, c.Host, c.Port)
}

type Connection struct {
	config *ConnectionConfig
	client *api.Client
	stats  *statsObserver
}

func (c *Connection) Host() string {
	return c.config.Host
}

func (c *Connection) Client() *api.Client {
	return c.client
}

func (c *Connection) ListSecrets(ctx context.Context, path string) (results map[string]string, err error) {
	err = trace.Operation(ctx, "vault."+statsApiListSecrets, func(ctx context.Context) error {
		return c.stats.Observe(statsApiListSecrets, path, func() error {
			results = make(map[string]string)

			if secrets, err := c.read(ctx, path); err != nil {
				return errors.Wrap(err, "Failed to list vault secrets")
			} else if secrets != nil {
				logger.Infof("Retrieved %d configs from vault (%s): %s", len(secrets.Data), c.Host(), path)
				for key, val := range secrets.Data {
					results[key] = val.(string)
				}
			} else {
				logger.Warningf("No secrets retrieved from vault (%s): %s", c.Host(), path)
			}

			return nil
		})
	})

	return
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

func (c *Connection) StoreSecrets(ctx context.Context, path string, secrets map[string]string) (err error) {
	err = trace.Operation(ctx, "vault."+statsApiStoreSecrets, func(ctx context.Context) error {
		return c.stats.Observe(statsApiStoreSecrets, path, func() error {
			if _, err := c.write(ctx, path, secrets); err != nil {
				return errors.Wrap(err, "Failed to store vault secrets")
			}

			return nil
		})
	})

	return
}

func (c *Connection) write(ctx context.Context, path string, requestBody interface{}) (*api.Secret, error) {
	r := c.client.NewRequest("POST", "/v1/"+path)
	if err := r.SetJSONBody(requestBody); err != nil {
		return nil, err
	}

	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()
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
			return secret, err
		}
	}
	if err != nil {
		return nil, err
	}

	return api.ParseSecret(resp.Body)
}

func (c *Connection) CreateTransitKey(ctx context.Context, keyName string, request CreateTransitKeyRequest) (err error) {
	err = trace.Operation(ctx, "vault."+statsApiCreateTransitKey, func(ctx context.Context) error {
		return c.stats.Observe(statsApiCreateTransitKey, keyName, func() error {
			path := "transit/keys/" + keyName
			if _, err = c.write(ctx, path, request); err != nil {
				return errors.Wrap(err, "Failed to create transit key")
			}
			return nil
		})
	})

	return
}

func (c *Connection) TransitEncrypt(ctx context.Context, keyName string, plaintext string) (ciphertext string, err error) {
	err = trace.Operation(ctx, "vault."+statsApiTransitEncrypt, func(ctx context.Context) error {
		return c.stats.Observe(statsApiTransitEncrypt, keyName, func() error {
			path := "/transit/encrypt/" + keyName

			data := map[string]interface{}{
				"plaintext": base64.StdEncoding.EncodeToString([]byte(plaintext)),
			}

			result, err := c.write(ctx, path, data)
			if err != nil {
				return errors.Wrap(err, "Failed to encrypt data")
			}

			ciphertext = result.Data["ciphertext"].(string)
			return nil
		})
	})

	return
}

func (c *Connection) TransitDecrypt(ctx context.Context, keyName string, ciphertext string) (plaintext string, err error) {
	err = trace.Operation(ctx, "vault."+statsApiTransitDecrypt, func(ctx context.Context) error {
		return c.stats.Observe(statsApiTransitDecrypt, keyName, func() error {
			path := "/transit/decrypt/" + keyName

			data := map[string]interface{}{
				"ciphertext": ciphertext,
			}

			result, err := c.write(ctx, path, data)
			if err != nil {
				return errors.Wrap(err, "Failed to decrypt data")
			}

			plaintextBytes, err := base64.StdEncoding.DecodeString(result.Data["plaintext"].(string))
			if err != nil {
				return err
			}

			plaintext = string(plaintextBytes)
			return nil
		})
	})

	return
}

// Copied from vault/api to allow custom context
func (c *Connection) Health(ctx context.Context) (response *api.HealthResponse, err error) {
	err = trace.Operation(ctx, "vault."+statsApiHealth, func(ctx context.Context) error {
		return c.stats.Observe(statsApiHealth, "", func() error {
			r := c.client.NewRequest("GET", "/v1/sys/health")
			// If the code is 400 or above it will automatically turn into an error,
			// but the sys/health API defaults to returning 5xx when not sealed or
			// inited, so we force this code to be something else so we parse correctly
			r.Params.Add("uninitcode", "299")
			r.Params.Add("sealedcode", "299")
			r.Params.Add("standbycode", "299")
			r.Params.Add("drsecondarycode", "299")
			r.Params.Add("performancestandbycode", "299")

			ctx, cancelFunc := context.WithCancel(ctx)
			defer cancelFunc()
			resp, err := c.client.RawRequestWithContext(ctx, r)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			var result api.HealthResponse
			err = resp.DecodeJSON(&result)
			response = &result
			return err
		})
	})
	return
}

func NewConnection(connectionConfig *ConnectionConfig) (*Connection, error) {
	if !connectionConfig.Enabled {
		return nil, ErrDisabled
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
		stats:  new(statsObserver),
	}, nil
}

func NewConnectionFromConfig(cfg *config.Config) (*Connection, error) {
	connectionConfig := &ConnectionConfig{}
	if err := cfg.Populate(connectionConfig, configRootVaultConnection); err != nil {
		return nil, err
	}

	return NewConnection(connectionConfig)
}
