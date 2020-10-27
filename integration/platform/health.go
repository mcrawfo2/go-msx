package platform

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	platform "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
	"net/http"
)

type HealthApi interface {
	GetDevicesHealthList(ctx context.Context, ids []string) ([]platform.ResourceHealth, *http.Response, error)
	GetServicesHealthList(ctx context.Context, ids []string) ([]platform.ResourceHealth, *http.Response, error)
}

func NewHealthApiService(ctx context.Context) *platform.HealthApiService {
	cfg := newPlatformClientConfigFromContext(ctx, integration.ServiceNameMonitor)
	return platform.NewAPIClient(cfg).HealthApi
}
