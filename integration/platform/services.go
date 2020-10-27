package platform

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	platform "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
	"net/http"
)

type ServicesApi interface {
	DeleteService(ctx context.Context, id string) (*http.Response, error)
	GetService(ctx context.Context, id string) (platform.Service, *http.Response, error)
	GetServicesPage(ctx context.Context, page int32, pageSize int32, localVarOptionals *platform.GetServicesPageOpts) (platform.ServicesPage, *http.Response, error)
	SubmitOrder(ctx context.Context, productId string, offerId string, legacyServiceOrder platform.LegacyServiceOrder) (platform.LegacyServiceOrderResponse, *http.Response, error)
	UpdateOrder(ctx context.Context, productId string, offerId string, legacyServiceOrder platform.LegacyServiceOrder) (platform.LegacyServiceOrderResponse, *http.Response, error)
}

func NewServicesApiService(ctx context.Context) *platform.ServicesApiService {
	cfg := newPlatformClientConfigFromContext(ctx, integration.ServiceNameManage)
	return platform.NewAPIClient(cfg).ServicesApi
}
