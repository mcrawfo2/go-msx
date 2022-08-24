// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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
	if err := cfg.LatestValues().Populate(securityAccountsConfig, configRootIntegrationSecurityAccountsDefault); err != nil {
		return nil, err
	}

	return securityAccountsConfig, nil
}

func WithDefaultServiceAccount(ctx context.Context, action types.ActionFunc) (err error) {
	newUserContext, err := LoginDefaultServiceAccount(ctx)
	if err != nil {
		return err
	}

	newCtx := security.ContextWithUserContext(ctx, newUserContext)
	return action(newCtx)
}

func WithUserContext(ctx context.Context, userId string, action types.ActionFunc) (err error) {
	newUserContext, err := LoginWithUser(ctx, userId)
	if err != nil {
		return err
	}

	newCtx := security.ContextWithUserContext(ctx, newUserContext)
	return action(newCtx)
}

func DefaultServiceAccountDecorator(action types.ActionFunc) types.ActionFunc {
	return func(ctx context.Context) error {
		return WithDefaultServiceAccount(ctx, action)
	}
}

func SwitchUserAccountDecorator(userId string) types.ActionFuncDecorator {
	return func(action types.ActionFunc) types.ActionFunc {
		return func(ctx context.Context) error {
			return WithUserContext(ctx, userId, action)
		}
	}
}

func LoginDefaultServiceAccount(ctx context.Context) (*security.UserContext, error) {
	api, err := usermanagement.NewIntegration(ctx)
	if err != nil {
		return nil, err
	}

	cfg, err := NewSecurityAccountsDefaultSettings(config.FromContext(ctx))
	if err != nil {
		return nil, err
	}

	msxResponse, err := api.Login(cfg.Username, cfg.Password)
	if err != nil {
		return nil, err
	}

	loginResponse, ok := msxResponse.Payload.(*usermanagement.LoginResponse)
	if !ok {
		return nil, errors.New("Invalid login response object")
	}

	return security.NewUserContextFromToken(ctx, loginResponse.AccessToken)
}

func LoginWithUser(ctx context.Context, userId string) (*security.UserContext, error) {
	api, err := usermanagement.NewIntegration(ctx)
	if err != nil {
		return nil, err
	}

	cfg, err := NewSecurityAccountsDefaultSettings(config.FromContext(ctx))
	if err != nil {
		return nil, err
	}

	msxResponseSystemToken, err := api.Login(cfg.Username, cfg.Password)
	if err != nil {
		return nil, err
	}

	systemToken, ok := msxResponseSystemToken.Payload.(*usermanagement.LoginResponse)

	msxResponse, err := api.SwitchContext(systemToken.AccessToken, userId)
	if err != nil {
		return nil, err
	}

	loginResponse, ok := msxResponse.Payload.(*usermanagement.LoginResponse)
	if !ok {
		return nil, errors.New("Invalid login response object")
	}

	return security.NewUserContextFromToken(ctx, loginResponse.AccessToken)
}

var defaultSecurityAccountUserContextCache security.UserContextCacheApi

func DefaultSecurityAccountUserContextCache(ctx context.Context) security.UserContextCacheApi {
	if defaultSecurityAccountUserContextCache == nil {
		defaultSecurityAccountUserContextCache = security.NewUserContextCache(ctx, LoginDefaultServiceAccount)
	}

	return defaultSecurityAccountUserContextCache
}
