package oss

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"net/http"
)

type Integration struct {
	*integration.ExternalService
	ctx      context.Context
	endpoint integration.Endpoint
}

func (i *Integration) GetPricePlanOptions(serviceId, offerId types.UUID, options PricingOptionsRequest) (response PricingOptionsResponse, err error) {
	userContext := security.UserContextFromContext(i.ctx)
	userName := userContext.UserName
	tenantId := userContext.TenantId

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

	req, err := i.ExternalService.Request(
		i.endpoint,
		uriVariables,
		nil,
		nil,
		bodyBytes)
	if err != nil {
		return
	}

	_, _, err = i.ExternalService.Do(req, &response)

	return
}

func NewPricePlanOptionsIntegration(ctx context.Context, outboundApi OutboundApi) Api {
	externalService := integration.NewExternalService(ctx, "http", "oss")
	externalService.AddInterceptor(outboundApi.Interceptor)
	return &Integration{
		ExternalService: externalService,
		ctx:             ctx,
		endpoint: integration.Endpoint{
			Name:   OutboundApiPricingOptions,
			Method: http.MethodPost,
			Path:   outboundApi.Url,
		},
	}
}
