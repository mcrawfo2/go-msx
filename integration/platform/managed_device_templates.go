package platform

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	platform "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
	"net/http"
)

type ManagedDeviceTemplatesApi interface {
	CreateManagedDeviceTemplate(ctx context.Context, managedDeviceTemplateDto platform.ManagedDeviceTemplateDto) (platform.ManagedDeviceTemplateDto, *http.Response, error)
	GetManagedDeviceTemplatesList(ctx context.Context) ([]platform.ManagedDeviceTemplateDto, *http.Response, error)
}

func NewManagedDeviceTemplatesApiService(ctx context.Context) *platform.ManagedDeviceTemplatesApiService {
	cfg := newPlatformClientConfigFromContext(ctx, integration.ServiceNameManagedDevice)
	return platform.NewAPIClient(cfg).ManagedDeviceTemplatesApi
}
