package platform

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	platform "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
	"net/http"
)

type DeviceTemplatesApi interface {
	CreateDeviceTemplate(ctx context.Context, deviceTemplateCreate platform.DeviceTemplateCreate) (platform.DeviceTemplate, *http.Response, error)
	DeleteDeviceTemplate(ctx context.Context, id string) (*http.Response, error)
	GetDeviceTemplate(ctx context.Context, id string) (platform.DeviceTemplate, *http.Response, error)
	GetDeviceTemplatesList(ctx context.Context, localVarOptionals *platform.GetDeviceTemplatesListOpts) ([]platform.DeviceTemplate, *http.Response, error)
	UpdateDeviceTemplateAccess(ctx context.Context, id string, deviceTemplateAccess platform.DeviceTemplateAccess) (platform.DeviceTemplateAccessResponse, *http.Response, error)
}

func NewDeviceTemplatesApiService(ctx context.Context) *platform.DeviceTemplatesApiService {
	cfg := newPlatformClientConfigFromContext(ctx, integration.ServiceNameManage)
	return platform.NewAPIClient(cfg).DeviceTemplatesApi
}
