// Code generated by mockery v2.3.0. DO NOT EDIT.

package platform

import (
	context "context"
	http "net/http"

	mock "github.com/stretchr/testify/mock"

	msx_platform_client "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
)

// MockDevicesApi is an autogenerated mock type for the DevicesApi type
type MockDevicesApi struct {
	mock.Mock
}

// AttachDeviceTemplates provides a mock function with given fields: ctx, id, deviceTemplateAttachRequest
func (_m *MockDevicesApi) AttachDeviceTemplates(ctx context.Context, id string, deviceTemplateAttachRequest msx_platform_client.DeviceTemplateAttachRequest) ([]msx_platform_client.DeviceTemplateHistory, *http.Response, error) {
	ret := _m.Called(ctx, id, deviceTemplateAttachRequest)

	var r0 []msx_platform_client.DeviceTemplateHistory
	if rf, ok := ret.Get(0).(func(context.Context, string, msx_platform_client.DeviceTemplateAttachRequest) []msx_platform_client.DeviceTemplateHistory); ok {
		r0 = rf(ctx, id, deviceTemplateAttachRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]msx_platform_client.DeviceTemplateHistory)
		}
	}

	var r1 *http.Response
	if rf, ok := ret.Get(1).(func(context.Context, string, msx_platform_client.DeviceTemplateAttachRequest) *http.Response); ok {
		r1 = rf(ctx, id, deviceTemplateAttachRequest)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*http.Response)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, msx_platform_client.DeviceTemplateAttachRequest) error); ok {
		r2 = rf(ctx, id, deviceTemplateAttachRequest)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// CreateDevice provides a mock function with given fields: ctx, deviceCreate
func (_m *MockDevicesApi) CreateDevice(ctx context.Context, deviceCreate msx_platform_client.DeviceCreate) (msx_platform_client.Device, *http.Response, error) {
	ret := _m.Called(ctx, deviceCreate)

	var r0 msx_platform_client.Device
	if rf, ok := ret.Get(0).(func(context.Context, msx_platform_client.DeviceCreate) msx_platform_client.Device); ok {
		r0 = rf(ctx, deviceCreate)
	} else {
		r0 = ret.Get(0).(msx_platform_client.Device)
	}

	var r1 *http.Response
	if rf, ok := ret.Get(1).(func(context.Context, msx_platform_client.DeviceCreate) *http.Response); ok {
		r1 = rf(ctx, deviceCreate)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*http.Response)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, msx_platform_client.DeviceCreate) error); ok {
		r2 = rf(ctx, deviceCreate)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DeleteDevice provides a mock function with given fields: ctx, id
func (_m *MockDevicesApi) DeleteDevice(ctx context.Context, id string) (*http.Response, error) {
	ret := _m.Called(ctx, id)

	var r0 *http.Response
	if rf, ok := ret.Get(0).(func(context.Context, string) *http.Response); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DetachDeviceTemplate provides a mock function with given fields: ctx, id, templateId
func (_m *MockDevicesApi) DetachDeviceTemplate(ctx context.Context, id string, templateId string) ([]msx_platform_client.DeviceTemplateHistory, *http.Response, error) {
	ret := _m.Called(ctx, id, templateId)

	var r0 []msx_platform_client.DeviceTemplateHistory
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []msx_platform_client.DeviceTemplateHistory); ok {
		r0 = rf(ctx, id, templateId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]msx_platform_client.DeviceTemplateHistory)
		}
	}

	var r1 *http.Response
	if rf, ok := ret.Get(1).(func(context.Context, string, string) *http.Response); ok {
		r1 = rf(ctx, id, templateId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*http.Response)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctx, id, templateId)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DetachDeviceTemplates provides a mock function with given fields: ctx, id
func (_m *MockDevicesApi) DetachDeviceTemplates(ctx context.Context, id string) ([]msx_platform_client.DeviceTemplateHistory, *http.Response, error) {
	ret := _m.Called(ctx, id)

	var r0 []msx_platform_client.DeviceTemplateHistory
	if rf, ok := ret.Get(0).(func(context.Context, string) []msx_platform_client.DeviceTemplateHistory); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]msx_platform_client.DeviceTemplateHistory)
		}
	}

	var r1 *http.Response
	if rf, ok := ret.Get(1).(func(context.Context, string) *http.Response); ok {
		r1 = rf(ctx, id)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*http.Response)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, id)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetDevice provides a mock function with given fields: ctx, id
func (_m *MockDevicesApi) GetDevice(ctx context.Context, id string) (msx_platform_client.Device, *http.Response, error) {
	ret := _m.Called(ctx, id)

	var r0 msx_platform_client.Device
	if rf, ok := ret.Get(0).(func(context.Context, string) msx_platform_client.Device); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(msx_platform_client.Device)
	}

	var r1 *http.Response
	if rf, ok := ret.Get(1).(func(context.Context, string) *http.Response); ok {
		r1 = rf(ctx, id)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*http.Response)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, id)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetDeviceConfig provides a mock function with given fields: ctx, id
func (_m *MockDevicesApi) GetDeviceConfig(ctx context.Context, id string) (string, *http.Response, error) {
	ret := _m.Called(ctx, id)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 *http.Response
	if rf, ok := ret.Get(1).(func(context.Context, string) *http.Response); ok {
		r1 = rf(ctx, id)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*http.Response)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, id)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetDeviceTemplateHistory provides a mock function with given fields: ctx, id
func (_m *MockDevicesApi) GetDeviceTemplateHistory(ctx context.Context, id string) ([]msx_platform_client.DeviceTemplateHistory, *http.Response, error) {
	ret := _m.Called(ctx, id)

	var r0 []msx_platform_client.DeviceTemplateHistory
	if rf, ok := ret.Get(0).(func(context.Context, string) []msx_platform_client.DeviceTemplateHistory); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]msx_platform_client.DeviceTemplateHistory)
		}
	}

	var r1 *http.Response
	if rf, ok := ret.Get(1).(func(context.Context, string) *http.Response); ok {
		r1 = rf(ctx, id)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*http.Response)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, id)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetDevicesPage provides a mock function with given fields: ctx, page, pageSize, localVarOptionals
func (_m *MockDevicesApi) GetDevicesPage(ctx context.Context, page int32, pageSize int32, localVarOptionals *msx_platform_client.GetDevicesPageOpts) (msx_platform_client.DevicesPage, *http.Response, error) {
	ret := _m.Called(ctx, page, pageSize, localVarOptionals)

	var r0 msx_platform_client.DevicesPage
	if rf, ok := ret.Get(0).(func(context.Context, int32, int32, *msx_platform_client.GetDevicesPageOpts) msx_platform_client.DevicesPage); ok {
		r0 = rf(ctx, page, pageSize, localVarOptionals)
	} else {
		r0 = ret.Get(0).(msx_platform_client.DevicesPage)
	}

	var r1 *http.Response
	if rf, ok := ret.Get(1).(func(context.Context, int32, int32, *msx_platform_client.GetDevicesPageOpts) *http.Response); ok {
		r1 = rf(ctx, page, pageSize, localVarOptionals)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*http.Response)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, int32, int32, *msx_platform_client.GetDevicesPageOpts) error); ok {
		r2 = rf(ctx, page, pageSize, localVarOptionals)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// RedeployDevice provides a mock function with given fields: ctx, id
func (_m *MockDevicesApi) RedeployDevice(ctx context.Context, id string) (*http.Response, error) {
	ret := _m.Called(ctx, id)

	var r0 *http.Response
	if rf, ok := ret.Get(0).(func(context.Context, string) *http.Response); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateDeviceTemplates provides a mock function with given fields: ctx, id, deviceTemplateUpdateRequest
func (_m *MockDevicesApi) UpdateDeviceTemplates(ctx context.Context, id string, deviceTemplateUpdateRequest msx_platform_client.DeviceTemplateUpdateRequest) ([]msx_platform_client.DeviceTemplateHistory, *http.Response, error) {
	ret := _m.Called(ctx, id, deviceTemplateUpdateRequest)

	var r0 []msx_platform_client.DeviceTemplateHistory
	if rf, ok := ret.Get(0).(func(context.Context, string, msx_platform_client.DeviceTemplateUpdateRequest) []msx_platform_client.DeviceTemplateHistory); ok {
		r0 = rf(ctx, id, deviceTemplateUpdateRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]msx_platform_client.DeviceTemplateHistory)
		}
	}

	var r1 *http.Response
	if rf, ok := ret.Get(1).(func(context.Context, string, msx_platform_client.DeviceTemplateUpdateRequest) *http.Response); ok {
		r1 = rf(ctx, id, deviceTemplateUpdateRequest)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*http.Response)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, msx_platform_client.DeviceTemplateUpdateRequest) error); ok {
		r2 = rf(ctx, id, deviceTemplateUpdateRequest)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
