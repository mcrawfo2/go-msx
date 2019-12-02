package jwttokenprovider

import (
	"context"
	"crypto/x509"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/vault"
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	jwtClaimUserName    = "user_name"
	jwtClaimScope       = "scope"
	jwtClaimTenantId    = "tenantId"
	jwtClaimRoles       = "roles"
	jwtClaimAuthorities = "authorities"

	configRootJwtTokenProvider = "security.keys.jwt"
)

var ErrNotFound = errors.New("Token not found in security context")

type TokenProviderConfig struct {
	KeyPath string `config:"default=secret/phi_pnp"`
	KeyName string `config:"default=key"`
}

type TokenProvider struct {
	ctx context.Context
	cfg *TokenProviderConfig
}

func (j *TokenProvider) SecurityContextFromToken(token string) (userContext *security.UserContext, err error) {
	parser := &jwt.Parser{}
	jwtClaims := jwt.MapClaims{}
	_, err = parser.ParseWithClaims(token, jwtClaims, j.signingKey)
	if err != nil {
		return nil, err
	}

	return &security.UserContext{
		UserName:    jwtClaims[jwtClaimUserName].(string),
		Roles:       types.InterfaceSliceToStringSlice(jwtClaims[jwtClaimRoles].([]interface{})),
		TenantId:    jwtClaims[jwtClaimTenantId].(string),
		Scopes:      types.InterfaceSliceToStringSlice(jwtClaims[jwtClaimScope].([]interface{})),
		Authorities: types.InterfaceSliceToStringSlice(jwtClaims[jwtClaimAuthorities].([]interface{})),
		Token:       token,
	}, nil
}

func (j *TokenProvider) TokenFromSecurityContext(userContext *security.UserContext) (token string, err error) {
	if userContext.Token == "" {
		return "", ErrNotFound
	}
	return userContext.Token, nil
}

func (j *TokenProvider) signingKey(token *jwt.Token) (interface{}, error) {
	vaultPool := vault.PoolFromContext(j.ctx)

	if j.cfg.KeyPath == "" {
		return nil, errors.New("JWT Key Path not configured")
	} else if j.cfg.KeyName == "" {
		return nil, errors.New("JWT Key Name not configured")
	}

	var keyEncoded string
	err := vaultPool.WithConnection(func(connection *vault.Connection) error {
		results, err := connection.ListSecrets(j.ctx, j.cfg.KeyPath)
		if err != nil {
			return err
		}

		keyEncoded, _ = results[j.cfg.KeyName]
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

func RegisterTokenProvider(ctx context.Context) error {
	cfg := config.FromContext(ctx)
	jwtTokenProviderConfig := new(TokenProviderConfig)
	if err := cfg.Populate(jwtTokenProviderConfig, configRootJwtTokenProvider); err != nil {
		return err
	}

	security.SetTokenProvider(&TokenProvider{
		ctx: ctx,
		cfg: jwtTokenProviderConfig,
	})

	return nil
}
