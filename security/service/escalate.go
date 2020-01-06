package service

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
)

const (
	configRootIntegrationSecurityAccountsDefault = "integration.security.accounts.default"
)

type SecurityAccountsDefaultSettings struct {
	Username string `config:"default=system"`
	Password string `config:"default=system"`
}

func NewSecurityAccountsDefaultSettings(ctx context.Context) (cfg *SecurityAccountsDefaultSettings, err error) {
	cfg = &SecurityAccountsDefaultSettings{}
	err = config.FromContext(ctx).Populate(cfg, configRootIntegrationSecurityAccountsDefault)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func WithDefaultServiceAccount(ctx context.Context, action types.ActionFunc) (err error) {
	api, err := usermanagement.NewIntegration(ctx)
	if err != nil {
		return err
	}

	cfg, err := NewSecurityAccountsDefaultSettings(ctx)
	if err != nil {
		return err
	}

	msxResponse, err := api.Login(cfg.Username, cfg.Password)
	if err != nil {
		return err
	}

	loginResponse, ok := msxResponse.Payload.(*usermanagement.LoginResponse)
	if !ok {
		return errors.New("Invalid login response object")
	}

	token := loginResponse.AccessToken
	newUserContext, err := security.NewUserContextFromToken(ctx, token)
	if err != nil {
		return err
	}

	newCtx := security.ContextWithUserContext(ctx, newUserContext)
	return action(newCtx)
}
