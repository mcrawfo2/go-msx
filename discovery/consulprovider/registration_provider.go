// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package consulprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/consul"
	"cto-github.cisco.com/NFV-BU/go-msx/discovery"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"path"
	"strings"
)

const (
	configRootRegistrationProvider = "spring.cloud.consul.discovery"
)

var (
	logger      = log.NewLogger("msx.discovery.consulprovider")
	ErrDisabled = errors.New("Consul registration provider disabled")
)

type RegistrationProvider struct {
	config  *discovery.RegistrationConfig
	conn    *consul.Connection
	details *discovery.RegistrationDetails
}

func (c *RegistrationProvider) tags() []string {
	var tags []string

	if c.config.Tags != "" {
		tags = append(tags, strings.Split(c.config.Tags, ",")...)
	}

	if c.config.HiddenApiListing {
		tags = append(tags, "hiddenApiListing=true")
	}
	tags = append(tags, c.details.Tags()...)

	return tags
}

func (c *RegistrationProvider) meta() map[string]string {
	return c.details.Meta()
}

func (c *RegistrationProvider) healthCheck() *api.AgentServiceCheck {
	if c.config.RegisterHealthCheck == false {
		return nil
	}

	checkPath := path.Clean(path.Join(c.details.ContextPath, c.config.HealthCheckPath))

	return &api.AgentServiceCheck{
		Interval:      c.config.HealthCheckInterval.String(),
		Timeout:       c.config.HealthCheckTimeout.String(),
		TLSSkipVerify: true,
		HTTP:          fmt.Sprintf("%s://%s%s", c.config.Scheme, c.details.SocketAddress(), checkPath),
	}
}

func (c *RegistrationProvider) serviceRegistration() *api.AgentServiceRegistration {
	registration := &api.AgentServiceRegistration{
		ID:      c.details.InstanceId,
		Name:    c.details.Name,
		Address: c.details.ServiceAddress,
		Port:    c.config.Port,
		Tags:    c.tags(),
		Meta:    c.meta(),
	}

	if c.config.RegisterHealthCheck {
		registration.Check = c.healthCheck()
	}

	return registration
}

func (c *RegistrationProvider) Register(ctx context.Context) error {
	if c.details != nil {
		logger.Infof("Registering service in consul: %v", *c.details)
		return c.conn.RegisterService(ctx, c.serviceRegistration())
	} else {
		return nil
	}
}

func (c *RegistrationProvider) Deregister(ctx context.Context) error {
	if c.details != nil {
		logger.Infof("De-registering service in consul: %v", *c.details)
		return c.conn.DeregisterService(ctx, c.serviceRegistration())
	} else {
		return nil
	}
}

func NewRegistrationProvider(ctx context.Context) (*RegistrationProvider, error) {
	providerConfig, err := discovery.NewRegistrationConfig(ctx, configRootRegistrationProvider)
	if err != nil {
		return nil, err
	}

	if !providerConfig.Enabled || !providerConfig.Register {
		logger.Warn(ErrDisabled)
		return nil, nil
	}

	details, err := discovery.RegistrationFactory{}.NewRegistrationDetails(ctx, providerConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create registration details")
	}

	var conn *consul.Connection
	if conn, err = consul.NewConnectionFromConfig(config.FromContext(ctx)); err != nil {
		return nil, err
	} else if conn == nil {
		return nil, nil
	}

	return &RegistrationProvider{
		config:  providerConfig,
		details: details,
		conn:    conn,
	}, nil
}
