// Code generated by mockery v2.1.0. DO NOT EDIT.

package awsresourceprovider

import (
	integration "cto-github.cisco.com/NFV-BU/go-msx/integration"
	mock "github.com/stretchr/testify/mock"

	types "cto-github.cisco.com/NFV-BU/go-msx/types"
)

// MockAwsResourceProvider is an autogenerated mock type for the Api type
type MockAwsResourceProvider struct {
	mock.Mock
}

// Connect provides a mock function with given fields: request
func (_m *MockAwsResourceProvider) Connect(request AwsConnectRequest) (*integration.MsxResponse, error) {
	ret := _m.Called(request)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(AwsConnectRequest) *integration.MsxResponse); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(AwsConnectRequest) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAvailabilityZones provides a mock function with given fields: controlPlaneId, region
func (_m *MockAwsResourceProvider) GetAvailabilityZones(controlPlaneId types.UUID, region string) (*integration.MsxResponse, error) {
	ret := _m.Called(controlPlaneId, region)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(types.UUID, string) *integration.MsxResponse); ok {
		r0 = rf(controlPlaneId, region)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.UUID, string) error); ok {
		r1 = rf(controlPlaneId, region)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEc2InstanceStatus provides a mock function with given fields: controlPlaneId, region, instanceId
func (_m *MockAwsResourceProvider) GetEc2InstanceStatus(controlPlaneId types.UUID, region string, instanceId string) (*integration.MsxResponse, error) {
	ret := _m.Called(controlPlaneId, region, instanceId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(types.UUID, string, string) *integration.MsxResponse); ok {
		r0 = rf(controlPlaneId, region, instanceId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.UUID, string, string) error); ok {
		r1 = rf(controlPlaneId, region, instanceId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRegions provides a mock function with given fields: controlPlaneId
func (_m *MockAwsResourceProvider) GetRegions(controlPlaneId types.UUID) (*integration.MsxResponse, error) {
	ret := _m.Called(controlPlaneId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(types.UUID) *integration.MsxResponse); ok {
		r0 = rf(controlPlaneId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.UUID) error); ok {
		r1 = rf(controlPlaneId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetResources provides a mock function with given fields: serviceConfigurationApplicationId
func (_m *MockAwsResourceProvider) GetResources(serviceConfigurationApplicationId types.UUID) (*integration.MsxResponse, error) {
	ret := _m.Called(serviceConfigurationApplicationId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(types.UUID) *integration.MsxResponse); ok {
		r0 = rf(serviceConfigurationApplicationId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.UUID) error); ok {
		r1 = rf(serviceConfigurationApplicationId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransitGatewayAttachmentStatus provides a mock function with given fields: controlPlaneId, region, transitGatewayAttachmentIds
func (_m *MockAwsResourceProvider) GetTransitGatewayAttachmentStatus(controlPlaneId types.UUID, region string, transitGatewayAttachmentIds []string) (*integration.MsxResponse, error) {
	ret := _m.Called(controlPlaneId, region, transitGatewayAttachmentIds)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(types.UUID, string, []string) *integration.MsxResponse); ok {
		r0 = rf(controlPlaneId, region, transitGatewayAttachmentIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.UUID, string, []string) error); ok {
		r1 = rf(controlPlaneId, region, transitGatewayAttachmentIds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransitGatewayStatus provides a mock function with given fields: controlPlaneId, region, transitGatewayIds
func (_m *MockAwsResourceProvider) GetTransitGatewayStatus(controlPlaneId types.UUID, region string, transitGatewayIds []string) (*integration.MsxResponse, error) {
	ret := _m.Called(controlPlaneId, region, transitGatewayIds)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(types.UUID, string, []string) *integration.MsxResponse); ok {
		r0 = rf(controlPlaneId, region, transitGatewayIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.UUID, string, []string) error); ok {
		r1 = rf(controlPlaneId, region, transitGatewayIds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVpnConnectionDetails provides a mock function with given fields: controlPlaneId, vpnConnectionIds, region
func (_m *MockAwsResourceProvider) GetVpnConnectionDetails(controlPlaneId types.UUID, vpnConnectionIds []string, region string) (*integration.MsxResponse, error) {
	ret := _m.Called(controlPlaneId, vpnConnectionIds, region)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(types.UUID, []string, string) *integration.MsxResponse); ok {
		r0 = rf(controlPlaneId, vpnConnectionIds, region)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.UUID, []string, string) error); ok {
		r1 = rf(controlPlaneId, vpnConnectionIds, region)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
