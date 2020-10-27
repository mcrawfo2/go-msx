package platform

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	platform "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
	"net/http"
)

type DevicesApi interface {
	AttachDeviceTemplates(ctx context.Context, id string, deviceTemplateAttachRequest platform.DeviceTemplateAttachRequest) ([]platform.DeviceTemplateHistory, *http.Response, error)
	CreateDevice(ctx context.Context, deviceCreate platform.DeviceCreate) (platform.Device, *http.Response, error)
	DeleteDevice(ctx context.Context, id string) (*http.Response, error)
	DetachDeviceTemplate(ctx context.Context, id string, templateId string) ([]platform.DeviceTemplateHistory, *http.Response, error)
	DetachDeviceTemplates(ctx context.Context, id string) ([]platform.DeviceTemplateHistory, *http.Response, error)
	GetDevice(ctx context.Context, id string) (platform.Device, *http.Response, error)
	GetDeviceConfig(ctx context.Context, id string) (string, *http.Response, error)
	GetDeviceTemplateHistory(ctx context.Context, id string) ([]platform.DeviceTemplateHistory, *http.Response, error)
	GetDevicesPage(ctx context.Context, page int32, pageSize int32, localVarOptionals *platform.GetDevicesPageOpts) (platform.DevicesPage, *http.Response, error)
	RedeployDevice(ctx context.Context, id string) (*http.Response, error)
	UpdateDeviceTemplates(ctx context.Context, id string, deviceTemplateUpdateRequest platform.DeviceTemplateUpdateRequest) ([]platform.DeviceTemplateHistory, *http.Response, error)
}

func NewDevicesApiService(ctx context.Context) *platform.DevicesApiService {
	cfg := newPlatformClientConfigFromContext(ctx, integration.ServiceNameManage)
	return platform.NewAPIClient(cfg).DevicesApi
}
