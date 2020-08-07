package serviceconfigmanager

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"net/url"
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

	endpointNameGetServiceConfigurationAssignmentByAssignmentId                               = "getServiceConfigurationAssignmentsByAssignmentId"
	endpointNameGetServiceConfigurationAssignmentAll                                          = "getServiceConfigurationAssignmentsAll"
	endpointNameCreateServiceConfigurationAssignment                                          = "createServiceConfigurationAssignment"
	endpointNameDeleteServiceConfigurationAssignment                                          = "deleteServiceConfigurationAssignment"
	endpointNameGetServiceConfigurationAssignmentsByServiceConfigurationId                    = "getServiceConfigurationAssignmentsByServiceConfigurationId"
	endpointNameUpdateServiceConfigurationAssignmentStatusByServiceConfigurationIdAndTenantId = "updateServiceConfigurationAssignmentStatusByServiceConfigurationIdAndTenantId"
	endpointNameGetTenantAssignmentsByServiceConfigurationId                                  = "getTenantAssignmentsByServiceConfigurationId"

	endpointNameCreateServiceConfigurationApplication                                                = "createServiceConfigurationApplication"
	endpointNameUpdateServiceConfigurationApplicationStatus                                          = "updateServiceConfigurationApplicationStatus"
	endpointNameDeleteServiceConfigurationApplication                                                = "deleteServiceConfigurationApplicationStatus"
	endpointNameGetServiceConfigurationApplications                                                  = "getServiceConfigurationApplications"
	endpointNameGetServiceConfigurationApplicationById                                               = "getServiceConfigurationApplicationById"
	endpointNameGetServiceConfigurationApplicationByServiceConfigIdTargetEntityTypeAndTargetEntityId = "getServiceConfigurationApplicationByServiceConfigIdTargetEntityTypeAndTargetEntityId"

	serviceName = integration.ServiceNameServiceConfig
)

var (
	endpoints = map[string]integration.MsxServiceEndpoint{
		endpointNameGetAdminHealth: {Method: "GET", Path: "/admin/health"},

		endpointNameGetServiceConfigurations:                 {Method: "GET", Path: "/api/v1/serviceconfigurations"},
		endpointNameGetServiceConfigurationByServiceConfigId: {Method: "GET", Path: "/api/v1/serviceconfigurations/{{.serviceConfigId}}"},
		endpointNameCreateServiceConfiguration:               {Method: "POST", Path: "/api/v1/serviceconfigurations"},
		endpointNameUpdateServiceConfiguration:               {Method: "PUT", Path: "/api/v1/serviceconfigurations"},
		endpointNameDeleteServiceConfiguration:               {Method: "DELETE", Path: "/api/v1/serviceconfigurations/{{.serviceConfigId}}"},
		endpointNameUpdateServiceConfigurationStatus:         {Method: "POST", Path: "/api/v1/serviceconfigurations/events/{{.serviceConfigId}}"},

		endpointNameGetServiceConfigurationAssignmentByAssignmentId:                               {Method: "GET", Path: "/api/v1/serviceconfigurations/assign/{{.id}}"},
		endpointNameGetServiceConfigurationAssignmentAll:                                          {Method: "GET", Path: "/api/v1/serviceconfigurations/assign/all"},
		endpointNameCreateServiceConfigurationAssignment:                                          {Method: "POST", Path: "/api/v1/serviceconfigurations/assign/{{.serviceConfigId}}"},
		endpointNameDeleteServiceConfigurationAssignment:                                          {Method: "DELETE", Path: "/api/v1/serviceconfigurations/assign/{{.serviceConfigId}}"},
		endpointNameGetServiceConfigurationAssignmentsByServiceConfigurationId:                    {Method: "GET", Path: "/api/v1/serviceconfigurations/assign/config/{{.serviceConfigId}}"},
		endpointNameUpdateServiceConfigurationAssignmentStatusByServiceConfigurationIdAndTenantId: {Method: "POST", Path: "/api/v1/serviceconfigurations/assign/events/{{.serviceConfigId}}/{{.assignedTenantId}}"},
		endpointNameGetTenantAssignmentsByServiceConfigurationId:                                  {Method: "GET", Path: "/api/v1/serviceconfigurations/assign/tenants/{{.serviceConfigId}}"},

		endpointNameCreateServiceConfigurationApplication:                                                {Method: "POST", Path: "/api/v1/serviceconfigurations/applications"},
		endpointNameUpdateServiceConfigurationApplicationStatus:                                          {Method: "POST", Path: "/api/v1/serviceconfigurations/{{.serviceConfigId}}/applications/{{.applicationId}}/status"},
		endpointNameDeleteServiceConfigurationApplication:                                                {Method: "DELETE", Path: "/api/v1/serviceconfigurations/applications/{{.id}}"},
		endpointNameGetServiceConfigurationApplications:                                                  {Method: "GET", Path: "/api/v1/serviceconfigurations/applications"},
		endpointNameGetServiceConfigurationApplicationById:                                               {Method: "GET", Path: "/api/v1/serviceconfigurations/applications/{{.id}}"},
		endpointNameGetServiceConfigurationApplicationByServiceConfigIdTargetEntityTypeAndTargetEntityId: {Method: "GET", Path: "/api/v1/serviceconfigurations/{(.serviceConfigId}}/applications/target/{(.targetEntityType}}/{{.targetEntityId}}"},
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

func (i *Integration) GetServiceConfigurationByServiceConfigId(serviceConfigId types.UUID) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetServiceConfigurationByServiceConfigId,
		EndpointParameters: map[string]string{
			"serviceConfigId": serviceConfigId.String(),
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

func (i *Integration) DeleteServiceConfiguration(serviceConfigId types.UUID) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameDeleteServiceConfiguration,
		EndpointParameters: map[string]string{
			"serviceConfigId": serviceConfigId.String(),
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) UpdateServiceConfigurationStatus(serviceConfigId types.UUID, serviceConfigurationStatus StatusUpdateRequest) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(serviceConfigurationStatus)

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameUpdateServiceConfigurationStatus,
		EndpointParameters: map[string]string{
			"serviceConfigId": serviceConfigId.String(),
		},
		Body:           bodyBytes,
		Payload:        new(ServiceConfigurationResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetAllServiceConfigurationAssignments(page, pageSize int, filterTenantId types.UUID) (*integration.MsxResponse, error) {
	queryMap := map[string][]string{
		"page":     {strconv.Itoa(page)},
		"pageSize": {strconv.Itoa(pageSize)},
	}

	if filterTenantId != nil {
		queryMap["assignedTenantId"] = []string{filterTenantId.String()}
	}
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:       endpointNameGetServiceConfigurationAssignmentAll,
		EndpointParameters: map[string]string{},
		QueryParameters:    queryMap,
		Payload: &paging.PaginatedResponse{
			Content: []ServiceConfigurationAssignmentResponse{},
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetServiceConfigurationAssignmentByAssignmentId(assignmentId types.UUID) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetServiceConfigurationAssignmentByAssignmentId,
		EndpointParameters: map[string]string{
			"id": assignmentId.String(),
		},
		Payload:        &ServiceConfigurationAssignmentResponse{},
		ExpectEnvelope: true,
	})
}

func (i *Integration) CreateServiceConfigurationAssignment(serviceConfigId types.UUID, tenantIdList []types.UUID) (*integration.MsxResponse, error) {

	bodyBytes, err := json.Marshal(ServiceConfigurationAssignmentRequest{
		Tenants: tenantIdList,
	})

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameCreateServiceConfigurationAssignment,
		EndpointParameters: map[string]string{
			"serviceConfigId": serviceConfigId.String(),
		},
		Body:           bodyBytes,
		ExpectEnvelope: true,
	})
}

func (i *Integration) DeleteServiceConfigurationAssignments(serviceConfigId types.UUID, tenantIdList []types.UUID) (*integration.MsxResponse, error) {

	bodyBytes, err := json.Marshal(ServiceConfigurationAssignmentRequest{
		Tenants: tenantIdList,
	})

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameDeleteServiceConfigurationAssignment,
		EndpointParameters: map[string]string{
			"serviceConfigId": serviceConfigId.String(),
		},
		Body:           bodyBytes,
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetServiceConfigurationAssignmentsByServiceConfigurationId(serviceConfigurationId types.UUID) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetServiceConfigurationAssignmentsByServiceConfigurationId,
		EndpointParameters: map[string]string{
			"serviceConfigId": serviceConfigurationId.String(),
		},
		Payload: &paging.PaginatedResponse{
			Content: []ServiceConfigurationAssignmentResponse{},
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) UpdateServiceConfigurationAssignmentStatusByServiceConfigurationIdAndTenantId(serviceConfigId types.UUID, tenantId types.UUID, status StatusUpdateRequest) (*integration.MsxResponse, error) {

	bodyBytes, err := json.Marshal(status)

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameUpdateServiceConfigurationAssignmentStatusByServiceConfigurationIdAndTenantId,
		EndpointParameters: map[string]string{
			"serviceConfigId":  serviceConfigId.String(),
			"assignedTenantId": tenantId.String(),
		},
		Body:           bodyBytes,
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetTenantAssignmentsByServiceConfigurationId(serviceConfigurationId types.UUID) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetTenantAssignmentsByServiceConfigurationId,
		EndpointParameters: map[string]string{
			"serviceConfigId": serviceConfigurationId.String(),
		},
		Payload:        &Pojo{},
		ExpectEnvelope: true,
	})
}

func (i *Integration) CreateServiceConfigurationApplication(applicationRequest ServiceConfigurationApplicationRequest) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(applicationRequest)

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:       endpointNameCreateServiceConfigurationApplication,
		EndpointParameters: map[string]string{},
		Body:               bodyBytes,
		Payload:            &ServiceConfigurationApplicationResponse{},
		ExpectEnvelope:     true,
	})
}

func (i *Integration) UpdateServiceConfigurationApplicationStatus(applicationId types.UUID, serviceConfigId types.UUID, applicationStatusUpdateRequest ServiceConfigurationApplicationStatusUpdateRequest) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(applicationStatusUpdateRequest)

	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameUpdateServiceConfigurationApplicationStatus,
		EndpointParameters: map[string]string{
			"applicationId":   applicationId.String(),
			"serviceConfigId": serviceConfigId.String(),
		},
		Body:           bodyBytes,
		ExpectEnvelope: true,
	})
}

func (i *Integration) DeleteServiceConfigurationApplication(applicationId types.UUID) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameDeleteServiceConfigurationApplication,
		EndpointParameters: map[string]string{
			"id": applicationId.String(),
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetServiceConfigurationApplications(page, pageSize int, tenantId types.UUID, serviceConfigId types.UUID, sortBy, sortOrder, targetEntityId, targetEntityType *string) (*integration.MsxResponse, error) {
	pageString := strconv.Itoa(page)
	pageSizeString := strconv.Itoa(pageSize)

	searchParameters := map[string]*string{
		"sortBy":           sortBy,
		"sortOrder":        sortOrder,
		"targetEntityId":   targetEntityId,
		"targetEntityType": targetEntityType,
		"page":             &pageString,
		"pageSize":         &pageSizeString,
	}

	// Convert optional search queries into query parameters
	queryParameters := make(url.Values)
	for k, v := range searchParameters {
		if v != nil && *v != "" {
			queryParameters[k] = []string{*v}
		}
	}
	queryParameters.Add("tenantId", tenantId.String())

	if serviceConfigId != nil {
		queryParameters.Add("serviceConfigId", serviceConfigId.String())
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:       endpointNameGetServiceConfigurationApplications,
		EndpointParameters: map[string]string{},
		QueryParameters:    queryParameters,
		ExpectEnvelope:     true,
	})
}

func (i *Integration) GetServiceConfigurationApplicationById(applicationId types.UUID) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetServiceConfigurationApplicationById,
		EndpointParameters: map[string]string{
			"id": applicationId.String(),
		},
		Payload:        &ServiceConfigurationApplicationResponse{},
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetServiceConfigurationApplicationByServiceConfigIdTargetEntityTypeAndTargetEntityId(serviceConfigId types.UUID, targetEntityType string, targetEntityId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetServiceConfigurationApplicationByServiceConfigIdTargetEntityTypeAndTargetEntityId,
		EndpointParameters: map[string]string{
			"serviceConfigId":  serviceConfigId.String(),
			"targetEntityType": targetEntityType,
			"targetEntityId":   targetEntityId,
		},
		ExpectEnvelope: true,
	})
}
