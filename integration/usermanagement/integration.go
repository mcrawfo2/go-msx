// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package usermanagement

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/auth"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/idm"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/secrets"
)

func NewIntegration(ctx context.Context) (Api, error) {
	integrationInstance := IntegrationFromContext(ctx)
	if integrationInstance == nil {
		authInt, err := auth.NewIntegration(ctx)
		if err != nil {
			return nil, err
		}
		idmInt, err := idm.NewIntegration(ctx)
		if err != nil {
			return nil, err
		}
		secretsInt, err := secrets.NewIntegration(ctx)
		if err != nil {
			return nil, err
		}

		integrationInstance = &Integration{
			AuthApi:    authInt,
			IdmApi:     idmInt,
			SecretsApi: secretsInt,
		}
	}
	return integrationInstance, nil
}

func NewIntegrationWithExecutor(executor integration.MsxContextServiceExecutor) Api {
	ctx := executor.Context()

	authInt, err := auth.NewIntegration(ctx)
	if err != nil {
		return nil
	}
	idmInt, err := idm.NewIntegration(ctx)
	if err != nil {
		return nil
	}
	secretsInt, err := secrets.NewIntegration(ctx)
	if err != nil {
		return nil
	}

	return &Integration{
		AuthApi:    authInt,
		IdmApi:     idmInt,
		SecretsApi: secretsInt,
	}
}

type Integration struct {
	auth.AuthApi
	idm.IdmApi
	secrets.SecretsApi
}
