package populate

import (
	"context"
	api "cto-github.cisco.com/NFV-BU/go-msx/integration/serviceconfigmanager"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/populate"
	"cto-github.cisco.com/NFV-BU/go-msx/resource"
	"cto-github.cisco.com/NFV-BU/go-msx/security/service"
	"github.com/pkg/errors"
	"path"
)

const (
	manifestDir  = "/platform-common/serviceconfig"
	manifestFile = "manifest.json"

	artifactKeyServiceConfigs = "serviceconfigs"
)

var logger = log.NewLogger("msx.integration.serviceconfigmanager.populate")

type ServiceConfigPopulator struct{}

func (p ServiceConfigPopulator) Populate(ctx context.Context) error {
	return service.WithDefaultServiceAccount(ctx, func(ctx context.Context) error {
		logger.WithContext(ctx).Info("Populating service configs")

		var manifest populate.Manifest
		err := resource.
			Reference(path.Join(manifestDir, manifestFile)).
			Unmarshal(&manifest)
		if err != nil {
			return err
		}

		scm, err := api.NewIntegration(ctx)
		if err != nil {
			return err
		}

		for _, serviceDefinitionArtifact := range manifest[artifactKeyServiceConfigs] {
			err = p.populateServiceConfig(ctx, scm, serviceDefinitionArtifact)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (p ServiceConfigPopulator) populateServiceConfig(ctx context.Context, scm api.Api, artifact populate.Artifact) (err error) {
	logger.WithContext(ctx).Infof("Populating service config %q", artifact.TemplateFileName)

	var request api.ServiceConfigurationRequest
	err = resource.Reference(
		path.Join(manifestDir, artifact.TemplateFileName)).
		Unmarshal(&request)
	if err != nil {
		return errors.Wrapf(err, "Failed to load service config %q", artifact.TemplateFileName)
	}

	_, err = scm.CreateServiceConfiguration(request)
	if err != nil {
		return errors.Wrapf(err, "Failed to populate service config %q", artifact.TemplateFileName)
	}

	logger.WithContext(ctx).Infof("Successfully populated service config %q", artifact.TemplateFileName)
	return nil
}

func init() {
	populate.RegisterPopulationTask(
		populate.NewPopulatorTask(
			"Populate service configurations",
			1000,
			[]string{"all", "serviceConfig", "serviceMetadata"},
			func(ctx context.Context) populate.Populator {
				return &ServiceConfigPopulator{}
			}))
}
