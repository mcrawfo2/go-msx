package authprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/bmatcuk/doublestar"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

const (
	configRootAuthenticationProvider = "security.resources.patterns"
	oAuthRoleClient                  = "ROLE_CLIENT"
	oAuthScopeRead                   = "read"
	oAuthScopeWrite                  = "write"
)

var ErrUserForbidden = webservice.NewStatusError(
	errors.New("User does not have required identity"),
	http.StatusForbidden)

var ErrUserExpired = webservice.NewStatusError(
	errors.New("User token has expired"),
	http.StatusForbidden)

type ResourcePatternAuthenticationConfig struct {
	Blacklist []string `config:"default=/api/**;/admin;/admin/**"`
	Whitelist []string `config:"default=/admin/health;/admin/info;/admin/alive"`
}

type ResourcePatternAuthenticationProvider struct {
	cfg *ResourcePatternAuthenticationConfig
}

func (f *ResourcePatternAuthenticationProvider) requiresAuthentication(path string) (bool, error) {
	blacklisted := false

	// Check blacklist for path match
	for _, rule := range f.cfg.Blacklist {
		if matches, err := doublestar.Match(rule, path); err != nil {
			return true, err
		} else if matches {
			blacklisted = true
			break
		}
	}

	// Blacklist did not contain a match
	if !blacklisted {
		return false, nil
	}

	// Check whitelist for path match
	for _, rule := range f.cfg.Whitelist {
		if matches, err := doublestar.Match(rule, path); err != nil {
			return true, err
		} else if matches {
			return false, nil
		}
	}

	// Whitelist did not match, but blacklist did
	return true, nil
}

func (f *ResourcePatternAuthenticationProvider) Authenticate(req *restful.Request) (err error) {
	serverContextPath := webservice.WebServerFromContext(req.Request.Context()).ContextPath()
	path := strings.TrimPrefix(req.Request.URL.Path, serverContextPath)

	var required bool
	if required, err = f.requiresAuthentication(path); err != nil || !required {
		return err
	}

	userContext := security.UserContextFromContext(req.Request.Context())
	if !types.StringStack(userContext.Authorities).Contains(oAuthRoleClient) ||
		!types.StringStack(userContext.Scopes).Contains(oAuthScopeRead) ||
		!types.StringStack(userContext.Scopes).Contains(oAuthScopeWrite) {
		return ErrUserForbidden
	}

	if userContext.Token != "" {
		active, err := security.IsTokenActive(req.Request.Context())
		if err != nil {
			return err
		}
		if !active {
			return ErrUserExpired
		}
	}

	return nil
}

func NewResourcePatternAuthenticationProvider(cfg *ResourcePatternAuthenticationConfig) *ResourcePatternAuthenticationProvider {
	return &ResourcePatternAuthenticationProvider{cfg: cfg}
}

func NewResourcePatternAuthenticationProviderFromConfig(cfg *config.Config) (*ResourcePatternAuthenticationProvider, error) {
	jwtSecurityProviderConfig := new(ResourcePatternAuthenticationConfig)
	if err := cfg.Populate(jwtSecurityProviderConfig, configRootAuthenticationProvider); err != nil {
		return nil, err
	}

	return NewResourcePatternAuthenticationProvider(jwtSecurityProviderConfig), nil
}

func RegisterAuthenticationProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server == nil {
		// Server disabled
		return nil
	}

	cfg := config.MustFromContext(ctx)
	jwtSecurityProvider, err := NewResourcePatternAuthenticationProviderFromConfig(cfg)
	if err != nil {
		return errors.Wrap(err, "Failed to create JWT web security provider")
	}

	server.SetAuthenticationProvider(jwtSecurityProvider)
	return nil
}
