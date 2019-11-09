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
