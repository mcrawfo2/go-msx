package idmdetailsprovider

import (
	"context"

	"cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
)

type detailsFetcher interface {
	FetchDetails(ctx context.Context) (*security.UserContextDetails, error)
	FetchActive(ctx context.Context) (bool, error)
}

type fastDetailsFetcher struct{}

func (t *fastDetailsFetcher) FetchDetails(ctx context.Context) (*security.UserContextDetails, error) {
	return t.fetchUserContextDetails(ctx)
}

func (t *fastDetailsFetcher) FetchActive(ctx context.Context) (bool, error) {
	details, err := t.fetchUserContextDetails(ctx)
	if err != nil {
		return false, err
	}
	return details.Active, err
}

func (t *fastDetailsFetcher) fetchUserContextDetails(ctx context.Context) (*security.UserContextDetails, error) {
	userManagementApi, err := usermanagement.NewIntegration(ctx)
	if err != nil {
		return nil, err
	}

	response, err := userManagementApi.GetTokenDetails(false)
	if err != nil {
		return nil, err
	}

	payload := response.Payload.(*usermanagement.TokenDetails)
	return t.toUserContextDetails(payload), nil
}

func (t *fastDetailsFetcher) toUserContextDetails(payload *usermanagement.TokenDetails) *security.UserContextDetails {
	return &security.UserContextDetails{
		Active:       payload.Active,
		Issuer:       payload.Issuer,
		Subject:      payload.Subject,
		Expires:      payload.Expires,
		IssuedAt:     payload.IssuedAt,
		Jti:          payload.Jti,
		AuthTime:     payload.AuthTime,
		GivenName:    payload.GivenName,
		FamilyName:   payload.FamilyName,
		Email:        payload.Email,
		Locale:       payload.Locale,
		Scopes:       payload.Scopes,
		ClientId:     payload.ClientId,
		Username:     payload.Username,
		UserId:       payload.UserId,
		Currency:     payload.Currency,
		TenantId:     payload.TenantId,
		TenantName:   payload.TenantName,
		ProviderId:   payload.ProviderId,
		ProviderName: payload.ProviderName,
		Tenants:      payload.Tenants,
		Roles:        payload.Roles,
		Permissions:  payload.Permissions,
	}
}
