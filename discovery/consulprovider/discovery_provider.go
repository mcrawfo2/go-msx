package consulprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/consul"
	"cto-github.cisco.com/NFV-BU/go-msx/discovery"
	"github.com/hashicorp/consul/api"
)

const configRootDiscoveryProvider = "spring.cloud.consul.discovery"

type DiscoveryProviderConfig struct {
	DefaultQueryTag string `config:"default="`
}

type DiscoveryProvider struct {
	cfg  *DiscoveryProviderConfig
	conn *consul.Connection
}

func (p *DiscoveryProvider) Discover(ctx context.Context, service string, passingOnly bool, tags ...string) (result discovery.ServiceInstances, err error) {
	var serviceEntries []*api.ServiceEntry
	if p.cfg.DefaultQueryTag != "" {
		tags = append(tags, p.cfg.DefaultQueryTag)
	}
	if serviceEntries, err = p.conn.GetServiceInstances(ctx, service, passingOnly, tags...); err != nil {
		return nil, err
	}

	return convertToServiceInstances(serviceEntries), nil
}

func (p *DiscoveryProvider) DiscoverAll(ctx context.Context, passingOnly bool, tags ...string) (result discovery.ServiceInstances, err error) {
	var serviceEntries []*api.ServiceEntry
	if p.cfg.DefaultQueryTag != "" {
		tags = append(tags, p.cfg.DefaultQueryTag)
	}
	if serviceEntries, err = p.conn.GetAllServiceInstances(ctx, passingOnly, tags...); err != nil {
		return nil, err
	}

	return convertToServiceInstances(serviceEntries), nil
}

func convertToServiceInstances(sourceEntries []*api.ServiceEntry) (result discovery.ServiceInstances) {
	for _, sourceEntry := range sourceEntries {
		result = append(result, convertToServiceInstance(sourceEntry))
	}
	return result
}

func convertToServiceInstance(sourceEntry *api.ServiceEntry) *discovery.ServiceInstance {
	return &discovery.ServiceInstance{
		ID:   sourceEntry.Service.ID,
		Name: sourceEntry.Service.Service,
		Host: sourceEntry.Service.Address,
		Tags: sourceEntry.Service.Tags,
		Meta: sourceEntry.Service.Meta,
		Port: sourceEntry.Service.Port,
	}
}

func NewDiscoveryProviderFromConfig(cfg *config.Config) (provider *DiscoveryProvider, err error) {
	var discoveryConfig DiscoveryProviderConfig
	if err := cfg.Populate(&discoveryConfig, configRootDiscoveryProvider); err != nil {
		return nil, err
	}

	var conn *consul.Connection
	if conn, err = consul.NewConnectionFromConfig(cfg); err != nil && err != consul.ErrDisabled {
		return nil, err
	} else if err == consul.ErrDisabled {
		return nil, nil
	} else if conn == nil {
		return nil, nil
	}

	return &DiscoveryProvider{
		cfg:  &discoveryConfig,
		conn: conn,
	}, nil
}
