// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package service

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/auth"
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

// WithDefaultServiceAccount is an Action wrapper to escalate to the system account
func WithDefaultServiceAccount(ctx context.Context, action types.ActionFunc) (err error) {
	newUserContext, err := LoginDefaultServiceAccount(ctx)
	if err != nil {
		return err
	}

	newCtx := security.ContextWithUserContext(ctx, newUserContext)
	return action(newCtx)
}

// DefaultServiceAccount is an Action decorator to escalate to the system account
func DefaultServiceAccount(action types.ActionFunc) types.ActionFunc {
	return func(ctx context.Context) error {
		return WithDefaultServiceAccount(ctx, action)
	}
}

// LoginDefaultServiceAccount returns a new security.UserContext by escalating to the system account
func LoginDefaultServiceAccount(ctx context.Context) (*security.UserContext, error) {
	api, err := auth.NewIntegration(ctx)
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

	loginResponse, ok := msxResponse.Payload.(*auth.LoginResponse)
	if !ok {
		return nil, errors.New("Invalid login response object")
	}

	return security.NewUserContextFromToken(ctx, loginResponse.AccessToken)
}

// WithUserAccount is an Action wrapper to switch user to the specified user
func WithUserAccount(ctx context.Context, userId types.UUID, action types.ActionFunc) (err error) {
	newUserContext, err := LoginUserAccount(ctx, userId)
	if err != nil {
		return err
	}

	newCtx := security.ContextWithUserContext(ctx, newUserContext)
	return action(newCtx)
}

// SwitchUserAccountDecorator is an Action decorator factory to switch user to the specified user
func SwitchUserAccountDecorator(userId types.UUID) types.ActionFuncDecorator {
	return func(action types.ActionFunc) types.ActionFunc {
		return func(ctx context.Context) error {
			return WithUserAccount(ctx, userId, action)
		}
	}
}

// LoginUserAccount returns a new security.UserContext by switching to the specified user
func LoginUserAccount(ctx context.Context, userId types.UUID) (*security.UserContext, error) {
	api, err := auth.NewIntegration(ctx)
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

	systemToken, ok := msxResponseSystemToken.Payload.(*auth.LoginResponse)

	msxResponse, err := api.SwitchContext(systemToken.AccessToken, userId)
	if err != nil {
		return nil, err
	}

	loginResponse, ok := msxResponse.Payload.(*auth.LoginResponse)
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
