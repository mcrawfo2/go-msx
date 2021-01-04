package jwttokenprovider

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"reflect"
	"testing"
)

func loadFile(t *testing.T, filename string, target interface{}) []byte {
	data, err := ioutil.ReadFile("testdata/" + filename + ".json")
	if err != nil {
		assert.FailNow(t, "failed to read file", err.Error())
	}
	if target != nil {
		if err = json.Unmarshal(data, target); err != nil {
			assert.FailNow(t, "failed to unmarshal file contents", err.Error())
		}
	}
	return data
}

func createTokenProviderKeystore() *TokenProvider {
	return &TokenProvider{cfg: &TokenProviderConfig{
		KeySource:   keySourceKeystore,
		KeyPath:     "testdata/msxjwtkeystore.jks",
		KeyName:     "jwt",
		KeyPassword: "AwesomEKEyStorePWD4jWt",
	}}
}

func createTokenProviderPem() *TokenProvider {
	return &TokenProvider{cfg: &TokenProviderConfig{
		KeySource: keySourcePem,
		KeyPath:   "testdata/jwt-pubkey.pem",
	}}
}

func TestTokenProvider_keystoreSigningKey(t *testing.T) {
	tokenProvider := createTokenProviderKeystore()
	keyStore, err := tokenProvider.keystoreSigningKey(nil)
	assert.NoError(t, err)
	assert.NotNil(t, keyStore)
}

func TestTokenProvider_pemSigningKey(t *testing.T) {
	tokenProvider := createTokenProviderPem()
	keyStore, err := tokenProvider.pemSigningKey(nil)
	assert.NoError(t, err)
	assert.NotNil(t, keyStore)
}

// TODO: Implement using mock clock
//func TestTokenProvider_SecurityContextFromToken(t *testing.T) {
//	tokenProvider := createTokenProviderPem()
//	var token *string
//	loadFile(t, "token", &token)
//
//	userContext, err := tokenProvider.UserContextFromToken(context.Background(), *token)
//	assert.NoError(t, err)
//	assert.NotNil(t, userContext)
//}

func TestNewTokenProviderConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *TokenProviderConfig
		wantErr bool
	}{
		{
			name: "Defaults",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{}),
			},
			want: &TokenProviderConfig{
				KeySource:   "vault",
				KeyPath:     "secret/phi_pnp",
				KeyName:     "key",
				KeyPassword: "",
			},
		},
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"security.keys.jwt.key-source":   "pem",
					"security.keys.jwt.key-path":     "testdata/jwt-pubkey.pem",
					"security.keys.jwt.key-name":     "ignored1",
					"security.keys.jwt.key-password": "ignored2",
				}),
			},
			want: &TokenProviderConfig{
				KeySource:   "pem",
				KeyPath:     "testdata/jwt-pubkey.pem",
				KeyName:     "ignored1",
				KeyPassword: "ignored2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTokenProviderConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTokenProviderConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTokenProviderConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
