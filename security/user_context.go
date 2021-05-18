package security

import (
	"crypto/x509"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

const contextDefaultUserName = "anonymous"

var (
	defaultUserContext = &UserContext{
		UserName:    contextDefaultUserName,
		Roles:       nil,
		TenantId:    nil,
		Scopes:      nil,
		Authorities: nil,
		FirstName:   "",
		LastName:    "",
		Issuer:      "",
		Subject:     "",
		Exp:         0,
		IssuedAt:    0,
		Jti:         "",
		Email:       "",
		Token:       "",
		Certificate: nil,
	}
)

type UserContext struct {
	UserName    string            `json:"user_name"`
	Roles       []string          `json:"roles"`
	TenantId    types.UUID        `json:"tenant_id"`
	Scopes      []string          `json:"scope"`
	Authorities []string          `json:"authorities"`
	FirstName   string            `json:"firstName"`
	LastName    string            `json:"lastName"`
	Issuer      string            `json:"iss"`
	Subject     string            `json:"sub"`
	Exp         int               `json:"exp"`
	IssuedAt    int               `json:"iat"`
	Jti         string            `json:"jti"`
	Email       string            `json:"email"`
	Token       string            `json:"-"`
	Certificate *x509.Certificate `json:"-"`
	ClientId    string            `json:"client_id"`
}

func (c *UserContext) Clone() *UserContext {
	return &UserContext{
		UserName:    c.UserName,
		Roles:       c.Roles[:],
		TenantId:    c.TenantId,
		Scopes:      c.Scopes[:],
		Authorities: c.Authorities[:],
		FirstName:   c.FirstName,
		LastName:    c.LastName,
		Issuer:      c.Issuer,
		Subject:     c.Subject,
		Exp:         c.Exp,
		IssuedAt:    c.IssuedAt,
		Jti:         c.Jti,
		Email:       c.Email,
		Token:       c.Token,
		Certificate: c.Certificate,
		ClientId:    c.ClientId,
	}
}
