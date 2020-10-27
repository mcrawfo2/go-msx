package platform

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	platform "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
	"net/http"
)

type TenantsApi interface {
	CreateTenant(ctx context.Context, tenantCreate platform.TenantCreate) (platform.Tenant, *http.Response, error)
	DeleteTenant(ctx context.Context, id string) (*http.Response, error)
	GetTenant(ctx context.Context, id string) (platform.Tenant, *http.Response, error)
	GetTenantsList(ctx context.Context, ids []string, localVarOptionals *platform.GetTenantsListOpts) ([]platform.Tenant, *http.Response, error)
	GetTenantsPage(ctx context.Context, page int32, pageSize int32, localVarOptionals *platform.GetTenantsPageOpts) (platform.TenantsPage, *http.Response, error)
	UpdateTenant(ctx context.Context, id string, tenantUpdate platform.TenantUpdate) (platform.Tenant, *http.Response, error)
}

func NewTenantsApiService(ctx context.Context) *platform.TenantsApiService {
	cfg := newPlatformClientConfigFromContext(ctx, integration.ServiceNameUserManagement)
	return platform.NewAPIClient(cfg).TenantsApi
}
