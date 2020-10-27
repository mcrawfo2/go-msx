package platform

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	platform "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
	"net/http"
)

type ProductsApi interface {
	CreateProduct(ctx context.Context, productCreate platform.ProductCreate) (platform.Product, *http.Response, error)
	DeleteProduct(ctx context.Context, id string) (*http.Response, error)
	GetProduct(ctx context.Context, id string) (platform.Product, *http.Response, error)
	GetProductsCount(ctx context.Context) (int64, *http.Response, error)
	GetProductsPage(ctx context.Context, page int32, pageSize int32) (platform.ProductsPage, *http.Response, error)
	UpdateProduct(ctx context.Context, id string, productUpdate platform.ProductUpdate) (platform.Product, *http.Response, error)
}

func NewProductsApiService(ctx context.Context) *platform.ProductsApiService {
	cfg := newPlatformClientConfigFromContext(ctx, integration.ServiceNameConsume)
	return platform.NewAPIClient(cfg).ProductsApi
}
