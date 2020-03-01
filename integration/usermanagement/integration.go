package usermanagement

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

const (
	endpointNameGetAdminHealth = "getAdminHealth"

	endpointNameLogin  = "login"
	endpointNameLogout = "logout"

	endpointNameGetMyCapabilities   = "getMyCapabilities"
	endpointNameGetUserCapabilities = "getUserCapabilities"

	endpointNameGetTenantIds   = "getTenantIds"
	endpointNameGetMyTenants   = "getMyTenants"
	endpointNameGetUserTenants = "getUserTenants"

	endpointNameGetTenantById   = "getTenantById"
	endpointNameGetTenantByName = "getTenantByName"

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

	serviceName = integration.ServiceNameUserManagement
)

var (
	logger    = log.NewLogger("msx.integration.usermanagement")
	endpoints = map[string]integration.MsxServiceEndpoint{
		endpointNameGetAdminHealth: {Method: "GET", Path: "/admin/health"},

		endpointNameLogin:  {Method: "POST", Path: "/v2/token"},
		endpointNameLogout: {Method: "GET", Path: "/v2/logout"},

		endpointNameGetMyCapabilities:   {Method: "GET", Path: "/api/v1/users/capabilities"},
		endpointNameGetUserCapabilities: {Method: "GET", Path: "/api/v1/users/{{.userId}}/capabilities"},

		endpointNameGetTenantIds:   {Method: "GET", Path: "/api/v1/tenantIds"},
		endpointNameGetMyTenants:   {Method: "GET", Path: "/api/v1/users/tenants"},
		endpointNameGetUserTenants: {Method: "GET", Path: "/api/v1/users/{{.userId}}/tenants"},

		endpointNameGetTenantById:   {Method: "GET", Path: "/api/v3/tenants/{{.tenantId}}"},
		endpointNameGetTenantByName: {Method: "GET", Path: "/api/v1/tenants/{{.tenantName}}"},

		endpointNameGetSystemSecrets:      {Method: "GET", Path: "/api/v2/secrets/scope/{{.scope}}"},
		endpointNameAddSystemSecrets:      {Method: "POST", Path: "/api/v2/secrets/scope/{{.scope}}"},
		endpointNameReplaceSystemSecrets:  {Method: "PUT", Path: "/api/v2/secrets/scope/{{.scope}}"},
		endpointNameRemoveSystemSecrets:   {Method: "DELETE", Path: "/api/v2/secrets/scope/{{.scope}}"},
		endpointNameEncryptSystemSecrets:  {Method: "POST", Path: "/api/v2/secrets/scope/{{.scope}}/encrypt"},
		endpointNameGenerateSystemSecrets: {Method: "POST", Path: "/api/v2/secrets/scope/{{.scope}}/generate"},

		endpointNameGetTenantSecrets:      {Method: "GET", Path: "/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}"},
		endpointNameAddTenantSecrets:      {Method: "POST", Path: "/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}"},
		endpointNameReplaceTenantSecrets:  {Method: "PUT", Path: "/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}"},
		endpointNameRemoveTenantSecrets:   {Method: "DELETE", Path: "/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}"},
		endpointNameEncryptTenantSecrets:  {Method: "POST", Path: "/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}/encrypt"},
		endpointNameGenerateTenantSecrets: {Method: "POST", Path: "/api/v2/secrets/tenant/{{.tenantId}}/scope/{{.scope}}/generate"},
	}
)

func NewIntegration(ctx context.Context) (Api, error) {
	return &Integration{
		MsxService: integration.NewMsxService(ctx, serviceName, endpoints),
	}, nil
}

type Integration struct {
	*integration.MsxService
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

func (i *Integration) Login(user, password string) (result *integration.MsxResponse, err error) {
	securityClientSettings, err := integration.NewSecurityClientSettings(i.Context())
	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameLogin,
		Headers: http.Header(map[string][]string{
			"Authorization": {securityClientSettings.Authorization()},
			"Content-Type":  {"application/x-www-form-urlencoded"},
		}),
		Body: []byte(url.Values(map[string][]string{
			"grant_type": {"password"},
			"username":   {user},
			"password":   {password},
		}).Encode()),
		Payload:      new(LoginResponse),
		ErrorPayload: new(integration.OAuthErrorDTO),
		NoToken:      true,
	})
}

func (i *Integration) Logout() (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameLogout,
	})
}

func (i *Integration) GetMyCapabilities() (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameGetMyCapabilities,
		Payload:      &UserCapabilityListResponse{},
		ErrorPayload: &ErrorDTO{},
	})
}

func (i *Integration) GetUserCapabilities(userId string) (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameGetUserCapabilities,
		EndpointParameters: map[string]string{
			"userId": userId,
		},
		Payload:      &UserCapabilityListResponse{},
		ErrorPayload: &ErrorDTO{},
	})
}

func (i *Integration) GetMyTenants() (result *integration.MsxResponse, err error) {
	result, err = i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameGetMyTenants,
		Payload:      new(TenantListResponse),
		ErrorPayload: &ErrorDTO{},
	})

	if result != nil && result.StatusCode == 404 {
		logger.Info("Converting 404 on list to 200")
		result = &integration.MsxResponse{
			StatusCode: 200,
			Status:     "200 OK",
			Headers:    result.Headers,
			Envelope: &integration.MsxEnvelope{
				Success:    true,
				HttpStatus: "OK",
				Payload:    new(TenantListResponse),
			},
		}
		result.Payload = result.Envelope.Payload
		result.Body, _ = json.Marshal(result.Envelope)
		result.BodyString = string(result.Body)
		err = nil
	}

	return result, err
}

func (i *Integration) GetTenantIds() (result *integration.MsxResponse, err error) {
	result, err = i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameGetTenantIds,
		Payload:      new(TenantIdList),
		ErrorPayload: &ErrorDTO2{},
	})

	return result, err
}

func (i *Integration) GetUserTenants(userId string) (result *integration.MsxResponse, err error) {
	result, err = i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameGetUserTenants,
		EndpointParameters: map[string]string{
			"userId": userId,
		},
		Payload:      new(TenantListResponse),
		ErrorPayload: &ErrorDTO{},
	})

	if result != nil && result.StatusCode == 404 {
		logger.Info("Converting 404 on list to 200")
		result = &integration.MsxResponse{
			StatusCode: 200,
			Status:     "200 OK",
			Headers:    result.Headers,
			Envelope: &integration.MsxEnvelope{
				Success:    true,
				HttpStatus: "OK",
				Payload:    new(TenantListResponse),
			},
		}
		result.Payload = result.Envelope.Payload
		result.Body, _ = json.Marshal(result.Envelope)
		result.BodyString = string(result.Body)
		err = nil
	}

	return result, err
}

func (i *Integration) GetTenantById(tenantId string) (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameGetTenantById,
		EndpointParameters: map[string]string{
			"tenantId": tenantId,
		},
		ExpectEnvelope: true,
		Payload:        new(TenantResponse),
	})
}

func (i *Integration) GetTenantByName(tenantName string) (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameGetTenantByName,
		EndpointParameters: map[string]string{
			"tenantName": tenantName,
		},
		Payload: new(TenantResponse),
	})
}

func (i *Integration) GetSystemSecrets(scope string) (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameGetSystemSecrets,
		EndpointParameters: map[string]string{
			"scope": scope,
		},
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) AddSystemSecrets(scope string, secrets map[string]string) (result *integration.MsxResponse, err error) {
	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(secrets); err != nil {
		return nil, errors.Wrap(err, "Failed to serialize body")
	}
	return i.Execute(&integration.MsxRequest{
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
		return nil, errors.Wrap(err, "Failed to serialize body")
	}
	return i.Execute(&integration.MsxRequest{
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
		return nil, errors.Wrap(err, "Failed to serialize body")
	}

	return i.Execute(&integration.MsxRequest{
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
	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameRemoveSystemSecrets,
		EndpointParameters: map[string]string{
			"scope": scope,
		},
		ExpectEnvelope: true,
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
		return nil, errors.Wrap(err, "Failed to serialize body")
	}

	return i.Execute(&integration.MsxRequest{
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
	return i.Execute(&integration.MsxRequest{
		EndpointName: endpointNameGetTenantSecrets,
		EndpointParameters: map[string]string{
			"tenantId": tenantId,
			"scope":    scope,
		},
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}

func (i *Integration) AddTenantSecrets(tenantId, scope string, secrets map[string]string) (result *integration.MsxResponse, err error) {
	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(secrets); err != nil {
		return nil, errors.Wrap(err, "Failed to serialize body")
	}
	return i.Execute(&integration.MsxRequest{
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
		return nil, errors.Wrap(err, "Failed to serialize body")
	}
	return i.Execute(&integration.MsxRequest{
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
		return nil, errors.Wrap(err, "Failed to serialize body")
	}

	return i.Execute(&integration.MsxRequest{
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
	return i.Execute(&integration.MsxRequest{
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
		return nil, errors.Wrap(err, "Failed to serialize body")
	}

	return i.Execute(&integration.MsxRequest{
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
