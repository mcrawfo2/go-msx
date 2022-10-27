// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package secrets

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	integrationClientTest "cto-github.cisco.com/NFV-BU/go-msx/integration/testhelpers/clienttest"
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
			"remoteservice.secretsservice.service": "secretsservice",
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

type SecretsIntegrationTest struct {
	*integrationClientTest.EndpointTest
}

func NewSecretsIntegrationTest() *SecretsIntegrationTest {
	return &SecretsIntegrationTest{
		EndpointTest: new(integrationClientTest.EndpointTest).WithEndpoints(endpoints),
	}
}

type ManageCall func(t *testing.T, api Api) (*integration.MsxResponse, error)
type AuthCall func(t *testing.T, api Api) (*integration.MsxResponse, []types.UUID, error)

func (m *SecretsIntegrationTest) WithCall(call ManageCall) *SecretsIntegrationTest {
	m.EndpointTest.WithCall(func(t *testing.T, executor integration.MsxContextServiceExecutor) (*integration.MsxResponse, error) {
		return call(t, NewIntegrationWithExecutor(executor))
	})
	return m
}

func (m *SecretsIntegrationTest) WithMultiTenantResultCall(call AuthCall) *SecretsIntegrationTest {
	m.EndpointTest.WithMultiTenantResultCall(func(t *testing.T, executor integration.MsxContextServiceExecutor) (*integration.MsxResponse, []types.UUID, error) {
		return call(t, NewIntegrationWithExecutor(executor))
	})
	return m
}

func TestIntegration_GetSystemSecrets(t *testing.T) {
	const scope = "scope"

	NewSecretsIntegrationTest().
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

	NewSecretsIntegrationTest().
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

	NewSecretsIntegrationTest().
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

	NewSecretsIntegrationTest().
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

	NewSecretsIntegrationTest().
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

func TestIntegration_RemoveSystemSecretsPermanent(t *testing.T) {
	const scope = "scope"

	NewSecretsIntegrationTest().
		WithCall(func(t *testing.T, api Api) (*integration.MsxResponse, error) {
			permanent := true
			return api.RemoveSystemSecretsPermanent(scope, &permanent)
		}).
		WithResponseStatus(http.StatusOK).
		WithResponseEnvelope().
		WithRequestPredicates(
			clienttest.EndpointRequestHasName(endpointNameRemoveSystemSecrets),
			clienttest.EndpointRequestHasEndpointParameter("scope", scope),
			clienttest.EndpointRequestHasQueryParam("permanent", "true"),
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

	NewSecretsIntegrationTest().
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

	NewSecretsIntegrationTest().
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

	NewSecretsIntegrationTest().
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

	NewSecretsIntegrationTest().
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

	NewSecretsIntegrationTest().
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

	NewSecretsIntegrationTest().
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

	NewSecretsIntegrationTest().
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

func TestIntegration_GetSecretPolicy(t *testing.T) {
	const policyName = "policy-name"

	NewSecretsIntegrationTest().
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

	NewSecretsIntegrationTest().
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

	NewSecretsIntegrationTest().
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
