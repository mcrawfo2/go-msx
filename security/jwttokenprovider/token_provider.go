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
	"github.com/pavel-v-chernykh/keystore-go"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"strings"
)

const (
	jwtClaimUserName    = "user_name"
	jwtClaimScope       = "scope"
	jwtClaimTenantId    = "tenantId"
	jwtClaimRoles       = "roles"
	jwtClaimAuthorities = "authorities"

	configRootJwtTokenProvider = "security.keys.jwt"

	keySourcePem      = "pem"
	keySourceKeystore = "keystore"
	keySourceVault    = "vault"

	pemBeginPublicKey = `-----BEGIN PUBLIC KEY-----`
	pemEndPublicKey   = `-----END PUBLIC KEY-----`
)

var ErrNotFound = errors.New("Token not found in security context")

type TokenProviderConfig struct {
	KeySource   string `config:"default=vault"`
	KeyPath     string `config:"default=secret/phi_pnp"`
	KeyName     string `config:"default=key"`
	KeyPassword string `config:"default="`
}

type TokenProvider struct {
	cfg *TokenProviderConfig
}

func (j *TokenProvider) SecurityContextFromToken(ctx context.Context, token string) (userContext *security.UserContext, err error) {
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

func (j *TokenProvider) signingKeyFunc(ctx context.Context) (jwt.Keyfunc, error) {
	switch j.cfg.KeySource {
	case keySourcePem:
		return j.pemSigningKey, nil
	case keySourceKeystore:
		return j.keystoreSigningKey, nil
	case keySourceVault:
		return j.vaultSigningKeyFunc(ctx), nil
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

	pemLines := strings.Fields(string(pemBytes))

	for firstLine := 0; firstLine < len(pemLines); firstLine++ {
		if pemLines[firstLine] == pemBeginPublicKey {
			pemLines = pemLines[firstLine+1:]
			break
		}
	}

	for lastLine := 0; lastLine < len(pemLines); lastLine++ {
		if pemLines[lastLine] == pemEndPublicKey {
			pemLines = pemLines[:lastLine]
			break
		}
	}

	pem := strings.Join(pemLines, "")
	re, err := x509.ParsePKIXPublicKey([]byte(pem))
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
	if !exists  {
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
		vaultPool := vault.PoolFromContext(ctx)

		if j.cfg.KeyPath == "" {
			return nil, errors.New("JWT Key Path not configured")
		} else if j.cfg.KeyName == "" {
			return nil, errors.New("JWT Key Name not configured")
		}

		var keyEncoded string
		err := vaultPool.WithConnection(func(connection *vault.Connection) error {
			results, err := connection.ListSecrets(ctx, j.cfg.KeyPath)
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
		re, err := x509.ParsePKIXPublicKey(keyBytes)
		return re, err
	}
}

func RegisterTokenProvider(ctx context.Context) error {
	cfg := config.FromContext(ctx)
	jwtTokenProviderConfig := new(TokenProviderConfig)
	if err := cfg.Populate(jwtTokenProviderConfig, configRootJwtTokenProvider); err != nil {
		return err
	}

	security.SetTokenProvider(&TokenProvider{
		cfg: jwtTokenProviderConfig,
	})

	return nil
}