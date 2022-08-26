// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package service

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

type SecurityService interface {
	// WithDefaultServiceAccount is an action wrapper for system account
	WithDefaultServiceAccount(ctx context.Context, action types.ActionFunc) (err error)
	// DefaultServiceAccount is an action decorator for system account
	DefaultServiceAccount(action types.ActionFunc) types.ActionFunc
	// WithUserAccount is an action wrapper for user account
	WithUserAccount(ctx context.Context, userId types.UUID, action types.ActionFunc) (err error)
	// SwitchUserAccountDecorator is an action decorator factory for user account
	SwitchUserAccountDecorator(userId types.UUID) types.ActionFuncDecorator
}

type productionSecurityService struct{}

func (p productionSecurityService) WithDefaultServiceAccount(ctx context.Context, action types.ActionFunc) (err error) {
	return WithDefaultServiceAccount(ctx, action)
}

func (p productionSecurityService) DefaultServiceAccount(action types.ActionFunc) types.ActionFunc {
	return DefaultServiceAccount(action)
}

func (p productionSecurityService) WithUserAccount(ctx context.Context, userId types.UUID, action types.ActionFunc) (err error) {
	return WithUserAccount(ctx, userId, action)
}

func (p productionSecurityService) SwitchUserAccountDecorator(userId types.UUID) types.ActionFuncDecorator {
	return SwitchUserAccountDecorator(userId)
}

func NewSecurityService(ctx context.Context) SecurityService {
	svc := ContextSecurityService.Get(ctx)
	if svc == nil {
		svc = &productionSecurityService{}
	}
	return svc
}

type contextKey string

const contextKeySecurityService = contextKey("SecurityService")

var ContextSecurityService = types.NewContextKeyAccessor[SecurityService](contextKeySecurityService)
