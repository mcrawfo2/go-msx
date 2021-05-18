package security

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

type UserContextDetails struct {
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
	Tenants      []types.UUID `json:"tenants"`
	Roles        []string     `json:"roles"`
	Permissions  []string     `json:"permissions"`
}

func (u UserContextDetails) HasTenantId(tenantId types.UUID) bool {
	for _, id := range u.Tenants {
		if tenantId.Equals(id) {
			return true
		}
	}
	return false
}

func NewUserContextDetailsFromUserContext(userContext *UserContext) *UserContextDetails {
	return &UserContextDetails{
		Active:       true,
		Issuer:       &userContext.Issuer,
		Subject:      &userContext.Subject,
		Expires:      &userContext.Exp,
		IssuedAt:     &userContext.IssuedAt,
		GivenName:    &userContext.FirstName,
		FamilyName:   &userContext.LastName,
		Email:        &userContext.Email,
		Scopes:       userContext.Scopes,
		Username:     &userContext.UserName,
		UserId:       types.EmptyUUID(),
		TenantId:     types.EmptyUUID(),
		Tenants:      []types.UUID{},
		Roles:        userContext.Roles,
		Permissions:  []string{},
	}
}