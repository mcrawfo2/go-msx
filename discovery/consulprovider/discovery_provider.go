package consulprovider

import (
	"cto-github.cisco.com/NFV-BU/go-msx/discovery"
	"cto-github.cisco.com/NFV-BU/go-msx/support/config"
	"cto-github.cisco.com/NFV-BU/go-msx/support/consul"
	"github.com/hashicorp/consul/api"
)

type DiscoveryProvider struct {
	conn *consul.Connection
}

func (p *DiscoveryProvider) Discover(service string, passingOnly bool, tags ...string) (result discovery.ServiceInstances, err error) {
	var serviceEntries []*api.ServiceEntry
	if serviceEntries, err = p.conn.GetServiceInstances(service, passingOnly, tags...); err != nil {
		return nil, err
	}

	for _, sourceEntry := range serviceEntries {
		resultEntry := &discovery.ServiceInstance{
			ID:   sourceEntry.Service.ID,
			Name: sourceEntry.Service.Service,
			Host: sourceEntry.Service.Address,
			Tags: sourceEntry.Service.Tags,
			Meta: sourceEntry.Service.Meta,
			Port: sourceEntry.Service.Port,
		}

		result = append(result, resultEntry)
	}

	return
}

func NewDiscoveryProviderFromConfig(cfg *config.Config) (provider *DiscoveryProvider, err error) {
	var conn *consul.Connection
	if conn, err = consul.NewConnectionFromConfig(cfg); err != nil {
		return nil, err
	} else if conn == nil {
		return nil, nil
	}

	return &DiscoveryProvider{
		conn: conn,
	}, nil
}
