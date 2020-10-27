package platform

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	platform "cto-github.cisco.com/NFV-BU/msx-platform-go-client"
	"net/http"
)

type UsersApi interface {
	CreateUser(ctx context.Context, userCreate platform.UserCreate) (platform.User, *http.Response, error)
	DeleteUser(ctx context.Context, id string) (*http.Response, error)
	GetCurrentUser(ctx context.Context) (platform.User, *http.Response, error)
	GetUser(ctx context.Context, id string) (platform.User, *http.Response, error)
	GetUsersPage(ctx context.Context, page int32, pageSize int32, localVarOptionals *platform.GetUsersPageOpts) (platform.UsersPage, *http.Response, error)
	UpdateUser(ctx context.Context, id string, userUpdate platform.UserUpdate) (platform.User, *http.Response, error)
	UpdateUserPassword(ctx context.Context, updatePassword platform.UpdatePassword) (*http.Response, error)
}

func NewUsersApiService(ctx context.Context) *platform.UsersApiService {
	cfg := newPlatformClientConfigFromContext(ctx, integration.ServiceNameUserManagement)
	return platform.NewAPIClient(cfg).UsersApi
}
