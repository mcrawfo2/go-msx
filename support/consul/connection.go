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
	ErrNoInstances    = errors.New("No matching service instances found")
	logger            = log.NewLogger("msx.support.consul")
)

type ConnectionConfig struct {
	Enabled bool   `config:"default=true"`
	Host    string `config:"default=localhost"`
	Port    int    `config:"default=8500"`
	Scheme  string `config:"default=http"`
	Config  struct {
		AclToken string `config:"default="`
	}
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

func (c *Connection) GetServiceInstances(service string, passingOnly bool, tags ...string) ([]*api.ServiceEntry, error) {
	if serviceEntries, _, err := c.client.Health().ServiceMultipleTags(service, tags, passingOnly, nil); err != nil {
		return nil, err
	} else if len(serviceEntries) == 0 {
		return nil, errors.Wrap(ErrNoInstances, service)
	} else {
		// Add a quick walk to fix results that have no address to deal with kube2consul entries
		for _, v := range serviceEntries {
			if v.Service.Address == "" {
				v.Service.Address = v.Node.Address
			}
		}
		return serviceEntries, nil
	}
}

func (c *Connection) RegisterService(ctx context.Context, registration *api.AgentServiceRegistration) error {
	return c.client.Agent().ServiceRegister(registration)
}

func (c *Connection) DeregisterService(ctx context.Context, registration *api.AgentServiceRegistration) error {
	return c.client.Agent().ServiceDeregister(registration.ID)
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
