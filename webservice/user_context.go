package webservice

import (
	"context"
	"crypto/x509"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/vault"
	"encoding/base64"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"net/http/httptest"
)

const (
	jwtClaimUserName = "user_name"
	jwtClaimScope    = "scope"
	jwtClaimTenantId = "tenantId"
	jwtClaimRoles    = "roles"

	configRootUserContextFilter = "server.jwt"
)

type UserContextFilterConfig struct {
	KeyPath string `config:"default=secret/phi_pnp"`
	KeyName string `config:"default=key"`
}

type UserContextFilter struct {
	cfg *UserContextFilterConfig
}

func (f *UserContextFilter) Filter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	// 0. Retrieve the context
	ctx := req.Request.Context()
	if ctx == nil {
		logger.Error("Failed to retrieve context from request")
		return
	}

	middleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: f.jwtSigningKey(ctx),
		SigningMethod:       jwt.SigningMethodRS256,
		Debug:               false,
	})

	// 1. Validate the token signature
	recorder := httptest.NewRecorder()
	if err := middleware.CheckJWT(recorder, req.Request); err != nil {
		logger.WithError(err).Error("JWT validation failed")
		err := WriteErrorEnvelope(req, resp, 401, errors.Wrap(err, "JWT validation failed"))
		if err != nil {
			logger.WithError(err).Error("Failed to write error to response")
		}
		return
	}

	// 2. Validate the session
	tokenString, err := jwtmiddleware.FromAuthHeader(req.Request)
	if err != nil {
		logger.WithError(err).Error("Failed to retrieve token from request")
		return
	}

	// TODO: Validate the session
	// if !auth.Validate(j.consul, token) {
	//	logger.Error("Token is expired")
	//	_ = resp.WriteErrorString(http.StatusUnauthorized,  "Token not Valid")
	//}

	// 3. Inject the security context
	jwtParser := &jwt.Parser{SkipClaimsValidation: true}
	jwtClaims := jwt.MapClaims{}
	_, _, err = jwtParser.ParseUnverified(tokenString, jwtClaims)
	if err != nil {
		logger.WithError(err).Error("Failed to parse JWT")
	}

	userContext := &security.UserContext{
		UserName: jwtClaims[jwtClaimUserName].(string),
		Roles:    types.InterfaceSliceToStringSlice(jwtClaims[jwtClaimRoles].([]interface{})),
		TenantId: jwtClaims[jwtClaimTenantId].(string),
		Scopes:   types.InterfaceSliceToStringSlice(jwtClaims[jwtClaimScope].([]interface{})),
		Token:    tokenString,
	}

	ctx = security.ContextWithUserContext(ctx, userContext)
	req.Request = req.Request.WithContext(ctx)

	// 4. Continue with the filter chain
	chain.ProcessFilter(req, resp)
}

func (f *UserContextFilter) jwtSigningKey(ctx context.Context) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		vaultPool := vault.PoolFromContext(ctx)

		if f.cfg.KeyPath == "" {
			return nil, errors.New("JWT Key Path not configured")
		} else if f.cfg.KeyName == "" {
			return nil, errors.New("JWT Key Name not configured")
		}

		var keyEncoded string
		err := vaultPool.WithConnection(func(connection *vault.Connection) error {
			results, err := connection.ListSecrets(ctx, f.cfg.KeyPath)
			if err != nil {
				return err
			}

			keyEncoded, _ = results[f.cfg.KeyName]
			if keyEncoded == "" {
				return errors.New("Missing JWT key in vault")
			}

			return nil
		})

		if err != nil {
			return nil, err
		}

		keyBytes, _ := base64.StdEncoding.DecodeString(keyEncoded)
		re, _ := x509.ParsePKIXPublicKey(keyBytes)
		return re, nil
	}
}

func NewUserContextFilter(cfg *UserContextFilterConfig) *UserContextFilter {
	return &UserContextFilter{cfg: cfg}
}

func NewUserContextFilterFromConfig(cfg *config.Config) (*UserContextFilter, error) {
	userContextFilterConfig := new(UserContextFilterConfig)
	if err := cfg.Populate(userContextFilterConfig, configRootUserContextFilter); err != nil {
		return nil, err
	}

	return NewUserContextFilter(userContextFilterConfig), nil
}

func RequireAuthorizedFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	cfg := UserContextFilterConfigFromContext(req.Request.Context())
	userContextFilter := NewUserContextFilter(cfg)
	userContextFilter.Filter(req, resp, chain)
}

type userContextContextKey int

const contextKeyUserContextFilterConfig userContextContextKey = iota

func ContextWithUserContextFilterConfig(ctx context.Context, cfg *UserContextFilterConfig) context.Context {
	return context.WithValue(ctx, contextKeyUserContextFilterConfig, cfg)
}

func UserContextFilterConfigFromContext(ctx context.Context) *UserContextFilterConfig {
	return ctx.Value(contextKeyUserContextFilterConfig).(*UserContextFilterConfig)
}
