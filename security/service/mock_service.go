package service

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/securitytest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

const systemUserName = "system"
const userUserName = "user"

type MockSecurityService struct {
	SystemTokenDetailsProvider *securitytest.MockTokenDetailsProvider
	SystemEscalationError      error
	UserTokenDetailsProvider   *securitytest.MockTokenDetailsProvider
	UserSwitchError            error
}

func (m *MockSecurityService) WithDefaultServiceAccount(ctx context.Context, action types.ActionFunc) (err error) {
	if m.SystemEscalationError != nil {
		return m.SystemEscalationError
	}

	// Override token details provider
	ctx = m.SystemTokenDetailsProvider.Inject(ctx)

	// Override user context
	userContext := m.SystemTokenDetailsProvider.UserContext()
	ctx = security.ContextWithUserContext(ctx, userContext)

	// Call the action
	return action(ctx)
}

func (m *MockSecurityService) DefaultServiceAccount(action types.ActionFunc) types.ActionFunc {
	return func(ctx context.Context) error {
		return WithDefaultServiceAccount(ctx, action)
	}
}

func (m *MockSecurityService) WithUserAccount(ctx context.Context, userId types.UUID, action types.ActionFunc) (err error) {
	if m.UserSwitchError != nil {
		return m.UserSwitchError
	}

	// Override token details provider
	m.UserTokenDetailsProvider.UserId = userId
	ctx = m.UserTokenDetailsProvider.Inject(ctx)

	// Override user context
	userContext := m.UserTokenDetailsProvider.UserContext()
	ctx = security.ContextWithUserContext(ctx, userContext)

	// Call the action
	return action(ctx)
}

func (m *MockSecurityService) SwitchUserAccountDecorator(userId types.UUID) types.ActionFuncDecorator {
	return func(action types.ActionFunc) types.ActionFunc {
		return func(ctx context.Context) error {
			return m.WithUserAccount(ctx, userId, action)
		}
	}
}

func (m *MockSecurityService) Inject(ctx context.Context) context.Context {
	return ContextSecurityService.Set(ctx, m)
}

func NewMockSecurityService() *MockSecurityService {
	systemTokenDetailsProvider := securitytest.NewMockTokenDetailsProvider()
	systemTokenDetailsProvider.UserId = securitytest.DefaultUserId
	systemTokenDetailsProvider.UserName = systemUserName

	userTokenDetailsProvider := securitytest.NewMockTokenDetailsProvider()
	userTokenDetailsProvider.UserName = userUserName

	return &MockSecurityService{
		SystemTokenDetailsProvider: systemTokenDetailsProvider,
		UserTokenDetailsProvider:   userTokenDetailsProvider,
	}
}
