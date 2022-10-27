// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package auth

import (
	"crypto/rsa"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/base64"
	"github.com/pkg/errors"
	"math/big"
)

type TokenDetails struct {
	Active       bool         `json:"active"`
	Issuer       *string      `json:"iss"`
	Subject      *string      `json:"sub"`
	Expires      *int         `json:"exp"`
	IssuedAt     *int         `json:"iat"`
	Jti          *string      `json:"jti"`
	AuthTime     *int         `json:"auth_time"`
	GivenName    *string      `json:"given_name"`
	FamilyName   *string      `json:"family_name"`
	Email        *string      `json:"email"`
	Locale       *string      `json:"locale"`
	Scopes       []string     `json:"scope"`
	ClientId     *string      `json:"client_id"`
	Username     *string      `json:"username"`
	UserId       types.UUID   `json:"user_id"`
	Currency     *string      `json:"currency"`
	TenantId     types.UUID   `json:"tenant_id"`
	TenantName   *string      `json:"tenant_name"`
	ProviderId   types.UUID   `json:"provider_id"`
	ProviderName *string      `json:"provider_name"`
	Tenants      []types.UUID `json:"assigned_tenants"`
	Roles        []string     `json:"roles"`
	Permissions  []string     `json:"permissions"`
}

type LoginResponse struct {
	AccessToken string      `json:"access_token"`
	TokenType   string      `json:"token_type"`
	ExpiresIn   int         `json:"expires_in"`
	Scope       string      `json:"scope"`
	FirstName   string      `json:"firstName"`
	LastName    string      `json:"lastName"`
	Roles       []string    `json:"roles"`
	IdToken     string      `json:"id_token"`
	TenantId    *types.UUID `json:"tenantId"`
	Email       string      `json:"email"`
	Username    string      `json:"username"`
}

type JsonWebKey struct {
	KeyId              string `json:"kid"`
	KeyType            string `json:"kty"`
	RsaModulus         string `json:"n"`
	RsaPublicExponent  string `json:"e"`
	RsaPrivateExponent string `json:"d"`
}

func (j JsonWebKey) RsaPublicKey() (*rsa.PublicKey, error) {
	modulusBytes, err := base64.RawURLEncoding.DecodeString(j.RsaModulus)
	if err != nil {
		return nil, err
	}
	exponentBytes, err := base64.RawURLEncoding.DecodeString(j.RsaPublicExponent)
	if err != nil {
		return nil, err
	}

	var modulus = new(big.Int)
	modulus.SetBytes(modulusBytes)

	var exponent = new(big.Int)
	exponent.SetBytes(exponentBytes)

	return &rsa.PublicKey{
		N: modulus,
		E: int(exponent.Int64()),
	}, nil
}

type JsonWebKeys struct {
	Keys []JsonWebKey `json:"keys"`
}

func (k JsonWebKeys) KeyById(kid string) (JsonWebKey, error) {
	for _, v := range k.Keys {
		if v.KeyId == kid {
			return v, nil
		}
	}
	return JsonWebKey{}, errors.Errorf("Json Web key not found: %q", kid)
}
