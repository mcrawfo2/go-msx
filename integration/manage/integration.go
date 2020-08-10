package manage

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

const (
	endpointNameGetAdminHealth = "getAdminHealth"

	endpointNameGetSubscription    = "getSubscription"
	endpointNameGetSubscriptionsV3 = "getSubscriptionsV3"
	endpointNameCreateSubscription = "createSubscription"
	endpointNameUpdateSubscription = "updateSubscription"
	endpointNameDeleteSubscription = "deleteSubscription"

	endpointNameGetServiceInstance              = "getServiceInstance"
	endpointNameGetSubscriptionServiceInstances = "getSubscriptionServiceInstances"
	endpointNameCreateServiceInstance           = "createServiceInstance"
	endpointNameUpdateServiceInstance           = "updateServiceInstance"
	endpointNameDeleteServiceInstance           = "deleteServiceInstance"

	endpointNameGetSite    = "getSite"
	endpointNameCreateSite = "createSite"
	endpointNameUpdateSite = "updateSite"
	endpointNameDeleteSite = "deleteSite"

	endpointNameCreateSiteV3           = "createSiteV3"
	endpointNameGetSitesV3             = "getSitesV3"
	endpointNameGetSiteV3              = "getSiteV3"
	endpointNameUpdateSiteV3           = "updateSiteV3"
	endpointNameDeleteSiteV3           = "deleteSiteV3"
	endpointNameAddDevicetoSiteV3      = "addDeviceToSiteV3"
	endpointNameDeleteDeviceFromSiteV3 = "deleteDeviceFromSiteV3"
	endpointNameUpdateSiteStatusV3     = "updateSiteStatusV3"

	endpointNameGetDevice    = "getDevice"
	endpointNameGetDevices   = "getDevices"
	endpointNameCreateDevice = "createDevice"
	endpointNameUpdateDevice = "updateDevice"
	endpointNameDeleteDevice = "deleteDevice"

	endpointNameCreateManagedDevice = "createManagedDevice"
	endpointNameDeleteManagedDevice = "deleteManagedDevice"
	endpointNameGetDeviceConfig     = "getDeviceConfig"

	endpointNameGetDevicesV4         = "getDevicesV4"
	endpointNameGetDeviceV4          = "getDeviceV4"
	endpointNameCreateDeviceV4       = "createDeviceV4"
	endpointNameDeleteDeviceV4       = "deleteDeviceV4"
	endpointNameUpdateDeviceStatusV4 = "updateDeviceStatusV4"

	endpointNameGetDeviceTemplateHistory = "getDeviceTemplateHistory"
	endpointNameAttachDeviceTemplates    = "attachDeviceTemplates"
	endpointNameUpdateDeviceTemplates    = "updateDeviceTemplates"
	endpointNameDetachDeviceTemplates    = "detachDeviceTemplates"
	endpointNameDetachDeviceTemplate     = "detachDeviceTemplate"

	endpointNameGetAllControlPlanes          = "getAllControlPlanes"
	endpointNameCreateControlPlane           = "createControlPlane"
	endpointNameGetControlPlane              = "getControlPlane"
	endpointNameUpdateControlPlane           = "updateControlPlane"
	endpointNameDeleteControlPlane           = "deleteControlPlane"
	endpointNameConnectControlPlane          = "connectControlPlane"
	endpointNameConnectUnmanagedControlPlane = "connectUnmanagedControlPlane"

	endpointNameGetEntityShard = "getEntityShard"

	serviceName = integration.ServiceNameManage
)

var (
	logger    = log.NewLogger("msx.integration.manage")
	endpoints = map[string]integration.MsxServiceEndpoint{
		endpointNameGetAdminHealth: {Method: "GET", Path: "/admin/health"},

		endpointNameGetSubscription:    {Method: "GET", Path: "/api/v2/subscriptions/{{.subscriptionId}}"},
		endpointNameGetSubscriptionsV3: {Method: "GET", Path: "/api/v3/subscriptions"},

		endpointNameCreateSubscription: {Method: "POST", Path: "/api/v2/subscriptions/tenants/{{.tenantId}}"},
		endpointNameUpdateSubscription: {Method: "PUT", Path: "/api/v2/subscriptions/{{.subscriptionId}}"},
		endpointNameDeleteSubscription: {Method: "DELETE", Path: "/api/v2/subscriptions/{{.subscriptionId}}"},

		endpointNameGetServiceInstance:              {Method: "GET", Path: "/api/v2/serviceinstances/{{.serviceInstanceId}}"},
		endpointNameGetSubscriptionServiceInstances: {Method: "GET", Path: "/api/v2/serviceinstances/subscriptions/{{.subscriptionId}}"},
		endpointNameCreateServiceInstance:           {Method: "POST", Path: "/api/v1/serviceinstances/subscriptions/{{.subscriptionId}}"},
		endpointNameUpdateServiceInstance:           {Method: "PUT", Path: "/api/v1/serviceinstances/{{.serviceInstanceId}}"},
		endpointNameDeleteServiceInstance:           {Method: "DELETE", Path: "/api/v1/serviceinstances/{{.serviceInstanceId}}"},

		endpointNameGetSite:    {Method: "GET", Path: "/api/v1/sites/{{.siteId}}"},
		endpointNameCreateSite: {Method: "POST", Path: "/api/v1/sites/subscriptions/{{.subscriptionId}}"},
		endpointNameUpdateSite: {Method: "PUT", Path: "/api/v1/sites/{{.siteId}}"},
		endpointNameDeleteSite: {Method: "DELETE", Path: "/api/v1/site/{{.siteId}}"},

		endpointNameCreateSiteV3:           {Method: "POST", Path: "/api/v3/sites"},
		endpointNameGetSitesV3:             {Method: "GET", Path: "/api/v3/sites"},
		endpointNameGetSiteV3:              {Method: "GET", Path: "/api/v3/sites/{{.siteId}}"},
		endpointNameUpdateSiteV3:           {Method: "PUT", Path: "/api/v3/sites/{{.siteId}}"},
		endpointNameDeleteSiteV3:           {Method: "DELETE", Path: "/api/v3/sites/{{.siteId}}"},
		endpointNameAddDevicetoSiteV3:      {Method: "PUT", Path: "/api/v3/sites/{{.siteId}}/devices/{{.deviceId}}"},
		endpointNameDeleteDeviceFromSiteV3: {Method: "DELETE", Path: "/api/v3/sites/{{.siteId}}/devices/{{.deviceId}}"},
		endpointNameUpdateSiteStatusV3:     {Method: "PUT", Path: "/api/v3/sites/{{.siteId}}/status"},

		endpointNameGetDevice:           {Method: "GET", Path: "/api/v1/devices/{{.deviceInstanceId}}"},
		endpointNameGetDevices:          {Method: "GET", Path: "/api/v2/devices"},
		endpointNameCreateDevice:        {Method: "POST", Path: "/api/v1/devices/subscriptions/{{.subscriptionId}}"},
		endpointNameUpdateDevice:        {Method: "PUT", Path: "/api/v1/devices/{{.deviceInstanceId}}"},
		endpointNameDeleteDevice:        {Method: "DELETE", Path: "/api/v1/devices/{{.deviceInstanceId}}"},
		endpointNameCreateManagedDevice: {Method: "POST", Path: "/api/v3/devices"},
		endpointNameDeleteManagedDevice: {Method: "DELETE", Path: "/api/v3/devices/{{.deviceInstanceId}}"},
		endpointNameGetDeviceConfig:     {Method: "GET", Path: "/api/v3/devices/{{.deviceInstanceId}}/config"},

		endpointNameGetDevicesV4:         {Method: "GET", Path: "/api/v4/devices"},
		endpointNameGetDeviceV4:          {Method: "GET", Path: "/api/v4/devices/{{.deviceId}}"},
		endpointNameCreateDeviceV4:       {Method: "POST", Path: "/api/v4/devices"},
		endpointNameDeleteDeviceV4:       {Method: "DELETE", Path: "/api/v4/devices/{{.deviceId}}"},
		endpointNameUpdateDeviceStatusV4: {Method: "PUT", Path: "/api/v4/devices/{{.deviceId}}/status"},

		endpointNameGetDeviceTemplateHistory: {Method: "GET", Path: "/api/v3/devices/{{.deviceInstanceId}}/templates"},
		endpointNameAttachDeviceTemplates:    {Method: "POST", Path: "/api/v3/devices/{{.deviceInstanceId}}/templates"},
		endpointNameUpdateDeviceTemplates:    {Method: "PUT", Path: "/api/v3/devices/{{.deviceInstanceId}}/templates"},
		endpointNameDetachDeviceTemplates:    {Method: "DELETE", Path: "/api/v3/devices/{{.deviceInstanceId}}/templates"},
		endpointNameDetachDeviceTemplate:     {Method: "DELETE", Path: "/api/v3/devices/{{.deviceInstanceId}}/templates/{{.templateId}}"},

		endpointNameGetAllControlPlanes:          {Method: "GET", Path: "/api/v1/controlplanes"},
		endpointNameCreateControlPlane:           {Method: "POST", Path: "/api/v1/controlplanes"},
		endpointNameGetControlPlane:              {Method: "GET", Path: "/api/v1/controlplanes/{{.controlPlaneId}}"},
		endpointNameUpdateControlPlane:           {Method: "PUT", Path: "/api/v1/controlplanes/{{.controlPlaneId}}"},
		endpointNameDeleteControlPlane:           {Method: "DELETE", Path: "/api/v1/controlplanes/{{.controlPlaneId}}"},
		endpointNameConnectControlPlane:          {Method: "POST", Path: "/api/v1/controlplanes/{{.controlPlaneId}}/connect"},
		endpointNameConnectUnmanagedControlPlane: {Method: "POST", Path: "/api/v1/controlplanes/connect"},

		endpointNameGetEntityShard: {Method: "GET", Path: "/api/v2/shardmanagers/entity/{{.entityId}}"},
	}
)

type Integration struct {
	*integration.MsxService
}

func NewIntegration(ctx context.Context) (Api, error) {
	integrationInstance := IntegrationFromContext(ctx)
	if integrationInstance == nil {
		integrationInstance = &Integration{
			MsxService: integration.NewMsxService(ctx, serviceName, endpoints),
		}
	}
	return integrationInstance, nil
}

func (i *Integration) GetAdminHealth() (result *HealthResult, err error) {
	result = &HealthResult{
		Payload: &integration.HealthDTO{},
	}

	result.Response, err = i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetAdminHealth,
		Payload:      result.Payload,
		NoToken:      true,
	})

	return result, err
}

func (i *Integration) GetSubscription(subscriptionId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetSubscription,
		EndpointParameters: map[string]string{
			"subscriptionId": subscriptionId,
		},
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetSubscriptionsV3(serviceType string, page, pageSize int) (*integration.MsxResponse, error) {
	pageString := strconv.Itoa(page)
	pageSizeString := strconv.Itoa(pageSize)

	searchParameters := map[string]*string{
		"serviceType": &serviceType,
		"page":        &pageString,
		"pageSize":    &pageSizeString,
	}

	// Convert optional search queries into query parameters
	queryParameters := make(url.Values)
	for k, v := range searchParameters {
		if v != nil && *v != "" {
			queryParameters[k] = []string{*v}
		}
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:    endpointNameGetSubscriptionsV3,
		QueryParameters: queryParameters,
		Payload:         new(Pojo),
		ExpectEnvelope:  true,
	})
}

func (i *Integration) CreateSubscription(tenantId, serviceType string, subscriptionName *string,
	subscriptionAttribute, offerDefAttribute, offerSelectionDetail, costAttribute map[string]string) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(&Pojo{
		"serviceType":           serviceType,
		"subscriptionName":      subscriptionName,
		"subscriptionAttribute": subscriptionAttribute,
		"offerDefAttribute":     offerDefAttribute,
		"offerSelectionDetail":  offerSelectionDetail,
		"costAttribute":         costAttribute,
	})

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameCreateSubscription,
		EndpointParameters: map[string]string{
			"tenantId": tenantId,
		},
		Body:           bodyBytes,
		Payload:        new(CreateSubscriptionResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) UpdateSubscription(subscriptionId, serviceType string, subscriptionName *string,
	subscriptionAttribute, offerDefAttribute, offerSelectionDetail, costAttribute map[string]string) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(&Pojo{
		"serviceType":           serviceType,
		"subscriptionName":      subscriptionName,
		"subscriptionAttribute": subscriptionAttribute,
		"offerDefAttribute":     offerDefAttribute,
		"offerSelectionDetail":  offerSelectionDetail,
		"costAttribute":         costAttribute,
	})

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameUpdateSubscription,
		EndpointParameters: map[string]string{
			"subscriptionId": subscriptionId,
		},
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) DeleteSubscription(subscriptionId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameDeleteSubscription,
		EndpointParameters: map[string]string{
			"subscriptionId": subscriptionId,
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetServiceInstance(serviceInstanceId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetServiceInstance,
		EndpointParameters: map[string]string{
			"serviceInstanceId": serviceInstanceId,
		},
		Payload:        new(ServiceInstanceResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetSubscriptionServiceInstances(subscriptionId string, page, pageSize int) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetSubscriptionServiceInstances,
		EndpointParameters: map[string]string{
			"subscriptionId": subscriptionId,
		},
		QueryParameters: map[string][]string{
			"page":     {strconv.Itoa(page)},
			"pageSize": {strconv.Itoa(pageSize)},
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) CreateServiceInstance(subscriptionId, serviceInstanceId string, serviceAttribute, serviceDefAttribute, status map[string]string) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(&Pojo{
		"serviceInstanceId":   serviceInstanceId,
		"serviceAttribute":    serviceAttribute,
		"serviceDefAttribute": serviceDefAttribute,
		"status":              status,
	})

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameCreateServiceInstance,
		EndpointParameters: map[string]string{
			"subscriptionId": subscriptionId,
		},
		Body:           bodyBytes,
		Payload:        new(ServiceInstanceResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) UpdateServiceInstance(serviceInstanceId string, serviceAttribute, serviceDefAttribute, status map[string]string) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(&Pojo{
		"serviceInstanceId":   serviceInstanceId,
		"serviceAttribute":    serviceAttribute,
		"serviceDefAttribute": serviceDefAttribute,
		"status":              status,
	})

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameUpdateServiceInstance,
		EndpointParameters: map[string]string{
			"serviceInstanceId": serviceInstanceId,
		},
		Body:           bodyBytes,
		Payload:        new(ServiceInstanceResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) DeleteServiceInstance(serviceInstanceId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameDeleteServiceInstance,
		EndpointParameters: map[string]string{
			"serviceInstanceId": serviceInstanceId,
		},
		ExpectEnvelope: true,
	})
}

// Deprecated: Use v3 Endpoint Instead
func (i *Integration) GetSite(siteId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetSite,
		EndpointParameters: map[string]string{
			"siteId": siteId,
		},
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

// Deprecated: Use v3 Endpoint Instead
func (i *Integration) CreateSite(subscriptionId, serviceInstanceId string, siteId, siteName, siteType, displayName *string,
	siteAttributes, siteDefAttributes map[string]string, devices []string) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(&Pojo{
		"serviceInstanceId": serviceInstanceId,
		"siteId":            siteId,
		"siteName":          siteName,
		"siteType":          siteType,
		"displayName":       displayName,
		"siteAttributes":    siteAttributes,
		"siteDefAttributes": siteDefAttributes,
		"devices":           devices,
	})

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameCreateSite,
		EndpointParameters: map[string]string{
			"subscriptionId": subscriptionId,
		},
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

// Deprecated: Use v3 Endpoint Instead
func (i *Integration) UpdateSite(siteId string, siteType, displayName *string,
	siteAttributes, siteDefAttributes map[string]string, devices []string) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(&Pojo{
		"siteId":            siteId,
		"siteType":          siteType,
		"displayName":       displayName,
		"siteAttributes":    siteAttributes,
		"siteDefAttributes": siteDefAttributes,
		"devices":           devices,
	})
	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameUpdateSite,
		EndpointParameters: map[string]string{
			"siteId": siteId,
		},
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

// Deprecated: Use v3 Endpoint Instead
func (i *Integration) DeleteSite(siteId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameDeleteSite,
		EndpointParameters: map[string]string{
			"siteId": siteId,
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetSitesV3(siteFilters SiteQueryFilter, page, pageSize int) (*integration.MsxResponse, error) {
	pageString := strconv.Itoa(page)
	pageSizeString := strconv.Itoa(pageSize)

	searchParameters := map[string]*string{
		"deviceInstanceId":  siteFilters.DeviceInstanceId,
		"parentId":          siteFilters.ParentId,
		"serviceInstanceId": siteFilters.ServiceInstanceId,
		"serviceType":       siteFilters.ServiceType,
		"showImage":         siteFilters.ShowImage,
		"tenantId":          siteFilters.TenantId,
		"type":              siteFilters.Type,
		"page":              &pageString,
		"pageSize":          &pageSizeString,
	}

	// Convert optional search queries into query parameters
	queryParameters := make(url.Values)
	for k, v := range searchParameters {
		if v != nil && *v != "" {
			queryParameters[k] = []string{*v}
		}
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:    endpointNameGetSitesV3,
		QueryParameters: queryParameters,
		Payload:         new(Pojo),
		ExpectEnvelope:  true,
	})
}

func (i *Integration) GetSiteV3(siteId string, showImage string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetSiteV3,
		EndpointParameters: map[string]string{
			"siteId": siteId,
		},
		QueryParameters: map[string][]string{
			"showImage": {showImage},
		},
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) CreateSiteV3(siteRequest SiteCreateRequest) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(siteRequest)
	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:   endpointNameCreateSiteV3,
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) UpdateSiteV3(siteRequest SiteUpdateRequest, siteId string, notification string) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(siteRequest)
	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameUpdateSiteV3,
		EndpointParameters: map[string]string{
			"siteId": siteId,
		},
		QueryParameters: map[string][]string{
			"notification": {notification},
		},
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) DeleteSiteV3(siteId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameDeleteSiteV3,
		EndpointParameters: map[string]string{
			"siteId": siteId,
		},
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) AddDeviceToSiteV3(deviceId string, siteId string, notification string) (*integration.MsxResponse, error) {

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameAddDevicetoSiteV3,
		EndpointParameters: map[string]string{
			"siteId":   siteId,
			"deviceId": deviceId,
		},
		QueryParameters: map[string][]string{
			"notification": {notification},
		},
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) DeleteDeviceFromSiteV3(deviceId string, siteId string) (*integration.MsxResponse, error) {

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameDeleteDeviceFromSiteV3,
		EndpointParameters: map[string]string{
			"siteId":   siteId,
			"deviceId": deviceId,
		},
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) UpdateSiteStatusV3(siteStatus SiteStatusUpdateRequest, siteId string) (*integration.MsxResponse, error) {

	bodyBytes, err := json.Marshal(siteStatus)
	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameUpdateSiteStatusV3,
		EndpointParameters: map[string]string{
			"siteId": siteId,
		},
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

// Deprecated: Use v4 Endpoint Instead
func (i *Integration) CreateManagedDevice(tenantId string, deviceModel, deviceOnboardType string, deviceOnboardInfo map[string]string) (*integration.MsxResponse, error) {
	deviceUuid, _ := types.NewUUID()
	deviceInstanceId := "CPE-" + strings.ToLower(deviceUuid.MustMarshalText())

	bodyBytes, err := json.Marshal(&Pojo{
		"tenantId":             tenantId,
		"deviceInstanceId":     deviceInstanceId,
		"deviceName":           deviceInstanceId,
		"deviceModel":          deviceModel,
		"deviceOnboardingType": deviceOnboardType,
		"deviceOnboardInfo":    deviceOnboardInfo,
	})

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:   endpointNameCreateManagedDevice,
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

// Deprecated: Use v4 Endpoint Instead
func (i *Integration) DeleteManagedDevice(deviceInstanceId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameDeleteManagedDevice,
		EndpointParameters: map[string]string{
			"deviceInstanceId": deviceInstanceId,
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetDeviceConfig(deviceInstanceId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetDeviceConfig,
		EndpointParameters: map[string]string{
			"deviceInstanceId": deviceInstanceId,
		},
		ExpectEnvelope: true,
	})
}

// Deprecated: Use v4 Endpoint Instead
func (i *Integration) GetDevice(deviceInstanceId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetDevice,
		EndpointParameters: map[string]string{
			"deviceInstanceId": deviceInstanceId,
		},
		ExpectEnvelope: true,
	})
}

// Deprecated: Use v4 Endpoint Instead
func (i *Integration) GetDevices(deviceInstanceId, subscriptionId, serialKey, tenantId *string, page, pageSize int) (*integration.MsxResponse, error) {
	pageString := strconv.Itoa(page)
	pageSizeString := strconv.Itoa(pageSize)

	searchParameters := map[string]*string{
		"deviceID":       deviceInstanceId,
		"serialKey":      serialKey,
		"subscriptionID": subscriptionId,
		"tenantId":       tenantId,
		"page":           &pageString,
		"pageSize":       &pageSizeString,
	}

	// Convert optional search queries into query parameters
	queryParameters := make(url.Values)
	for k, v := range searchParameters {
		if v != nil && *v != "" {
			queryParameters[k] = []string{*v}
		}
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:    endpointNameGetDevices,
		QueryParameters: queryParameters,
		ExpectEnvelope:  true,
	})
}

// Deprecated: Use v4 Endpoint Instead
func (i *Integration) CreateDevice(subscriptionId string, deviceInstanceId *string, deviceAttribute, deviceDefAttribute, status map[string]string) (*integration.MsxResponse, error) {
	if deviceInstanceId == nil {
		deviceUuid, _ := types.NewUUID()
		deviceInstanceIdString := "CPE-" + strings.ToLower(deviceUuid.MustMarshalText())
		deviceInstanceId = &deviceInstanceIdString
	}

	bodyBytes, err := json.Marshal(&Pojo{
		"deviceInstanceId":   deviceInstanceId,
		"deviceAttribute":    deviceAttribute,
		"deviceDefAttribute": deviceDefAttribute,
		"status":             status,
	})

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameCreateDevice,
		EndpointParameters: map[string]string{
			"subscriptionId": subscriptionId,
		},
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

// Deprecated: Use v4 Endpoint Instead
func (i *Integration) UpdateDevice(deviceInstanceId string, deviceAttribute, deviceDefAttribute, status map[string]string) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(&Pojo{
		"deviceAttribute":    deviceAttribute,
		"deviceDefAttribute": deviceDefAttribute,
		"status":             status,
	})

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameUpdateDevice,
		EndpointParameters: map[string]string{
			"deviceInstanceId": deviceInstanceId,
		},
		Body:           bodyBytes,
		ExpectEnvelope: true,
	})
}

// Deprecated: Use v4 Endpoint Instead
func (i *Integration) DeleteDevice(deviceInstanceId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameDeleteDevice,
		EndpointParameters: map[string]string{
			"deviceInstanceId": deviceInstanceId,
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) CreateDeviceV4(deviceRequest DeviceCreateRequest) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(deviceRequest)
	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:       endpointNameCreateDeviceV4,
		EndpointParameters: map[string]string{},
		Body:               bodyBytes,
		Payload:            new(DeviceResponse),
		ExpectEnvelope:     true,
	})
}

func (i *Integration) DeleteDeviceV4(deviceId string, force string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameDeleteDeviceV4,
		EndpointParameters: map[string]string{
			"deviceId": deviceId,
		},
		QueryParameters: map[string][]string{
			"force": {force},
		},
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetDevicesV4(requestQuery map[string][]string, page, pageSize int) (*integration.MsxResponse, error) {

	// Convert optional search queries into query parameters
	queryParameters := make(url.Values)
	queryParameters["page"] = []string{strconv.Itoa(page)}
	queryParameters["pageSize"] = []string{strconv.Itoa(pageSize)}

	for k, v := range requestQuery {
		if v != nil && len(v) > 0 {
			queryParameters[k] = v
		}
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:    endpointNameGetDevicesV4,
		QueryParameters: queryParameters,
		Payload: paging.PaginatedResponse{
			Content: []DeviceResponse{},
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetDeviceV4(deviceId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetDeviceV4,
		EndpointParameters: map[string]string{
			"deviceId": deviceId,
		},
		Payload:        new(DeviceResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) UpdateDeviceStatusV4(deviceStatus DeviceStatusUpdateRequest, deviceId string) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(&Pojo{
		"message": deviceStatus.Message,
		"type":    deviceStatus.Type,
		"value":   deviceStatus.Value,
	})

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameUpdateDevice,
		EndpointParameters: map[string]string{
			"deviceId": deviceId,
		},
		Body:           bodyBytes,
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetDeviceTemplateHistory(deviceInstanceId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetDeviceTemplateHistory,
		EndpointParameters: map[string]string{
			"deviceInstanceId": deviceInstanceId,
		},
		Payload:        new(PojoArray),
		ExpectEnvelope: true,
	})
}

/*
	TODO: v3 device templates
	endpointNameAttachDeviceTemplates    = "attachDeviceTemplates"
	endpointNameUpdateDeviceTemplates    = "updateDeviceTemplates"
	endpointNameDetachDeviceTemplates    = "detachDeviceTemplates"
	endpointNameDetachDeviceTemplate     = "detachDeviceTemplate"
*/

func (i *Integration) GetAllControlPlanes(tenantId *string) (*integration.MsxResponse, error) {
	queryParameters := url.Values{}
	if tenantId != nil {
		queryParameters["tenantId"] = []string{*tenantId}
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:    endpointNameGetAllControlPlanes,
		QueryParameters: queryParameters,
		Payload:         new(PojoArray),
		ExpectEnvelope:  true,
	})
}

func (i *Integration) CreateControlPlane(tenantId, name, url, resourceProvider, authenticationType string, tlsInsecure bool, attributes map[string]string) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(&Pojo{
		"tenantId":           tenantId,
		"name":               name,
		"url":                url,
		"resourceProvider":   resourceProvider,
		"authenticationType": authenticationType,
		"tlsInsecure":        tlsInsecure,
		"attributes":         attributes,
	})

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:   endpointNameCreateControlPlane,
		Body:           bodyBytes,
		Payload:        new(ControlPlaneResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetControlPlane(controlPlaneId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:       endpointNameGetControlPlane,
		EndpointParameters: map[string]string{"controlPlaneId": controlPlaneId},
		Payload:            new(Pojo),
		ExpectEnvelope:     true,
	})
}

func (i *Integration) UpdateControlPlane(controlPlaneId, tenantId, name, url, resourceProvider, authenticationType string, tlsInsecure bool, attributes map[string]string) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(&Pojo{
		"tenantId":           tenantId,
		"name":               name,
		"url":                url,
		"resourceProvider":   resourceProvider,
		"authenticationType": authenticationType,
		"tlsInsecure":        tlsInsecure,
		"attributes":         attributes,
	})

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:       endpointNameUpdateControlPlane,
		EndpointParameters: map[string]string{"controlPlaneId": controlPlaneId},
		Body:               bodyBytes,
		Payload:            new(Pojo),
		ExpectEnvelope:     true,
	})
}

func (i *Integration) DeleteControlPlane(controlPlaneId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:       endpointNameDeleteControlPlane,
		EndpointParameters: map[string]string{"controlPlaneId": controlPlaneId},
		Payload:            new(Pojo),
		ExpectEnvelope:     true,
	})
}

func (i *Integration) ConnectControlPlane(controlPlaneId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:       endpointNameConnectControlPlane,
		EndpointParameters: map[string]string{"controlPlaneId": controlPlaneId},
		Payload:            new(Pojo),
		ExpectEnvelope:     true,
	})
}

func (i *Integration) ConnectUnmanagedControlPlane(username, password, url, resourceProvider string, tlsInsecure bool) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(&Pojo{
		"username":         username,
		"password":         password,
		"url":              url,
		"resourceProvider": resourceProvider,
		"tlsInsecure":      tlsInsecure,
	})

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:   endpointNameConnectUnmanagedControlPlane,
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetEntityShard(entityId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:       endpointNameGetEntityShard,
		EndpointParameters: map[string]string{"entityId": entityId},
		Payload:            new(EntityShard),
		ExpectEnvelope:     true,
	})
}
