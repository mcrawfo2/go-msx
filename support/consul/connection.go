package consul

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/support/config"
	"cto-github.cisco.com/NFV-BU/go-msx/support/log"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"strings"
)

const (
	configRootConsulConnection = "spring.cloud.consul"
)

var (
	ErrConsulDisabled = errors.New("Consul connection disabled")
	logger            = log.NewLogger("msx.support.consul")
)

type ConnectionConfig struct {
	Enabled bool   `properties:"enabled,default=true"`
	Host    string `properties:"host,default=localhost"`
	Port    int    `properties:"port,default=8500"`
	Scheme  string `properties:"scheme,default=http"`
	Config  struct {
		AclToken string `properties:"acltoken,default="`
	} `properties:"config"`
}

func (c ConnectionConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type Connection struct {
	config *ConnectionConfig
	client *api.Client
}

func (c *Connection) Host() string {
	return c.config.Host
}

func (c *Connection) ListKeyValuePairs(ctx context.Context, path string) (map[string]string, error) {
	queryOptions := &api.QueryOptions{}
	entries, _, err := c.client.KV().List(path, queryOptions.WithContext(ctx))
	if err != nil {
		return nil, err
	} else if entries == nil {
		logger.Warningf("No config retrieved from consul (%s): %s", c.Host(), path)
	} else {
		logger.Infof("Retrieved %d configs from consul (%s): %s", len(entries), c.Host(), path)
	}

	prefix := path + "/"
	results := make(map[string]string)
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Key, prefix) {
			continue
		}

		propName := strings.TrimPrefix(entry.Key, prefix)
		results[propName] = string(entry.Value)
	}

	return results, nil
}

func NewConnection(connectionConfig *ConnectionConfig) (*Connection, error) {
	if !connectionConfig.Enabled {
		return nil, ErrConsulDisabled
	}

	clientConfig := api.DefaultConfig()
	clientConfig.Address = connectionConfig.Address()
	clientConfig.Scheme = connectionConfig.Scheme
	clientConfig.Token = connectionConfig.Config.AclToken
	clientConfig.TLSConfig.InsecureSkipVerify = true

	if client, err := api.NewClient(clientConfig); err != nil {
		return nil, err
	} else {
		return &Connection{
			config: connectionConfig,
			client: client,
		}, nil
	}
}

func NewConnectionFromConfig(cfg *config.Config) (*Connection, error) {
	connectionConfig := &ConnectionConfig{}
	if err := cfg.Populate(connectionConfig, configRootConsulConnection); err != nil {
		return nil, err
	}

	return NewConnection(connectionConfig)
}
