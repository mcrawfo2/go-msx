package jwtprovider

import (
	"context"
	"crypto/x509"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/vault"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"encoding/base64"
	"github.com/auth0/go-jwt-middleware"
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

	configRootJwtSecurityProvider = "server.jwt"
)

var (
	logger = log.NewLogger("webservice.jwtprovider")
)

type JwtSecurityProviderConfig struct {
	KeyPath string `config:"default=secret/phi_pnp"`
	KeyName string `config:"default=key"`
}

type JwtSecurityProvider struct {
	cfg *JwtSecurityProviderConfig
}

func (f *JwtSecurityProvider) Authentication(req *restful.Request) error {
	// 0. Retrieve the context
	ctx := req.Request.Context()
	if ctx == nil {
		return errors.New("Failed to retrieve context from request")
	}

	middleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: f.jwtSigningKey(ctx),
		SigningMethod:       jwt.SigningMethodRS256,
		Debug:               false,
	})

	// 1. Validate the token signature
	recorder := httptest.NewRecorder()
	if err := middleware.CheckJWT(recorder, req.Request); err != nil {
		return errors.Wrap(err, "JWT validation failed")
	}

	// 2. Validate the session
	tokenString, err := jwtmiddleware.FromAuthHeader(req.Request)
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve token from request")
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
		return errors.Wrap(err, "Failed to parse JWT")
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
	return nil
}

func (f *JwtSecurityProvider) jwtSigningKey(ctx context.Context) jwt.Keyfunc {
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

func NewJwtSecurityProvider(cfg *JwtSecurityProviderConfig) *JwtSecurityProvider {
	return &JwtSecurityProvider{cfg: cfg}
}

func NewJwtSecurityProviderFromConfig(cfg *config.Config) (*JwtSecurityProvider, error) {
	jwtSecurityProviderConfig := new(JwtSecurityProviderConfig)
	if err := cfg.Populate(jwtSecurityProviderConfig, configRootJwtSecurityProvider); err != nil {
		return nil, err
	}

	return NewJwtSecurityProvider(jwtSecurityProviderConfig), nil
}

func RegisterSecurityProvider(ctx context.Context) error {
	cfg := config.MustFromContext(ctx)
	jwtSecurityProvider, err := NewJwtSecurityProviderFromConfig(cfg)
	if err != nil {
		return errors.Wrap(err, "Failed to create JWT web security provider")
	}

	webservice.RegisterSecurityProvider(jwtSecurityProvider)
	return nil
}
