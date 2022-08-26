// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package securitytest

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"time"
)

var DefaultUserId = types.MustParseUUID("67f9b089-532e-4b54-9a06-8e4eade2114e")
var DefaultTenantId = types.MustParseUUID("960272b3-e800-43e6-86ce-7d51672bd80d")
var DefaultProviderId = types.MustParseUUID("30b62544-860e-42fb-93ba-bc7e771dff61")
var DefaultRootTenantId = types.MustParseUUID("9f7972d0-ea58-4562-aa62-3dc480ae759b")

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

func (m *MockTokenDetailsProvider) WithUserName(userName string) *MockTokenDetailsProvider {
	m.UserName = userName
	return m
}

func (m *MockTokenDetailsProvider) WithUserId(userId types.UUID) *MockTokenDetailsProvider {
	m.UserId = userId
	return m
}

func (m *MockTokenDetailsProvider) WithClientId(clientId string) *MockTokenDetailsProvider {
	m.ClientId = clientId
	return m
}

func (m *MockTokenDetailsProvider) WithActive(active bool) *MockTokenDetailsProvider {
	m.Active = active
	return m
}

func (m *MockTokenDetailsProvider) WithTenantId(tenantId types.UUID) *MockTokenDetailsProvider {
	m.TenantId = tenantId
	return m
}

func (m *MockTokenDetailsProvider) WithTenantName(tenantName string) *MockTokenDetailsProvider {
	m.TenantName = tenantName
	return m
}

func (m *MockTokenDetailsProvider) WithProviderId(providerId types.UUID) *MockTokenDetailsProvider {
	m.ProviderId = providerId
	return m
}

func (m *MockTokenDetailsProvider) WithProviderName(providerName string) *MockTokenDetailsProvider {
	m.ProviderName = providerName
	return m
}

func (m *MockTokenDetailsProvider) WithRole(roleName string) *MockTokenDetailsProvider {
	roles := types.StringStack(m.Roles)
	if !roles.Contains(roleName) {
		m.Roles = append(m.Roles, roleName)
	}
	return m
}

func (m *MockTokenDetailsProvider) WithoutRole(roleName string) *MockTokenDetailsProvider {
	var roles []string
	for _, role := range m.Roles {
		if role != roleName {
			roles = append(roles, role)
		}
	}
	m.Roles = roles
	return m
}

func (m *MockTokenDetailsProvider) WithPermission(permissionName string) *MockTokenDetailsProvider {
	permissions := types.StringStack(m.Permissions)
	if !permissions.Contains(permissionName) {
		m.Permissions = append(m.Permissions, permissionName)
	}
	return m
}

func (m *MockTokenDetailsProvider) WithoutPermission(permissionName string) *MockTokenDetailsProvider {
	var permissions []string
	for _, permission := range m.Permissions {
		if permission != permissionName {
			permissions = append(permissions, permission)
		}
	}
	m.Permissions = permissions
	return m
}

func (m *MockTokenDetailsProvider) WithTenantAssociation(tenantId types.UUID) *MockTokenDetailsProvider {
	for _, tenant := range m.Tenants {
		if tenant.Equals(tenantId) {
			return m
		}
	}
	m.Tenants = append(m.Tenants, tenantId)
	return m
}

func (m *MockTokenDetailsProvider) WithoutTenantAssociation(tenantId types.UUID) *MockTokenDetailsProvider {
	for n, tenant := range m.Tenants {
		if tenant.Equals(tenantId) {
			m.Tenants = append(m.Tenants[:n], m.Tenants[n+1:]...)
			break
		}
	}
	return m
}

func (m *MockTokenDetailsProvider) TokenDetails(_ context.Context) (*security.UserContextDetails, error) {
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

func (m *MockTokenDetailsProvider) UserContext() *security.UserContext {
	return &security.UserContext{
		UserName:    m.UserName,
		Roles:       m.Roles,
		TenantId:    m.TenantId,
		Scopes:      []string{"read", "write"},
		Authorities: []string{"ROLE_CLIENT"},
		FirstName:   "first-name",
		LastName:    "last-name",
		Issuer:      "",
		Subject:     m.UserName,
		Exp:         int(time.Now().UTC().Add(1 * time.Hour).Unix()),
		IssuedAt:    int(time.Now().UTC().Add(-1 * time.Hour).Unix()),
		Jti:         types.MustNewUUID().String(),
		Email:       "user@ciscomsx.com",
		Token:       "token",
		Certificate: nil,
		ClientId:    "client-id",
	}
}

func (m *MockTokenDetailsProvider) IsTokenActive(_ context.Context) (bool, error) {
	return m.Active, nil
}

func (m *MockTokenDetailsProvider) Inject(ctx context.Context) context.Context {
	return security.ContextWithTokenDetailsProvider(ctx, m)
}

func NewMockTokenDetailsProvider() *MockTokenDetailsProvider {
	return &MockTokenDetailsProvider{
		UserName:     "tester",
		ClientId:     "client-id",
		UserId:       DefaultUserId,
		TenantId:     DefaultTenantId,
		TenantName:   "test-tenant",
		ProviderId:   DefaultProviderId,
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
