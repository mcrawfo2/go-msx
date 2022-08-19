// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package discovery

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"time"
)

type RegistrationConfig struct {
	Enabled             bool          `config:"default=false"`
	Name                string        `config:"default=${info.app.name}"`
	Register            bool          `config:"default=true"`
	IpAddress           string        `config:"default="`
	Interface           string        `config:"default="`
	Port                int           `config:"default=0"`
	Scheme              string        `config:"default=http"`
	RegisterHealthCheck bool          `config:"default=true"`
	HealthCheckPath     string        `config:"default=/admin/health"`
	HealthCheckInterval time.Duration `config:"default=10s"`
	HealthCheckTimeout  time.Duration `config:"default=10s"`
	Tags                string        `config:"default="`
	InstanceId          string        `config:"default=local"` // uuid, hostname, or any static string
	InstanceName        string        `config:"default=${info.app.name}"`
	HiddenApiListing    bool          `config:"default=false"`
}

func NewRegistrationConfigFromConfig(cfg *config.Config, root string) (*RegistrationConfig, error) {
	var providerConfig = RegistrationConfig{}
	var err = cfg.Populate(&providerConfig, root)
	if err != nil {
		return nil, err
	}

	return &providerConfig, nil
}

func NewRegistrationConfig(ctx context.Context, root string) (*RegistrationConfig, error) {
	return NewRegistrationConfigFromConfig(
		config.FromContext(ctx),
		root)
}
