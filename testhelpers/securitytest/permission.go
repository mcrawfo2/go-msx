package securitytest

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

type MockTokenDetailsProvider struct {
	UserName    string
	Active      bool
	Roles       []string
	Permissions []string
	Tenants     []types.UUID
}

func (m *MockTokenDetailsProvider) TokenDetails(ctx context.Context) (*security.UserContextDetails, error) {
	return &security.UserContextDetails{
		Active:       m.Active,
		Username:     types.NewOptionalStringFromString(m.UserName).Ptr(),
		Roles:        m.Roles,
		Permissions:  m.Permissions,
		Tenants:      m.Tenants,
	}, nil
}

func (m *MockTokenDetailsProvider) IsTokenActive(ctx context.Context) (bool, error) {
	return m.Active, nil
}

func NewMockTokenDetailsProvider() *MockTokenDetailsProvider {
	return &MockTokenDetailsProvider{
		UserName:    "tester",
		Active:      true,
		Roles:       []string{"TESTER"},
		Permissions: []string{},
		Tenants:     []types.UUID{},
	}
}

func MockTokenDetailsProviderFromContext(ctx context.Context) *MockTokenDetailsProvider {
	tokenDetailsProvider := security.TokenDetailsProviderFromContext(ctx)

	// Ensure the latest provider is our mock
	var ok bool
	if tokenDetailsProvider != nil {
		_, ok = tokenDetailsProvider.(*MockTokenDetailsProvider)
	}

	if !ok {
		return nil
	}

	return tokenDetailsProvider.(*MockTokenDetailsProvider)
}

func TokenDetailsProviderInjector(ctx context.Context) context.Context {
	mockTokenDetailsProvider := MockTokenDetailsProviderFromContext(ctx)

	if mockTokenDetailsProvider == nil {
		mockTokenDetailsProvider = NewMockTokenDetailsProvider()
		ctx = security.ContextWithTokenDetailsProvider(ctx, mockTokenDetailsProvider)
	}

	return ctx
}

func PermissionInjector(permissions ...string) types.ContextInjector {
	return func(ctx context.Context) context.Context {
		ctx = TokenDetailsProviderInjector(ctx)
		tokenDetailsProvider := MockTokenDetailsProviderFromContext(ctx)
		tokenDetailsProvider.Permissions = append(tokenDetailsProvider.Permissions, permissions...)
		return ctx
	}
}

func TenantAssignmentInjector(tenantIds ...types.UUID) types.ContextInjector {
	return func(ctx context.Context) context.Context {
		ctx = TokenDetailsProviderInjector(ctx)
		tokenDetailsProvider := MockTokenDetailsProviderFromContext(ctx)
		tokenDetailsProvider.Tenants = append(tokenDetailsProvider.Tenants, tenantIds...)
		return ctx
	}
}
