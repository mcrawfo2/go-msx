// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package populate

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/manage"
	"cto-github.cisco.com/NFV-BU/go-msx/populate"
	"cto-github.cisco.com/NFV-BU/go-msx/resource"
	"cto-github.cisco.com/NFV-BU/go-msx/security/service"
	"net/http"
	"path"
)

const (
	deviceActionPopulatorConfigRoot = "populate.manage.device-action"
)

type deviceActionManifest struct {
	ServiceType  string                         `json:"serviceType"`
	DeviceAction []deviceActionManifestArtifact `json:"deviceAction"`
}

type deviceActionManifestArtifact struct {
	FileName string `json:"fileName"`
}

type deviceActionInstance struct {
	manage.DeviceActionCreateRequest
}

type DeviceActionPopulatorConfig struct {
	Enabled bool   `config:"default=false"`
	Root    string `config:"default=${populate.root}/manage"`
}

func NewDeviceActionPopulatorConfigFromConfig(cfg *config.Config) (*DeviceActionPopulatorConfig, error) {
	var populatorConfig DeviceActionPopulatorConfig
	if err := cfg.Populate(&populatorConfig, deviceActionPopulatorConfigRoot); err != nil {
		return nil, err
	}

	return &populatorConfig, nil
}

type DeviceActionPopulator struct {
	cfg DeviceActionPopulatorConfig
}

func (p DeviceActionPopulator) populateDeviceAction(ctx context.Context, api manage.Api, serviceType string, deviceAction deviceActionInstance) error {
	request := deviceAction.DeviceActionCreateRequest

	logger.WithContext(ctx).Infof("Populating action %q", request.Name)

	if request.Owner == "" {
		request.Owner = serviceType
	}

	response, err := api.CreateDeviceActions([]manage.DeviceActionCreateRequest{request})
	if err != nil && response != nil {
		logger.WithContext(ctx).Error(response.BodyString)
	}

	if response != nil && (response.StatusCode == http.StatusConflict || response.StatusCode == http.StatusInternalServerError) {
		logger.WithContext(ctx).Info("Device Action already exists.  Updating...")
		response, err = api.UpdateDeviceActions([]manage.DeviceActionCreateRequest{request})
	}

	return err
}

func (p DeviceActionPopulator) Populate(ctx context.Context) error {
	if !p.cfg.Enabled {
		logger.WithContext(ctx).Warn("Device Action populator disabled.")
		return nil
	}

	return service.WithDefaultServiceAccount(ctx, func(ctx context.Context) error {
		var m deviceActionManifest
		err := resource.
			Reference(path.Join(p.cfg.Root, manifestFile)).
			Unmarshal(&m)
		if err != nil {
			return err
		}

		api, _ := manage.NewIntegration(ctx)

		logger.WithContext(ctx).Info("Populating device actions")

		for _, artifact := range m.DeviceAction {
			var deviceAction deviceActionInstance
			err := resource.
				Reference(path.Join(p.cfg.Root, artifact.FileName)).
				Unmarshal(&deviceAction)
			if err != nil {
				return err
			}

			err = p.populateDeviceAction(ctx, api, m.ServiceType, deviceAction)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func init() {
	populate.RegisterPopulationTask(
		populate.NewPopulatorTask(
			"Populate device actions",
			1000,
			[]string{"all", "deviceActions", "serviceMetadata"},
			func(ctx context.Context) (populate.Populator, error) {
				cfg, err := NewDeviceActionPopulatorConfigFromConfig(config.MustFromContext(ctx))
				if err != nil {
					return nil, err
				}
				return &DeviceActionPopulator{
					cfg: *cfg,
				}, nil
			}))
}
