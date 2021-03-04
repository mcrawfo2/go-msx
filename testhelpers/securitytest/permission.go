package securitytest

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

var defaultUserId = types.MustParseUUID("67f9b089-532e-4b54-9a06-8e4eade2114e")
var defaultTenantId = types.MustParseUUID("960272b3-e800-43e6-86ce-7d51672bd80d")
var defaultProviderId = types.MustParseUUID("30b62544-860e-42fb-93ba-bc7e771dff61")

type MockTokenDetailsProvider struct {
	UserName     string
	UserId       types.UUID
	ClientId     string
	Active       bool
	TenantId     types.UUID
	TenantName   string
	ProviderId   types.UUID
	ProviderName string
	Roles        []string
	Permissions  []string
	Tenants      []types.UUID
}

func (m *MockTokenDetailsProvider) TokenDetails(ctx context.Context) (*security.UserContextDetails, error) {
	return &security.UserContextDetails{
		Active:       m.Active,
		ClientId:     &m.ClientId,
		Username:     &m.UserName,
		UserId:       m.UserId,
		Roles:        m.Roles,
		Permissions:  m.Permissions,
		Tenants:      m.Tenants,
		TenantId:     m.TenantId,
		TenantName:   &m.TenantName,
		ProviderId:   m.ProviderId,
		ProviderName: &m.ProviderName,
	}, nil
}

func (m *MockTokenDetailsProvider) IsTokenActive(ctx context.Context) (bool, error) {
	return m.Active, nil
}

func NewMockTokenDetailsProvider() *MockTokenDetailsProvider {
	return &MockTokenDetailsProvider{
		UserName:     "tester",
		ClientId:     "client-id",
		UserId:       defaultUserId,
		TenantId:     defaultTenantId,
		TenantName:   "test-tenant",
		ProviderId:   defaultProviderId,
		ProviderName: "cisco",
		Active:       true,
		Roles:        []string{"TESTER"},
		Permissions:  []string{},
		Tenants:      []types.UUID{},
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

func TokenDetailsProviderCustomizer(fn func(*MockTokenDetailsProvider)) types.ContextInjector {
	return func(ctx context.Context) context.Context {
		ctx = TokenDetailsProviderInjector(ctx)
		tokenDetailsProvider := MockTokenDetailsProviderFromContext(ctx)
		fn(tokenDetailsProvider)
		return ctx
	}
}

func ClientIdInjector(clientId string) types.ContextInjector {
	return TokenDetailsProviderCustomizer(func(provider *MockTokenDetailsProvider) {
		provider.ClientId = clientId
	})
}

func PermissionInjector(permissions ...string) types.ContextInjector {
	return TokenDetailsProviderCustomizer(func(provider *MockTokenDetailsProvider) {
		provider.Permissions = append(provider.Permissions, permissions...)
	})
}

func TenantAssignmentInjector(tenantIds ...types.UUID) types.ContextInjector {
	return TokenDetailsProviderCustomizer(func(provider *MockTokenDetailsProvider) {
		provider.Tenants = append(provider.Tenants, tenantIds...)
	})
}

func RolesInjector(roles ...string) types.ContextInjector {
	return TokenDetailsProviderCustomizer(func(provider *MockTokenDetailsProvider) {
		provider.Roles = append(provider.Roles, roles...)
	})
}

func AuthoritiesInjector(authorities ...string) types.ContextInjector {
	return func(ctx context.Context) context.Context {
		userContextPtr := security.UserContextFromContext(ctx)
		userContext := *userContextPtr
		userContext.Authorities = append(userContext.Authorities, authorities...)
		ctx = security.ContextWithUserContext(ctx, &userContext)
		return ctx
	}
}
