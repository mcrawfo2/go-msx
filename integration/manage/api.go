// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

//go:generate mockery --inpackage --name=Api --structname=MockManage

package manage

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

type Api interface {
	GetAdminHealth() (*integration.MsxResponse, error)

	GetSubscription(subscriptionId string) (*integration.MsxResponse, error)
	GetSubscriptionsV3(serviceType string, page, pageSize int) (*integration.MsxResponse, error)
	CreateSubscription(tenantId, serviceType string, subscriptionName *string,
		subscriptionAttribute, offerDefAttribute, offerSelectionDetail, costAttribute map[string]string) (*integration.MsxResponse, error)
	UpdateSubscription(subscriptionId, serviceType string, subscriptionName *string,
		subscriptionAttribute, offerDefAttribute, offerSelectionDetail, costAttribute map[string]string) (*integration.MsxResponse, error)
	DeleteSubscription(subscriptionId string) (*integration.MsxResponse, error)

	// CreateServiceOrder
	// UpdateServiceOrder

	// Deprecated: Use v8/ServicesApiService (or newer)
	GetServiceInstance(serviceInstanceId string) (*integration.MsxResponse, error)
	// Deprecated: REST API was deprecated in 3.10.0
	GetSubscriptionServiceInstances(subscriptionId string, page, pageSize int) (*integration.MsxResponse, error)
	// Deprecated: Use v8/ServicesApiService (or newer)
	CreateServiceInstance(subscriptionId, serviceInstanceId string, serviceAttribute, serviceDefAttribute, status map[string]string) (*integration.MsxResponse, error)
	// Deprecated: Use v8/ServicesApiService (or newer)
	UpdateServiceInstance(serviceInstanceId string, serviceAttribute, serviceDefAttribute, status map[string]string) (*integration.MsxResponse, error)
	// Deprecated: Use v8/ServicesApiService (or newer)
	DeleteServiceInstance(serviceInstanceId string) (*integration.MsxResponse, error)

	// Deprecated: Use v8/SiteApiService (or newer)
	GetSitesV3(siteFilters SiteQueryFilter, page, pageSize int) (*integration.MsxResponse, error)
	// Deprecated: Use v8/SiteApiService (or newer)
	GetSiteV3(siteId string, showImage string) (*integration.MsxResponse, error)
	// Deprecated: Use v8/SiteApiService (or newer)
	CreateSiteV3(siteRequest SiteCreateRequest) (*integration.MsxResponse, error)
	// Deprecated: Use v8/SiteApiService (or newer)
	UpdateSiteV3(siteRequest SiteUpdateRequest, siteId string, notification string) (*integration.MsxResponse, error)
	// Deprecated: Use v8/SiteApiService (or newer)
	DeleteSiteV3(siteId string) (*integration.MsxResponse, error)
	// Deprecated: Use v8/SiteApiService (or newer)
	AddDeviceToSiteV3(deviceId string, siteId string, notification string) (*integration.MsxResponse, error)
	// Deprecated: Use v8/SiteApiService (or newer)
	DeleteDeviceFromSiteV3(deviceId string, siteId string) (*integration.MsxResponse, error)
	// Deprecated: No replacement planned (site status will be calculated only)
	UpdateSiteStatusV3(siteStatus SiteStatusUpdateRequest, siteId string) (*integration.MsxResponse, error)

	// Deprecated: Use v8/DevicesApiService (or newer)
	GetDeviceConfig(deviceInstanceId string) (*integration.MsxResponse, error)

	// Deprecated: Use v8/DevicesApiService (or newer)
	CreateDeviceV4(deviceRequest DeviceCreateRequest) (*integration.MsxResponse, error)
	// Deprecated: Use v8/DevicesApiService (or newer)
	DeleteDeviceV4(deviceId string, force string) (*integration.MsxResponse, error)
	// Deprecated: Use v8/DevicesApiService (or newer)
	GetDevicesV4(requestQuery map[string][]string, page int, pageSize int) (*integration.MsxResponse, error)
	// Deprecated: Use v8/DevicesApiService (or newer)
	GetDeviceV4(deviceId string) (*integration.MsxResponse, error)
	// Deprecated: Use v8/DevicesApiService (or newer)
	UpdateDeviceV4(deviceRequest DeviceUpdateRequest, deviceId string) (*integration.MsxResponse, error)
	UpdateDeviceStatusV4(deviceStatus DeviceStatusUpdateRequest, deviceId string) (*integration.MsxResponse, error)

	ListDeviceTemplates(serviceType string, tenantId *types.UUID) (*integration.MsxResponse, error)
	GetDeviceTemplate(templateId types.UUID) (*integration.MsxResponse, error)
	AddDeviceTemplate(deviceTemplateCreateRequest DeviceTemplateCreateRequest) (*integration.MsxResponse, error)
	DeleteDeviceTemplate(templateId types.UUID) (*integration.MsxResponse, error)

	// Deprecated: Use v8/DevicesApiService (or newer)
	GetDeviceTemplateHistory(deviceInstanceId string) (*integration.MsxResponse, error)
	// Deprecated: Use v8/DevicesApiService (or newer)
	AttachDeviceTemplates(deviceId string, attachTemplateRequest AttachTemplateRequest) (*integration.MsxResponse, error)
	// Deprecated: Use v8/DevicesApiService (or newer)
	DetachDeviceTemplates(deviceId string) (*integration.MsxResponse, error)
	// Deprecated: Use v8/DevicesApiService (or newer)
	DetachDeviceTemplate(deviceId string, templateId types.UUID) (*integration.MsxResponse, error)

	UpdateTemplateAccess(templateId string, deviceTemplateAccess DeviceTemplateAccess) (*integration.MsxResponse, error)

	CreateDeviceActions(deviceActionList DeviceActionCreateRequests) (*integration.MsxResponse, error)
	UpdateDeviceActions(deviceActionList DeviceActionCreateRequests) (*integration.MsxResponse, error)

	GetAllControlPlanes(tenantId *string) (*integration.MsxResponse, error)
	CreateControlPlane(tenantId, name, url, resourceProvider, authenticationType string, tlsInsecure bool, attributes map[string]string) (*integration.MsxResponse, error)
	GetControlPlane(controlPlaneId string) (*integration.MsxResponse, error)
	UpdateControlPlane(controlPlaneId, tenantId, name, url, resourceProvider, authenticationType string, tlsInsecure bool, attributes map[string]string) (*integration.MsxResponse, error)
	DeleteControlPlane(controlPlaneId string) (*integration.MsxResponse, error)
	ConnectControlPlane(controlPlaneId string) (*integration.MsxResponse, error)
	ConnectUnmanagedControlPlane(username, password, url, resourceProvider string, tlsInsecure bool) (*integration.MsxResponse, error)

	CreateDeviceConnection(deviceConnection DeviceConnectionCreateRequest) (*integration.MsxResponse, *DeviceConnectionResponse, error)
	DeleteDeviceConnection(deviceConnectionId string) (*integration.MsxResponse, error)

	GetEntityShard(entityId string) (*integration.MsxResponse, error)

	GetLocationGeocode(location string) (*integration.MsxResponse, error)
}
