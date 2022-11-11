// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package idmdetailsprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
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
	userManagementApi, err := usermanagement.NewIntegration(ctx)
	if err != nil {
		return false, err
	}

	response, err := userManagementApi.GetTokenDetails(true)
	if err != nil {
		return false, err
	}

	payload := response.Payload.(*usermanagement.TokenDetails)
	return payload.Active, nil
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
	accessibleTenantsSet := make(types.UUIDSet)
	accessibleTenantsSet.Add(payload.Tenants...)

	for _, tenantId := range payload.Tenants {
		_, descendants, e := userManagementApi.GetTenantHierarchyDescendants(tenantId)
		if e != nil {
			return nil, e
		}
		accessibleTenantsSet.Add(descendants...)
	}

	accessibleTenants := make([]types.UUID, len(accessibleTenantsSet))
	i := 0
	for k := range accessibleTenantsSet {
		accessibleTenants[i] = types.MustFormatUUID(k[:])
		i++
	}

	userContext := t.toUserContextDetails(payload)
	userContext.Tenants = accessibleTenants //for backward compatibility, set userContext.tenants to the full accessible tenants list
	return userContext, nil
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
