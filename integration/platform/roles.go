package platform

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	platform "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
	"net/http"
)

type RolesApi interface {
	GetRoleByName(ctx context.Context, name string) (platform.Role, *http.Response, error)
	GetRolesList(ctx context.Context, ids []string) ([]platform.Role, *http.Response, error)
}

func NewRolesApiService(ctx context.Context) *platform.RolesApiService {
	cfg := newPlatformClientConfigFromContext(ctx, integration.ServiceNameUserManagement)
	return platform.NewAPIClient(cfg).RolesApi
}
