//go:generate mockery --inpackage --name=Api --structname=MockServiceConfigManagerApi

package serviceconfigmanager

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

type Api interface {
	GetAdminHealth() (*HealthResult, error)

	// Service Configuration
	GetServiceConfigurations(page, pageSize int) (*integration.MsxResponse, error)
	GetServiceConfigurationByServiceConfigId(serviceConfigId types.UUID) (*integration.MsxResponse, error)
	CreateServiceConfiguration(configuration ServiceConfigurationRequest) (*integration.MsxResponse, error)
	UpdateServiceConfiguration(configuration ServiceConfigurationUpdateRequest) (*integration.MsxResponse, error)
	DeleteServiceConfiguration(serviceConfigId types.UUID) (*integration.MsxResponse, error)
	UpdateServiceConfigurationStatus(serviceConfigId types.UUID, serviceConfigurationStatus StatusUpdateRequest) (*integration.MsxResponse, error)

	//Assignment to tenants
	GetAllServiceConfigurationAssignments(page, pageSize int, filterTenantId types.UUID) (*integration.MsxResponse, error)
	GetServiceConfigurationAssignmentByAssignmentId(assignmentId types.UUID) (*integration.MsxResponse, error)
	CreateServiceConfigurationAssignment(serviceConfigId types.UUID, tenantIdList []types.UUID) (*integration.MsxResponse, error)
	DeleteServiceConfigurationAssignments(serviceConfigId types.UUID, tenantIdList []types.UUID) (*integration.MsxResponse, error)
	GetServiceConfigurationAssignmentsByServiceConfigurationId(serviceConfigurationId types.UUID) (*integration.MsxResponse, error)
	UpdateServiceConfigurationAssignmentStatusByServiceConfigurationIdAndTenantId(serviceConfigId types.UUID, tenantId types.UUID, status StatusUpdateRequest) (*integration.MsxResponse, error)
	GetTenantAssignmentsByServiceConfigurationId(serviceConfigurationId types.UUID) (*integration.MsxResponse, error)

	// Application of service configurations
	CreateServiceConfigurationApplication(applicationRequest ServiceConfigurationApplicationRequest) (*integration.MsxResponse, error)
	UpdateServiceConfigurationApplicationStatus(applicationId types.UUID, serviceConfigId types.UUID, applicationStatusUpdateRequest ServiceConfigurationApplicationStatusUpdateRequest) (*integration.MsxResponse, error)
	DeleteServiceConfigurationApplication(applicationId types.UUID) (*integration.MsxResponse, error)
	GetServiceConfigurationApplications(page, pageSize int, tenantId types.UUID, serviceConfigId types.UUID, sortBy, sortOrder, targetEntityId, targetEntityType *string) (*integration.MsxResponse, error)
	GetServiceConfigurationApplicationById(applicationId types.UUID) (*integration.MsxResponse, error)
	GetServiceConfigurationApplicationByServiceConfigIdTargetEntityTypeAndTargetEntityId(serviceConfigId types.UUID, targetEntityType string, targetEntityId string) (*integration.MsxResponse, error)
}
