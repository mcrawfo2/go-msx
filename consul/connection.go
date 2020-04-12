package consul

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"strings"
)

const (
	configRootConsulConnection = "spring.cloud.consul"
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
	stats  *statsObserver
}

func (c *Connection) Client() *api.Client {
	return c.client
}

func (c *Connection) Host() string {
	return c.config.Host
}

func (c *Connection) ListKeyValuePairs(ctx context.Context, path string) (results map[string]string, err error) {
	ctx, span := trace.NewSpan(ctx, "consul."+statsApiListKeyValuePairs)
	defer span.Finish()

	err = c.stats.Observe(statsApiListKeyValuePairs, path, func() error {
		queryOptions := &api.QueryOptions{}
		entries, _, err := c.client.KV().List(path, queryOptions.WithContext(ctx))
		if err != nil {
			return err
		} else if entries == nil {
			logger.Warningf("No config retrieved from consul (%s): %s", c.Host(), path)
		} else {
			logger.Infof("Retrieved %d configs from consul (%s): %s", len(entries), c.Host(), path)
		}

		prefix := path + "/"
		results = make(map[string]string)
		for _, entry := range entries {
			if !strings.HasPrefix(entry.Key, prefix) {
				continue
			}

			propName := strings.TrimPrefix(entry.Key, prefix)
			results[propName] = string(entry.Value)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return results, nil
}

func (c *Connection) GetKeyValue(ctx context.Context, path string) (value []byte, err error) {
	ctx, span := trace.NewSpan(ctx, "consul."+statsApiGetKeyValue)
	defer span.Finish()

	err = c.stats.Observe(statsApiGetKeyValue, path, func() error {
		queryOptions := &api.QueryOptions{}
		data, _, err := c.client.KV().Get(path, queryOptions.WithContext(ctx))
		if err != nil {
			return err
		} else if data == nil {
			logger.Warningf("No kv pair retrieved from consul %q: %s", c.Host(), path)
			value = nil
		} else {
			logger.Infof("Retrieved kv pair from consul %q: %s", c.Host(), path)
			value = data.Value
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return
}

func (c *Connection) SetKeyValue(ctx context.Context, path string, value []byte) (err error) {
	ctx, span := trace.NewSpan(ctx, "consul."+statsApiSetKeyValue)
	defer span.Finish()

	return c.stats.Observe(statsApiSetKeyValue, path, func() error {
		kvPair := &api.KVPair{
			Key:   path,
			Value: value,
		}

		writeOptions := &api.WriteOptions{}
		_, err := c.client.KV().Put(kvPair, writeOptions.WithContext(ctx))
		if err != nil {
			return err
		}

		logger.Infof("Stored kv pair to consul %q: %s", c.Host(), path)
		return nil
	})
}

func (c *Connection) GetServiceInstances(ctx context.Context, service string, passingOnly bool, tags ...string) (serviceEntries []*api.ServiceEntry, err error) {
	ctx, span := trace.NewSpan(ctx, "consul."+statsApiGetServiceInstances)
	defer span.Finish()

	err = c.stats.Observe(statsApiGetServiceInstances, service, func() error {
		queryOptions := &api.QueryOptions{}
		queryOptions = queryOptions.WithContext(ctx)
		if serviceEntries, _, err = c.client.Health().ServiceMultipleTags(service, tags, passingOnly, queryOptions); err != nil {
			return err
		} else if len(serviceEntries) == 0 {
			err = errors.Wrap(ErrNoInstances, service)
			return err
		} else {
			// Fix results that have no address to deal with kube2consul entries
			for _, v := range serviceEntries {
				if v.Service.Address == "" {
					v.Service.Address = v.Node.Address
				}
			}
			return nil
		}
	})

	if err != nil {
		return nil, err
	}
	return serviceEntries, nil
}

func (c *Connection) GetAllServiceInstances(ctx context.Context, passingOnly bool, tags ...string) (serviceEntries []*api.ServiceEntry, err error) {
	ctx, span := trace.NewSpan(ctx, "consul."+statsApiGetAllServiceInstances)
	defer span.Finish()

	err = c.stats.Observe(statsApiGetAllServiceInstances, "", func() error {
		queryOptions := &api.QueryOptions{}
		queryOptions = queryOptions.WithContext(ctx)

		var serviceMap map[string][]string
		if serviceMap, _, err = c.client.Catalog().Services(queryOptions); err != nil {
			return err
		} else if len(serviceMap) == 0 {
			return ErrNoInstances
		} else {
			for serviceName := range serviceMap {
				var serviceSpecificEntries []*api.ServiceEntry
				if serviceSpecificEntries, _, err = c.client.Health().ServiceMultipleTags(serviceName, tags, passingOnly, queryOptions); err != nil {
					return err
				} else if len(serviceSpecificEntries) > 0 {
					serviceEntries = append(serviceEntries, serviceSpecificEntries...)
				}
			}

			if len(serviceEntries) == 0 {
				return ErrNoInstances
			}

			// Fix results that have no address to deal with kube2consul entries
			for _, v := range serviceEntries {
				if v.Service.Address == "" {
					v.Service.Address = v.Node.Address
				}
			}
			return nil
		}
	})

	if err != nil {
		return nil, err
	}
	return serviceEntries, nil
}

func (c *Connection) RegisterService(ctx context.Context, registration *api.AgentServiceRegistration) error {
	return trace.Operation(ctx, "consul."+statsApiRegisterService, func(ctx context.Context) error {
		return c.stats.Observe(statsApiRegisterService, "", func() error {
			return c.client.Agent().ServiceRegister(registration)
		})
	})
}

func (c *Connection) DeregisterService(ctx context.Context, registration *api.AgentServiceRegistration) error {
	return trace.Operation(ctx, "consul."+statsApiDeregisterService, func(ctx context.Context) error {
		return c.stats.Observe(statsApiDeregisterService, "", func() error {
			return c.client.Agent().ServiceDeregister(registration.ID)
		})
	})
}

func (c *Connection) NodeHealth(ctx context.Context) (healthChecks api.HealthChecks, err error) {
	err = trace.Operation(ctx, "consul."+statsApiNodeHealth, func(ctx context.Context) error {
		err = c.stats.Observe(statsApiNodeHealth, "", func() error {
			var nodeName string
			nodeName, err = c.client.Agent().NodeName()
			if err != nil {
				return err
			}

			q := &api.QueryOptions{}
			q = q.WithContext(ctx)
			healthChecks, _, err = c.client.Health().Node(nodeName, q)

			return nil
		})

		if err != nil {
			healthChecks = nil
		}

		return err
	})

	return
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
			stats:  &statsObserver{},
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
