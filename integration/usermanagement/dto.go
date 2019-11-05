package usermanagement

type UserCapabilityListDTO struct {
	Capabilities []UserCapabilityDTO `json:"capabilities"`
}

type UserCapabilityDTO struct {
	Name string `json:"name"`
}

type TenantDTO struct {
	TenantId     string `json:"tenantId"`
	ProviderId   string `json:"providerId"`
	ProviderName string `json:"providerName"`
	TenantName   string `json:"tenantName"`
	DisplayName  string `json:"displayName"`
	CreatedOn    int32  `json:"createdOn"`
	ModifiedOn   int32  `json:"lastUpdated"`
	Suspended    bool   `json:"suspended"`
}

type TenantListDTO struct {
	Tenants []TenantDTO `json:"tenants"`
}

type TenantIdList []string
