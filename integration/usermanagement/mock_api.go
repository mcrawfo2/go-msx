// Code generated by mockery v2.9.4. DO NOT EDIT.

package usermanagement

import (
	integration "cto-github.cisco.com/NFV-BU/go-msx/integration"
	mock "github.com/stretchr/testify/mock"

	paging "cto-github.cisco.com/NFV-BU/go-msx/paging"

	types "cto-github.cisco.com/NFV-BU/go-msx/types"
)

// MockUserManagement is an autogenerated mock type for the Api type
type MockUserManagement struct {
	mock.Mock
}

// AddSystemSecrets provides a mock function with given fields: scope, secrets
func (_m *MockUserManagement) AddSystemSecrets(scope string, secrets map[string]string) (*integration.MsxResponse, error) {
	ret := _m.Called(scope, secrets)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, map[string]string) *integration.MsxResponse); ok {
		r0 = rf(scope, secrets)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, map[string]string) error); ok {
		r1 = rf(scope, secrets)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddTenantSecrets provides a mock function with given fields: tenantId, scope, secrets
func (_m *MockUserManagement) AddTenantSecrets(tenantId string, scope string, secrets map[string]string) (*integration.MsxResponse, error) {
	ret := _m.Called(tenantId, scope, secrets)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, string, map[string]string) *integration.MsxResponse); ok {
		r0 = rf(tenantId, scope, secrets)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, map[string]string) error); ok {
		r1 = rf(tenantId, scope, secrets)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BatchCreateCapabilities provides a mock function with given fields: populator, owner, capabilities
func (_m *MockUserManagement) BatchCreateCapabilities(populator bool, owner string, capabilities []CapabilityCreateRequest) (*integration.MsxResponse, error) {
	ret := _m.Called(populator, owner, capabilities)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(bool, string, []CapabilityCreateRequest) *integration.MsxResponse); ok {
		r0 = rf(populator, owner, capabilities)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool, string, []CapabilityCreateRequest) error); ok {
		r1 = rf(populator, owner, capabilities)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BatchUpdateCapabilities provides a mock function with given fields: populator, owner, capabilities
func (_m *MockUserManagement) BatchUpdateCapabilities(populator bool, owner string, capabilities []CapabilityUpdateRequest) (*integration.MsxResponse, error) {
	ret := _m.Called(populator, owner, capabilities)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(bool, string, []CapabilityUpdateRequest) *integration.MsxResponse); ok {
		r0 = rf(populator, owner, capabilities)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool, string, []CapabilityUpdateRequest) error); ok {
		r1 = rf(populator, owner, capabilities)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateRole provides a mock function with given fields: dbinstaller, body
func (_m *MockUserManagement) CreateRole(dbinstaller bool, body RoleCreateRequest) (*integration.MsxResponse, error) {
	ret := _m.Called(dbinstaller, body)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(bool, RoleCreateRequest) *integration.MsxResponse); ok {
		r0 = rf(dbinstaller, body)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool, RoleCreateRequest) error); ok {
		r1 = rf(dbinstaller, body)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteCapability provides a mock function with given fields: populator, owner, name
func (_m *MockUserManagement) DeleteCapability(populator bool, owner string, name string) (*integration.MsxResponse, error) {
	ret := _m.Called(populator, owner, name)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(bool, string, string) *integration.MsxResponse); ok {
		r0 = rf(populator, owner, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool, string, string) error); ok {
		r1 = rf(populator, owner, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteRole provides a mock function with given fields: roleName
func (_m *MockUserManagement) DeleteRole(roleName string) (*integration.MsxResponse, error) {
	ret := _m.Called(roleName)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(roleName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(roleName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteSecretPolicy provides a mock function with given fields: name
func (_m *MockUserManagement) DeleteSecretPolicy(name string) (*integration.MsxResponse, error) {
	ret := _m.Called(name)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EncryptSystemSecrets provides a mock function with given fields: scope, names, encrypt
func (_m *MockUserManagement) EncryptSystemSecrets(scope string, names []string, encrypt EncryptSecretsDTO) (*integration.MsxResponse, error) {
	ret := _m.Called(scope, names, encrypt)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, []string, EncryptSecretsDTO) *integration.MsxResponse); ok {
		r0 = rf(scope, names, encrypt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []string, EncryptSecretsDTO) error); ok {
		r1 = rf(scope, names, encrypt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EncryptTenantSecrets provides a mock function with given fields: tenantId, scope, names, encrypt
func (_m *MockUserManagement) EncryptTenantSecrets(tenantId string, scope string, names []string, encrypt EncryptSecretsDTO) (*integration.MsxResponse, error) {
	ret := _m.Called(tenantId, scope, names, encrypt)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, string, []string, EncryptSecretsDTO) *integration.MsxResponse); ok {
		r0 = rf(tenantId, scope, names, encrypt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, []string, EncryptSecretsDTO) error); ok {
		r1 = rf(tenantId, scope, names, encrypt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateSystemSecrets provides a mock function with given fields: scope, names, save
func (_m *MockUserManagement) GenerateSystemSecrets(scope string, names []string, save bool) (*integration.MsxResponse, error) {
	ret := _m.Called(scope, names, save)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, []string, bool) *integration.MsxResponse); ok {
		r0 = rf(scope, names, save)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []string, bool) error); ok {
		r1 = rf(scope, names, save)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateTenantSecrets provides a mock function with given fields: tenantId, scope, names, save
func (_m *MockUserManagement) GenerateTenantSecrets(tenantId string, scope string, names []string, save bool) (*integration.MsxResponse, error) {
	ret := _m.Called(tenantId, scope, names, save)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, string, []string, bool) *integration.MsxResponse); ok {
		r0 = rf(tenantId, scope, names, save)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, []string, bool) error); ok {
		r1 = rf(tenantId, scope, names, save)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAdminHealth provides a mock function with given fields:
func (_m *MockUserManagement) GetAdminHealth() (*integration.MsxResponse, error) {
	ret := _m.Called()

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func() *integration.MsxResponse); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCapabilities provides a mock function with given fields: p
func (_m *MockUserManagement) GetCapabilities(p paging.Request) (*integration.MsxResponse, error) {
	ret := _m.Called(p)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(paging.Request) *integration.MsxResponse); ok {
		r0 = rf(p)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(paging.Request) error); ok {
		r1 = rf(p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMyProvider provides a mock function with given fields:
func (_m *MockUserManagement) GetMyProvider() (*integration.MsxResponse, error) {
	ret := _m.Called()

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func() *integration.MsxResponse); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProviderByName provides a mock function with given fields: providerName
func (_m *MockUserManagement) GetProviderByName(providerName string) (*integration.MsxResponse, error) {
	ret := _m.Called(providerName)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(providerName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(providerName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProviderExtensionByName provides a mock function with given fields: name
func (_m *MockUserManagement) GetProviderExtensionByName(name string) (*integration.MsxResponse, error) {
	ret := _m.Called(name)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRoles provides a mock function with given fields: resolvePermissionNames, p
func (_m *MockUserManagement) GetRoles(resolvePermissionNames bool, p paging.Request) (*integration.MsxResponse, error) {
	ret := _m.Called(resolvePermissionNames, p)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(bool, paging.Request) *integration.MsxResponse); ok {
		r0 = rf(resolvePermissionNames, p)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool, paging.Request) error); ok {
		r1 = rf(resolvePermissionNames, p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSecretPolicy provides a mock function with given fields: name
func (_m *MockUserManagement) GetSecretPolicy(name string) (*integration.MsxResponse, error) {
	ret := _m.Called(name)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSystemSecrets provides a mock function with given fields: scope
func (_m *MockUserManagement) GetSystemSecrets(scope string) (*integration.MsxResponse, error) {
	ret := _m.Called(scope)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(scope)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(scope)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTenantById provides a mock function with given fields: tenantId
func (_m *MockUserManagement) GetTenantById(tenantId string) (*integration.MsxResponse, error) {
	ret := _m.Called(tenantId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(tenantId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(tenantId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTenantByIdV8 provides a mock function with given fields: tenantId
func (_m *MockUserManagement) GetTenantByIdV8(tenantId string) (*integration.MsxResponse, error) {
	ret := _m.Called(tenantId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(tenantId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(tenantId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTenantByName provides a mock function with given fields: tenantName
func (_m *MockUserManagement) GetTenantByName(tenantName string) (*integration.MsxResponse, error) {
	ret := _m.Called(tenantName)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(tenantName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(tenantName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTenantHierarchyAncestors provides a mock function with given fields: tenantId
func (_m *MockUserManagement) GetTenantHierarchyAncestors(tenantId types.UUID) (*integration.MsxResponse, []types.UUID, error) {
	ret := _m.Called(tenantId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(types.UUID) *integration.MsxResponse); ok {
		r0 = rf(tenantId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 []types.UUID
	if rf, ok := ret.Get(1).(func(types.UUID) []types.UUID); ok {
		r1 = rf(tenantId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]types.UUID)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(types.UUID) error); ok {
		r2 = rf(tenantId)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetTenantHierarchyChildren provides a mock function with given fields: tenantId
func (_m *MockUserManagement) GetTenantHierarchyChildren(tenantId types.UUID) (*integration.MsxResponse, []types.UUID, error) {
	ret := _m.Called(tenantId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(types.UUID) *integration.MsxResponse); ok {
		r0 = rf(tenantId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 []types.UUID
	if rf, ok := ret.Get(1).(func(types.UUID) []types.UUID); ok {
		r1 = rf(tenantId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]types.UUID)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(types.UUID) error); ok {
		r2 = rf(tenantId)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetTenantHierarchyDescendants provides a mock function with given fields: tenantId
func (_m *MockUserManagement) GetTenantHierarchyDescendants(tenantId types.UUID) (*integration.MsxResponse, []types.UUID, error) {
	ret := _m.Called(tenantId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(types.UUID) *integration.MsxResponse); ok {
		r0 = rf(tenantId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 []types.UUID
	if rf, ok := ret.Get(1).(func(types.UUID) []types.UUID); ok {
		r1 = rf(tenantId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]types.UUID)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(types.UUID) error); ok {
		r2 = rf(tenantId)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetTenantHierarchyParent provides a mock function with given fields: tenantId
func (_m *MockUserManagement) GetTenantHierarchyParent(tenantId types.UUID) (*integration.MsxResponse, error) {
	ret := _m.Called(tenantId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(types.UUID) *integration.MsxResponse); ok {
		r0 = rf(tenantId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.UUID) error); ok {
		r1 = rf(tenantId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTenantHierarchyRoot provides a mock function with given fields:
func (_m *MockUserManagement) GetTenantHierarchyRoot() (*integration.MsxResponse, error) {
	ret := _m.Called()

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func() *integration.MsxResponse); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTenantSecrets provides a mock function with given fields: tenantId, scope
func (_m *MockUserManagement) GetTenantSecrets(tenantId string, scope string) (*integration.MsxResponse, error) {
	ret := _m.Called(tenantId, scope)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, string) *integration.MsxResponse); ok {
		r0 = rf(tenantId, scope)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(tenantId, scope)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTokenDetails provides a mock function with given fields: noDetails
func (_m *MockUserManagement) GetTokenDetails(noDetails bool) (*integration.MsxResponse, error) {
	ret := _m.Called(noDetails)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(bool) *integration.MsxResponse); ok {
		r0 = rf(noDetails)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool) error); ok {
		r1 = rf(noDetails)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTokenKeys provides a mock function with given fields:
func (_m *MockUserManagement) GetTokenKeys() (JsonWebKeys, *integration.MsxResponse, error) {
	ret := _m.Called()

	var r0 JsonWebKeys
	if rf, ok := ret.Get(0).(func() JsonWebKeys); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(JsonWebKeys)
	}

	var r1 *integration.MsxResponse
	if rf, ok := ret.Get(1).(func() *integration.MsxResponse); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*integration.MsxResponse)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetUserById provides a mock function with given fields: userId
func (_m *MockUserManagement) GetUserById(userId string) (*integration.MsxResponse, error) {
	ret := _m.Called(userId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByIdV8 provides a mock function with given fields: userId
func (_m *MockUserManagement) GetUserByIdV8(userId string) (*integration.MsxResponse, error) {
	ret := _m.Called(userId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsTokenActive provides a mock function with given fields:
func (_m *MockUserManagement) IsTokenActive() (*integration.MsxResponse, error) {
	ret := _m.Called()

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func() *integration.MsxResponse); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: user, password
func (_m *MockUserManagement) Login(user string, password string) (*integration.MsxResponse, error) {
	ret := _m.Called(user, password)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, string) *integration.MsxResponse); ok {
		r0 = rf(user, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(user, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Logout provides a mock function with given fields:
func (_m *MockUserManagement) Logout() (*integration.MsxResponse, error) {
	ret := _m.Called()

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func() *integration.MsxResponse); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveSystemSecrets provides a mock function with given fields: scope
func (_m *MockUserManagement) RemoveSystemSecrets(scope string) (*integration.MsxResponse, error) {
	ret := _m.Called(scope)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(scope)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(scope)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveSystemSecretsPermanent provides a mock function with given fields: scope, permanent
func (_m *MockUserManagement) RemoveSystemSecretsPermanent(scope string, permanent *bool) (*integration.MsxResponse, error) {
	ret := _m.Called(scope, permanent)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, *bool) *integration.MsxResponse); ok {
		r0 = rf(scope, permanent)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *bool) error); ok {
		r1 = rf(scope, permanent)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveTenantSecrets provides a mock function with given fields: tenantId, scope
func (_m *MockUserManagement) RemoveTenantSecrets(tenantId string, scope string) (*integration.MsxResponse, error) {
	ret := _m.Called(tenantId, scope)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, string) *integration.MsxResponse); ok {
		r0 = rf(tenantId, scope)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(tenantId, scope)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReplaceSystemSecrets provides a mock function with given fields: scope, secrets
func (_m *MockUserManagement) ReplaceSystemSecrets(scope string, secrets map[string]string) (*integration.MsxResponse, error) {
	ret := _m.Called(scope, secrets)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, map[string]string) *integration.MsxResponse); ok {
		r0 = rf(scope, secrets)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, map[string]string) error); ok {
		r1 = rf(scope, secrets)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReplaceTenantSecrets provides a mock function with given fields: tenantId, scope, secrets
func (_m *MockUserManagement) ReplaceTenantSecrets(tenantId string, scope string, secrets map[string]string) (*integration.MsxResponse, error) {
	ret := _m.Called(tenantId, scope, secrets)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, string, map[string]string) *integration.MsxResponse); ok {
		r0 = rf(tenantId, scope, secrets)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, map[string]string) error); ok {
		r1 = rf(tenantId, scope, secrets)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StoreSecretPolicy provides a mock function with given fields: name, policy
func (_m *MockUserManagement) StoreSecretPolicy(name string, policy SecretPolicySetRequest) (*integration.MsxResponse, error) {
	ret := _m.Called(name, policy)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, SecretPolicySetRequest) *integration.MsxResponse); ok {
		r0 = rf(name, policy)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, SecretPolicySetRequest) error); ok {
		r1 = rf(name, policy)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateRole provides a mock function with given fields: dbinstaller, body
func (_m *MockUserManagement) UpdateRole(dbinstaller bool, body RoleUpdateRequest) (*integration.MsxResponse, error) {
	ret := _m.Called(dbinstaller, body)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(bool, RoleUpdateRequest) *integration.MsxResponse); ok {
		r0 = rf(dbinstaller, body)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool, RoleUpdateRequest) error); ok {
		r1 = rf(dbinstaller, body)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
