package manage

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

const (
	endpointNameGetAdminHealth = "getAdminHealth"

	endpointNameGetSubscription    = "getSubscription"
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

	endpointNameGetDevice    = "getDevice"
	endpointNameGetDevices   = "getDevices"
	endpointNameCreateDevice = "createDevice"
	endpointNameUpdateDevice = "updateDevice"
	endpointNameDeleteDevice = "deleteDevice"

	endpointNameCreateManagedDevice = "createManagedDevice"
	endpointNameDeleteManagedDevice = "deleteManagedDevice"
	endpointNameGetDeviceConfig     = "getDeviceConfig"

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

	serviceName        = integration.ServiceNameManage
)

var (
	logger    = log.NewLogger("msx.integration.usermanagement")
	endpoints = map[string]integration.MsxServiceEndpoint{
		endpointNameGetAdminHealth: {Method: "GET", Path: "/admin/health"},

		endpointNameGetSubscription:    {Method: "GET", Path: "/api/v2/subscriptions/{{.subscriptionId}}"},
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

		endpointNameGetDevice:    {Method: "GET", Path: "/api/v1/devices/{{.deviceInstanceId}}"},
		endpointNameGetDevices:   {Method: "GET", Path: "/api/v2/devices"},
		endpointNameCreateDevice: {Method: "POST", Path: "/api/v1/devices/subscriptions/{{.subscriptionId}}"},
		endpointNameUpdateDevice: {Method: "PUT", Path: "/api/v1/devices/{{.deviceInstanceId}}"},
		endpointNameDeleteDevice: {Method: "DELETE", Path: "/api/v1/devices/{{.deviceInstanceId}}"},

		endpointNameCreateManagedDevice: {Method: "POST", Path: "/api/v3/devices"},
		endpointNameDeleteManagedDevice: {Method: "DELETE", Path: "/api/v3/devices/{{.deviceInstanceId}}"},
		endpointNameGetDeviceConfig:     {Method: "GET", Path: "/api/v3/devices/{{.deviceInstanceId}}/config"},

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
	} 
)

type Integration struct {
	*integration.MsxService
}

func NewIntegration(ctx context.Context) (Api, error) {
	return &Integration{
		MsxService: integration.NewMsxService(ctx, serviceName, endpoints),
	}, nil
}

func (i *Integration) GetAdminHealth() (result *HealthResult, err error) {
	result = &HealthResult{
		Payload: &integration.HealthDTO{},
	}

	result.Response, err = i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameGetAdminHealth,
		Payload:      result.Payload,
		NoToken:      true,
	})

	return result, err
}

func (i *Integration) GetSubscription(subscriptionId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameGetSubscription,
		EndpointParameters: map[string]string{
			"subscriptionId": subscriptionId,
		},
		Payload:        new(Pojo),
		ExpectEnvelope: true,
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

	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameCreateSubscription,
		EndpointParameters: map[string]string{
			"tenantId": tenantId,
		},
		Body:           bodyBytes,
		Payload:        new(Pojo),
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

	return i.Execute(&integration.MsxRequest{
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
	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameDeleteSubscription,
		EndpointParameters: map[string]string{
			"subscriptionId": subscriptionId,
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetServiceInstance(serviceInstanceId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameGetServiceInstance,
		EndpointParameters: map[string]string{
			"serviceInstanceId": serviceInstanceId,
		},
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetSubscriptionServiceInstances(subscriptionId string, page, pageSize int) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxRequest{
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

	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameCreateServiceInstance,
		EndpointParameters: map[string]string{
			"subscriptionId": subscriptionId,
		},
		Body:           bodyBytes,
		Payload:        new(Pojo),
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

	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameUpdateServiceInstance,
		EndpointParameters: map[string]string{
			"serviceInstanceId": serviceInstanceId,
		},
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) DeleteServiceInstance(serviceInstanceId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameDeleteServiceInstance,
		EndpointParameters: map[string]string{
			"serviceInstanceId": serviceInstanceId,
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetSite(siteId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameGetSite,
		EndpointParameters: map[string]string{
			"siteId": siteId,
		},
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

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

	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameCreateSite,
		EndpointParameters: map[string]string{
			"subscriptionId": subscriptionId,
		},
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

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

	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameUpdateSite,
		EndpointParameters: map[string]string{
			"siteId": siteId,
		},
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) DeleteSite(siteId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameDeleteSite,
		EndpointParameters: map[string]string{
			"siteId": siteId,
		},
		ExpectEnvelope: true,
	})
}

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

	return i.Execute(&integration.MsxRequest{
		EndpointName:   endpointNameCreateManagedDevice,
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) DeleteManagedDevice(deviceInstanceId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameDeleteManagedDevice,
		EndpointParameters: map[string]string{
			"deviceInstanceId": deviceInstanceId,
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetDeviceConfig(deviceInstanceId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameGetDeviceConfig,
		EndpointParameters: map[string]string{
			"deviceInstanceId": deviceInstanceId,
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetDevice(deviceInstanceId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameGetDevice,
		EndpointParameters: map[string]string{
			"deviceInstanceId": deviceInstanceId,
		},
		ExpectEnvelope: true,
	})
}

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

	return i.Execute(&integration.MsxRequest{
		EndpointName:    endpointNameGetDevices,
		QueryParameters: queryParameters,
		ExpectEnvelope:  true,
	})
}

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

	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameCreateDevice,
		EndpointParameters: map[string]string{
			"subscriptionId": subscriptionId,
		},
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) UpdateDevice(deviceInstanceId string, deviceAttribute, deviceDefAttribute, status map[string]string) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(&Pojo{
		"deviceAttribute":    deviceAttribute,
		"deviceDefAttribute": deviceDefAttribute,
		"status":             status,
	})

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameUpdateDevice,
		EndpointParameters: map[string]string{
			"deviceInstanceId": deviceInstanceId,
		},
		Body:           bodyBytes,
		ExpectEnvelope: true,
	})
}

func (i *Integration) DeleteDevice(deviceInstanceId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameDeleteDevice,
		EndpointParameters: map[string]string{
			"deviceInstanceId": deviceInstanceId,
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetDeviceTemplateHistory(deviceInstanceId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxRequest{
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
		queryParameters[*tenantId] = []string{*tenantId}
	}

	return i.Execute(&integration.MsxRequest{
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

	return i.Execute(&integration.MsxRequest{
		EndpointName:   endpointNameCreateControlPlane,
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetControlPlane(controlPlaneId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxRequest{
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

	return i.Execute(&integration.MsxRequest{
		EndpointName:       endpointNameUpdateControlPlane,
		EndpointParameters: map[string]string{"controlPlaneId": controlPlaneId},
		Body:               bodyBytes,
		Payload:            new(Pojo),
		ExpectEnvelope:     true,
	})
}

func (i *Integration) DeleteControlPlane(controlPlaneId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxRequest{
		EndpointName:       endpointNameDeleteControlPlane,
		EndpointParameters: map[string]string{"controlPlaneId": controlPlaneId},
		Payload:            new(Pojo),
		ExpectEnvelope:     true,
	})
}

func (i *Integration) ConnectControlPlane(controlPlaneId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxRequest{
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

	return i.Execute(&integration.MsxRequest{
		EndpointName:   endpointNameConnectUnmanagedControlPlane,
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}
