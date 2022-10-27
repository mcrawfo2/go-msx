// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package idm

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	integrationClientTest "cto-github.cisco.com/NFV-BU/go-msx/integration/testhelpers/clienttest"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/clienttest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

func TestNewIntegration(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	ctxWithConfig := configtest.ContextWithNewInMemoryConfig(
		context.Background(),
		map[string]string{
			"remoteservice.usermanagementservice.service": "usermanagementservice",
		})

	ctxWithConfigDifferentName := configtest.ContextWithNewInMemoryConfig(
		context.Background(),
		map[string]string{
			"remoteservice.usermanagementservice.service": "testservice1",
			"remoteservice.authservice.service":           "testservice2",
			"remoteservice.secretsservice.service":        "testservice3",
		})

	tests := []struct {
		name string
		args args
		want Api
	}{
		{
			name: "NonExisting",
			args: args{
				ctx: ctxWithConfig,
			},
			want: &Integration{
				MsxContextServiceExecutor: integration.NewMsxService(ctxWithConfig, serviceName, endpoints),
			},
		},
		{
			name: "Existing",
			args: args{
				ctx: ContextWithIntegration(ctxWithConfig, &Integration{}),
			},
			want: &Integration{},
		},
		{
			name: "ServiceName",
			args: args{
				ctx: ctxWithConfig,
			},
			want: &Integration{
				MsxContextServiceExecutor: integration.NewMsxService(ctxWithConfig, serviceName, endpoints),
			},
		},
		{
			name: "DifferentServiceName",
			args: args{
				ctx: ctxWithConfigDifferentName,
			},
			want: &Integration{
				MsxContextServiceExecutor: integration.NewMsxService(ctxWithConfigDifferentName, serviceName, endpoints),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := NewIntegration(tt.args.ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIntegration() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type IdmIntegrationTest struct {
	*integrationClientTest.EndpointTest
}

func NewIdmIntegrationTest() *IdmIntegrationTest {
	return &IdmIntegrationTest{
		EndpointTest: new(integrationClientTest.EndpointTest).WithEndpoints(endpoints),
	}
}

type ManageCall func(t *testing.T, api Api) (*integration.MsxResponse, error)
type AuthCall func(t *testing.T, api Api) (*integration.MsxResponse, []types.UUID, error)

func (m *IdmIntegrationTest) WithCall(call ManageCall) *IdmIntegrationTest {
	m.EndpointTest.WithCall(func(t *testing.T, executor integration.MsxContextServiceExecutor) (*integration.MsxResponse, error) {
		return call(t, NewIntegrationWithExecutor(executor))
	})
	return m
}

func (m *IdmIntegrationTest) WithMultiTenantResultCall(call AuthCall) *IdmIntegrationTest {
	m.EndpointTest.WithMultiTenantResultCall(func(t *testing.T, executor integration.MsxContextServiceExecutor) (*integration.MsxResponse, []types.UUID, error) {
		return call(t, NewIntegrationWithExecutor(executor))
	})
	return m
}

func TestIntegration_IsTokenActive(t *testing.T) {
	NewIdmIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.IsTokenActive()
		}).
		WithResponseStatus(http.StatusOK).
		WithRequestPredicate(clienttest.EndpointRequestHasName(endpointNameIsTokenValid)).
		WithRequestPredicate(clienttest.EndpointRequestHasToken(true)).
		WithEndpointPredicate(clienttest.ServiceEndpointHasMethod(http.MethodGet)).
		WithEndpointPredicate(clienttest.ServiceEndpointHasPath("/api/v1/isTokenValid")).
		Test(t)
}

func TestIntegration_GetMyProvider(t *testing.T) {
	NewIdmIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetMyProvider()
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(ProviderResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetMyProvider),
			clienttest.EndpointRequestHasToken(true)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v1/providers")).
		Test(t)
}

func TestIntegration_GetProviderByName(t *testing.T) {
	const providerName = "provider-name"

	NewIdmIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetProviderByName(providerName)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(ProviderResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetProviderByName),
			clienttest.EndpointRequestHasEndpointParameter("providerName", providerName),
			clienttest.EndpointRequestHasToken(true)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v1/providers/{{.providerName}}")).
		Test(t)
}

func TestIntegration_GetProviderExtensionByName(t *testing.T) {
	const providerExtensionName = "provider-extension-name"

	NewIdmIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetProviderExtensionByName(providerExtensionName)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(ProviderResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetProviderExtensionByName),
			clienttest.EndpointRequestHasEndpointParameter("providerExtensionName", providerExtensionName),
			clienttest.EndpointRequestHasToken(true)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v1/providers/providerextension/parameters/{{.providerExtensionName}}")).
		Test(t)
}

func TestIntegration_GetUserById(t *testing.T) {
	const userId = "user-id"

	NewIdmIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetUserById(userId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(UserResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetUserById),
			clienttest.EndpointRequestHasEndpointParameter("userId", userId),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(true)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v2/users/{{.userId}}")).
		Test(t)
}

func TestIntegration_GetUserByIdV8(t *testing.T) {
	const userId = "user-id"

	NewIdmIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetUserByIdV8(userId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(UserResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetUserByIdV8),
			clienttest.EndpointRequestHasEndpointParameter("userId", userId),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(false)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v8/users/{{.userId}}")).
		Test(t)
}

func TestIntegration_GetTenantById(t *testing.T) {
	const tenantId = "tenant-id"

	NewIdmIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetTenantById(tenantId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(TenantResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetTenantById),
			clienttest.EndpointRequestHasEndpointParameter("tenantId", tenantId),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(true)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v3/tenants/{{.tenantId}}")).
		Test(t)
}

func TestIntegration_GetTenantByIdV8(t *testing.T) {
	const tenantId = "tenant-id"

	NewIdmIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetTenantByIdV8(tenantId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(TenantResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetTenantByIdV8),
			clienttest.EndpointRequestHasEndpointParameter("tenantId", tenantId),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(false)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v8/tenants/{{.tenantId}}")).
		Test(t)
}

func TestIntegration_GetTenantByName(t *testing.T) {
	const tenantName = "tenant-name"

	NewIdmIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetTenantByName(tenantName)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(TenantResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetTenantByName),
			clienttest.EndpointRequestHasEndpointParameter("tenantName", tenantName),
			clienttest.EndpointRequestHasToken(true)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v1/tenants/{{.tenantName}}")).
		Test(t)
}

func TestIntegration_GetRoles(t *testing.T) {
	const resolve = true
	preq := paging.Request{
		Page: 0,
		Size: 10,
	}

	NewIdmIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetRoles(resolve, preq)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(RoleListResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetRoles),
			clienttest.EndpointRequestHasQueryParam("resolvepermissionname", strconv.FormatBool(resolve)),
			clienttest.EndpointRequestHasQueryParam("page", strconv.Itoa(int(preq.Page))),
			clienttest.EndpointRequestHasQueryParam("pageSize", strconv.Itoa(int(preq.Size))),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(false)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v1/roles")).
		Test(t)
}

func TestIntegration_CreateRole(t *testing.T) {
	const installer = true
	var request = RoleCreateRequest{
		Description: "description",
		Owner:       "owner",
	}

	NewIdmIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.CreateRole(installer, request)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(RoleResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameCreateRole),
			clienttest.EndpointRequestHasQueryParam("owner", request.Owner),
			clienttest.EndpointRequestHasQueryParam("dbinstaller", strconv.FormatBool(installer)),
			clienttest.EndpointRequestHasBodyJsonValue("description", request.Description),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(false)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v1/roles")).
		Test(t)
}

func TestIntegration_UpdateRole(t *testing.T) {
	const installer = true
	var request = RoleUpdateRequest{
		Description: "description",
		Owner:       "owner",
		RoleName:    "role-name",
	}

	NewIdmIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.UpdateRole(installer, request)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(RoleResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameUpdateRole),
			clienttest.EndpointRequestHasQueryParam("owner", request.Owner),
			clienttest.EndpointRequestHasQueryParam("dbinstaller", strconv.FormatBool(installer)),
			clienttest.EndpointRequestHasBodyJsonValue("description", request.Description),
			clienttest.EndpointRequestHasEndpointParameter("roleName", request.RoleName),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(false)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPut),
			clienttest.ServiceEndpointHasPath("/api/v1/roles/{roleName}")).
		Test(t)
}

func TestIntegration_DeleteRole(t *testing.T) {
	const roleName = "role-name"

	NewIdmIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.DeleteRole(roleName)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(RoleResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameDeleteRole),
			clienttest.EndpointRequestHasEndpointParameter("roleName", roleName),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(false)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodDelete),
			clienttest.ServiceEndpointHasPath("/api/v1/roles/{roleName}")).
		Test(t)
}

func TestIntegration_GetCapabilities(t *testing.T) {
	preq := paging.Request{
		Page: 0,
		Size: 10,
	}

	NewIdmIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetCapabilities(preq)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(CapabilityListResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetCapabilities),
			clienttest.EndpointRequestHasQueryParam("page", strconv.Itoa(int(preq.Page))),
			clienttest.EndpointRequestHasQueryParam("pageSize", strconv.Itoa(int(preq.Size))),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(false)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v1/roles/capabilities")).
		Test(t)
}

func TestIntegration_BatchCreateCapabilities(t *testing.T) {
	const populator = true
	const owner = "owner"
	var capabilities = []CapabilityCreateRequest{
		{
			Name: "capability-name",
		},
	}

	NewIdmIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.BatchCreateCapabilities(populator, owner, capabilities)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(CapabilityListResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameBatchCreateCapabilities),
			clienttest.EndpointRequestHasQueryParam("dbinstaller", strconv.FormatBool(populator)),
			clienttest.EndpointRequestHasQueryParam("owner", owner),
			clienttest.EndpointRequestHasBodyJsonValue("capabilities.0.name", capabilities[0].Name),
			clienttest.EndpointRequestHasBodyJsonValue("capabilities.#", float64(1)),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(false)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v1/roles/capabilities")).
		Test(t)
}

func TestIntegration_BatchUpdateCapabilities(t *testing.T) {
	const populator = true
	const owner = "owner"
	var capabilities = []CapabilityUpdateRequest{
		{
			Name: "capability-name",
		},
	}

	NewIdmIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.BatchUpdateCapabilities(populator, owner, capabilities)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(CapabilityListResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameBatchUpdateCapabilities),
			clienttest.EndpointRequestHasQueryParam("dbinstaller", strconv.FormatBool(populator)),
			clienttest.EndpointRequestHasQueryParam("owner", owner),
			clienttest.EndpointRequestHasBodyJsonValue("capabilities.0.name", capabilities[0].Name),
			clienttest.EndpointRequestHasBodyJsonValue("capabilities.#", float64(1)),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(false)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPut),
			clienttest.ServiceEndpointHasPath("/api/v1/roles/capabilities")).
		Test(t)
}

func TestIntegration_DeleteCapability(t *testing.T) {
	const populator = true
	const owner = "owner"
	const capabilityName = "capability-name"

	NewIdmIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.DeleteCapability(populator, owner, capabilityName)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(CapabilityListResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameDeleteCapability),
			clienttest.EndpointRequestHasQueryParam("dbinstaller", strconv.FormatBool(populator)),
			clienttest.EndpointRequestHasQueryParam("owner", owner),
			clienttest.EndpointRequestHasEndpointParameter("capabilityName", capabilityName),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(false)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodDelete),
			clienttest.ServiceEndpointHasPath("/api/v1/roles/capabilities/{capabilityName}")).
		Test(t)
}
