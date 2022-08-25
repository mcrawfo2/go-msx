// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package service

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient/mocks"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/jarcoal/httpmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestNewSecurityAccountsDefaultSettings(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *SecurityAccountsDefaultSettings
		wantErr bool
	}{
		{
			name: "Defaults",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{}),
			},
			want: &SecurityAccountsDefaultSettings{
				Username: "system",
				Password: "system",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSecurityAccountsDefaultSettings(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSecurityAccountsDefaultSettings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSecurityAccountsDefaultSettings() gotCfg = %v, want %v", got, tt.want)
			}
		})
	}
}

func createSampleToken(payloadType string) mocks.HttpRouteBuilder {
	if payloadType == "401" {
		return mocks.HttpRouteBuilder{
			Url:        "http://authservice/v2/token",
			Method:     "POST",
			StatusCode: 401,
			Status:     "FORBIDDEN",
			BodyPath:   "./tokenResponse.json",
			Headers:    nil,
		}
	}
	return mocks.HttpRouteBuilder{
		Url:        "http://authservice/v2/token",
		Method:     "POST",
		StatusCode: 200,
		Status:     "OK",
		BodyPath:   "./tokenResponse.json",
		Headers:    nil,
	}
}

func setup(getToken mocks.HttpRouteBuilder) (context.Context, *security.MockTokenProvider) {
	ctx := context.Background()
	cfg := configtest.NewInMemoryConfig(map[string]string{
		"Username":  "system",
		"Password":  "system",
		"KeySource": "pem",
	})
	ctx = config.ContextWithConfig(ctx, cfg)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockFactory := mocks.NewMockHttpClientFactory()

	cfg.Load(ctx)
	mockFactory.AddRouteDefinition(getToken)
	ctx = httpclient.ContextWithFactory(ctx, mockFactory)
	mockTokenProvider := new(security.MockTokenProvider)
	security.SetTokenProvider(mockTokenProvider)
	return ctx, mockTokenProvider
}

func TestLoginWithUser(t *testing.T) {
	var userId = types.MustNewUUID()

	mockSystemUserContext := security.UserContext{
		UserName:    "system",
		Roles:       []string{"API_ADMIN"},
		TenantId:    types.UUID("534e159b-6c58-4901-a862-8363a25da0f2"),
		Scopes:      []string{"phone", "token_details", "openid"},
		Authorities: nil,
		FirstName:   "test",
		LastName:    "test",
		Issuer:      "test",
		Subject:     "test",
		Exp:         0,
		IssuedAt:    0,
		Jti:         "",
		Email:       "",
		Token:       "dGVzdHRva2Vu",
		Certificate: nil,
		ClientId:    "nfv-client-id",
	}
	mockUserContext := security.UserContext{
		UserName:    "username",
		Roles:       []string{"API_ADMIN"},
		TenantId:    types.UUID("534e159b-6c58-4901-a862-8363a25da0f2"),
		Scopes:      []string{"phone", "token_details", "openid"},
		Authorities: nil,
		FirstName:   "test",
		LastName:    "test",
		Issuer:      "test",
		Subject:     "test",
		Exp:         0,
		IssuedAt:    0,
		Jti:         "",
		Email:       "",
		Token:       "dGVzdHRva2Vu",
		Certificate: nil,
		ClientId:    "nfv-client-id",
	}
	ctx, mockTokenProvider := setup(createSampleToken("success"))
	mockTokenProvider.
		On("UserContextFromToken", ctx, mock.Anything).
		Return(&mockSystemUserContext, nil).Return(&mockUserContext, nil)

	userContextDetails, err := LoginWithUser(ctx, userId)
	assert.NoError(t, err)
	assert.Equal(t, userContextDetails.UserName, "username")
}

func TestLoginWithUser_Error(t *testing.T) {
	var userId = types.MustNewUUID()
	ctx, mockTokenProvider := setup(createSampleToken("success"))
	mockTokenProvider.
		On("UserContextFromToken", ctx, mock.Anything).
		Return(nil, errors.New("Connection refused"))

	_, err := LoginWithUser(ctx, userId)
	assert.ErrorContains(t, err, "Connection refused")
}

func TestWithUserContext(t *testing.T) {
	var userId = types.MustNewUUID()

	mockSystemUserContext := security.UserContext{
		UserName:    "system",
		Roles:       []string{"API_ADMIN"},
		TenantId:    types.UUID("534e159b-6c58-4901-a862-8363a25da0f2"),
		Scopes:      []string{"phone", "token_details", "openid"},
		Authorities: nil,
		FirstName:   "test",
		LastName:    "test",
		Issuer:      "test",
		Subject:     "test",
		Exp:         0,
		IssuedAt:    0,
		Jti:         "",
		Email:       "",
		Token:       "dGVzdHRva2Vu",
		Certificate: nil,
		ClientId:    "nfv-client-id",
	}
	mockUserContext := security.UserContext{
		UserName:    "username",
		Roles:       []string{"API_ADMIN"},
		TenantId:    types.UUID("534e159b-6c58-4901-a862-8363a25da0f2"),
		Scopes:      []string{"phone", "token_details", "openid"},
		Authorities: nil,
		FirstName:   "test",
		LastName:    "test",
		Issuer:      "test",
		Subject:     "test",
		Exp:         0,
		IssuedAt:    0,
		Jti:         "",
		Email:       "",
		Token:       "dGVzdHRva2Vu",
		Certificate: nil,
		ClientId:    "nfv-client-id",
	}
	ctx, mockTokenProvider := setup(createSampleToken("success"))
	mockTokenProvider.
		On("UserContextFromToken", ctx, mock.Anything).
		Return(&mockSystemUserContext, nil).Return(&mockUserContext, nil)
	err := WithUserContext(ctx, userId, func(executeContext context.Context) error {
		assert.Equal(t, security.UserNameFromContext(executeContext), "username")
		return nil
	})
	assert.NoError(t, err)
}

func TestWithUserContext_Error(t *testing.T) {
	var userId = types.MustNewUUID()
	ctx, mockTokenProvider := setup(createSampleToken("success"))
	mockTokenProvider.
		On("UserContextFromToken", ctx, mock.Anything).
		Return(nil, errors.New("Connection refused"))

	err := WithUserContext(ctx, userId, func(ctx context.Context) error {
		return nil
	})

	assert.ErrorContains(t, err, "Connection refused")
}

func TestFailSystemLogin_TokenCallFail(t *testing.T) {
	var userId = types.MustNewUUID()
	ctx, _ := setup(createSampleToken("401"))

	_, err := LoginWithUser(ctx, userId)
	assert.ErrorContains(t, err, "Unauthorized")
}
