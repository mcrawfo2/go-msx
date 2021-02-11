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
	"path/filepath"
)

const (
	deviceTemplatePopulatorConfigRoot = "populate.manage.device-template"
	manifestFile                      = "manifest.json"
)

type deviceTemplateManifest struct {
	ServiceType    string                           `json:"serviceType"`
	DeviceTemplate []deviceTemplateManifestArtifact `json:"deviceTemplate"`
}

type deviceTemplateManifestArtifact struct {
	FileName string `json:"fileName"`
}

type deviceTemplateInstance struct {
	manage.DeviceTemplateCreateRequest
	ConfigFileName string `json:"configFileName"`
}

type DeviceTemplatePopulatorConfig struct {
	Enabled bool   `config:"default=false"`
	Root    string `config:"default=${populate.root}/manage"`
}

func NewDeviceTemplatePopulatorConfigFromConfig(cfg *config.Config) (*DeviceTemplatePopulatorConfig, error) {
	var populatorConfig DeviceTemplatePopulatorConfig
	if err := cfg.Populate(&populatorConfig, deviceTemplatePopulatorConfigRoot); err != nil {
		return nil, err
	}

	return &populatorConfig, nil
}

type DeviceTemplatePopulator struct {
	cfg DeviceTemplatePopulatorConfig
}

func (p DeviceTemplatePopulator) populateDeviceTemplate(ctx context.Context, api manage.Api, serviceType string, deviceTemplate deviceTemplateInstance) error {
	request := deviceTemplate.DeviceTemplateCreateRequest

	logger.WithContext(ctx).Infof("Populating template %q", request.Name)

	if request.ServiceType == "" {
		request.ServiceType = serviceType
	}

	if request.ConfigContent == "" {
		if deviceTemplate.ConfigFileName != "" {
			configBytes, err := resource.
				Reference(path.Join(p.cfg.Root, deviceTemplate.ConfigFileName)).
				ReadAll()
			if err != nil {
				return err
			}
			request.ConfigContent = string(configBytes)
		}
	}

	response, err := api.AddDeviceTemplate(request)
	if response != nil && response.StatusCode == http.StatusConflict {
		logger.WithContext(ctx).Error("Device template already exists")
		return nil
	}

	if err != nil && response != nil {
		logger.WithContext(ctx).Error(response.BodyString)
	}

	return err
}

func (p DeviceTemplatePopulator) Populate(ctx context.Context) error {
	if !p.cfg.Enabled {
		logger.WithContext(ctx).Warn("Device Template populator disabled.")
		return nil
	}

	return service.WithDefaultServiceAccount(ctx, func(ctx context.Context) error {
		var m deviceTemplateManifest
		err := resource.
			Reference(path.Join(p.cfg.Root, manifestFile)).
			Unmarshal(&m)
		if err != nil {
			return err
		}

		api, _ := manage.NewIntegration(ctx)

		logger.WithContext(ctx).Info("Populating device templates")

		for _, artifact := range m.DeviceTemplate {
			var deviceTemplate deviceTemplateInstance
			err := resource.
				Reference(filepath.Join(p.cfg.Root, artifact.FileName)).
				Unmarshal(&deviceTemplate)
			if err != nil {
				return err
			}

			err = p.populateDeviceTemplate(ctx, api, m.ServiceType, deviceTemplate)
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
			"Populate device templates",
			1000,
			[]string{"all", "deviceTemplates", "serviceMetadata"},
			func(ctx context.Context) (populate.Populator, error) {
				cfg, err := NewDeviceTemplatePopulatorConfigFromConfig(config.MustFromContext(ctx))
				if err != nil {
					return nil, err
				}
				return &DeviceTemplatePopulator{
					cfg: *cfg,
				}, nil
			}))
}
