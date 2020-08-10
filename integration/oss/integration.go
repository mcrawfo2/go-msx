package oss

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"github.com/pkg/errors"
)

var ErrEndpointNotFound = errors.New("Endpoint not defined")

type Integration struct {
	ctx       context.Context
	services  map[string]*integration.ExternalService
	endpoints map[string]integration.Endpoint
}

func (i *Integration) GetAccessibleServices() (response ServicesResponse, err error) {
	service, endpoint, err := i.getServiceAndEndpoint(OutboundApiServiceAccess)
	if err != nil {
		return
	}

	tenantId, userName := i.getTenantIdAndUserIdFromContext()

	uriVariables := map[string]string{
		"tenantId": tenantId.String(),
		"userId":   userName,
	}

	req, err := service.Request(
		endpoint,
		uriVariables,
		nil,
		nil,
		nil)

	if err != nil {
		return
	}

	_, _, err = service.Do(req, &response)
	return
}

func (i *Integration) GetPricePlanOptions(serviceId, offerId types.UUID, options PricingOptionsRequest) (response PricingOptionsResponse, err error) {
	service, endpoint, err := i.getServiceAndEndpoint(OutboundApiPricingOptions)
	if err != nil {
		return
	}

	tenantId, userName := i.getTenantIdAndUserIdFromContext()

	uriVariables := map[string]string{
		"serviceId": serviceId.String(),
		"offerId":   offerId.String(),
		"tenantId":  tenantId.String(),
		"userId":    userName,
	}

	bodyBytes, err := json.Marshal(options)
	if err != nil {
		return
	}

	req, err := service.Request(
		endpoint,
		uriVariables,
		nil,
		nil,
		bodyBytes)
	if err != nil {
		return
	}

	_, _, err = service.Do(req, &response)

	return
}

func (i *Integration) GetAllowedValues(serviceId, propertyName string) (response AllowedValuesResponse, err error) {
	service, endpoint, err := i.getServiceAndEndpoint(OutboundApiAllowedValues)
	if err != nil {
		return
	}

	tenantId, userName := i.getTenantIdAndUserIdFromContext()
	uriVariables := map[string]string{
		"serviceId":    serviceId,
		"propertyName": propertyName,
		"tenantId":     tenantId.String(),
		"userId":       userName,
	}
	req, err := service.Request(
		endpoint,
		uriVariables,
		nil,
		nil,
		nil)
	if err != nil {
		return
	}

	_, _, err = service.Do(req, &response)

	return
}

func (i *Integration) getTenantIdAndUserIdFromContext() (tenantId types.UUID, userName string) {
	userContext := security.UserContextFromContext(i.ctx)
	userName = userContext.UserName
	tenantId = userContext.TenantId
	return
}

func (i *Integration) getServiceAndEndpoint(apiName string) (*integration.ExternalService, integration.Endpoint, error) {
	service, ok := i.services[apiName]
	if !ok {
		return nil, integration.Endpoint{}, ErrEndpointNotFound
	}
	endpoint := i.endpoints[apiName]
	return service, endpoint, nil
}

func NewIntegration(ctx context.Context, outboundApi OutboundApi) Api {
	integrationInstance := IntegrationFromContext(ctx)
	if integrationInstance == nil {
		integrationInstance = NewOssIntegration(ctx, []OutboundApi{outboundApi})
	}
	return integrationInstance
}

func NewOssIntegration(ctx context.Context, outboundApis []OutboundApi) Api {
	integrationInstance := IntegrationFromContext(ctx)
	if integrationInstance == nil {
		services := map[string]*integration.ExternalService{}
		endpoints := map[string]integration.Endpoint{}

		for _, outboundApi := range outboundApis {
			apiName := outboundApi.ApiName

			externalService := integration.NewExternalService(ctx, "http", "oss")
			externalService.AddInterceptor(outboundApi.Interceptor)
			services[apiName] = externalService

			endpoint := integration.Endpoint{
				Name:   outboundApi.ApiName,
				Method: outboundApi.HttpMethod,
				Path:   outboundApi.Url,
			}
			endpoints[apiName] = endpoint
		}

		integrationInstance = &Integration{
			ctx:       ctx,
			services:  services,
			endpoints: endpoints,
		}
	}
	return integrationInstance
}
