// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package oauth2detailsprovider

import (
	"bytes"
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cache/lru"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"encoding/base64"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	configRootOauth2TokenDetailsProvider = "security.token.details.oauth2"
)

type ProviderConfig struct {
	Scheme       string `config:"default=http"`
	Authority    string `config:"default=authservice"`
	Path         string `config:"default=/v2/jwks"`
	ClientId     string `config:"default=${integration.security.client.client-id}"`
	ClientSecret string `config:"default=${integration.security.client.client-secret}"`
	ActiveCache  lru.CacheConfig
	DetailsCache lru.CacheConfig
}

type Provider struct {
	cfg          ProviderConfig
	detailsCache lru.Cache
	activeCache  lru.Cache
}

func (t *Provider) IsTokenActive(ctx context.Context) (active bool, err error) {
	token := security.UserContextFromContext(ctx).Token
	if token == "" {
		return false, errors.Wrap(security.ErrTokenNotFound, "Empty Token")
	}

	active, err = t.getActive(ctx, token)

	if err != nil {
		logger.WithContext(ctx).WithError(err).Error("Failed to load token details")
	}

	return active, nil
}

func (t *Provider) getActive(ctx context.Context, token string) (active bool, err error) {
	activeInterface, exists := t.activeCache.Get(token)
	if !exists {
		active, err = t.fetchActive(ctx, token)
		if err != nil {
			return false, err
		}
	} else {
		active = activeInterface.(bool)
	}

	return active, nil
}

func (t *Provider) fetchActive(ctx context.Context, token string) (active bool, err error) {
	details, err := t.fetch(ctx, token)
	if err != nil {
		return false, err
	}

	t.activeCache.Set(token, active)

	return details.Active, nil
}

func (t *Provider) TokenDetails(ctx context.Context) (details *security.UserContextDetails, err error) {
	token := security.UserContextFromContext(ctx).Token
	if token == "" {
		return nil, errors.Wrap(security.ErrTokenNotFound, "Empty Token")
	}

	details, err = t.getDetails(ctx, token)

	if err != nil {
		logger.WithContext(ctx).WithError(err).Error("Failed to load token details")
	}

	return details, nil
}

func (t *Provider) getDetails(ctx context.Context, token string) (details *security.UserContextDetails, err error) {
	detailsInterface, exists := t.detailsCache.Get(token)
	if !exists {
		details, err = t.fetch(ctx, token)
		if err != nil {
			return nil, err
		}
	} else {
		details = detailsInterface.(*security.UserContextDetails)
	}

	return details, nil
}

func (t *Provider) fetch(ctx context.Context, token string) (details *security.UserContextDetails, err error) {
	client, err := httpclient.New(ctx, nil)
	if err != nil {
		return nil, err
	}

	tokenDetailsUrl := new(url.URL)
	tokenDetailsUrl.Scheme = t.cfg.Scheme
	tokenDetailsUrl.Host = t.cfg.Authority
	tokenDetailsUrl.Path = t.cfg.Path

	bodyRequest := url.Values{
		"token": []string{token},
	}
	bodyBytes := bodyRequest.Encode()
	bodyReader := bytes.NewBufferString(bodyBytes)

	request, err := http.NewRequest(http.MethodPost, tokenDetailsUrl.String(), bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create HTTP request")
	}

	clientCredentials := []byte(t.cfg.ClientId + ":" + t.cfg.ClientSecret)
	request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(clientCredentials))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute HTTP request")
	}
	defer response.Body.Close()

	tokenDetailsBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read HTTP response body")
	}

	details, err = t.parseTokenDetails(tokenDetailsBytes)
	if err != nil {
		return nil, err
	}

	t.detailsCache.Set(token, details)

	return
}

func (t *Provider) parseTokenDetails(detailsBytes []byte) (result *security.UserContextDetails, err error) {
	var token map[string]interface{}
	if err = json.Unmarshal(detailsBytes, &token); err != nil {
		return nil, errors.Wrap(err, "Failed to deserialize raw token details")
	}

	var ext TokenDetailsExt
	if err = json.Unmarshal(detailsBytes, &ext); err != nil {
		return nil, errors.Wrap(err, "Failed to deserialize raw token extended details")
	}

	// Merge ext into body
	for k, v := range ext.Ext {
		if _, ok := token[k]; !ok {
			token[k] = v
		}
	}

	// Serialize modified token back to json
	detailsBytes, err = json.Marshal(token)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to serialize combined token details")
	}

	var tokenDetails security.UserContextDetails
	if err = json.Unmarshal(detailsBytes, &tokenDetails); err != nil {
		return nil, errors.Wrap(err, "Failed to deserialize combined token details")
	}

	return &tokenDetails, nil
}

func NewProviderConfig(ctx context.Context) (*ProviderConfig, error) {
	var providerConfig ProviderConfig
	if err := config.FromContext(ctx).Populate(&providerConfig, configRootOauth2TokenDetailsProvider); err != nil {
		return nil, err
	}

	return &providerConfig, nil
}

func RegisterTokenDetailsProvider(ctx context.Context) error {
	logger.Info("Registering OAuth2 token details provider")
	providerConfig, err := NewProviderConfig(ctx)
	if err != nil {
		return err
	}

	security.SetTokenDetailsProvider(&Provider{
		cfg:          *providerConfig,
		activeCache:  lru.NewCacheFromConfig(&providerConfig.ActiveCache),
		detailsCache: lru.NewCacheFromConfig(&providerConfig.DetailsCache),
	})

	return nil
}

type TokenDetailsExt struct {
	Ext map[string]interface{} `json:"ext"`
}
