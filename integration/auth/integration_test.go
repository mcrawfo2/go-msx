// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package auth

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	integrationClientTest "cto-github.cisco.com/NFV-BU/go-msx/integration/testhelpers/clienttest"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/clienttest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"net/http"
	"reflect"
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

type AuthIntegrationTest struct {
	*integrationClientTest.EndpointTest
}

func NewAuthIntegrationTest() *AuthIntegrationTest {
	return &AuthIntegrationTest{
		EndpointTest: new(integrationClientTest.EndpointTest).WithEndpoints(endpoints),
	}
}

type ManageCall func(t *testing.T, api Api) (*integration.MsxResponse, error)
type AuthCall func(t *testing.T, api Api) (*integration.MsxResponse, []types.UUID, error)

func (m *AuthIntegrationTest) WithCall(call ManageCall) *AuthIntegrationTest {
	m.EndpointTest.WithCall(func(t *testing.T, executor integration.MsxContextServiceExecutor) (*integration.MsxResponse, error) {
		return call(t, NewIntegrationWithExecutor(executor))
	})
	return m
}

func (m *AuthIntegrationTest) WithMultiTenantResultCall(call AuthCall) *AuthIntegrationTest {
	m.EndpointTest.WithMultiTenantResultCall(func(t *testing.T, executor integration.MsxContextServiceExecutor) (*integration.MsxResponse, []types.UUID, error) {
		return call(t, NewIntegrationWithExecutor(executor))
	})
	return m
}

func TestIntegration_Login(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), nil)
	securityClientSettings, _ := integration.NewSecurityClientSettings(ctx)

	NewAuthIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.Login("username", "password")
		}).
		WithInjector(func(ctx context.Context) context.Context {
			return configtest.ContextWithNewInMemoryConfig(ctx, nil)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(LoginResponse)).
		WithRequestPredicate(clienttest.EndpointRequestHasHeader("Authorization", securityClientSettings.Authorization())).
		WithRequestPredicate(clienttest.EndpointRequestHasHeader("Content-Type", "application/x-www-form-urlencoded")).
		WithRequestPredicate(clienttest.EndpointRequestHasName(endpointNameLogin)).
		WithRequestPredicate(clienttest.EndpointRequestHasToken(false)).
		WithRequestPredicate(clienttest.EndpointRequestHasBodySubstring("grant_type=password")).
		WithRequestPredicate(clienttest.EndpointRequestHasBodySubstring("username=username")).
		WithRequestPredicate(clienttest.EndpointRequestHasBodySubstring("password=password")).
		WithEndpointPredicate(clienttest.ServiceEndpointHasMethod(http.MethodPost)).
		WithEndpointPredicate(clienttest.ServiceEndpointHasPath("/v2/token")).
		Test(t)
}

func TestIntegration_SwitchContext(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), nil)
	securityClientSettings, _ := integration.NewSecurityClientSettings(ctx)
	userId := types.MustNewUUID()
	token := "user-token"

	NewAuthIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.SwitchContext(token, userId)
		}).
		WithInjector(func(ctx context.Context) context.Context {
			return configtest.ContextWithNewInMemoryConfig(ctx, nil)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(LoginResponse)).
		WithRequestPredicate(clienttest.EndpointRequestHasHeader("Authorization", securityClientSettings.Authorization())).
		WithRequestPredicate(clienttest.EndpointRequestHasHeader("Content-Type", "application/x-www-form-urlencoded")).
		WithRequestPredicate(clienttest.EndpointRequestHasName(endpointNameLogin)).
		WithRequestPredicate(clienttest.EndpointRequestHasToken(false)).
		WithRequestPredicate(clienttest.EndpointRequestHasBodySubstring("grant_type=urn%3Acisco%3Anfv%3Aoauth%3Agrant-type%3Aswitch-user")).
		WithRequestPredicate(clienttest.EndpointRequestHasBodySubstring("switch_user_id=" + userId.String())).
		WithRequestPredicate(clienttest.EndpointRequestHasBodySubstring("access_token=" + token)).
		WithEndpointPredicate(clienttest.ServiceEndpointHasMethod(http.MethodPost)).
		WithEndpointPredicate(clienttest.ServiceEndpointHasPath("/v2/token")).
		Test(t)
}

func TestIntegration_Logout(t *testing.T) {
	NewAuthIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.Logout()
		}).
		WithResponseStatus(http.StatusOK).
		WithRequestPredicate(clienttest.EndpointRequestHasName(endpointNameLogout)).
		WithRequestPredicate(clienttest.EndpointRequestHasToken(true)).
		WithEndpointPredicate(clienttest.ServiceEndpointHasMethod(http.MethodGet)).
		WithEndpointPredicate(clienttest.ServiceEndpointHasPath("/v2/logout")).
		Test(t)
}

func TestIntegration_GetTokenDetails(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), nil)
	securityClientSettings, _ := integration.NewSecurityClientSettings(ctx)

	NewAuthIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetTokenDetails(false)
		}).
		WithInjector(func(ctx context.Context) context.Context {
			ctx = configtest.ContextWithNewInMemoryConfig(ctx, nil)
			ctx = security.ContextWithUserContext(ctx, &security.UserContext{
				UserName: "username",
				Token:    "token-value",
			})
			return ctx
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(TokenDetails)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetTokenDetails),
			clienttest.EndpointRequestHasToken(false),
			clienttest.EndpointRequestHasExpectEnvelope(false),
			clienttest.EndpointRequestHasHeader("Authorization", securityClientSettings.Authorization()),
			clienttest.EndpointRequestHasHeader("Content-Type", httpclient.MimeTypeApplicationWwwFormUrlencoded),
			clienttest.EndpointRequestHasHeader("Accept", httpclient.MimeTypeApplicationJson),
			clienttest.EndpointRequestHasBodySubstring("token=token-value"),
		).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/v2/check_token"),
		).
		Test(t)
}

func TestIntegration_GetTenantHierarchyRoot(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), nil)
	securityClientSettings, _ := integration.NewSecurityClientSettings(ctx)

	NewAuthIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetTenantHierarchyRoot()
		}).
		WithInjector(func(ctx context.Context) context.Context {
			return configtest.ContextWithNewInMemoryConfig(ctx, nil)
		}).
		WithResponseStatus(http.StatusOK).
		WithRequestPredicate(clienttest.EndpointRequestHasHeader("Authorization", securityClientSettings.Authorization())).
		WithRequestPredicate(clienttest.EndpointRequestHasHeader("Content-Type", httpclient.MimeTypeApplicationWwwFormUrlencoded)).
		WithRequestPredicate(clienttest.EndpointRequestHasHeader("Accept", httpclient.MimeTypeApplicationJson)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointTenantHierarchyRoot),
			clienttest.EndpointRequestHasToken(false),
			clienttest.EndpointRequestHasExpectEnvelope(false)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/v2/tenant_hierarchy/root")).
		Test(t)
}

func TestIntegration_GetTenantHierarchyParent(t *testing.T) {
	var tenantId, _ = types.NewUUID()
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), nil)
	securityClientSettings, _ := integration.NewSecurityClientSettings(ctx)

	NewAuthIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetTenantHierarchyParent(tenantId)
		}).
		WithInjector(func(ctx context.Context) context.Context {
			return configtest.ContextWithNewInMemoryConfig(ctx, nil)
		}).
		WithResponseStatus(http.StatusOK).
		WithRequestPredicate(clienttest.EndpointRequestHasHeader("Authorization", securityClientSettings.Authorization())).
		WithRequestPredicate(clienttest.EndpointRequestHasHeader("Content-Type", httpclient.MimeTypeApplicationWwwFormUrlencoded)).
		WithRequestPredicate(clienttest.EndpointRequestHasHeader("Accept", httpclient.MimeTypeApplicationJson)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointTenantHierarchyParent),
			clienttest.EndpointRequestHasQueryParam("tenantId", tenantId.String()),
			clienttest.EndpointRequestHasToken(false),
			clienttest.EndpointRequestHasExpectEnvelope(false)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/v2/tenant_hierarchy/parent")).
		Test(t)
}

func TestIntegration_GetTenantHierarchyChildren(t *testing.T) {
	tenantId, _ := types.NewUUID()
	child1, _ := types.NewUUID()
	child2, _ := types.NewUUID()

	children := []types.UUID{child1, child2}

	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), nil)
	securityClientSettings, _ := integration.NewSecurityClientSettings(ctx)

	NewAuthIntegrationTest().
		WithMultiTenantResultCall(func(t *testing.T, api Api) (*integration.MsxResponse, []types.UUID, error) {
			return api.GetTenantHierarchyChildren(tenantId)
		}).
		WithInjector(func(ctx context.Context) context.Context {
			return configtest.ContextWithNewInMemoryConfig(ctx, nil)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(children).
		WithRequestPredicate(clienttest.EndpointRequestHasHeader("Authorization", securityClientSettings.Authorization())).
		WithRequestPredicate(clienttest.EndpointRequestHasHeader("Content-Type", httpclient.MimeTypeApplicationWwwFormUrlencoded)).
		WithRequestPredicate(clienttest.EndpointRequestHasHeader("Accept", httpclient.MimeTypeApplicationJson)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointTenantHierarchyChildren),
			clienttest.EndpointRequestHasQueryParam("tenantId", tenantId.String()),
			clienttest.EndpointRequestHasToken(false),
			clienttest.EndpointRequestHasExpectEnvelope(false)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/v2/tenant_hierarchy/children")).
		WithTenants(children).
		TestMultiTenantResult(t)
}

func TestIntegration_GetTenantHierarchyDescendants(t *testing.T) {
	tenantId, _ := types.NewUUID()
	desc1, _ := types.NewUUID()
	desc2, _ := types.NewUUID()

	descendants := []types.UUID{desc1, desc2}

	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), nil)
	securityClientSettings, _ := integration.NewSecurityClientSettings(ctx)

	NewAuthIntegrationTest().
		WithMultiTenantResultCall(func(t *testing.T, api Api) (*integration.MsxResponse, []types.UUID, error) {
			return api.GetTenantHierarchyDescendants(tenantId)
		}).
		WithInjector(func(ctx context.Context) context.Context {
			return configtest.ContextWithNewInMemoryConfig(ctx, nil)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(descendants).
		WithRequestPredicate(clienttest.EndpointRequestHasHeader("Authorization", securityClientSettings.Authorization())).
		WithRequestPredicate(clienttest.EndpointRequestHasHeader("Content-Type", httpclient.MimeTypeApplicationWwwFormUrlencoded)).
		WithRequestPredicate(clienttest.EndpointRequestHasHeader("Accept", httpclient.MimeTypeApplicationJson)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointTenantHierarchyDescendants),
			clienttest.EndpointRequestHasQueryParam("tenantId", tenantId.String()),
			clienttest.EndpointRequestHasToken(false),
			clienttest.EndpointRequestHasExpectEnvelope(false)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/v2/tenant_hierarchy/descendants")).
		WithTenants(descendants).
		TestMultiTenantResult(t)
}

func TestIntegration_GetTenantHierarchyAncestors(t *testing.T) {
	tenantId, _ := types.NewUUID()
	ancestor1, _ := types.NewUUID()
	ancestor2, _ := types.NewUUID()

	ancestors := []types.UUID{ancestor1, ancestor2}

	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), nil)
	securityClientSettings, _ := integration.NewSecurityClientSettings(ctx)

	NewAuthIntegrationTest().
		WithMultiTenantResultCall(func(t *testing.T, api Api) (*integration.MsxResponse, []types.UUID, error) {
			return api.GetTenantHierarchyAncestors(tenantId)
		}).
		WithInjector(func(ctx context.Context) context.Context {
			return configtest.ContextWithNewInMemoryConfig(ctx, nil)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(ancestors).
		WithRequestPredicate(clienttest.EndpointRequestHasHeader("Authorization", securityClientSettings.Authorization())).
		WithRequestPredicate(clienttest.EndpointRequestHasHeader("Content-Type", httpclient.MimeTypeApplicationWwwFormUrlencoded)).
		WithRequestPredicate(clienttest.EndpointRequestHasHeader("Accept", httpclient.MimeTypeApplicationJson)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointTenantHierarchyAncestors),
			clienttest.EndpointRequestHasQueryParam("tenantId", tenantId.String()),
			clienttest.EndpointRequestHasToken(false),
			clienttest.EndpointRequestHasExpectEnvelope(false)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/v2/tenant_hierarchy/ancestors")).
		WithTenants(ancestors).
		TestMultiTenantResult(t)
}
