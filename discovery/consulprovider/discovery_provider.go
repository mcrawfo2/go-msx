package consulprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/discovery"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/consul"
	"github.com/hashicorp/consul/api"
)

type DiscoveryProvider struct {
	conn *consul.Connection
}

func (p *DiscoveryProvider) Discover(ctx context.Context, service string, passingOnly bool, tags ...string) (result discovery.ServiceInstances, err error) {
	var serviceEntries []*api.ServiceEntry
	if serviceEntries, err = p.conn.GetServiceInstances(ctx, service, passingOnly, tags...); err != nil {
		return nil, err
	}

	return convertToServiceInstances(serviceEntries), nil
}

func (p *DiscoveryProvider) DiscoverAll(ctx context.Context, passingOnly bool, tags ...string) (result discovery.ServiceInstances, err error) {
	var serviceEntries []*api.ServiceEntry
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
	var conn *consul.Connection
	if conn, err = consul.NewConnectionFromConfig(cfg); err != nil && err != consul.ErrDisabled {
		return nil, err
	} else if err == consul.ErrDisabled {
		return nil, nil
	} else if conn == nil {
		return nil, nil
	}

	return &DiscoveryProvider{
		conn: conn,
	}, nil
}
