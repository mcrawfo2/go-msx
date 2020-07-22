package populate

import (
	"context"
	api "cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement"
	"cto-github.cisco.com/NFV-BU/go-msx/populate"
	"cto-github.cisco.com/NFV-BU/go-msx/resource"
	"cto-github.cisco.com/NFV-BU/go-msx/security/service"
	"github.com/pkg/errors"
	"path"
)

const (
	artifactKeySecretPolicies = "secretPolicies"
)

type SecretPolicyPopulator struct{}

func (p SecretPolicyPopulator) populateSecretPolicy(ctx context.Context, idm api.Api, artifact populate.Artifact) error {
	logger.WithContext(ctx).Infof("Loading policy from %q", artifact.TemplateFileName)

	var policy struct {
		api.SecretPolicySetRequest
		Name string
	}

	err := resource.
		Reference(path.Join(manifestDir, artifact.TemplateFileName)).
		Unmarshal(&policy)
	if err != nil {
		return errors.Wrapf(err, "Unable to load policy %q", artifact.TemplateFileName)
	}

	_, err = idm.StoreSecretPolicy(policy.Name, policy.SecretPolicySetRequest)
	if err != nil {
		return errors.Wrapf(err, "Failed to store secret policy %q", artifact.TemplateFileName)
	}

	return nil
}

func (p SecretPolicyPopulator) Populate(ctx context.Context) error {
	return service.WithDefaultServiceAccount(ctx, func(ctx context.Context) error {

		var manifest populate.Manifest
		err := resource.
			Reference(path.Join(manifestDir, manifestFile)).
			Unmarshal(&manifest)
		if err != nil {
			return err
		}

		idm, _ := api.NewIntegration(ctx)

		logger.WithContext(ctx).Info("Populating capabilities")

		for _, artifact := range manifest[artifactKeySecretPolicies] {
			err = p.populateSecretPolicy(ctx, idm, artifact)
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
			"Populate roles and capabilities",
			1000,
			[]string{"all", "customRolesAndCapabilities", "serviceMetadata"},
			func(ctx context.Context) populate.Populator {
				return &RoleCapabilityPopulator{}
			}))
}
