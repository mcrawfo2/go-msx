// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package manage

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
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

	endpointNameCreateSiteV3           = "createSiteV3"
	endpointNameGetSitesV3             = "getSitesV3"
	endpointNameGetSiteV3              = "getSiteV3"
	endpointNameUpdateSiteV3           = "updateSiteV3"
	endpointNameDeleteSiteV3           = "deleteSiteV3"
	endpointNameAddDevicetoSiteV3      = "addDeviceToSiteV3"
	endpointNameDeleteDeviceFromSiteV3 = "deleteDeviceFromSiteV3"
	endpointNameUpdateSiteStatusV3     = "updateSiteStatusV3"

	endpointNameCreateDeviceActions = "createDeviceActions"
	endpointNameUpdateDeviceActions = "updateDeviceActions"

	endpointNameGetDeviceConfig = "getDeviceConfig"

	endpointNameGetDevicesV4         = "getDevicesV4"
	endpointNameGetDeviceV4          = "getDeviceV4"
	endpointNameCreateDeviceV4       = "createDeviceV4"
	endpointNameDeleteDeviceV4       = "deleteDeviceV4"
	endpointNameUpdateDeviceV4       = "updateDeviceV4"
	endpointNameUpdateDeviceStatusV4 = "updateDeviceStatusV4"

	endpointNameGetDeviceTemplateHistory = "getDeviceTemplateHistory"
	endpointNameAttachDeviceTemplates    = "attachDeviceTemplates"
	endpointNameUpdateDeviceTemplates    = "updateDeviceTemplates"
	endpointNameDetachDeviceTemplates    = "detachDeviceTemplates"
	endpointNameDetachDeviceTemplate     = "detachDeviceTemplate"

	endpointNameListDeviceTemplates  = "listDeviceTemplates"
	endpointNameGetDeviceTemplate    = "getDeviceTemplate"
	endpointNameSetDeviceTemplate    = "setDeviceTemplate"
	endpointNameDeleteDeviceTemplate = "deleteDeviceTemplate"

	endpointNameGetAllControlPlanes          = "getAllControlPlanes"
	endpointNameCreateControlPlane           = "createControlPlane"
	endpointNameGetControlPlane              = "getControlPlane"
	endpointNameUpdateControlPlane           = "updateControlPlane"
	endpointNameDeleteControlPlane           = "deleteControlPlane"
	endpointNameConnectControlPlane          = "connectControlPlane"
	endpointNameConnectUnmanagedControlPlane = "connectUnmanagedControlPlane"

	endpointNameGetEntityShard = "getEntityShard"

	endpointNameCreateDeviceConnection = "createDeviceConnection"
	endpointNameDeleteDeviceConnection = "deleteDeviceConnection"

	endpointNameUpdateTemplateAccess = "updateAccessTemplate"

	endpointNameGetLocationGeocode = "getLocationGeocode"

	serviceName = integration.ServiceNameManage
)

var (
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

		endpointNameCreateSiteV3:           {Method: "POST", Path: "/api/v3/sites"},
		endpointNameGetSitesV3:             {Method: "GET", Path: "/api/v3/sites"},
		endpointNameGetSiteV3:              {Method: "GET", Path: "/api/v3/sites/{{.siteId}}"},
		endpointNameUpdateSiteV3:           {Method: "PUT", Path: "/api/v3/sites/{{.siteId}}"},
		endpointNameDeleteSiteV3:           {Method: "DELETE", Path: "/api/v3/sites/{{.siteId}}"},
		endpointNameAddDevicetoSiteV3:      {Method: "PUT", Path: "/api/v3/sites/{{.siteId}}/devices/{{.deviceId}}"},
		endpointNameDeleteDeviceFromSiteV3: {Method: "DELETE", Path: "/api/v3/sites/{{.siteId}}/devices/{{.deviceId}}"},
		endpointNameUpdateSiteStatusV3:     {Method: "PUT", Path: "/api/v3/sites/{{.siteId}}/status"},

		endpointNameGetDeviceConfig: {Method: "GET", Path: "/api/v3/devices/{{.deviceInstanceId}}/config"},

		endpointNameGetDevicesV4:         {Method: "GET", Path: "/api/v4/devices"},
		endpointNameGetDeviceV4:          {Method: "GET", Path: "/api/v4/devices/{{.deviceId}}"},
		endpointNameCreateDeviceV4:       {Method: "POST", Path: "/api/v4/devices"},
		endpointNameDeleteDeviceV4:       {Method: "DELETE", Path: "/api/v4/devices/{{.deviceId}}"},
		endpointNameUpdateDeviceV4:       {Method: "PUT", Path: "/api/v4/devices/{{.deviceId}}"},
		endpointNameUpdateDeviceStatusV4: {Method: "PUT", Path: "/api/v4/devices/{{.deviceId}}/status"},

		endpointNameCreateDeviceActions: {Method: "POST", Path: "/api/v1/deviceActions"},
		endpointNameUpdateDeviceActions: {Method: "PUT", Path: "/api/v1/deviceActions"},

		endpointNameGetDeviceTemplateHistory: {Method: "GET", Path: "/api/v3/devices/{{.deviceInstanceId}}/templates"},
		endpointNameAttachDeviceTemplates:    {Method: "POST", Path: "/api/v3/devices/{{.deviceInstanceId}}/templates"},
		endpointNameUpdateDeviceTemplates:    {Method: "PUT", Path: "/api/v3/devices/{{.deviceInstanceId}}/templates"},
		endpointNameDetachDeviceTemplates:    {Method: "DELETE", Path: "/api/v3/devices/{{.deviceInstanceId}}/templates"},
		endpointNameDetachDeviceTemplate:     {Method: "DELETE", Path: "/api/v3/devices/{{.deviceInstanceId}}/templates/{{.templateId}}"},

		endpointNameListDeviceTemplates:  {Method: "GET", Path: "/api/v1/devicetemplates"},
		endpointNameGetDeviceTemplate:    {Method: "GET", Path: "/api/v1/devicetemplates/{{.id}}"},
		endpointNameSetDeviceTemplate:    {Method: "POST", Path: "/api/v1/devicetemplates"},
		endpointNameDeleteDeviceTemplate: {Method: "DELETE", Path: "/api/v1/devicetemplates/{{.id}}"},

		endpointNameGetAllControlPlanes:          {Method: "GET", Path: "/api/v1/controlplanes"},
		endpointNameCreateControlPlane:           {Method: "POST", Path: "/api/v1/controlplanes"},
		endpointNameGetControlPlane:              {Method: "GET", Path: "/api/v1/controlplanes/{{.controlPlaneId}}"},
		endpointNameUpdateControlPlane:           {Method: "PUT", Path: "/api/v1/controlplanes/{{.controlPlaneId}}"},
		endpointNameDeleteControlPlane:           {Method: "DELETE", Path: "/api/v1/controlplanes/{{.controlPlaneId}}"},
		endpointNameConnectControlPlane:          {Method: "POST", Path: "/api/v1/controlplanes/{{.controlPlaneId}}/connect"},
		endpointNameConnectUnmanagedControlPlane: {Method: "POST", Path: "/api/v1/controlplanes/connect"},

		endpointNameGetEntityShard:       {Method: "GET", Path: "/api/v2/shardmanagers/entity/{{.entityId}}"},
		endpointNameUpdateTemplateAccess: {Method: "PUT", Path: "/api/v1/devicetemplates/{{.templateId}}"},

		endpointNameCreateDeviceConnection: {Method: "POST", Path: "/api/v2/devices/connections"},
		endpointNameDeleteDeviceConnection: {Method: "DELETE", Path: "/api/v2/devices/connections/{{.deviceConnectionId}}"},
		endpointNameGetLocationGeocode:     {Method: "GET", Path: "/api/v1/location/geocode"},
	}
)

type Integration struct {
	integration.MsxServiceExecutor
}

func NewIntegration(ctx context.Context) (Api, error) {
	integrationInstance := IntegrationFromContext(ctx)
	if integrationInstance == nil {
		integrationInstance = &Integration{
			MsxServiceExecutor: integration.NewMsxService(ctx, serviceName, endpoints),
		}
	}
	return integrationInstance, nil
}

func NewIntegrationWithExecutor(executor integration.MsxServiceExecutor) Api {
	return &Integration{
		MsxServiceExecutor: executor,
	}
}

func (i *Integration) GetAdminHealth() (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetAdminHealth,
		Payload:      &integration.HealthDTO{},
		NoToken:      true,
	})
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

func (i *Integration) GetDeviceConfig(deviceInstanceId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetDeviceConfig,
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

func (i *Integration) CreateDeviceActions(deviceActionCreateRequests DeviceActionCreateRequests) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(deviceActionCreateRequests)
	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:       endpointNameCreateDeviceActions,
		EndpointParameters: map[string]string{},
		Body:               bodyBytes,
		ExpectEnvelope:     true,
	})
}

func (i *Integration) UpdateDeviceActions(deviceActionCreateRequests DeviceActionCreateRequests) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(deviceActionCreateRequests)
	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:       endpointNameUpdateDeviceActions,
		EndpointParameters: map[string]string{},
		Body:               bodyBytes,
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
		Payload: &paging.PaginatedResponse{
			Content: new(DeviceListResponse),
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

func (i *Integration) UpdateDeviceV4(deviceRequest DeviceUpdateRequest, deviceId string) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(deviceRequest)
	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameUpdateDeviceV4,
		EndpointParameters: map[string]string{
			"deviceId": deviceId,
		},
		Body:           bodyBytes,
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
		EndpointName: endpointNameUpdateDeviceStatusV4,
		EndpointParameters: map[string]string{
			"deviceId": deviceId,
		},
		Body:           bodyBytes,
		ExpectEnvelope: true,
	})
}

func (i *Integration) ListDeviceTemplates(serviceType string, tenantId *types.UUID) (*integration.MsxResponse, error) {
	queryParameters := make(url.Values)
	queryParameters.Set("serviceType", serviceType)
	if tenantId != nil {
		queryParameters.Set("tenantId", tenantId.String())
	}

	responsePayload := []DeviceTemplateListItemResponse{}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:    endpointNameListDeviceTemplates,
		QueryParameters: queryParameters,
		Payload:         &responsePayload,
		ExpectEnvelope:  true,
	})

}

func (i *Integration) GetDeviceTemplate(templateId types.UUID) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetDeviceTemplate,
		EndpointParameters: map[string]string{
			"id": templateId.String(),
		},
		Payload:        new(DeviceTemplateResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) AddDeviceTemplate(deviceTemplateCreateRequest DeviceTemplateCreateRequest) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(deviceTemplateCreateRequest)
	if err != nil {
		return nil, err
	}
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:   endpointNameSetDeviceTemplate,
		Body:           bodyBytes,
		Payload:        new(DeviceTemplateCreateResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) DeleteDeviceTemplate(templateId types.UUID) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameDeleteDeviceTemplate,
		EndpointParameters: map[string]string{
			"id": templateId.String(),
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetDeviceTemplateHistory(deviceInstanceId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetDeviceTemplateHistory,
		EndpointParameters: map[string]string{
			"deviceInstanceId": deviceInstanceId,
		},
		Payload:        new(AttachTemplateResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) AttachDeviceTemplates(deviceId string, attachTemplateRequest AttachTemplateRequest) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(attachTemplateRequest)
	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameAttachDeviceTemplates,
		EndpointParameters: map[string]string{
			"deviceInstanceId": deviceId,
		},
		Body:           bodyBytes,
		Payload:        new(AttachTemplateResponse),
		ExpectEnvelope: true,
	})
}

/*
	TODO: "updateDeviceTemplates"
*/

func (i *Integration) DetachDeviceTemplates(deviceId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameDetachDeviceTemplates,
		EndpointParameters: map[string]string{
			"deviceInstanceId": deviceId,
		},
		Payload:        new(AttachTemplateResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) DetachDeviceTemplate(deviceId string, templateId types.UUID) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameDetachDeviceTemplate,
		EndpointParameters: map[string]string{
			"deviceInstanceId": deviceId,
			"templateId":       templateId.String(),
		},
		Payload:        new(AttachTemplateResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) UpdateTemplateAccess(templateId string, deviceTemplateAccess DeviceTemplateAccess) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(deviceTemplateAccess)
	if err != nil {
		return nil, err
	}
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameUpdateTemplateAccess,
		EndpointParameters: map[string]string{
			"templateId": templateId,
		},
		Body:           bodyBytes,
		Payload:        new(DeviceTemplateAccessResponse),
		ExpectEnvelope: true,
	})
}

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

func (i *Integration) CreateDeviceConnection(deviceConnection DeviceConnectionCreateRequest) (*integration.MsxResponse, *DeviceConnectionResponse, error) {
	bodyBytes, err := json.Marshal(deviceConnection)

	if err != nil {
		return nil, nil, err
	}

	response, err := i.Execute(&integration.MsxEndpointRequest{
		EndpointName:   endpointNameCreateDeviceConnection,
		Body:           bodyBytes,
		Payload:        new(DeviceConnectionResponse),
		ExpectEnvelope: true,
	})

	if err != nil {
		return response, nil, err
	}

	return response, response.Payload.(*DeviceConnectionResponse), err
}

func (i *Integration) DeleteDeviceConnection(deviceConnectionId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointParameters: map[string]string{
			"deviceConnectionId": deviceConnectionId,
		},
		EndpointName:   endpointNameDeleteDeviceConnection,
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetLocationGeocode(address string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetLocationGeocode,
		QueryParameters: map[string][]string{
			"address": {address},
		},
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}
