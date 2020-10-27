package platform

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	platform "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
	"net/http"
)

type SecurityApi interface {
	GetAccessToken(ctx context.Context, authorization string, grantType string, localVarOptionals *platform.GetAccessTokenOpts) (platform.AccessToken, *http.Response, error)
}

func NewSecurityApiService(ctx context.Context) *platform.SecurityApiService {
	cfg := newSecurityClientConfigFromContext(ctx, integration.ServiceNameUserManagement)
	return platform.NewAPIClient(cfg).SecurityApi
}
