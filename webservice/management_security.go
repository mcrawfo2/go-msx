package webservice

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/rbac"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

const (
	configRootManagementSecurity = "management.security"
)

type ManagementSecurityConfig struct {
	EnabledByDefault bool     `config:"default=true"`
	Permissions      []string `config:"default=IS_API_ADMIN"`
	Roles            []string `config:"default=ROLE_CLIENT"`
	Endpoint         map[string]EndpointConfig
}

func (s ManagementSecurityConfig) EndpointSecurityEnabled(endpoint string) bool {
	if endpointOverride, ok := s.Endpoint[endpoint]; !ok {
		return s.EnabledByDefault
	} else if strings.ToLower(endpointOverride.Enabled) == "false" {
		return false
	} else {
		return true
	}
}

type EndpointConfig struct {
	Enabled string `config:"default="`
}

func NewManagementSecurityConfig(ctx context.Context) (*ManagementSecurityConfig, error) {
	var cfg ManagementSecurityConfig
	if err := config.FromContext(ctx).Populate(&cfg, configRootManagementSecurity); err != nil {
		return nil, err
	}

	return &cfg, nil
}

type ManagementSecurityFilter struct {
	cfg *ManagementSecurityConfig
}

func (s ManagementSecurityFilter) Filter(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	ctx := request.Request.Context()

	if err := s.roles(ctx); err != nil {
		WriteError(request, response, http.StatusUnauthorized, err)
		return
	}

	if err := rbac.HasPermission(ctx, s.cfg.Permissions); err != nil {
		WriteError(request, response, http.StatusUnauthorized, err)
		return
	}

	chain.ProcessFilter(request, response)
}

func (s ManagementSecurityFilter) roles(ctx context.Context) error {
	userContext := security.UserContextFromContext(ctx)
	for _, role := range s.cfg.Roles {
		if types.StringStack(userContext.Authorities).Contains(role) {
			return nil
		}
	}

	return errors.New("Token does not contain required roles")
}

func NewManagementSecurityFilter(cfg *ManagementSecurityConfig) restful.FilterFunction {
	return ManagementSecurityFilter{cfg: cfg}.Filter
}
