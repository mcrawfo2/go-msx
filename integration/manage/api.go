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

	CreateManagedDevice(tenantId string, deviceModel, deviceOnboardType string, deviceOnboardInfo map[string]string) (*integration.MsxResponse, error)
	DeleteManagedDevice(deviceInstanceId string) (*integration.MsxResponse, error)
	GetDeviceConfig(deviceInstanceId string) (*integration.MsxResponse, error)

	GetDevice(deviceInstanceId string) (*integration.MsxResponse, error)
	GetDevices(deviceInstanceId, subscriptionId, serialKey, tenantId *string, page, pageSize int) (*integration.MsxResponse, error)
	CreateDevice(subscriptionId string, deviceInstanceId *string, deviceAttribute, deviceDefAttribute, status map[string]string) (*integration.MsxResponse, error)
	UpdateDevice(deviceInstanceId string, deviceAttribute, deviceDefAttribute, status map[string]string) (*integration.MsxResponse, error)
	DeleteDevice(deviceInstanceId string) (*integration.MsxResponse, error)

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
