package jwttokenprovider

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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

func createTokenProvider() *TokenProvider {
	return &TokenProvider{cfg: &TokenProviderConfig{
		KeySource:   keySourceKeystore,
		KeyPath:     "testdata/msxjwtkeystore.jks",
		KeyName:     "jwt",
		KeyPassword: "AwesomEKEyStorePWD4jWt",
	}}
}

func TestTokenProvider_keystoreSigningKey(t *testing.T) {
	tokenProvider := createTokenProvider()
	keyStore, err := tokenProvider.keystoreSigningKey(nil)
	assert.NoError(t, err)
	assert.NotNil(t, keyStore)
}

func TestTokenProvider_SecurityContextFromToken(t *testing.T) {
	tokenProvider := createTokenProvider()
	var token *string
	loadFile(t, "token", &token)

	userContext, err := tokenProvider.SecurityContextFromToken(context.Background(), *token)
	assert.NoError(t, err)
	assert.NotNil(t, userContext)
}
