package serviceconfigmanager

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"encoding/json"
	"strconv"
)

const (
	endpointNameGetAdminHealth = "getAdminHealth"

	endpointNameGetServiceConfigurations                 = "getServiceConfigurations"
	endpointNameGetServiceConfigurationByServiceConfigId = "getServiceConfigurationByServiceConfigId"
	endpointNameCreateServiceConfiguration               = "createServiceConfiguration"
	endpointNameUpdateServiceConfiguration               = "updateServiceConfiguration"
	endpointNameDeleteServiceConfiguration               = "deleteServiceConfiguration"
	endpointNameUpdateServiceConfigurationStatus         = "updateServiceConfigurationStatus"

	serviceName = integration.ServiceNameServiceConfig
)

var (
	logger    = log.NewLogger("msx.integration.usermanagement")
	endpoints = map[string]integration.MsxServiceEndpoint{
		endpointNameGetAdminHealth: {Method: "GET", Path: "/admin/health"},

		endpointNameGetServiceConfigurations:                 {Method: "GET", Path: "/api/v1/serviceconfigurations"},
		endpointNameGetServiceConfigurationByServiceConfigId: {Method: "GET", Path: "/api/v1/serviceconfigurations/{{.serviceConfigId}}"},
		endpointNameCreateServiceConfiguration:               {Method: "POST", Path: "/api/v1/serviceconfigurations"},
		endpointNameUpdateServiceConfiguration:               {Method: "PUT", Path: "/api/v1/serviceconfigurations"},
		endpointNameDeleteServiceConfiguration:               {Method: "DELETE", Path: "/api/v1/serviceconfigurations/{{.serviceConfigId}}"},
		endpointNameUpdateServiceConfigurationStatus:         {Method: "POST", Path: "/api/v1/serviceconfigurations/events/{{.serviceConfigId}}"},
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

	result.Response, err = i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetAdminHealth,
		Payload:      result.Payload,
		NoToken:      true,
	})

	return result, err
}

func (i *Integration) GetServiceConfigurations(page, pageSize int) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:       endpointNameGetServiceConfigurations,
		EndpointParameters: map[string]string{},
		QueryParameters: map[string][]string{
			"page":     {strconv.Itoa(page)},
			"pageSize": {strconv.Itoa(pageSize)},
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetServiceConfigurationByServiceConfigId(serviceConfigId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetServiceConfigurationByServiceConfigId,
		EndpointParameters: map[string]string{
			"serviceConfigId": serviceConfigId,
		},
		Payload:        new(ServiceConfigurationResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) CreateServiceConfiguration(configuration ServiceConfigurationRequest) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(configuration)

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:       endpointNameCreateServiceConfiguration,
		EndpointParameters: map[string]string{},
		Body:               bodyBytes,
		Payload:            new(ServiceConfigurationResponse),
		ExpectEnvelope:     true,
	})
}

func (i *Integration) UpdateServiceConfiguration(configuration ServiceConfigurationUpdateRequest) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(configuration)

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:       endpointNameUpdateServiceConfiguration,
		EndpointParameters: map[string]string{},
		Body:               bodyBytes,
		Payload:            new(ServiceConfigurationResponse),
		ExpectEnvelope:     true,
	})
}

func (i *Integration) DeleteServiceConfiguration(serviceConfigId string) (*integration.MsxResponse, error) {

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameDeleteServiceConfiguration,
		EndpointParameters: map[string]string{
			"serviceConfigId": serviceConfigId,
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) UpdateServiceConfigurationStatus(serviceConfigId string, serviceConfigurationStatus ServiceConfigurationStatusRequest) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(serviceConfigurationStatus)

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameUpdateServiceConfigurationStatus,
		EndpointParameters: map[string]string{
			"serviceConfigId": serviceConfigId,
		},
		Body:           bodyBytes,
		Payload:        new(ServiceConfigurationResponse),
		ExpectEnvelope: true,
	})
}
