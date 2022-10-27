// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package auth

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	endpointNameLogin  = "login"
	endpointNameLogout = "logout"

	endpointNameGetTokenDetails = "getTokenDetails"
	endpointNameGetTokenKeys    = "getTokenKeys"

	endpointTenantHierarchyRoot        = "getTenantHierarchyRoot"
	endpointTenantHierarchyParent      = "getTenantHierarchyParent"
	endpointTenantHierarchyAncestors   = "getTenantHierarchyAncestors"
	endpointTenantHierarchyDescendants = "getTenantHierarchyDescendants"
	endpointTenantHierarchyChildren    = "getTenantHierarchyChildren"

	serviceName       = integration.ServiceNameAuth
	HeaderContentType = "Content-Type"
)

var (
	endpoints = map[string]integration.MsxServiceEndpoint{
		endpointNameLogin:  {Method: "POST", Path: "/v2/token"},
		endpointNameLogout: {Method: "GET", Path: "/v2/logout"},

		endpointNameGetTokenDetails: {Method: "POST", Path: "/v2/check_token"},
		endpointNameGetTokenKeys:    {Method: "GET", Path: "/v2/jwks"},

		endpointTenantHierarchyRoot:        {Method: "GET", Path: "/v2/tenant_hierarchy/root"},
		endpointTenantHierarchyParent:      {Method: "GET", Path: "/v2/tenant_hierarchy/parent"},
		endpointTenantHierarchyAncestors:   {Method: "GET", Path: "/v2/tenant_hierarchy/ancestors"},
		endpointTenantHierarchyDescendants: {Method: "GET", Path: "/v2/tenant_hierarchy/descendants"},
		endpointTenantHierarchyChildren:    {Method: "GET", Path: "/v2/tenant_hierarchy/children"},
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

func (i *Integration) Login(user, password string) (result *integration.MsxResponse, err error) {
	securityClientSettings, err := integration.NewSecurityClientSettings(i.Context())
	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameLogin,
		Headers: http.Header(map[string][]string{
			"Authorization":   {securityClientSettings.Authorization()},
			HeaderContentType: {httpclient.MimeTypeApplicationWwwFormUrlencoded},
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

func (i *Integration) SwitchContext(accessToken string, userId types.UUID) (result *integration.MsxResponse, err error) {
	securityClientSettings, err := integration.NewSecurityClientSettings(i.Context())
	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameLogin,
		Headers: http.Header(map[string][]string{
			"Authorization":   {securityClientSettings.Authorization()},
			HeaderContentType: {httpclient.MimeTypeApplicationWwwFormUrlencoded},
		}),
		Body: []byte(url.Values(map[string][]string{
			"grant_type":     {"urn:cisco:nfv:oauth:grant-type:switch-user"},
			"switch_user_id": {userId.String()},
			"access_token":   {accessToken},
		}).Encode()),
		Payload:      new(LoginResponse),
		ErrorPayload: new(integration.OAuthErrorDTO),
		NoToken:      true,
	})
}

func (i *Integration) Logout() (result *integration.MsxResponse, err error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameLogout,
	})
}

func (i *Integration) GetTokenKeys() (keys JsonWebKeys, response *integration.MsxResponse, err error) {
	response, err = i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetTokenKeys,
		Payload:      &keys,
		NoToken:      true,
		Headers: http.Header{
			"Accept": []string{
				"application/jwk+json",
				"application/json",
			},
		},
	})
	return
}

func (i *Integration) GetTokenDetails(noDetails bool) (*integration.MsxResponse, error) {
	securityClientSettings, err := integration.NewSecurityClientSettings(i.Context())
	if err != nil {
		return nil, err
	}

	var headers = make(http.Header)
	headers.Set("Authorization", securityClientSettings.Authorization())
	headers.Set(HeaderContentType, httpclient.MimeTypeApplicationWwwFormUrlencoded)
	headers.Set("Accept", httpclient.MimeTypeApplicationJson)

	var body = make(url.Values)
	userContext := security.UserContextFromContext(i.Context())
	body.Set("token", userContext.Token)
	if noDetails {
		body.Set("no_details", "true")
	}
	var bodyBytes = []byte(body.Encode())

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:   endpointNameGetTokenDetails,
		Headers:        headers,
		Body:           bodyBytes,
		ExpectEnvelope: false,
		NoToken:        true,
		Payload:        new(TokenDetails),
		ErrorPayload:   new(integration.ErrorDTO3),
	})
}

func (i *Integration) GetTenantHierarchyRoot() (*integration.MsxResponse, error) {
	msxEndpointRequest, err := i.buildTenantHierarchyMsxEndpointRequest(endpointTenantHierarchyRoot)
	if err != nil {
		return nil, err
	}

	return i.Execute(msxEndpointRequest)
}

func (i *Integration) GetTenantHierarchyParent(tenantId types.UUID) (*integration.MsxResponse, error) {
	msxEndpointRequest, err := i.buildTenantHierarchyMsxEndpointRequest(endpointTenantHierarchyParent)
	if err != nil {
		return nil, err
	}

	qp := url.Values{}
	qp.Set("tenantId", tenantId.String())
	msxEndpointRequest.QueryParameters = qp

	return i.Execute(msxEndpointRequest)
}

func (i *Integration) GetTenantHierarchyAncestors(tenantId types.UUID) (*integration.MsxResponse, []types.UUID, error) {
	request, err := i.buildTenantHierarchyMsxEndpointRequest(endpointTenantHierarchyAncestors)
	if err != nil {
		return nil, nil, err
	}

	request.QueryParameters = url.Values{
		"tenantId": []string{tenantId.String()},
	}

	response, err := i.Execute(request)
	if err != nil {
		return nil, nil, err
	}

	var result []types.UUID
	if err = json.Unmarshal(response.Body, &result); err != nil {
		return nil, nil, err
	}

	return response, result, err
}

func (i *Integration) GetTenantHierarchyDescendants(tenantId types.UUID) (*integration.MsxResponse, []types.UUID, error) {
	request, err := i.buildTenantHierarchyMsxEndpointRequest(endpointTenantHierarchyDescendants)
	if err != nil {
		return nil, nil, err
	}

	request.QueryParameters = url.Values{
		"tenantId": []string{tenantId.String()},
	}

	response, err := i.Execute(request)
	if err != nil {
		return nil, nil, err
	}

	var result []types.UUID
	if err = json.Unmarshal(response.Body, &result); err != nil {
		return nil, nil, err
	}

	return response, result, err
}

func (i *Integration) GetTenantHierarchyChildren(tenantId types.UUID) (*integration.MsxResponse, []types.UUID, error) {
	request, err := i.buildTenantHierarchyMsxEndpointRequest(endpointTenantHierarchyChildren)
	if err != nil {
		return nil, nil, err
	}

	request.QueryParameters = url.Values{
		"tenantId": []string{tenantId.String()},
	}

	response, err := i.Execute(request)
	if err != nil {
		return nil, nil, err
	}

	var result []types.UUID
	if err = json.Unmarshal(response.Body, &result); err != nil {
		return nil, nil, err
	}

	return response, result, err
}

func (i *Integration) buildTenantHierarchyMsxEndpointRequest(endpointName string) (*integration.MsxEndpointRequest, error) {

	securityClientSettings, err := integration.NewSecurityClientSettings(i.Context())
	if err != nil {
		return nil, err
	}

	var headers = make(http.Header)
	headers.Set("Authorization", securityClientSettings.Authorization())
	headers.Set(HeaderContentType, httpclient.MimeTypeApplicationWwwFormUrlencoded)
	headers.Set("Accept", httpclient.MimeTypeApplicationJson)

	return &integration.MsxEndpointRequest{
		EndpointName:   endpointName,
		Headers:        headers,
		ExpectEnvelope: false,
		NoToken:        true,
		Payload:        nil,
		ErrorPayload:   new(integration.ErrorDTO3),
	}, nil
}
