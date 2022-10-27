// Code generated by mockery v2.14.0. DO NOT EDIT.

package auth

import (
	integration "cto-github.cisco.com/NFV-BU/go-msx/integration"
	mock "github.com/stretchr/testify/mock"

	types "cto-github.cisco.com/NFV-BU/go-msx/types"
)

// MockAuth is an autogenerated mock type for the Api type
type MockAuth struct {
	mock.Mock
}

// GetAdminHealth provides a mock function with given fields:
func (_m *MockAuth) GetAdminHealth() (*integration.MsxResponse, error) {
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

// GetTenantHierarchyAncestors provides a mock function with given fields: tenantId
func (_m *MockAuth) GetTenantHierarchyAncestors(tenantId types.UUID) (*integration.MsxResponse, []types.UUID, error) {
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
func (_m *MockAuth) GetTenantHierarchyChildren(tenantId types.UUID) (*integration.MsxResponse, []types.UUID, error) {
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
func (_m *MockAuth) GetTenantHierarchyDescendants(tenantId types.UUID) (*integration.MsxResponse, []types.UUID, error) {
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
func (_m *MockAuth) GetTenantHierarchyParent(tenantId types.UUID) (*integration.MsxResponse, error) {
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
func (_m *MockAuth) GetTenantHierarchyRoot() (*integration.MsxResponse, error) {
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

// GetTokenDetails provides a mock function with given fields: noDetails
func (_m *MockAuth) GetTokenDetails(noDetails bool) (*integration.MsxResponse, error) {
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
func (_m *MockAuth) GetTokenKeys() (JsonWebKeys, *integration.MsxResponse, error) {
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

// Login provides a mock function with given fields: user, password
func (_m *MockAuth) Login(user string, password string) (*integration.MsxResponse, error) {
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
func (_m *MockAuth) Logout() (*integration.MsxResponse, error) {
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

// SwitchContext provides a mock function with given fields: accessToken, userId
func (_m *MockAuth) SwitchContext(accessToken string, userId types.UUID) (*integration.MsxResponse, error) {
	ret := _m.Called(accessToken, userId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, types.UUID) *integration.MsxResponse); ok {
		r0 = rf(accessToken, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, types.UUID) error); ok {
		r1 = rf(accessToken, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockAuth interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockAuth creates a new instance of MockAuth. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockAuth(t mockConstructorTestingTNewMockAuth) *MockAuth {
	mock := &MockAuth{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
