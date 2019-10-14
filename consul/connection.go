package consul

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/stats"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"strings"
	"time"
)

const (
	configRootConsulConnection = "spring.cloud.consul"

	statsCounterPrefixConsulCallRequests = "consul.calls"
	statsCounterListKeyValuePairs = "list-kv-pairs"

	statsGetServiceInstances = "get-service-instances"
	statsRegisterService     = "register-service"
	statsDeregisterService   = "deregister-service"

	statsCounterConsulRegisteredServices = "consul.registration"

	statsTimerPrefixConsulCallTimer = "consul.timer"
)

var (
	ErrDisabled    = errors.New("Consul connection disabled")
	ErrNoInstances = errors.New("No matching service instances found")
	logger         = log.NewLogger("msx.consul")
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

func (c *Connection) Client() *api.Client {
	return c.client
}

func (c *Connection) Host() string {
	return c.config.Host
}

func (c *Connection) ListKeyValuePairs(ctx context.Context, path string) (map[string]string, error) {
	stats.Incr(
		strings.Join([]string{statsCounterPrefixConsulCallRequests, statsCounterListKeyValuePairs, path}, "."),
		1)

	start := time.Now()
	defer func() {
		stats.PrecisionTiming(
			strings.Join([]string{statsTimerPrefixConsulCallTimer, statsCounterListKeyValuePairs, path}, "."),
			time.Since(start))
	}()

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

func (c *Connection) GetServiceInstances(ctx context.Context, service string, passingOnly bool, tags ...string) ([]*api.ServiceEntry, error) {
	stats.Incr(strings.Join([]string{statsCounterPrefixConsulCallRequests, statsGetServiceInstances, service}, "."), 1)

	start := time.Now()
	defer func() {
		stats.PrecisionTiming(
			strings.Join([]string{statsTimerPrefixConsulCallTimer, statsGetServiceInstances, service}, "."),
			time.Since(start))
	}()

	queryOptions := &api.QueryOptions{}
	queryOptions = queryOptions.WithContext(ctx)
	if serviceEntries, _, err := c.client.Health().ServiceMultipleTags(service, tags, passingOnly, queryOptions); err != nil {
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
	stats.Incr(strings.Join([]string{statsCounterPrefixConsulCallRequests, statsRegisterService}, "."), 1)
	stats.Incr(statsCounterConsulRegisteredServices, 1)
	return c.client.Agent().ServiceRegister(registration)
}

func (c *Connection) DeregisterService(ctx context.Context, registration *api.AgentServiceRegistration) error {
	stats.Incr(strings.Join([]string{statsCounterPrefixConsulCallRequests, statsDeregisterService}, "."), 1)
	stats.Decr(statsCounterConsulRegisteredServices, 1)
	return c.client.Agent().ServiceDeregister(registration.ID)
}

func NewConnection(connectionConfig *ConnectionConfig) (*Connection, error) {
	if !connectionConfig.Enabled {
		return nil, ErrDisabled
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
