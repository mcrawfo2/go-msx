// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package idm

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"encoding/json"
	"github.com/pkg/errors"
	"net/url"
	"strconv"
)

const (
	endpointNameIsTokenValid = "isTokenValid"

	endpointNameGetMyProvider = "getMyProvider"

	endpointNameGetProviderByName = "getProviderByName"

	endpointNameGetProviderExtensionByName = "getProviderExtensionByName"

	endpointNameGetTenantById   = "getTenantById"
	endpointNameGetTenantByIdV8 = "getTenantByIdV8"
	endpointNameGetTenantByName = "getTenantByName"

	endpointNameGetUserById   = "getUserById"
	endpointNameGetUserByIdV8 = "getUserByIdV8"

	endpointNameGetRoles   = "getRoles"
	endpointNameCreateRole = "createRole"
	endpointNameUpdateRole = "updateRole"
	endpointNameDeleteRole = "deleteRole"

	endpointNameGetCapabilities         = "getCapabilities"
	endpointNameBatchCreateCapabilities = "batchCreateCapabilities"
	endpointNameBatchUpdateCapabilities = "batchUpdateCapabilities"
	endpointNameDeleteCapability        = "deleteCapability"

	serviceName                       = integration.ServiceNameUserManagement
	ErrorMessageFailedToSerializeBody = "Failed to serialize body"
	pathApiV1RolesCapabilities        = "/api/v1/roles/capabilities"
)

var (
	endpoints = map[string]integration.MsxServiceEndpoint{
		endpointNameIsTokenValid: {Method: "GET", Path: "/api/v1/isTokenValid"},

		endpointNameGetMyProvider:     {Method: "GET", Path: "/api/v1/providers"},
		endpointNameGetProviderByName: {Method: "GET", Path: "/api/v1/providers/{{.providerName}}"},

		endpointNameGetProviderExtensionByName: {Method: "GET", Path: "/api/v1/providers/providerextension/parameters/{{.providerExtensionName}}"},

		endpointNameGetTenantById:   {Method: "GET", Path: "/api/v3/tenants/{{.tenantId}}"},
		endpointNameGetTenantByIdV8: {Method: "GET", Path: "/api/v8/tenants/{{.tenantId}}"},
		endpointNameGetTenantByName: {Method: "GET", Path: "/api/v1/tenants/{{.tenantName}}"},

		endpointNameGetUserById:   {Method: "GET", Path: "/api/v2/users/{{.userId}}"},
		endpointNameGetUserByIdV8: {Method: "GET", Path: "/api/v8/users/{{.userId}}"},

		endpointNameGetRoles:   {Method: "GET", Path: "/api/v1/roles"},
		endpointNameCreateRole: {Method: "POST", Path: "/api/v1/roles"},
		endpointNameUpdateRole: {Method: "PUT", Path: "/api/v1/roles/{roleName}"},
		endpointNameDeleteRole: {Method: "DELETE", Path: "/api/v1/roles/{roleName}"},

		endpointNameGetCapabilities:         {Method: "GET", Path: pathApiV1RolesCapabilities},
		endpointNameBatchCreateCapabilities: {Method: "POST", Path: pathApiV1RolesCapabilities},
		endpointNameBatchUpdateCapabilities: {Method: "PUT", Path: pathApiV1RolesCapabilities},
		endpointNameDeleteCapability:        {Method: "DELETE", Path: "/api/v1/roles/capabilities/{capabilityName}"},
	}
)

func NewIntegration(ctx context.Context) (Api, error) {
	integrationInstance := IntegrationFromContext(ctx)
	if integrationInstance == nil {
		integrationInstance = &Integration{
			MsxContextServiceExecutor: integration.NewMsxService(ctx, serviceName, endpoints),
		}
	}
	return integrationInstance, nil
}

func NewIntegrationWithExecutor(executor integration.MsxContextServiceExecutor) Api {
	return &Integration{
		MsxContextServiceExecutor: executor,
	}
}

type Integration struct {
	integration.MsxContextServiceExecutor
}

func (i *Integration) IsTokenActive() (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameIsTokenValid,
		ErrorPayload: new(integration.OAuthErrorDTO),
	})
}

func (i *Integration) GetMyProvider() (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetMyProvider,
		Payload:      &ProviderResponse{},
		ErrorPayload: new(integration.ErrorDTO),
	})
}

func (i *Integration) GetProviderByName(name string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetProviderByName,
		EndpointParameters: map[string]string{
			"providerName": name,
		},
		Payload:      &ProviderResponse{},
		ErrorPayload: new(integration.ErrorDTO),
	})
}

func (i *Integration) GetProviderExtensionByName(name string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetProviderExtensionByName,
		EndpointParameters: map[string]string{
			"providerExtensionName": name,
		},
		Payload:      &ProviderExtensionResponse{},
		ErrorPayload: new(integration.ErrorDTO),
	})
}

// Deprecated: The underlying REST API was deprecated in 3.10.0.  v8 (or newer) API should be used instead.
func (i *Integration) GetUserById(userId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetUserById,
		EndpointParameters: map[string]string{
			"userId": userId,
		},
		ExpectEnvelope: true,
		Payload:        new(UserResponse),
	})
}

func (i *Integration) GetUserByIdV8(userId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetUserByIdV8,
		EndpointParameters: map[string]string{
			"userId": userId,
		},
		ExpectEnvelope: false,
		Payload:        new(UserResponseV8),
	})
}

// Deprecated: The underlying REST API was deprecated in 3.10.0.  v8 (or newer) API should be used instead.
func (i *Integration) GetTenantById(tenantId string) (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetTenantById,
		EndpointParameters: map[string]string{
			"tenantId": tenantId,
		},
		ExpectEnvelope: true,
		Payload:        new(TenantResponse),
	})
}

func (i *Integration) GetTenantByIdV8(tenantId string) (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetTenantByIdV8,
		EndpointParameters: map[string]string{
			"tenantId": tenantId,
		},
		ExpectEnvelope: false,
		Payload:        new(TenantResponseV8),
	})
}

// Deprecated: Tenants should generally be access by ID, not tenantName.  The REST API is retired
// and due for decomissioning.
func (i *Integration) GetTenantByName(tenantName string) (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetTenantByName,
		EndpointParameters: map[string]string{
			"tenantName": tenantName,
		},
		Payload: new(TenantResponse),
	})
}

func (i *Integration) GetRoles(resolvePermissionNames bool, p paging.Request) (result *integration.MsxResponse, err error) {
	qp := p.QueryParameters()
	qp.Set("resolvepermissionname", strconv.FormatBool(resolvePermissionNames))

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:    endpointNameGetRoles,
		QueryParameters: qp,
		Payload:         new(RoleListResponse),
		ExpectEnvelope:  false,
	})
}

func (i *Integration) CreateRole(populator bool, body RoleCreateRequest) (result *integration.MsxResponse, err error) {
	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(body); err != nil {
		return nil, errors.Wrap(err, ErrorMessageFailedToSerializeBody)
	}

	qp := url.Values{}
	if body.Owner != "" {
		qp.Set("owner", body.Owner)
	}
	if populator {
		qp.Set("dbinstaller", strconv.FormatBool(populator))
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:    endpointNameCreateRole,
		QueryParameters: qp,
		Body:            bodyBytes,
		Payload:         new(RoleResponse),
		ExpectEnvelope:  false,
	})
}

func (i *Integration) UpdateRole(populator bool, body RoleUpdateRequest) (result *integration.MsxResponse, err error) {
	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(body); err != nil {
		return nil, errors.Wrap(err, ErrorMessageFailedToSerializeBody)
	}

	qp := url.Values{}
	if body.Owner != "" {
		qp.Set("owner", body.Owner)
	}
	if populator {
		qp.Set("dbinstaller", strconv.FormatBool(populator))
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameUpdateRole,
		EndpointParameters: map[string]string{
			"roleName": body.RoleName,
		},
		QueryParameters: qp,
		Body:            bodyBytes,
		Payload:         new(RoleResponse),
		ExpectEnvelope:  false,
	})
}

func (i *Integration) DeleteRole(roleName string) (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameDeleteRole,
		EndpointParameters: map[string]string{
			"roleName": roleName,
		},
		Payload:        new(RoleResponse),
		ExpectEnvelope: false,
	})
}

func (i *Integration) GetCapabilities(p paging.Request) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:    endpointNameGetCapabilities,
		QueryParameters: p.QueryParameters(),
		Payload:         new(CapabilityListResponse),
		ExpectEnvelope:  false,
	})
}

func (i *Integration) BatchCreateCapabilities(populator bool, owner string, capabilities []CapabilityCreateRequest) (result *integration.MsxResponse, err error) {
	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(CapabilityBatchCreateRequest{capabilities}); err != nil {
		return nil, errors.Wrap(err, ErrorMessageFailedToSerializeBody)
	}

	qp := url.Values{}
	if owner != "" {
		qp.Set("owner", owner)
	}
	if populator {
		qp.Set("dbinstaller", strconv.FormatBool(populator))
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:    endpointNameBatchCreateCapabilities,
		QueryParameters: qp,
		Body:            bodyBytes,
		Payload:         new(CapabilityListResponse),
		ExpectEnvelope:  false,
	})
}

func (i *Integration) BatchUpdateCapabilities(populator bool, owner string, capabilities []CapabilityUpdateRequest) (result *integration.MsxResponse, err error) {
	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(CapabilityBatchUpdateRequest{capabilities}); err != nil {
		return nil, errors.Wrap(err, ErrorMessageFailedToSerializeBody)
	}

	qp := url.Values{}
	if owner != "" {
		qp.Set("owner", owner)
	}
	if populator {
		qp.Set("dbinstaller", strconv.FormatBool(populator))
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:    endpointNameBatchUpdateCapabilities,
		QueryParameters: qp,
		Body:            bodyBytes,
		Payload:         new(CapabilityListResponse),
		ExpectEnvelope:  false,
	})
}

func (i *Integration) DeleteCapability(populator bool, owner string, name string) (*integration.MsxResponse, error) {
	qp := url.Values{}
	if owner != "" {
		qp.Set("owner", owner)
	}
	if populator {
		qp.Set("dbinstaller", strconv.FormatBool(populator))
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameDeleteCapability,
		EndpointParameters: map[string]string{
			"capabilityName": name,
		},
		QueryParameters: qp,
		Payload:         new(CapabilityResponse),
		ExpectEnvelope:  false,
	})
}
