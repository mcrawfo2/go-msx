// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package manage

import (
	integration "cto-github.cisco.com/NFV-BU/go-msx/integration"
	mock "github.com/stretchr/testify/mock"
)

// MockManage is an autogenerated mock type for the Api type
type MockManage struct {
	mock.Mock
}

// AddDeviceToSiteV3 provides a mock function with given fields: deviceId, siteId, notification
func (_m *MockManage) AddDeviceToSiteV3(deviceId string, siteId string, notification string) (*integration.MsxResponse, error) {
	ret := _m.Called(deviceId, siteId, notification)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, string, string) *integration.MsxResponse); ok {
		r0 = rf(deviceId, siteId, notification)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(deviceId, siteId, notification)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AttachDeviceTemplates provides a mock function with given fields: deviceId, attachTemplateRequest
func (_m *MockManage) AttachDeviceTemplates(deviceId string, attachTemplateRequest AttachTemplateRequest) (*integration.MsxResponse, error) {
	ret := _m.Called(deviceId, attachTemplateRequest)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, AttachTemplateRequest) *integration.MsxResponse); ok {
		r0 = rf(deviceId, attachTemplateRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, AttachTemplateRequest) error); ok {
		r1 = rf(deviceId, attachTemplateRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ConnectControlPlane provides a mock function with given fields: controlPlaneId
func (_m *MockManage) ConnectControlPlane(controlPlaneId string) (*integration.MsxResponse, error) {
	ret := _m.Called(controlPlaneId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(controlPlaneId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(controlPlaneId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ConnectUnmanagedControlPlane provides a mock function with given fields: username, password, url, resourceProvider, tlsInsecure
func (_m *MockManage) ConnectUnmanagedControlPlane(username string, password string, url string, resourceProvider string, tlsInsecure bool) (*integration.MsxResponse, error) {
	ret := _m.Called(username, password, url, resourceProvider, tlsInsecure)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, string, string, string, bool) *integration.MsxResponse); ok {
		r0 = rf(username, password, url, resourceProvider, tlsInsecure)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, string, bool) error); ok {
		r1 = rf(username, password, url, resourceProvider, tlsInsecure)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateControlPlane provides a mock function with given fields: tenantId, name, url, resourceProvider, authenticationType, tlsInsecure, attributes
func (_m *MockManage) CreateControlPlane(tenantId string, name string, url string, resourceProvider string, authenticationType string, tlsInsecure bool, attributes map[string]string) (*integration.MsxResponse, error) {
	ret := _m.Called(tenantId, name, url, resourceProvider, authenticationType, tlsInsecure, attributes)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, string, string, string, string, bool, map[string]string) *integration.MsxResponse); ok {
		r0 = rf(tenantId, name, url, resourceProvider, authenticationType, tlsInsecure, attributes)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, string, string, bool, map[string]string) error); ok {
		r1 = rf(tenantId, name, url, resourceProvider, authenticationType, tlsInsecure, attributes)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateDevice provides a mock function with given fields: subscriptionId, deviceInstanceId, deviceAttribute, deviceDefAttribute, status
func (_m *MockManage) CreateDevice(subscriptionId string, deviceInstanceId *string, deviceAttribute map[string]string, deviceDefAttribute map[string]string, status map[string]string) (*integration.MsxResponse, error) {
	ret := _m.Called(subscriptionId, deviceInstanceId, deviceAttribute, deviceDefAttribute, status)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, *string, map[string]string, map[string]string, map[string]string) *integration.MsxResponse); ok {
		r0 = rf(subscriptionId, deviceInstanceId, deviceAttribute, deviceDefAttribute, status)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *string, map[string]string, map[string]string, map[string]string) error); ok {
		r1 = rf(subscriptionId, deviceInstanceId, deviceAttribute, deviceDefAttribute, status)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateDeviceV4 provides a mock function with given fields: deviceRequest
func (_m *MockManage) CreateDeviceV4(deviceRequest DeviceCreateRequest) (*integration.MsxResponse, error) {
	ret := _m.Called(deviceRequest)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(DeviceCreateRequest) *integration.MsxResponse); ok {
		r0 = rf(deviceRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(DeviceCreateRequest) error); ok {
		r1 = rf(deviceRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateManagedDevice provides a mock function with given fields: tenantId, deviceModel, deviceOnboardType, deviceOnboardInfo
func (_m *MockManage) CreateManagedDevice(tenantId string, deviceModel string, deviceOnboardType string, deviceOnboardInfo map[string]string) (*integration.MsxResponse, error) {
	ret := _m.Called(tenantId, deviceModel, deviceOnboardType, deviceOnboardInfo)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, string, string, map[string]string) *integration.MsxResponse); ok {
		r0 = rf(tenantId, deviceModel, deviceOnboardType, deviceOnboardInfo)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, map[string]string) error); ok {
		r1 = rf(tenantId, deviceModel, deviceOnboardType, deviceOnboardInfo)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateServiceInstance provides a mock function with given fields: subscriptionId, serviceInstanceId, serviceAttribute, serviceDefAttribute, status
func (_m *MockManage) CreateServiceInstance(subscriptionId string, serviceInstanceId string, serviceAttribute map[string]string, serviceDefAttribute map[string]string, status map[string]string) (*integration.MsxResponse, error) {
	ret := _m.Called(subscriptionId, serviceInstanceId, serviceAttribute, serviceDefAttribute, status)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, string, map[string]string, map[string]string, map[string]string) *integration.MsxResponse); ok {
		r0 = rf(subscriptionId, serviceInstanceId, serviceAttribute, serviceDefAttribute, status)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, map[string]string, map[string]string, map[string]string) error); ok {
		r1 = rf(subscriptionId, serviceInstanceId, serviceAttribute, serviceDefAttribute, status)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateSite provides a mock function with given fields: subscriptionId, serviceInstanceId, siteId, siteName, siteType, displayName, siteAttributes, siteDefAttributes, devices
func (_m *MockManage) CreateSite(subscriptionId string, serviceInstanceId string, siteId *string, siteName *string, siteType *string, displayName *string, siteAttributes map[string]string, siteDefAttributes map[string]string, devices []string) (*integration.MsxResponse, error) {
	ret := _m.Called(subscriptionId, serviceInstanceId, siteId, siteName, siteType, displayName, siteAttributes, siteDefAttributes, devices)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, string, *string, *string, *string, *string, map[string]string, map[string]string, []string) *integration.MsxResponse); ok {
		r0 = rf(subscriptionId, serviceInstanceId, siteId, siteName, siteType, displayName, siteAttributes, siteDefAttributes, devices)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, *string, *string, *string, *string, map[string]string, map[string]string, []string) error); ok {
		r1 = rf(subscriptionId, serviceInstanceId, siteId, siteName, siteType, displayName, siteAttributes, siteDefAttributes, devices)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateSiteV3 provides a mock function with given fields: siteRequest
func (_m *MockManage) CreateSiteV3(siteRequest SiteCreateRequest) (*integration.MsxResponse, error) {
	ret := _m.Called(siteRequest)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(SiteCreateRequest) *integration.MsxResponse); ok {
		r0 = rf(siteRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(SiteCreateRequest) error); ok {
		r1 = rf(siteRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateSubscription provides a mock function with given fields: tenantId, serviceType, subscriptionName, subscriptionAttribute, offerDefAttribute, offerSelectionDetail, costAttribute
func (_m *MockManage) CreateSubscription(tenantId string, serviceType string, subscriptionName *string, subscriptionAttribute map[string]string, offerDefAttribute map[string]string, offerSelectionDetail map[string]string, costAttribute map[string]string) (*integration.MsxResponse, error) {
	ret := _m.Called(tenantId, serviceType, subscriptionName, subscriptionAttribute, offerDefAttribute, offerSelectionDetail, costAttribute)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, string, *string, map[string]string, map[string]string, map[string]string, map[string]string) *integration.MsxResponse); ok {
		r0 = rf(tenantId, serviceType, subscriptionName, subscriptionAttribute, offerDefAttribute, offerSelectionDetail, costAttribute)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, *string, map[string]string, map[string]string, map[string]string, map[string]string) error); ok {
		r1 = rf(tenantId, serviceType, subscriptionName, subscriptionAttribute, offerDefAttribute, offerSelectionDetail, costAttribute)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteControlPlane provides a mock function with given fields: controlPlaneId
func (_m *MockManage) DeleteControlPlane(controlPlaneId string) (*integration.MsxResponse, error) {
	ret := _m.Called(controlPlaneId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(controlPlaneId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(controlPlaneId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteDevice provides a mock function with given fields: deviceInstanceId
func (_m *MockManage) DeleteDevice(deviceInstanceId string) (*integration.MsxResponse, error) {
	ret := _m.Called(deviceInstanceId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(deviceInstanceId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(deviceInstanceId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteDeviceFromSiteV3 provides a mock function with given fields: deviceId, siteId
func (_m *MockManage) DeleteDeviceFromSiteV3(deviceId string, siteId string) (*integration.MsxResponse, error) {
	ret := _m.Called(deviceId, siteId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, string) *integration.MsxResponse); ok {
		r0 = rf(deviceId, siteId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(deviceId, siteId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteDeviceV4 provides a mock function with given fields: deviceId, force
func (_m *MockManage) DeleteDeviceV4(deviceId string, force string) (*integration.MsxResponse, error) {
	ret := _m.Called(deviceId, force)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, string) *integration.MsxResponse); ok {
		r0 = rf(deviceId, force)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(deviceId, force)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteManagedDevice provides a mock function with given fields: deviceInstanceId
func (_m *MockManage) DeleteManagedDevice(deviceInstanceId string) (*integration.MsxResponse, error) {
	ret := _m.Called(deviceInstanceId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(deviceInstanceId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(deviceInstanceId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteServiceInstance provides a mock function with given fields: serviceInstanceId
func (_m *MockManage) DeleteServiceInstance(serviceInstanceId string) (*integration.MsxResponse, error) {
	ret := _m.Called(serviceInstanceId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(serviceInstanceId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(serviceInstanceId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteSite provides a mock function with given fields: siteId
func (_m *MockManage) DeleteSite(siteId string) (*integration.MsxResponse, error) {
	ret := _m.Called(siteId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(siteId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(siteId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteSiteV3 provides a mock function with given fields: siteId
func (_m *MockManage) DeleteSiteV3(siteId string) (*integration.MsxResponse, error) {
	ret := _m.Called(siteId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(siteId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(siteId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteSubscription provides a mock function with given fields: subscriptionId
func (_m *MockManage) DeleteSubscription(subscriptionId string) (*integration.MsxResponse, error) {
	ret := _m.Called(subscriptionId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(subscriptionId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(subscriptionId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAdminHealth provides a mock function with given fields:
func (_m *MockManage) GetAdminHealth() (*HealthResult, error) {
	ret := _m.Called()

	var r0 *HealthResult
	if rf, ok := ret.Get(0).(func() *HealthResult); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*HealthResult)
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

// GetAllControlPlanes provides a mock function with given fields: tenantId
func (_m *MockManage) GetAllControlPlanes(tenantId *string) (*integration.MsxResponse, error) {
	ret := _m.Called(tenantId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(*string) *integration.MsxResponse); ok {
		r0 = rf(tenantId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*string) error); ok {
		r1 = rf(tenantId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetControlPlane provides a mock function with given fields: controlPlaneId
func (_m *MockManage) GetControlPlane(controlPlaneId string) (*integration.MsxResponse, error) {
	ret := _m.Called(controlPlaneId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(controlPlaneId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(controlPlaneId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDevice provides a mock function with given fields: deviceInstanceId
func (_m *MockManage) GetDevice(deviceInstanceId string) (*integration.MsxResponse, error) {
	ret := _m.Called(deviceInstanceId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(deviceInstanceId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(deviceInstanceId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDeviceConfig provides a mock function with given fields: deviceInstanceId
func (_m *MockManage) GetDeviceConfig(deviceInstanceId string) (*integration.MsxResponse, error) {
	ret := _m.Called(deviceInstanceId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(deviceInstanceId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(deviceInstanceId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDeviceTemplateHistory provides a mock function with given fields: deviceInstanceId
func (_m *MockManage) GetDeviceTemplateHistory(deviceInstanceId string) (*integration.MsxResponse, error) {
	ret := _m.Called(deviceInstanceId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(deviceInstanceId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(deviceInstanceId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDeviceV4 provides a mock function with given fields: deviceId
func (_m *MockManage) GetDeviceV4(deviceId string) (*integration.MsxResponse, error) {
	ret := _m.Called(deviceId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(deviceId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(deviceId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDevices provides a mock function with given fields: deviceInstanceId, subscriptionId, serialKey, tenantId, page, pageSize
func (_m *MockManage) GetDevices(deviceInstanceId *string, subscriptionId *string, serialKey *string, tenantId *string, page int, pageSize int) (*integration.MsxResponse, error) {
	ret := _m.Called(deviceInstanceId, subscriptionId, serialKey, tenantId, page, pageSize)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(*string, *string, *string, *string, int, int) *integration.MsxResponse); ok {
		r0 = rf(deviceInstanceId, subscriptionId, serialKey, tenantId, page, pageSize)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*string, *string, *string, *string, int, int) error); ok {
		r1 = rf(deviceInstanceId, subscriptionId, serialKey, tenantId, page, pageSize)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDevicesV4 provides a mock function with given fields: requestQuery, page, pageSize
func (_m *MockManage) GetDevicesV4(requestQuery map[string][]string, page int, pageSize int) (*integration.MsxResponse, error) {
	ret := _m.Called(requestQuery, page, pageSize)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(map[string][]string, int, int) *integration.MsxResponse); ok {
		r0 = rf(requestQuery, page, pageSize)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(map[string][]string, int, int) error); ok {
		r1 = rf(requestQuery, page, pageSize)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEntityShard provides a mock function with given fields: entityId
func (_m *MockManage) GetEntityShard(entityId string) (*integration.MsxResponse, error) {
	ret := _m.Called(entityId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(entityId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(entityId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetServiceInstance provides a mock function with given fields: serviceInstanceId
func (_m *MockManage) GetServiceInstance(serviceInstanceId string) (*integration.MsxResponse, error) {
	ret := _m.Called(serviceInstanceId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(serviceInstanceId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(serviceInstanceId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSite provides a mock function with given fields: siteId
func (_m *MockManage) GetSite(siteId string) (*integration.MsxResponse, error) {
	ret := _m.Called(siteId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(siteId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(siteId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSiteV3 provides a mock function with given fields: siteId, showImage
func (_m *MockManage) GetSiteV3(siteId string, showImage string) (*integration.MsxResponse, error) {
	ret := _m.Called(siteId, showImage)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, string) *integration.MsxResponse); ok {
		r0 = rf(siteId, showImage)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(siteId, showImage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSitesV3 provides a mock function with given fields: siteFilters, page, pageSize
func (_m *MockManage) GetSitesV3(siteFilters SiteQueryFilter, page int, pageSize int) (*integration.MsxResponse, error) {
	ret := _m.Called(siteFilters, page, pageSize)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(SiteQueryFilter, int, int) *integration.MsxResponse); ok {
		r0 = rf(siteFilters, page, pageSize)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(SiteQueryFilter, int, int) error); ok {
		r1 = rf(siteFilters, page, pageSize)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSubscription provides a mock function with given fields: subscriptionId
func (_m *MockManage) GetSubscription(subscriptionId string) (*integration.MsxResponse, error) {
	ret := _m.Called(subscriptionId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string) *integration.MsxResponse); ok {
		r0 = rf(subscriptionId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(subscriptionId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSubscriptionServiceInstances provides a mock function with given fields: subscriptionId, page, pageSize
func (_m *MockManage) GetSubscriptionServiceInstances(subscriptionId string, page int, pageSize int) (*integration.MsxResponse, error) {
	ret := _m.Called(subscriptionId, page, pageSize)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, int, int) *integration.MsxResponse); ok {
		r0 = rf(subscriptionId, page, pageSize)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int, int) error); ok {
		r1 = rf(subscriptionId, page, pageSize)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSubscriptionsV3 provides a mock function with given fields: serviceType, page, pageSize
func (_m *MockManage) GetSubscriptionsV3(serviceType string, page int, pageSize int) (*integration.MsxResponse, error) {
	ret := _m.Called(serviceType, page, pageSize)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, int, int) *integration.MsxResponse); ok {
		r0 = rf(serviceType, page, pageSize)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int, int) error); ok {
		r1 = rf(serviceType, page, pageSize)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateControlPlane provides a mock function with given fields: controlPlaneId, tenantId, name, url, resourceProvider, authenticationType, tlsInsecure, attributes
func (_m *MockManage) UpdateControlPlane(controlPlaneId string, tenantId string, name string, url string, resourceProvider string, authenticationType string, tlsInsecure bool, attributes map[string]string) (*integration.MsxResponse, error) {
	ret := _m.Called(controlPlaneId, tenantId, name, url, resourceProvider, authenticationType, tlsInsecure, attributes)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, string, string, string, string, string, bool, map[string]string) *integration.MsxResponse); ok {
		r0 = rf(controlPlaneId, tenantId, name, url, resourceProvider, authenticationType, tlsInsecure, attributes)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, string, string, string, bool, map[string]string) error); ok {
		r1 = rf(controlPlaneId, tenantId, name, url, resourceProvider, authenticationType, tlsInsecure, attributes)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateDevice provides a mock function with given fields: deviceInstanceId, deviceAttribute, deviceDefAttribute, status
func (_m *MockManage) UpdateDevice(deviceInstanceId string, deviceAttribute map[string]string, deviceDefAttribute map[string]string, status map[string]string) (*integration.MsxResponse, error) {
	ret := _m.Called(deviceInstanceId, deviceAttribute, deviceDefAttribute, status)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, map[string]string, map[string]string, map[string]string) *integration.MsxResponse); ok {
		r0 = rf(deviceInstanceId, deviceAttribute, deviceDefAttribute, status)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, map[string]string, map[string]string, map[string]string) error); ok {
		r1 = rf(deviceInstanceId, deviceAttribute, deviceDefAttribute, status)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateDeviceStatusV4 provides a mock function with given fields: deviceStatus, deviceId
func (_m *MockManage) UpdateDeviceStatusV4(deviceStatus DeviceStatusUpdateRequest, deviceId string) (*integration.MsxResponse, error) {
	ret := _m.Called(deviceStatus, deviceId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(DeviceStatusUpdateRequest, string) *integration.MsxResponse); ok {
		r0 = rf(deviceStatus, deviceId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(DeviceStatusUpdateRequest, string) error); ok {
		r1 = rf(deviceStatus, deviceId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateServiceInstance provides a mock function with given fields: serviceInstanceId, serviceAttribute, serviceDefAttribute, status
func (_m *MockManage) UpdateServiceInstance(serviceInstanceId string, serviceAttribute map[string]string, serviceDefAttribute map[string]string, status map[string]string) (*integration.MsxResponse, error) {
	ret := _m.Called(serviceInstanceId, serviceAttribute, serviceDefAttribute, status)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, map[string]string, map[string]string, map[string]string) *integration.MsxResponse); ok {
		r0 = rf(serviceInstanceId, serviceAttribute, serviceDefAttribute, status)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, map[string]string, map[string]string, map[string]string) error); ok {
		r1 = rf(serviceInstanceId, serviceAttribute, serviceDefAttribute, status)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateSite provides a mock function with given fields: siteId, siteType, displayName, siteAttributes, siteDefAttributes, devices
func (_m *MockManage) UpdateSite(siteId string, siteType *string, displayName *string, siteAttributes map[string]string, siteDefAttributes map[string]string, devices []string) (*integration.MsxResponse, error) {
	ret := _m.Called(siteId, siteType, displayName, siteAttributes, siteDefAttributes, devices)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, *string, *string, map[string]string, map[string]string, []string) *integration.MsxResponse); ok {
		r0 = rf(siteId, siteType, displayName, siteAttributes, siteDefAttributes, devices)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *string, *string, map[string]string, map[string]string, []string) error); ok {
		r1 = rf(siteId, siteType, displayName, siteAttributes, siteDefAttributes, devices)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateSiteStatusV3 provides a mock function with given fields: siteStatus, siteId
func (_m *MockManage) UpdateSiteStatusV3(siteStatus SiteStatusUpdateRequest, siteId string) (*integration.MsxResponse, error) {
	ret := _m.Called(siteStatus, siteId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(SiteStatusUpdateRequest, string) *integration.MsxResponse); ok {
		r0 = rf(siteStatus, siteId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(SiteStatusUpdateRequest, string) error); ok {
		r1 = rf(siteStatus, siteId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateSiteV3 provides a mock function with given fields: siteRequest, siteId, notification
func (_m *MockManage) UpdateSiteV3(siteRequest SiteUpdateRequest, siteId string, notification string) (*integration.MsxResponse, error) {
	ret := _m.Called(siteRequest, siteId, notification)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(SiteUpdateRequest, string, string) *integration.MsxResponse); ok {
		r0 = rf(siteRequest, siteId, notification)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(SiteUpdateRequest, string, string) error); ok {
		r1 = rf(siteRequest, siteId, notification)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateSubscription provides a mock function with given fields: subscriptionId, serviceType, subscriptionName, subscriptionAttribute, offerDefAttribute, offerSelectionDetail, costAttribute
func (_m *MockManage) UpdateSubscription(subscriptionId string, serviceType string, subscriptionName *string, subscriptionAttribute map[string]string, offerDefAttribute map[string]string, offerSelectionDetail map[string]string, costAttribute map[string]string) (*integration.MsxResponse, error) {
	ret := _m.Called(subscriptionId, serviceType, subscriptionName, subscriptionAttribute, offerDefAttribute, offerSelectionDetail, costAttribute)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, string, *string, map[string]string, map[string]string, map[string]string, map[string]string) *integration.MsxResponse); ok {
		r0 = rf(subscriptionId, serviceType, subscriptionName, subscriptionAttribute, offerDefAttribute, offerSelectionDetail, costAttribute)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, *string, map[string]string, map[string]string, map[string]string, map[string]string) error); ok {
		r1 = rf(subscriptionId, serviceType, subscriptionName, subscriptionAttribute, offerDefAttribute, offerSelectionDetail, costAttribute)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateTemplateAccess provides a mock function with given fields: templateId, deviceTemplateDTO
func (_m *MockManage) UpdateTemplateAccess(templateId string, deviceTemplateDTO DeviceTemplateAccessDTO) (*integration.MsxResponse, error) {
	ret := _m.Called(templateId, deviceTemplateDTO)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(string, DeviceTemplateAccessDTO) *integration.MsxResponse); ok {
		r0 = rf(templateId, deviceTemplateDTO)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, DeviceTemplateAccessDTO) error); ok {
		r1 = rf(templateId, deviceTemplateDTO)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
