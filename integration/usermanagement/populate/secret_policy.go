package populate

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	api "cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement"
	"cto-github.cisco.com/NFV-BU/go-msx/populate"
	"cto-github.cisco.com/NFV-BU/go-msx/resource"
	"cto-github.cisco.com/NFV-BU/go-msx/security/service"
	"github.com/pkg/errors"
	"path"
)

const (
	artifactKeySecretPolicies       = "secretPolicies"
	secretPolicyPopulatorConfigRoot = "populate.usermanagement.secret-policy"
)

type SecretPolicyPopulatorConfig struct {
	Enabled bool   `config:"default=false"`
	Root    string `config:"default=${populate.root}/usermanagement"`
}

type secretPolicyManifest struct {
	SecretPolicies []populate.Artifact `json:"secretPolicies"`
}

type SecretPolicyPopulator struct {
	cfg SecretPolicyPopulatorConfig
}

func (p SecretPolicyPopulator) populateSecretPolicy(ctx context.Context, idm api.Api, artifact populate.Artifact) error {
	logger.WithContext(ctx).Infof("Loading policy from %q", artifact.TemplateFileName)

	var policy struct {
		api.SecretPolicySetRequest
		Name string `json:"name"`
	}

	err := resource.
		Reference(path.Join(p.cfg.Root, artifact.TemplateFileName)).
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
	if !p.cfg.Enabled {
		logger.WithContext(ctx).Warn("Secret Policy populator disabled.")
		return nil
	}

	return service.WithDefaultServiceAccount(ctx, func(ctx context.Context) error {

		var manifest secretPolicyManifest
		err := resource.
			Reference(path.Join(p.cfg.Root, manifestFile)).
			Unmarshal(&manifest)
		if err != nil {
			return err
		}

		idm, _ := api.NewIntegration(ctx)

		logger.WithContext(ctx).Info("Populating secret policies")

		for _, artifact := range manifest.SecretPolicies {
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
			"Populate secret policies",
			1000,
			[]string{"all", "secretPolicies", "serviceMetadata"},
			func(ctx context.Context) (populate.Populator, error) {
				var cfg SecretPolicyPopulatorConfig
				err := config.MustFromContext(ctx).Populate(&cfg, secretPolicyPopulatorConfigRoot)
				if err != nil {
					return nil, err
				}
				return &SecretPolicyPopulator{
					cfg: cfg,
				}, nil
			}))
}
