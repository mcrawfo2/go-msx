package idmdetailsprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement"
	"cto-github.cisco.com/NFV-BU/go-msx/rbac"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
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

type slowDetailsFetcher struct{}

func (t *slowDetailsFetcher) FetchDetails(ctx context.Context) (*security.UserContextDetails, error) {
	userManagementApi, err := usermanagement.NewIntegration(ctx)
	if err != nil {
		return nil, err
	}

	details := new(security.UserContextDetails)

	// Fill in Permission details
	capabilitiesResponse, err := userManagementApi.GetMyCapabilities()
	if err != nil {
		return nil, err
	}

	capabilities := capabilitiesResponse.Payload.(*usermanagement.UserCapabilityListResponse)
	permissionNames := types.StringSet{}
	for _, capability := range capabilities.Capabilities {
		permissionNames.Add(capability.Name)
	}
	details.Permissions = permissionNames.Values()

	if types.StringStack(details.Permissions).Contains(rbac.PermissionAccessAllTenants) {
		tenantIdsResponse, err := userManagementApi.GetTenantIds()
		if err != nil {
			return nil, err
		}

		tenantIdList := tenantIdsResponse.Payload.(*usermanagement.TenantIdList)
		details.Tenants = *tenantIdList

		primaryTenantUuid := security.UserContextFromContext(ctx).TenantId
		primaryTenantResponse, err := userManagementApi.GetTenantById(primaryTenantUuid.String())
		if err != nil {
			return nil, err
		}

		tenantResponse := primaryTenantResponse.Payload.(*usermanagement.TenantResponse)
		details.TenantId = primaryTenantUuid
		details.TenantName = &tenantResponse.TenantName
		details.ProviderId = tenantResponse.ProviderId
		details.ProviderName = &tenantResponse.ProviderName
	} else {
		// Fill in Tenant details
		tenantResponse, err := userManagementApi.GetMyTenants()
		if err != nil {
			return nil, err
		}

		primaryTenantId := security.UserContextFromContext(ctx).TenantId
		tenants := tenantResponse.Payload.(*usermanagement.TenantListResponse)
		for _, tenant := range tenants.Tenants {
			tenantId := tenant.TenantId
			details.Tenants = append(details.Tenants, tenantId)
			if primaryTenantId.Equals(tenantId) {
				details.TenantId = tenantId
				tenantName := tenant.TenantName
				details.TenantName = &tenantName
				details.ProviderId = tenant.ProviderId
				providerName := tenant.ProviderName
				details.ProviderName = &providerName
			}
		}
	}

	if details.ProviderId == nil {
		providersResponse, err := userManagementApi.GetMyProvider()
		if err != nil {
			return nil, err
		}
		provider, ok := providersResponse.Payload.(*usermanagement.ProviderResponse)
		if !ok {
			return nil, errors.New("Incorrect response format for provider")
		}

		details.ProviderName = &provider.Name
		details.ProviderId = provider.ProvidersID
	}

	// Fill in User Details and Roles
	personalInfoResponse, err := userManagementApi.GetMyPersonalInfo()
	if err != nil {
		var locale = "en_US"
		userContext := security.UserContextFromContext(ctx)
		details.Username = &userContext.UserName
		details.Email = &userContext.Email
		details.GivenName = &userContext.FirstName
		details.FamilyName = &userContext.LastName
		details.Locale = &locale
		details.Roles = userContext.Roles
	} else {
		personalInfo := personalInfoResponse.Payload.(*usermanagement.UserPersonalInfoResponse)
		details.Username = &personalInfo.UserId
		details.Email = personalInfo.Email
		details.GivenName = personalInfo.FirstName
		details.FamilyName = personalInfo.LastName
		details.Locale = personalInfo.Locale
		details.Roles = personalInfo.Roles
	}

	// Fill in User Id
	userIdResponse, err := userManagementApi.GetMyUserId()
	if err != nil {
		return nil, err
	}

	details.UserId = userIdResponse.Payload.(*usermanagement.UserIdResponse).Uuid

	// Fill in Token
	userContext := security.UserContextFromContext(ctx)
	details.Issuer = &userContext.Issuer
	details.Subject = &userContext.Subject
	details.Expires = &userContext.Exp
	details.IssuedAt = &userContext.IssuedAt
	details.Jti = &userContext.Jti
	details.AuthTime = &userContext.IssuedAt
	details.ClientId = &userContext.ClientId

	// TODO: details.Currency
	// TODO: details.ClientId

	details.Active = true
	return details, nil
}

func (t *slowDetailsFetcher) FetchActive(ctx context.Context) (bool, error) {
	userManagementApi, err := usermanagement.NewIntegration(ctx)
	if err != nil {
		return false, err
	}

	if _, err = userManagementApi.IsTokenActive(); err != nil {
		return false, err
	}

	return true, nil
}
