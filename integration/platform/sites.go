package platform

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	platform "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
	"net/http"
)

type SitesApi interface {
	AddDeviceToSite(ctx context.Context, id string, deviceId string, localVarOptionals *platform.AddDeviceToSiteOpts) (platform.Site, *http.Response, error)
	CreateSite(ctx context.Context, siteCreate platform.SiteCreate) (platform.Site, *http.Response, error)
	DeleteSite(ctx context.Context, id string) (*http.Response, error)
	GetSite(ctx context.Context, id string, localVarOptionals *platform.GetSiteOpts) (platform.Site, *http.Response, error)
	GetSitesPage(ctx context.Context, page int32, pageSize int32, localVarOptionals *platform.GetSitesPageOpts) (platform.SitesPage, *http.Response, error)
	RemoveDeviceFromSite(ctx context.Context, id string, deviceId string) (*http.Response, error)
	UpdateSite(ctx context.Context, id string, siteUpdate platform.SiteUpdate, localVarOptionals *platform.UpdateSiteOpts) (platform.Site, *http.Response, error)
}

func NewSitesApiService(ctx context.Context) *platform.SitesApiService {
	cfg := newPlatformClientConfigFromContext(ctx, integration.ServiceNameManage)
	return platform.NewAPIClient(cfg).SitesApi
}
