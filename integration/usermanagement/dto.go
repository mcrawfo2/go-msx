package usermanagement

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

type Pojo integration.Pojo
type PojoArray integration.PojoArray
type HealthResult integration.HealthResult

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
	Tenants      []types.UUID `json:"tenants"`
	Roles        []string     `json:"roles"`
	Permissions  []string     `json:"permissions"`
}

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

type UserIdResponse struct {
	Uuid types.UUID `json:"uuid"`
}

type UserPersonalInfoResponse struct {
	FirstName *string  `json:"firstName"`
	LastName  *string  `json:"lastName"`
	UserId    string   `json:"userId"`
	Email     *string  `json:"email"`
	Locale    *string  `json:"locale"`
	Roles     []string `json:"roles"`
}

type ProviderResponse struct {
	Name             string     `json:"name"`
	DisplayName      *string    `json:"displayName"`
	Description      *string    `json:"description"`
	ProvidersID      types.UUID `json:"providersId"`
	Email            *string    `json:"email"`
	NotificationType *string    `json:"notificationType"`
	Locale           *string    `json:"locale"`
	AnyConnectURL    *string    `json:"anyConnectURL"`
}

type TenantResponse struct {
	TenantId          types.UUID  `json:"tenantId"`
	ParentId          *types.UUID `json:"parentId"`
	ProviderId        types.UUID  `json:"providerId"`
	ProviderName      string      `json:"providerName"`
	TenantName        string      `json:"tenantName"`
	DisplayName       string      `json:"displayName"`
	Image             string      `json:"image"`
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

type TenantIdList []types.UUID

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

type RoleCreateRequest struct {
	CapabilityAddList    []string `json:"capabilityAddList"`
	CapabilityRemoveList []string `json:"capabilityRemoveList"`
	CapabilityList       []string `json:"capabilitylist"`
	Description          string   `json:"description"`
	DisplayName          string   `json:"displayName"`
	IsSeeded             string   `json:"isSeeded"`
	Owner                string   `json:"owner"`
	ResourceDescriptor   string   `json:"resourceDescriptor"`
	RoleName             string   `json:"roleName"`
	Roleid               string   `json:"roleid"`
}

type RoleUpdateRequest RoleCreateRequest

type RoleResponse struct {
	CapabilityList     []string `json:"capabilitylist"`
	Description        string   `json:"description"`
	DisplayName        string   `json:"displayName"`
	Href               string   `json:"href"`
	IsSeeded           string   `json:"isSeeded"`
	Owner              string   `json:"owner"`
	ResourceDescriptor string   `json:"resourceDescriptor"`
	RoleName           string   `json:"roleName"`
	RoleId             string   `json:"roleid"`
	Status             string   `json:"status"`
}

type RoleListResponse struct {
	Roles []RoleResponse `json:"roles"`
}

type CapabilityCreateRequest struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	DisplayName string `json:"displayName"`
	IsDefault   string `json:"isDefault"`
	Name        string `json:"name"`
	ObjectName  string `json:"objectName"`
	Operation   string `json:"operation"`
}

type CapabilityBatchCreateRequest struct {
	Capabilities []CapabilityCreateRequest `json:"capabilities"`
}

type CapabilityUpdateRequest CapabilityCreateRequest

type CapabilityBatchUpdateRequest struct {
	Capabilities []CapabilityUpdateRequest `json:"capabilities"`
}

type CapabilityResponse struct {
	CapabilityCreateRequest
	Resources []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Method      string `json:"method"`
		Endpoint    string `json:"endpoint"`
	} `json:"resources"`
}

type CapabilityListResponse struct {
	Capabilities []CapabilityResponse `json:"capabilities"`
}

type SecretPolicySetRequest struct {
	AgingRule      AgingRule      `json:"agingRule"`
	CharacterRule  CharacterRule  `json:"characterRule"`
	DictionaryRule DictionaryRule `json:"dictionaryRule"`
	HistoryRule    HistoryRule    `json:"historyRule"`
	KeyRule        KeyRule        `json:"keyRule"`
	LengthRule     LengthRule     `json:"lengthRule"`
}

type SecretPolicyResponse struct {
	SecretPolicySetRequest
	Name string `json:"name"`
}

type AgingRule struct {
	Enabled          bool `json:"enabled"`
	ExpireWarningSec int  `json:"expireWarningSec"`
	GraceAuthNLimit  int  `json:"graceAuthNLimit"`
	MaxAgeSec        int  `json:"maxAgeSec"`
	MinAgeSec        int  `json:"minAgeSec"`
}

type CharacterRule struct {
	ASCIICharactersonly bool   `json:"asciiCharactersonly"`
	Enabled             bool   `json:"enabled"`
	MinDigit            int    `json:"minDigit"`
	MinLowercasechars   int    `json:"minLowercasechars"`
	MinSpecialchars     int    `json:"minSpecialchars"`
	MinUppercasechars   int    `json:"minUppercasechars"`
	SpecialCharacterSet string `json:"specialCharacterSet"`
}

type DictionaryRule struct {
	Enabled              bool `json:"enabled"`
	TestReversedPassword bool `json:"testReversedPassword"`
}

type HistoryRule struct {
	Enabled                    bool `json:"enabled"`
	Passwdhistorycount         int  `json:"passwdhistorycount"`
	PasswdhistorydurationMonth int  `json:"passwdhistorydurationMonth"`
}

type KeyRule struct {
	Algorithm string `json:"algorithm"`
	Enabled   bool   `json:"enabled"`
	Format    string `json:"format"`
}

type LengthRule struct {
	Enabled   bool `json:"enabled"`
	MaxLength int  `json:"maxLength"`
	MinLength int  `json:"minLength"`
}

type SecretsResponse map[string]string

func (s SecretsResponse) Value(key string) string {
	return s[key]
}
