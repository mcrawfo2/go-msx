package jwttokenprovider

import (
	"context"
	"crypto/x509"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/vault"
	"encoding/base64"
	"encoding/pem"
	"github.com/dgrijalva/jwt-go"
	"github.com/pavel-v-chernykh/keystore-go"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"time"
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

var (
	logger = log.NewLogger("msx.security.jwttokenprovider")
)

func init() {
	jwt.TimeFunc = func() time.Time {
		return time.Now().Add(5 * time.Second)
	}
}

type TokenProviderConfig struct {
	KeySource   string `config:"default=vault"`
	KeyPath     string `config:"default=secret/phi_pnp"`
	KeyName     string `config:"default=key"`
	KeyPassword string `config:"default="`
}

type TokenProvider struct {
	cfg *TokenProviderConfig
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

	return &security.UserContext{
		UserName:    jwtClaims[jwtClaimUserName].(string),
		Roles:       types.InterfaceSliceToStringSlice(jwtClaims[jwtClaimRoles].([]interface{})),
		TenantId:    tenantUuid,
		Scopes:      types.InterfaceSliceToStringSlice(jwtClaims[jwtClaimScope].([]interface{})),
		Authorities: types.InterfaceSliceToStringSlice(jwtClaims[jwtClaimAuthorities].([]interface{})),
		Token:       token,
	}, nil
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
	logger.Info("Registering JWT token provider")

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
