// Code generated by mockery v2.22.1. DO NOT EDIT.

package monitor

import (
	integration "cto-github.cisco.com/NFV-BU/go-msx/integration"
	mock "github.com/stretchr/testify/mock"
)

// MockMonitor is an autogenerated mock type for the Api type
type MockMonitor struct {
	mock.Mock
}

// GetDeviceHealth provides a mock function with given fields: deviceIds
func (_m *MockMonitor) GetDeviceHealth(deviceIds string) (*integration.MsxResponse, error) {
	ret := _m.Called(deviceIds)

	var r0 *integration.MsxResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*integration.MsxResponse, error)); ok {
		return rf(deviceIds)
	}
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(deviceIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(deviceIds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockMonitor interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockMonitor creates a new instance of MockMonitor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockMonitor(t mockConstructorTestingTNewMockMonitor) *MockMonitor {
	mock := &MockMonitor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
