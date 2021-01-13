package usermanagement

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/clienttest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

func TestNewIntegration(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want Api
	}{
		{
			name: "NonExisting",
			args: args{
				ctx: context.Background(),
			},
			want: &Integration{
				MsxContextServiceExecutor: integration.NewMsxService(context.Background(), serviceName, endpoints),
			},
		},
		{
			name: "Existing",
			args: args{
				ctx: ContextWithIntegration(context.Background(), &Integration{}),
			},
			want: &Integration{},
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

type UserManagementIntegrationTest struct {
	*EndpointTest
}

func NewUserManagementIntegrationTest() *UserManagementIntegrationTest {
	return &UserManagementIntegrationTest{
		EndpointTest: new(EndpointTest).WithEndpoints(endpoints),
	}
}

type ManageCall func(t *testing.T, api Api) (*integration.MsxResponse, error)

func (m *UserManagementIntegrationTest) WithCall(call ManageCall) *UserManagementIntegrationTest {
	m.EndpointTest.WithCall(func(t *testing.T, executor integration.MsxContextServiceExecutor) (*integration.MsxResponse, error) {
		return call(t, NewIntegrationWithExecutor(executor))
	})
	return m
}

func TestIntegration_GetAdminHealth(t *testing.T) {
	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetAdminHealth()
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(&integration.HealthDTO{
			Status: "Up",
		}).
		WithRequestPredicate(clienttest.EndpointRequestHasName(endpointNameGetAdminHealth)).
		WithRequestPredicate(clienttest.EndpointRequestHasToken(false)).
		WithEndpointPredicate(clienttest.ServiceEndpointHasMethod(http.MethodGet)).
		WithEndpointPredicate(clienttest.ServiceEndpointHasPath("/admin/health")).
		Test(t)
}

func TestIntegration_Login(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), nil)
	securityClientSettings, _ := integration.NewSecurityClientSettings(ctx)

	NewUserManagementIntegrationTest().
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

func TestIntegration_Logout(t *testing.T) {
	NewUserManagementIntegrationTest().
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

func TestIntegration_IsTokenActive(t *testing.T) {
	NewUserManagementIntegrationTest().
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

func TestIntegration_GetTokenDetails(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), nil)
	securityClientSettings, _ := integration.NewSecurityClientSettings(ctx)

	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetTokenDetails(false)
		}).
		WithInjector(func(ctx context.Context) context.Context {
			ctx = configtest.ContextWithNewInMemoryConfig(ctx, nil)
			ctx = security.ContextWithUserContext(ctx, &security.UserContext{
				UserName: "username",
				Token: "token-value",
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

func TestIntegration_GetMyCapabilities(t *testing.T) {
	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetMyCapabilities()
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(UserCapabilityListResponse)).
		WithRequestPredicate(clienttest.EndpointRequestHasName(endpointNameGetMyCapabilities)).
		WithRequestPredicate(clienttest.EndpointRequestHasToken(true)).
		WithEndpointPredicate(clienttest.ServiceEndpointHasMethod(http.MethodGet)).
		WithEndpointPredicate(clienttest.ServiceEndpointHasPath("/api/v1/users/capabilities")).
		Test(t)
}

func TestIntegration_GetUserCapabilities(t *testing.T) {
	const userId = "user-id"

	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetUserCapabilities(userId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(UserCapabilityListResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetUserCapabilities),
			clienttest.EndpointRequestHasEndpointParameter("userId", userId),
			clienttest.EndpointRequestHasToken(true)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v1/users/{{.userId}}/capabilities")).
		Test(t)
}

func TestIntegration_GetMyUserId(t *testing.T) {
	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetMyUserId()
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(UserIdResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetMyUserId),
			clienttest.EndpointRequestHasToken(true)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v1/currentuser")).
		Test(t)
}

func TestIntegration_GetMyPersonalInfo(t *testing.T) {
	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetMyPersonalInfo()
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(UserIdResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetMyPersonalInfo),
			clienttest.EndpointRequestHasToken(true)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v1/personalinfo")).
		Test(t)
}

func TestIntegration_GetMyProvider(t *testing.T) {
	NewUserManagementIntegrationTest().
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

	NewUserManagementIntegrationTest().
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

	NewUserManagementIntegrationTest().
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

func TestIntegration_GetTenantIds(t *testing.T) {
	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetTenantIds()
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(TenantIdList)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetTenantIds),
			clienttest.EndpointRequestHasToken(true)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v1/tenantids")).
		Test(t)
}

func TestIntegration_GetMyTenants(t *testing.T) {
	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetMyTenants()
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(TenantListResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetMyTenants),
			clienttest.EndpointRequestHasToken(true)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v1/users/tenants")).
		Test(t)
}

func TestIntegration_GetUserTenants(t *testing.T) {
	const userId = "user-id"

	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetUserTenants(userId)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(TenantListResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetUserTenants),
			clienttest.EndpointRequestHasEndpointParameter("userId", userId),
			clienttest.EndpointRequestHasToken(true)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v1/users/{{.userId}}/tenants")).
		Test(t)
}

func TestIntegration_GetUserById(t *testing.T) {
	const userId = "user-id"

	NewUserManagementIntegrationTest().
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

func TestIntegration_GetTenantById(t *testing.T) {
	const tenantId = "tenant-id"

	NewUserManagementIntegrationTest().
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

func TestIntegration_GetTenantByName(t *testing.T) {
	const tenantName = "tenant-name"

	NewUserManagementIntegrationTest().
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

func TestIntegration_GetSystemSecrets(t *testing.T) {
	const scope = "scope"

	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetSystemSecrets(scope)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(SecretsResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetSystemSecrets),
			clienttest.EndpointRequestHasEndpointParameter("scope", scope),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(true)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v2/secrets/scope/{{.scope}}")).
		Test(t)
}

func TestIntegration_AddSystemSecrets(t *testing.T) {
	const scope = "scope"
	var secrets = map[string]string{
		"secret-key-1": "secret-value-1",
	}

	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.AddSystemSecrets(scope, secrets)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(SecretsResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameAddSystemSecrets),
			clienttest.EndpointRequestHasEndpointParameter("scope", scope),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasBodyJsonValue("secret-key-1", "secret-value-1")).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v2/secrets/scope/{{.scope}}")).
		Test(t)
}

func TestIntegration_ReplaceSystemSecrets(t *testing.T) {
	const scope = "scope"
	var secrets = map[string]string{
		"secret-key-1": "secret-value-1",
	}

	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.ReplaceSystemSecrets(scope, secrets)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(SecretsResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameReplaceSystemSecrets),
			clienttest.EndpointRequestHasEndpointParameter("scope", scope),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasBodyJsonValue("secret-key-1", "secret-value-1")).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPut),
			clienttest.ServiceEndpointHasPath("/api/v2/secrets/scope/{{.scope}}")).
		Test(t)
}

func TestIntegration_EncryptSystemSecrets(t *testing.T) {
	const scope = "scope=scope-value"
	var names = []string{"secret-key-1"}
	var encrypt = EncryptSecretsDTO{
		Scope:  map[string]string{"scope": "scope-value"},
		Name:   "name",
		Method: "nso",
		SecretNames: []string{
			"secret-key-2",
		},
	}

	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.EncryptSystemSecrets(scope, names, encrypt)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(SecretsResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameEncryptSystemSecrets),
			clienttest.EndpointRequestHasEndpointParameter("scope", scope),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasBodyJsonValue("names.0", "secret-key-1"),
			clienttest.EndpointRequestHasBodyJsonValue("names.#", float64(1)),
			clienttest.EndpointRequestHasBodyJsonValue("encrypt.name", "name")).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v2/secrets/scope/{{.scope}}/encrypt")).
		Test(t)
}

func TestIntegration_RemoveSystemSecrets(t *testing.T) {
	const scope = "scope"

	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.RemoveSystemSecrets(scope)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameRemoveSystemSecrets),
			clienttest.EndpointRequestHasEndpointParameter("scope", scope),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(true)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodDelete),
			clienttest.ServiceEndpointHasPath("/api/v2/secrets/scope/{{.scope}}")).
		Test(t)
}

func TestIntegration_GenerateSystemSecrets(t *testing.T) {
	const scope = "scope-key=scope-value"
	const save = true
	var names = []string{"secret-key-1"}

	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GenerateSystemSecrets(scope, names, save)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(SecretsResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGenerateSystemSecrets),
			clienttest.EndpointRequestHasEndpointParameter("scope", scope),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasBodyJsonValue("names.0", "secret-key-1"),
			clienttest.EndpointRequestHasBodyJsonValue("names.#", float64(1)),
			clienttest.EndpointRequestHasBodyJsonValue("save", save)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v2/secrets/scope/{{.scope}}/generate")).
		Test(t)
}

////

func TestIntegration_GetTenantSecrets(t *testing.T) {
	const scope = "scope"
	const tenantId = "tenant-id"

	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetTenantSecrets(tenantId, scope)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(SecretsResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetTenantSecrets),
			clienttest.EndpointRequestHasEndpointParameter("scope", scope),
			clienttest.EndpointRequestHasEndpointParameter("tenantId", tenantId),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(true)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}")).
		Test(t)
}

func TestIntegration_AddTenantSecrets(t *testing.T) {
	const scope = "scope"
	const tenantId = "tenant-id"
	var secrets = map[string]string{
		"secret-key-1": "secret-value-1",
	}

	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.AddTenantSecrets(tenantId, scope, secrets)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(SecretsResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameAddTenantSecrets),
			clienttest.EndpointRequestHasEndpointParameter("scope", scope),
			clienttest.EndpointRequestHasEndpointParameter("tenantId", tenantId),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasBodyJsonValue("secret-key-1", "secret-value-1")).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}")).
		Test(t)
}

func TestIntegration_ReplaceTenantSecrets(t *testing.T) {
	const scope = "scope"
	const tenantId = "tenant-id"
	var secrets = map[string]string{
		"secret-key-1": "secret-value-1",
	}

	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.ReplaceTenantSecrets(tenantId, scope, secrets)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(SecretsResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameReplaceTenantSecrets),
			clienttest.EndpointRequestHasEndpointParameter("scope", scope),
			clienttest.EndpointRequestHasEndpointParameter("tenantId", tenantId),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasBodyJsonValue("secret-key-1", "secret-value-1")).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPut),
			clienttest.ServiceEndpointHasPath("/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}")).
		Test(t)
}

func TestIntegration_EncryptTenantSecrets(t *testing.T) {
	const scope = "scope=scope-value"
	const tenantId = "tenant-id"
	var names = []string{"secret-key-1"}
	var encrypt = EncryptSecretsDTO{
		Scope:  map[string]string{"scope": "scope-value"},
		Name:   "name",
		Method: "nso",
		SecretNames: []string{
			"secret-key-2",
		},
	}

	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.EncryptTenantSecrets(tenantId, scope, names, encrypt)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(SecretsResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameEncryptTenantSecrets),
			clienttest.EndpointRequestHasEndpointParameter("scope", scope),
			clienttest.EndpointRequestHasEndpointParameter("tenantId", tenantId),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasBodyJsonValue("names.0", "secret-key-1"),
			clienttest.EndpointRequestHasBodyJsonValue("names.#", float64(1)),
			clienttest.EndpointRequestHasBodyJsonValue("encrypt.name", "name")).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}/encrypt")).
		Test(t)
}

func TestIntegration_RemoveTenantSecrets(t *testing.T) {
	const scope = "scope"
	const tenantId = "tenant-id"

	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.RemoveTenantSecrets(tenantId, scope)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameRemoveTenantSecrets),
			clienttest.EndpointRequestHasEndpointParameter("scope", scope),
			clienttest.EndpointRequestHasEndpointParameter("tenantId", tenantId),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(true)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodDelete),
			clienttest.ServiceEndpointHasPath("/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}")).
		Test(t)
}

func TestIntegration_GenerateTenantSecrets(t *testing.T) {
	const scope = "scope-key=scope-value"
	const tenantId = "tenant-id"
	const save = true
	var names = []string{"secret-key-1"}

	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GenerateTenantSecrets(tenantId, scope, names, save)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithResponsePayload(new(SecretsResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGenerateTenantSecrets),
			clienttest.EndpointRequestHasEndpointParameter("scope", scope),
			clienttest.EndpointRequestHasEndpointParameter("tenantId", tenantId),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(true),
			clienttest.EndpointRequestHasBodyJsonValue("names.0", "secret-key-1"),
			clienttest.EndpointRequestHasBodyJsonValue("names.#", float64(1)),
			clienttest.EndpointRequestHasBodyJsonValue("save", save)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPost),
			clienttest.ServiceEndpointHasPath("/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}/generate")).
		Test(t)
}

func TestIntegration_GetRoles(t *testing.T) {
	const resolve = true
	preq := paging.Request{
		Page: 0,
		Size: 10,
	}

	NewUserManagementIntegrationTest().
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

	NewUserManagementIntegrationTest().
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

	NewUserManagementIntegrationTest().
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

	NewUserManagementIntegrationTest().
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

	NewUserManagementIntegrationTest().
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

	NewUserManagementIntegrationTest().
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

	NewUserManagementIntegrationTest().
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

	NewUserManagementIntegrationTest().
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

func TestIntegration_GetSecretPolicy(t *testing.T) {
	const policyName = "policy-name"

	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.GetSecretPolicy(policyName)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(SecretPolicyResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameGetSecretPolicy),
			clienttest.EndpointRequestHasEndpointParameter("policyName", policyName),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(true)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodGet),
			clienttest.ServiceEndpointHasPath("/api/v2/secrets/policy/{policyName}")).
		Test(t)
}

func TestIntegration_StoreSecretPolicy(t *testing.T) {
	const policyName = "policy-name"
	request := SecretPolicySetRequest{
		AgingRule: AgingRule{
			Enabled: true,
		},
	}

	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.StoreSecretPolicy(policyName, request)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(SecretPolicyResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameSetSecretPolicy),
			clienttest.EndpointRequestHasEndpointParameter("policyName", policyName),
			clienttest.EndpointRequestHasBodyJsonValue("agingRule.enabled", true),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(true)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodPut),
			clienttest.ServiceEndpointHasPath("/api/v2/secrets/policy/{policyName}")).
		Test(t)
}

func TestIntegration_DeleteSecretPolicy(t *testing.T) {
	const policyName = "policy-name"

	NewUserManagementIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			return api.DeleteSecretPolicy(policyName)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponsePayload(new(SecretPolicyResponse)).
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameUnsetSecretPolicy),
			clienttest.EndpointRequestHasEndpointParameter("policyName", policyName),
			clienttest.EndpointRequestHasToken(true),
			clienttest.EndpointRequestHasExpectEnvelope(true)).
		WithEndpointPredicates(
			clienttest.ServiceEndpointHasMethod(http.MethodDelete),
			clienttest.ServiceEndpointHasPath("/api/v2/secrets/policy/{policyName}")).
		Test(t)
}
