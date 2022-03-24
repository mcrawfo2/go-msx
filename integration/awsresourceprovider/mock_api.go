// Code generated by mockery v2.3.0. DO NOT EDIT.

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

// CheckStatus provides a mock function with given fields: applicationId, request
func (_m *MockAwsResourceProvider) CheckStatus(applicationId types.UUID, request *CheckStatusRequest) (*integration.MsxResponse, error) {
	ret := _m.Called(applicationId, request)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(types.UUID, *CheckStatusRequest) *integration.MsxResponse); ok {
		r0 = rf(applicationId, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.UUID, *CheckStatusRequest) error); ok {
		r1 = rf(applicationId, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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

// GetAmiInformation provides a mock function with given fields: controlPlaneId, amiName, region
func (_m *MockAwsResourceProvider) GetAmiInformation(controlPlaneId types.UUID, amiName string, region string) (*integration.MsxResponse, error) {
	ret := _m.Called(controlPlaneId, amiName, region)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(types.UUID, string, string) *integration.MsxResponse); ok {
		r0 = rf(controlPlaneId, amiName, region)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.UUID, string, string) error); ok {
		r1 = rf(controlPlaneId, amiName, region)
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

// GetInstanceType provides a mock function with given fields: controlPlaneId, region, availabilityZone, instanceType
func (_m *MockAwsResourceProvider) GetInstanceType(controlPlaneId types.UUID, region string, availabilityZone string, instanceType string) (*integration.MsxResponse, error) {
	ret := _m.Called(controlPlaneId, region, availabilityZone, instanceType)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(types.UUID, string, string, string) *integration.MsxResponse); ok {
		r0 = rf(controlPlaneId, region, availabilityZone, instanceType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.UUID, string, string, string) error); ok {
		r1 = rf(controlPlaneId, region, availabilityZone, instanceType)
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

// GetRegionsV2 provides a mock function with given fields: controlPlaneId, amiName
func (_m *MockAwsResourceProvider) GetRegionsV2(controlPlaneId types.UUID, amiName *string) (*integration.MsxResponse, error) {
	ret := _m.Called(controlPlaneId, amiName)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(types.UUID, *string) *integration.MsxResponse); ok {
		r0 = rf(controlPlaneId, amiName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.UUID, *string) error); ok {
		r1 = rf(controlPlaneId, amiName)
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

// GetRouteTableInformation provides a mock function with given fields: controlPlaneId, region, vpcId
func (_m *MockAwsResourceProvider) GetRouteTableInformation(controlPlaneId types.UUID, region string, vpcId string) (*integration.MsxResponse, error) {
	ret := _m.Called(controlPlaneId, region, vpcId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(types.UUID, string, string) *integration.MsxResponse); ok {
		r0 = rf(controlPlaneId, region, vpcId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.UUID, string, string) error); ok {
		r1 = rf(controlPlaneId, region, vpcId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSecrets provides a mock function with given fields: controlPlaneId, secretName, region
func (_m *MockAwsResourceProvider) GetSecrets(controlPlaneId types.UUID, secretName string, region string) (*integration.MsxResponse, error) {
	ret := _m.Called(controlPlaneId, secretName, region)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(types.UUID, string, string) *integration.MsxResponse); ok {
		r0 = rf(controlPlaneId, secretName, region)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.UUID, string, string) error); ok {
		r1 = rf(controlPlaneId, secretName, region)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStackOutputs provides a mock function with given fields: applicationId
func (_m *MockAwsResourceProvider) GetStackOutputs(applicationId types.UUID) (*integration.MsxResponse, error) {
	ret := _m.Called(applicationId)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(types.UUID) *integration.MsxResponse); ok {
		r0 = rf(applicationId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.UUID) error); ok {
		r1 = rf(applicationId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransitGatewayAttachmentStatus provides a mock function with given fields: controlPlaneId, region, transitGatewayAttachmentIds, resourceIds
func (_m *MockAwsResourceProvider) GetTransitGatewayAttachmentStatus(controlPlaneId types.UUID, region string, transitGatewayAttachmentIds []string, resourceIds []string) (*integration.MsxResponse, error) {
	ret := _m.Called(controlPlaneId, region, transitGatewayAttachmentIds, resourceIds)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(types.UUID, string, []string, []string) *integration.MsxResponse); ok {
		r0 = rf(controlPlaneId, region, transitGatewayAttachmentIds, resourceIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.UUID, string, []string, []string) error); ok {
		r1 = rf(controlPlaneId, region, transitGatewayAttachmentIds, resourceIds)
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

// GetTransitVPCStatus provides a mock function with given fields: controlPlaneId, region, transitVPCIds
func (_m *MockAwsResourceProvider) GetTransitVPCStatus(controlPlaneId types.UUID, region string, transitVPCIds []string) (*integration.MsxResponse, error) {
	ret := _m.Called(controlPlaneId, region, transitVPCIds)

	var r0 *integration.MsxResponse
	if rf, ok := ret.Get(0).(func(types.UUID, string, []string) *integration.MsxResponse); ok {
		r0 = rf(controlPlaneId, region, transitVPCIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integration.MsxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.UUID, string, []string) error); ok {
		r1 = rf(controlPlaneId, region, transitVPCIds)
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
