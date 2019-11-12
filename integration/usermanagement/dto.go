package usermanagement

type UserCapabilityListResponse struct {
	Capabilities []UserCapabilityResponse `json:"capabilities"`
}

type UserCapabilityResponse struct {
	Name string `json:"name"`
}

type TenantResponse struct {
	TenantId     string `json:"tenantId"`
	ProviderId   string `json:"providerId"`
	ProviderName string `json:"providerName"`
	TenantName   string `json:"tenantName"`
	DisplayName  string `json:"displayName"`
	CreatedOn    int64  `json:"createdOn"`
	ModifiedOn   int64  `json:"lastUpdated"`
	Suspended    bool   `json:"suspended"`
}

type TenantListResponse struct {
	Tenants []TenantResponse `json:"tenants"`
}

type TenantIdList []string

type LoginResponse struct {
	AccessToken string   `json:"access_token"`
	TokenType   string   `json:"token_type"`
	ExpiresIn   int      `json:"expires_in"`
	Scope       string   `json:"scope"`
	FirstName   string   `json:"firstName"`
	LastName    string   `json:"lastName"`
	Roles       []string `json:"roles"`
	IdToken     string   `json:"id_token"`
	TenantId    string   `json:"tenantId"`
	Email       string   `json:"email"`
	Username    string   `json:"username"`
}
