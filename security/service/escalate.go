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

func NewSecurityAccountsDefaultSettings(cfg *config.Config) (*SecurityAccountsDefaultSettings, error) {
	securityAccountsConfig := &SecurityAccountsDefaultSettings{}
	if err := cfg.Populate(securityAccountsConfig, configRootIntegrationSecurityAccountsDefault); err != nil {
		return nil, err
	}

	return securityAccountsConfig, nil
}

func WithDefaultServiceAccount(ctx context.Context, action types.ActionFunc) (err error) {
	api, err := usermanagement.NewIntegration(ctx)
	if err != nil {
		return err
	}

	cfg, err := NewSecurityAccountsDefaultSettings(config.FromContext(ctx))
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
