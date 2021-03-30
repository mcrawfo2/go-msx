//go:generate mockery --inpackage --name=Api --structname=MockUserManagement

package usermanagement

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

type Api interface {
	GetAdminHealth() (*integration.MsxResponse, error)

	Login(user, password string) (*integration.MsxResponse, error)
	Logout() (*integration.MsxResponse, error)
	IsTokenActive() (*integration.MsxResponse, error)
	GetTokenDetails(noDetails bool) (*integration.MsxResponse, error)

	GetMyProvider() (*integration.MsxResponse, error)
	GetProviderByName(providerName string) (*integration.MsxResponse, error)

	GetProviderExtensionByName(name string) (*integration.MsxResponse, error)

	// Deprecated: The underlying REST API was deprecated in 3.10.0.  v8 (or newer) API should be used instead
	GetTenantById(tenantId string) (*integration.MsxResponse, error)
	// Deprecated: Tenants should generally be retrieved by ID, not tenantName.  The underlying REST API is
	// current retired and due for decommissioning.
	GetTenantByName(tenantName string) (*integration.MsxResponse, error)

	// Deprecated: The underlying REST API was deprecated in 3.10.0.  v8 (or newer) API should be used instead.
	GetUserById(userId string) (*integration.MsxResponse, error)

	GetSystemSecrets(scope string) (*integration.MsxResponse, error)
	AddSystemSecrets(scope string, secrets map[string]string) (*integration.MsxResponse, error)
	ReplaceSystemSecrets(scope string, secrets map[string]string) (*integration.MsxResponse, error)
	EncryptSystemSecrets(scope string, names []string, encrypt EncryptSecretsDTO) (*integration.MsxResponse, error)
	RemoveSystemSecrets(scope string) (*integration.MsxResponse, error)
	GenerateSystemSecrets(scope string, names []string, save bool) (*integration.MsxResponse, error)

	GetTenantSecrets(tenantId, scope string) (*integration.MsxResponse, error)
	AddTenantSecrets(tenantId, scope string, secrets map[string]string) (*integration.MsxResponse, error)
	ReplaceTenantSecrets(tenantId, scope string, secrets map[string]string) (*integration.MsxResponse, error)
	EncryptTenantSecrets(tenantId, scope string, names []string, encrypt EncryptSecretsDTO) (*integration.MsxResponse, error)
	RemoveTenantSecrets(tenantId, scope string) (*integration.MsxResponse, error)
	GenerateTenantSecrets(tenantId, scope string, names []string, save bool) (*integration.MsxResponse, error)

	GetSecretPolicy(name string) (*integration.MsxResponse, error)
	StoreSecretPolicy(name string, policy SecretPolicySetRequest) (*integration.MsxResponse, error)
	DeleteSecretPolicy(name string) (*integration.MsxResponse, error)

	GetRoles(resolvePermissionNames bool, p paging.Request) (*integration.MsxResponse, error)
	CreateRole(dbinstaller bool, body RoleCreateRequest) (*integration.MsxResponse, error)
	UpdateRole(dbinstaller bool, body RoleUpdateRequest) (*integration.MsxResponse, error)
	DeleteRole(roleName string) (*integration.MsxResponse, error)

	GetCapabilities(p paging.Request) (*integration.MsxResponse, error)
	BatchCreateCapabilities(populator bool, owner string, capabilities []CapabilityCreateRequest) (*integration.MsxResponse, error)
	BatchUpdateCapabilities(populator bool, owner string, capabilities []CapabilityUpdateRequest) (*integration.MsxResponse, error)
	DeleteCapability(populator bool, owner string, name string) (*integration.MsxResponse, error)

	GetTenantHierarchyRoot() (*integration.MsxResponse, error)
	GetTenantHierarchyParent(tenantId types.UUID) (*integration.MsxResponse, error)
}
