// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package jwttokenprovider

import (
	"context"
	"crypto/x509"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient/discoveryinterceptor"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient/loginterceptor"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient/statsinterceptor"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient/traceinterceptor"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/vault"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pavel-v-chernykh/keystore-go"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

const (
	jwtClaimUserName    = "user_name"
	jwtClaimScope       = "scope"
	jwtClaimScopeAlt    = "scp"
	jwtClaimTenantId    = "tenantId"
	jwtClaimRoles       = "roles"
	jwtClaimAuthorities = "authorities"

	defaultJwtClaimAuthorities = "ROLE_CLIENT"

	configRootJwtTokenProvider = "security.keys.jwt"

	keySourcePem      = "pem"
	keySourceKeystore = "keystore"
	keySourceVault    = "vault"
	keySourceIdm      = "idm"
	keySourceJwks     = "jwks"

	pemBeginPublicKey = `-----BEGIN PUBLIC KEY-----`
	pemEndPublicKey   = `-----END PUBLIC KEY-----`

	jwtHeaderAlg   = "alg"
	jwtHeaderKeyId = "kid"
)

var (
	logger = log.NewLogger("msx.security.jwttokenprovider")

	ErrInvalidTokenKid    = errors.New("Missing or invalid 'kid' field in token header")
	ErrKeyNotExists       = errors.New("Key not found")
	ErrUnsupportedKeyType = errors.New("Unsupported key type")
)

func init() {
	jwt.TimeFunc = func() time.Time {
		return time.Now().Add(5 * time.Second)
	}
}

type TokenProviderConfig struct {
	KeySource    string `config:"default=jwks"`
	KeyPath      string `config:"default=/v2/jwks"`
	KeyScheme    string `config:"default=http"`
	KeyAuthority string `config:"default=authservice"`
	KeyName      string `config:"default=key"`
	KeyPassword  string `config:"default="`
}

type TokenProvider struct {
	cfg      *TokenProviderConfig
	keyCache sync.Map
}

func (j *TokenProvider) UserContextFromToken(ctx context.Context, token string) (userContext *security.UserContext, err error) {
	parser := &jwt.Parser{}
	jwtClaims := jwt.MapClaims{}
	keyFunc, err := j.signingKeyFunc(ctx)
	if err != nil {
		return nil, err
	}

	_, err = parser.ParseWithClaims(token, jwtClaims, keyFunc)
	if err != nil {
		return nil, err
	}

	tenantUuid := types.EmptyUUID()
	tenantId := jwtClaims[jwtClaimTenantId].(string)
	if tenantId != "" {
		tenantUuid, err = types.ParseUUID(tenantId)
		if err != nil {
			return nil, err
		}
	}

	scope := jwtClaims[jwtClaimScope]
	if scope == nil {
		scope = jwtClaims[jwtClaimScopeAlt]
	}
	if scope == nil {
		scope = []interface{}{}
	}

	uc := &security.UserContext{
		UserName: jwtClaims[jwtClaimUserName].(string),
		Roles:    types.InterfaceSliceToStringSlice(jwtClaims[jwtClaimRoles].([]interface{})),
		TenantId: tenantUuid,
		Scopes:   types.InterfaceSliceToStringSlice(scope.([]interface{})),
		Token:    token[:],
	}

	//jwtClaimAuthorities is deprecated and is not present in the token issued by auth service
	//so if it's not found in the token claim, use a default value.
	if claimAuthorities, ok := jwtClaims[jwtClaimAuthorities].([]interface{}); ok {
		uc.Authorities = types.InterfaceSliceToStringSlice(claimAuthorities)
	} else {
		uc.Authorities = []string{defaultJwtClaimAuthorities}
	}

	return uc, nil
}

func (j *TokenProvider) signingKeyFunc(ctx context.Context) (jwt.Keyfunc, error) {
	switch j.cfg.KeySource {
	case keySourcePem:
		return j.pemSigningKey, nil
	case keySourceKeystore:
		return j.keystoreSigningKey, nil
	case keySourceVault:
		return j.vaultSigningKeyFunc(ctx), nil
	case keySourceIdm:
		return j.cachedSigningKeyFunc(ctx, j.idmSigningKey), nil
	case keySourceJwks:
		return j.cachedSigningKeyFunc(ctx, j.jwksSigningKey), nil
	default:
		return nil, errors.Errorf("Unknown JWT Key Source: %s", j.cfg.KeySource)
	}
}

func (j *TokenProvider) pemSigningKey(token *jwt.Token) (interface{}, error) {
	if j.cfg.KeyPath == "" {
		return nil, errors.New("JWT Key Path not configured")
	}

	pemBytes, err := ioutil.ReadFile(j.cfg.KeyPath)
	if err != nil {
		return nil, err
	}

	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil || pemBlock.Type != "PUBLIC KEY" {
		return nil, errors.New("JWT PEM file does not contain valid public key")
	}

	re, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	return re, err
}

func (j *TokenProvider) keystoreSigningKey(token *jwt.Token) (interface{}, error) {
	f, err := os.Open(j.cfg.KeyPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	keyStore, err := keystore.Decode(f, []byte(j.cfg.KeyPassword))
	if err != nil {
		return nil, err
	}

	privateKeyInterface, exists := keyStore[j.cfg.KeyName]
	if !exists {
		return nil, errors.Errorf("No key named %s found in keystore", j.cfg.KeyName)
	}

	privateKey, ok := privateKeyInterface.(*keystore.PrivateKeyEntry)
	if !ok {
		return nil, errors.Errorf("Invalid private key format for %s in keystore", j.cfg.KeyName)
	}

	if len(privateKey.CertChain) == 0 {
		return nil, errors.Errorf("Certificate not found for private key %s in keystore", j.cfg.KeyName)
	}

	jksCert := privateKey.CertChain[0]
	if jksCert.Type != "X.509" {
		return nil, errors.Errorf("Invalid certificate format for certificate %s in keystore", j.cfg.KeyName)
	}

	cert, err := x509.ParseCertificate(jksCert.Content)
	if err != nil {
		return nil, err
	}

	return cert.PublicKey, nil
}

func (j *TokenProvider) vaultSigningKeyFunc(ctx context.Context) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		connection := vault.ConnectionFromContext(ctx)

		if j.cfg.KeyPath == "" {
			return nil, errors.New("JWT Key Path not configured")
		} else if j.cfg.KeyName == "" {
			return nil, errors.New("JWT Key Name not configured")
		}

		var keyEncoded string
		err := (func(connection vault.ConnectionApi) error {
			results, err := connection.ListSecrets(ctx, j.cfg.KeyPath)
			if err != nil {
				return err
			}

			keyEncoded, _ = results[j.cfg.KeyName]
			if keyEncoded == "" {
				return errors.New("Missing JWT key in vault")
			}

			return nil
		})(connection)

		if err != nil {
			return nil, err
		}

		keyBytes, _ := base64.StdEncoding.DecodeString(keyEncoded)
		re, err := x509.ParsePKIXPublicKey(keyBytes)
		return re, err
	}
}

func (j *TokenProvider) cachedSigningKey(kid string) (interface{}, error) {
	cachedValue, ok := j.keyCache.Load(kid)
	if !ok {
		return nil, ErrKeyNotExists
	}

	return cachedValue, nil
}

func (j *TokenProvider) idmSigningKey(ctx context.Context, kid string) (interface{}, error) {
	api, err := usermanagement.NewIntegration(ctx)
	if err != nil {
		return nil, err
	}

	jwks, _, err := api.GetTokenKeys()
	if err != nil {
		return nil, err
	}

	jwk, err := jwks.KeyById(kid)
	if err != nil {
		return nil, errors.Wrap(ErrKeyNotExists, err.Error())
	}

	switch jwk.KeyType {
	case "RSA":
		return jwk.RsaPublicKey()
	default:
		return nil, errors.Wrap(ErrUnsupportedKeyType, jwk.KeyType)
	}
}

func (j *TokenProvider) jwksSigningKey(ctx context.Context, kid string) (interface{}, error) {
	httpProdClientFactory, err := httpclient.NewProductionHttpClientFactory(ctx)
	if err != nil {
		return nil, err
	}

	prodClient := httpProdClientFactory.NewHttpClientWithConfigurer(ctx, httpclient.ClientConfigurer{
		ClientFuncs: []httpclient.ClientConfigurationFunc{
			httpclient.ApplyRecoveryErrorInterceptor,
			loginterceptor.ApplyInterceptor(),
			discoveryinterceptor.ApplyInterceptor(),
			statsinterceptor.ApplyInterceptor(),
			traceinterceptor.ApplyInterceptor(),
		},
		TransportFuncs: []httpclient.TransportConfigurationFunc{},
	})

	requestUrl := new(url.URL)
	requestUrl.Scheme = j.cfg.KeyScheme
	requestUrl.Host = j.cfg.KeyAuthority
	requestUrl.Path = j.cfg.KeyPath

	request, err := http.NewRequest("GET", requestUrl.String(), http.NoBody)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/jwk+json")
	request.Header.Add("Accept", "application/json")

	response, err := prodClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var jwks usermanagement.JsonWebKeys
	err = json.Unmarshal(responseBodyBytes, &jwks)
	if err != nil {
		return nil, err
	}

	jwk, err := jwks.KeyById(kid)
	if err != nil {
		return nil, errors.Wrap(ErrKeyNotExists, err.Error())
	}

	switch jwk.KeyType {
	case "RSA":
		return jwk.RsaPublicKey()
	default:
		return nil, errors.Wrap(ErrUnsupportedKeyType, jwk.KeyType)
	}
}

type signingKeyFunc func(context.Context, string) (interface{}, error)

func (j *TokenProvider) cachedSigningKeyFunc(ctx context.Context, fn signingKeyFunc) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header[jwtHeaderKeyId].(string)
		if !ok {
			return nil, ErrInvalidTokenKid
		}

		key, err := j.cachedSigningKey(kid)

		if err != nil && errors.Is(err, ErrKeyNotExists) {
			key, err = fn(ctx, kid)
			if err == nil {
				j.keyCache.Store(kid, key)
			}
		}

		if err != nil {
			return nil, err
		}

		return key, nil
	}
}

func NewTokenProviderConfig(cfg *config.Config) (*TokenProviderConfig, error) {
	jwtTokenProviderConfig := new(TokenProviderConfig)
	if err := cfg.Populate(jwtTokenProviderConfig, configRootJwtTokenProvider); err != nil {
		return nil, err
	}
	return jwtTokenProviderConfig, nil
}

func RegisterTokenProvider(ctx context.Context) error {
	logger.Info("Registering JWT token provider")

	jwtTokenProviderConfig, err := NewTokenProviderConfig(config.FromContext(ctx))
	if err != nil {
		return err
	}

	security.SetTokenProvider(&TokenProvider{
		cfg:      jwtTokenProviderConfig,
		keyCache: sync.Map{},
	})

	return nil
}
