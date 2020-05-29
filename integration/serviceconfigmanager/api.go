package serviceconfigmanager

import "cto-github.cisco.com/NFV-BU/go-msx/integration"

type Api interface {
	GetAdminHealth() (*HealthResult, error)

	GetServiceConfigurations(page, pageSize int) (*integration.MsxResponse, error)
	GetServiceConfigurationByServiceConfigId(serviceConfigId string) (*integration.MsxResponse, error)
	CreateServiceConfiguration(configuration ServiceConfigurationRequest) (*integration.MsxResponse, error)
	UpdateServiceConfiguration(configuration ServiceConfigurationUpdateRequest) (*integration.MsxResponse, error)
	DeleteServiceConfiguration(serviceConfigId string) (*integration.MsxResponse, error)
	UpdateServiceConfigurationStatus(serviceConfigId string, serviceConfigurationStatus ServiceConfigurationStatusRequest) (*integration.MsxResponse, error)
}
