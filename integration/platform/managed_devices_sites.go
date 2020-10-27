package platform

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	platform "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
	"net/http"
)

type ManagedDevicesSitesApi interface {
	GetManagedDeviceSites(ctx context.Context, localVarOptionals *platform.GetManagedDeviceSitesOpts) ([]platform.ManagedDeviceSiteDto, *http.Response, error)
}

func NewManagedDeviceSitesApiService(ctx context.Context) *platform.ManagedDevicesSitesApiService {
	cfg := newPlatformClientConfigFromContext(ctx, integration.ServiceNameManagedDevice)
	return platform.NewAPIClient(cfg).ManagedDevicesSitesApi
}
