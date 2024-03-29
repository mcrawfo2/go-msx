// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package populate

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/secrets"
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

func NewSecretPolicyPopulatorConfigFromConfig(cfg *config.Config) (*SecretPolicyPopulatorConfig, error) {
	var populatorConfig SecretPolicyPopulatorConfig
	if err := cfg.Populate(&populatorConfig, secretPolicyPopulatorConfigRoot); err != nil {
		return nil, err
	}

	return &populatorConfig, nil
}

type secretPolicyManifest struct {
	SecretPolicies []populate.Artifact `json:"secretPolicies"`
}

type SecretPolicyPopulator struct {
	cfg SecretPolicyPopulatorConfig
}

func (p SecretPolicyPopulator) populateSecretPolicy(ctx context.Context, api secrets.Api, artifact populate.Artifact) error {
	logger.WithContext(ctx).Infof("Loading policy from %q", artifact.TemplateFileName)

	var policy struct {
		secrets.SecretPolicySetRequest
		Name string `json:"name"`
	}

	err := resource.
		Reference(path.Join(p.cfg.Root, artifact.TemplateFileName)).
		Unmarshal(&policy)
	if err != nil {
		return errors.Wrapf(err, "Unable to load policy %q", artifact.TemplateFileName)
	}

	_, err = api.StoreSecretPolicy(policy.Name, policy.SecretPolicySetRequest)
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

		api, _ := secrets.NewIntegration(ctx)

		logger.WithContext(ctx).Info("Populating secret policies")

		for _, artifact := range manifest.SecretPolicies {
			err = p.populateSecretPolicy(ctx, api, artifact)
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
				cfg, err := NewSecretPolicyPopulatorConfigFromConfig(config.MustFromContext(ctx))
				if err != nil {
					return nil, err
				}
				return &SecretPolicyPopulator{
					cfg: *cfg,
				}, nil
			}))
}
