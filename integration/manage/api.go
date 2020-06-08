package manage

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
)

type Api interface {
	GetAdminHealth() (*HealthResult, error)

	GetSubscription(subscriptionId string) (*integration.MsxResponse, error)
	CreateSubscription(tenantId, serviceType string, subscriptionName *string,
		subscriptionAttribute, offerDefAttribute, offerSelectionDetail, costAttribute map[string]string) (*integration.MsxResponse, error)
	UpdateSubscription(subscriptionId, serviceType string, subscriptionName *string,
		subscriptionAttribute, offerDefAttribute, offerSelectionDetail, costAttribute map[string]string) (*integration.MsxResponse, error)
	DeleteSubscription(subscriptionId string) (*integration.MsxResponse, error)

	// CreateServiceOrder
	// UpdateServiceOrder

	GetServiceInstance(serviceInstanceId string) (*integration.MsxResponse, error)
	GetSubscriptionServiceInstances(subscriptionId string, page, pageSize int) (*integration.MsxResponse, error)
	CreateServiceInstance(subscriptionId, serviceInstanceId string, serviceAttribute, serviceDefAttribute, status map[string]string) (*integration.MsxResponse, error)
	UpdateServiceInstance(serviceInstanceId string, serviceAttribute, serviceDefAttribute, status map[string]string) (*integration.MsxResponse, error)
	DeleteServiceInstance(serviceInstanceId string) (*integration.MsxResponse, error)

	GetSite(siteId string) (*integration.MsxResponse, error)
	CreateSite(subscriptionId, serviceInstanceId string, siteId, siteName, siteType, displayName *string, siteAttributes, siteDefAttributes map[string]string, devices []string) (*integration.MsxResponse, error)
	UpdateSite(siteId string, siteType, displayName *string, siteAttributes, siteDefAttributes map[string]string, devices []string) (*integration.MsxResponse, error)
	DeleteSite(siteId string) (*integration.MsxResponse, error)

	//Deprecated: User v4 Endpoint Instead
	CreateManagedDevice(tenantId string, deviceModel, deviceOnboardType string, deviceOnboardInfo map[string]string) (*integration.MsxResponse, error)
	//Deprecated: User v4 Endpoint Instead
	DeleteManagedDevice(deviceInstanceId string) (*integration.MsxResponse, error)
	GetDeviceConfig(deviceInstanceId string) (*integration.MsxResponse, error)
	//Deprecated: User v4 Endpoint Instead
	GetDevice(deviceInstanceId string) (*integration.MsxResponse, error)
	//Deprecated: User v4 Endpoint Instead
	GetDevices(deviceInstanceId, subscriptionId, serialKey, tenantId *string, page, pageSize int) (*integration.MsxResponse, error)
	//Deprecated: User v4 Endpoint Instead
	CreateDevice(subscriptionId string, deviceInstanceId *string, deviceAttribute, deviceDefAttribute, status map[string]string) (*integration.MsxResponse, error)
	//Deprecated: User v4 Endpoint Instead
	UpdateDevice(deviceInstanceId string, deviceAttribute, deviceDefAttribute, status map[string]string) (*integration.MsxResponse, error)
	//Deprecated: User v4 Endpoint Instead
	DeleteDevice(deviceInstanceId string) (*integration.MsxResponse, error)


	CreateDeviceV4(deviceRequest DeviceCreateRequest) (*integration.MsxResponse, error)
	DeleteDeviceV4(deviceId string, force string) (*integration.MsxResponse, error)
	GetDevicesV4(requestQuery map[string][]string,  page int, pageSize int) (*integration.MsxResponse, error)
	GetDeviceV4(deviceId string) (*integration.MsxResponse, error)
	UpdateDeviceStatusV4(deviceStatus DeviceStatusUpdateRequest, deviceId string) (*integration.MsxResponse, error)

	GetDeviceTemplateHistory(deviceInstanceId string) (*integration.MsxResponse, error)
	// AttachDeviceTemplates
	// UpdateDeviceTemplates
	// DetachDeviceTemplates
	// DetachDeviceTemplate

	GetAllControlPlanes(tenantId *string) (*integration.MsxResponse, error)
	CreateControlPlane(tenantId, name, url, resourceProvider, authenticationType string, tlsInsecure bool, attributes map[string]string) (*integration.MsxResponse, error)
	GetControlPlane(controlPlaneId string) (*integration.MsxResponse, error)
	UpdateControlPlane(controlPlaneId, tenantId, name, url, resourceProvider, authenticationType string, tlsInsecure bool, attributes map[string]string) (*integration.MsxResponse, error)
	DeleteControlPlane(controlPlaneId string) (*integration.MsxResponse, error)
	ConnectControlPlane(controlPlaneId string) (*integration.MsxResponse, error)
	ConnectUnmanagedControlPlane(username, password, url, resourceProvider string, tlsInsecure bool) (*integration.MsxResponse, error)

	GetEntityShard(entityId string) (*integration.MsxResponse, error)
}
