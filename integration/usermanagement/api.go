package usermanagement

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
)

type Api interface {
	GetAdminHealth() (*HealthResult, error)

	Login(user, password string) (*integration.MsxResponse, error)
	Logout() (*integration.MsxResponse, error)

	GetMyCapabilities() (*integration.MsxResponse, error)
	GetUserCapabilities(userId string) (*integration.MsxResponse, error)

	GetTenantIds() (*integration.MsxResponse, error)
	GetMyTenants() (*integration.MsxResponse, error)
	GetUserTenants(userId string) (*integration.MsxResponse, error)
	GetTenantById(tenantId string) (*integration.MsxResponse, error)
	GetTenantByName(tenantName string) (*integration.MsxResponse, error)

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
}
