// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package secrets

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"encoding/json"
	"github.com/pkg/errors"
	"net/url"
	"strconv"
)

const (
	endpointNameGetSystemSecrets      = "getSystemSecrets"
	endpointNameEncryptSystemSecrets  = "encryptSystemSecrets"
	endpointNameAddSystemSecrets      = "addSystemSecrets"
	endpointNameReplaceSystemSecrets  = "replaceSystemSecrets"
	endpointNameRemoveSystemSecrets   = "removeSystemSecrets"
	endpointNameGenerateSystemSecrets = "generateSystemSecrets"

	endpointNameGetTenantSecrets      = "getTenantSecrets"
	endpointNameEncryptTenantSecrets  = "encryptTenantSecrets"
	endpointNameAddTenantSecrets      = "addTenantSecrets"
	endpointNameReplaceTenantSecrets  = "replaceTenantSecrets"
	endpointNameRemoveTenantSecrets   = "removeTenantSecrets"
	endpointNameGenerateTenantSecrets = "generateTenantSecrets"

	endpointNameGetSecretPolicy   = "getSecretPolicy"
	endpointNameSetSecretPolicy   = "setSecretPolicy"
	endpointNameUnsetSecretPolicy = "unsetSecretPolicy"

	serviceName                       = integration.ServiceNameSecrets
	ErrorMessageFailedToSerializeBody = "Failed to serialize body"
	pathApiV2SecretsScope             = "/api/v2/secrets/scope/{{.scope}}"
	pathApiV2SecretsTenantScope       = "/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}"
	pathApiV2SecretsPolicy            = "/api/v2/secrets/policy/{policyName}"
)

var (
	endpoints = map[string]integration.MsxServiceEndpoint{
		endpointNameGetSystemSecrets:      {Method: "GET", Path: pathApiV2SecretsScope},
		endpointNameAddSystemSecrets:      {Method: "POST", Path: pathApiV2SecretsScope},
		endpointNameReplaceSystemSecrets:  {Method: "PUT", Path: pathApiV2SecretsScope},
		endpointNameRemoveSystemSecrets:   {Method: "DELETE", Path: pathApiV2SecretsScope},
		endpointNameEncryptSystemSecrets:  {Method: "POST", Path: "/api/v2/secrets/scope/{{.scope}}/encrypt"},
		endpointNameGenerateSystemSecrets: {Method: "POST", Path: "/api/v2/secrets/scope/{{.scope}}/generate"},

		endpointNameGetTenantSecrets:      {Method: "GET", Path: pathApiV2SecretsTenantScope},
		endpointNameAddTenantSecrets:      {Method: "POST", Path: pathApiV2SecretsTenantScope},
		endpointNameReplaceTenantSecrets:  {Method: "PUT", Path: pathApiV2SecretsTenantScope},
		endpointNameRemoveTenantSecrets:   {Method: "DELETE", Path: pathApiV2SecretsTenantScope},
		endpointNameEncryptTenantSecrets:  {Method: "POST", Path: "/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}/encrypt"},
		endpointNameGenerateTenantSecrets: {Method: "POST", Path: "/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}/generate"},

		endpointNameGetSecretPolicy:   {Method: "GET", Path: pathApiV2SecretsPolicy},
		endpointNameSetSecretPolicy:   {Method: "PUT", Path: pathApiV2SecretsPolicy},
		endpointNameUnsetSecretPolicy: {Method: "DELETE", Path: pathApiV2SecretsPolicy},
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

func (i *Integration) GetSystemSecrets(scope string) (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetSystemSecrets,
		EndpointParameters: map[string]string{
			"scope": scope,
		},
		Payload:        new(SecretsResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) AddSystemSecrets(scope string, secrets map[string]string) (result *integration.MsxResponse, err error) {
	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(secrets); err != nil {
		return nil, errors.Wrap(err, ErrorMessageFailedToSerializeBody)
	}
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameAddSystemSecrets,
		EndpointParameters: map[string]string{
			"scope": scope,
		},
		Body:           bodyBytes,
		ExpectEnvelope: true,
	})
}

func (i *Integration) ReplaceSystemSecrets(scope string, secrets map[string]string) (result *integration.MsxResponse, err error) {
	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(secrets); err != nil {
		return nil, errors.Wrap(err, ErrorMessageFailedToSerializeBody)
	}
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameReplaceSystemSecrets,
		EndpointParameters: map[string]string{
			"scope": scope,
		},
		Body:           bodyBytes,
		ExpectEnvelope: true,
	})
}

func (i *Integration) EncryptSystemSecrets(scope string, names []string, encrypt EncryptSecretsDTO) (result *integration.MsxResponse, err error) {
	body := &GetSecretRequestDTO{
		Names:   names,
		Encrypt: encrypt,
	}

	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(body); err != nil {
		return nil, errors.Wrap(err, ErrorMessageFailedToSerializeBody)
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameEncryptSystemSecrets,
		EndpointParameters: map[string]string{
			"scope": scope,
		},
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) RemoveSystemSecrets(scope string) (result *integration.MsxResponse, err error) {
	return i.RemoveSystemSecretsPermanent(scope, nil)
}

func (i *Integration) RemoveSystemSecretsPermanent(scope string, permanent *bool) (result *integration.MsxResponse, err error) {
	var qp url.Values
	if permanent != nil {
		qp = make(url.Values)
		qp.Set("permanent", strconv.FormatBool(*permanent))
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameRemoveSystemSecrets,
		EndpointParameters: map[string]string{
			"scope": scope,
		},
		QueryParameters: qp,
		ExpectEnvelope:  true,
	})
}

func (i *Integration) GenerateSystemSecrets(scope string, names []string, save bool) (result *integration.MsxResponse, err error) {
	body := &GenerateSecretRequestDTO{
		Names:   names,
		Save:    save,
		Encrypt: nil,
	}

	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(body); err != nil {
		return nil, errors.Wrap(err, ErrorMessageFailedToSerializeBody)
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGenerateSystemSecrets,
		EndpointParameters: map[string]string{
			"scope": scope,
		},
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetTenantSecrets(tenantId, scope string) (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetTenantSecrets,
		EndpointParameters: map[string]string{
			"tenantId": tenantId,
			"scope":    scope,
		},
		Payload:        new(SecretsResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) AddTenantSecrets(tenantId, scope string, secrets map[string]string) (result *integration.MsxResponse, err error) {
	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(secrets); err != nil {
		return nil, errors.Wrap(err, ErrorMessageFailedToSerializeBody)
	}
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameAddTenantSecrets,
		EndpointParameters: map[string]string{
			"tenantId": tenantId,
			"scope":    scope,
		},
		Body:           bodyBytes,
		ExpectEnvelope: true,
	})
}

func (i *Integration) ReplaceTenantSecrets(tenantId, scope string, secrets map[string]string) (result *integration.MsxResponse, err error) {
	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(secrets); err != nil {
		return nil, errors.Wrap(err, ErrorMessageFailedToSerializeBody)
	}
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameReplaceTenantSecrets,
		EndpointParameters: map[string]string{
			"tenantId": tenantId,
			"scope":    scope,
		},
		Body:           bodyBytes,
		ExpectEnvelope: true,
	})
}

func (i *Integration) EncryptTenantSecrets(tenantId, scope string, names []string, encrypt EncryptSecretsDTO) (result *integration.MsxResponse, err error) {
	body := &GetSecretRequestDTO{
		Names:   names,
		Encrypt: encrypt,
	}

	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(body); err != nil {
		return nil, errors.Wrap(err, ErrorMessageFailedToSerializeBody)
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameEncryptTenantSecrets,
		EndpointParameters: map[string]string{
			"tenantId": tenantId,
			"scope":    scope,
		},
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) RemoveTenantSecrets(tenantId, scope string) (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameRemoveTenantSecrets,
		EndpointParameters: map[string]string{
			"tenantId": tenantId,
			"scope":    scope,
		},
		ExpectEnvelope: true,
	})
}

func (i *Integration) GenerateTenantSecrets(tenantId, scope string, names []string, save bool) (result *integration.MsxResponse, err error) {
	body := &GenerateSecretRequestDTO{
		Names:   names,
		Save:    save,
		Encrypt: nil,
	}

	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(body); err != nil {
		return nil, errors.Wrap(err, ErrorMessageFailedToSerializeBody)
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGenerateTenantSecrets,
		EndpointParameters: map[string]string{
			"tenantId": tenantId,
			"scope":    scope,
		},
		Body:           bodyBytes,
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) GetSecretPolicy(name string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetSecretPolicy,
		EndpointParameters: map[string]string{
			"policyName": name,
		},
		Payload:        new(SecretPolicyResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) StoreSecretPolicy(name string, policy SecretPolicySetRequest) (result *integration.MsxResponse, err error) {
	var body struct {
		SecretPolicySetRequest
		Name string `json:"name"`
	}

	body.SecretPolicySetRequest = policy
	body.Name = name

	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(body); err != nil {
		return nil, errors.Wrap(err, ErrorMessageFailedToSerializeBody)
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameSetSecretPolicy,
		EndpointParameters: map[string]string{
			"policyName": name,
		},
		Body:           bodyBytes,
		Payload:        new(SecretPolicyResponse),
		ExpectEnvelope: true,
	})
}

func (i *Integration) DeleteSecretPolicy(name string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameUnsetSecretPolicy,
		EndpointParameters: map[string]string{
			"policyName": name,
		},
		ExpectEnvelope: true,
	})
}
