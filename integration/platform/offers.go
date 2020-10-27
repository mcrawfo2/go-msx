package platform

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	platform "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
	"net/http"
)

type OffersApi interface {
	CreateOffer(ctx context.Context, offerCreate platform.OfferCreate) (platform.Offer, *http.Response, error)
	DeleteOffer(ctx context.Context, id string) (*http.Response, error)
	GetOffer(ctx context.Context, id string) (platform.Offer, *http.Response, error)
	GetOffersCount(ctx context.Context, localVarOptionals *platform.GetOffersCountOpts) (int64, *http.Response, error)
	GetOffersPage(ctx context.Context, page int32, pageSize int32, localVarOptionals *platform.GetOffersPageOpts) (platform.OffersPage, *http.Response, error)
	UpdateOffer(ctx context.Context, id string, offerUpdate platform.OfferUpdate) (platform.Offer, *http.Response, error)
}

func NewOffersApiService(ctx context.Context) *platform.OffersApiService {
	cfg := newPlatformClientConfigFromContext(ctx, integration.ServiceNameConsume)
	return platform.NewAPIClient(cfg).OffersApi
}
