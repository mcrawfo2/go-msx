package usermanagement

import "cto-github.cisco.com/NFV-BU/go-msx/integration"

type Pojo integration.Pojo
type PojoArray integration.PojoArray
type HealthResult integration.HealthResult
type ErrorDTO integration.ErrorDTO
type ErrorDTO2 integration.ErrorDTO2

type EncryptSecretsDTO struct {
	Scope       map[string]string `json:"scope"`
	Name        string            `json:"name"`
	Method      string            `json:"method"`
	SecretNames []string          `json:"secretNames"`
}

type GetSecretRequestDTO struct {
	Names   []string          `json:"names"`
	Encrypt EncryptSecretsDTO `json:"encrypt"`
}

type GenerateSecretRequestDTO struct {
	Names   []string `json:"names"`
	Save    bool     `json:"save"`
	Encrypt *Pojo    `json:"encrypt"`
}

type UserCapabilityListResponse struct {
	Capabilities []UserCapabilityResponse `json:"capabilities"`
}

type UserCapabilityResponse struct {
	Name string `json:"name"`
}

type TenantResponse struct {
	TenantId          string      `json:"tenantId"`
	ProviderId        string      `json:"providerId"`
	ProviderName      string      `json:"providerName"`
	TenantName        string      `json:"tenantName"`
	DisplayName       string      `json:"displayName"`
	CreatedOn         int64       `json:"createdOn"`
	ModifiedOn        int64       `json:"lastUpdated"`
	Suspended         bool        `json:"suspended"`
	TenantDescription string      `json:"tenantDescription"`
	URL               string      `json:"url"`
	TenantGroupName   interface{} `json:"tenantGroupName"`
	TenantExtension   struct {
		Parameters interface{} `json:"parameters"`
	} `json:"tenantExtension"`
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
