package serviceconfigmanager

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"time"
)

type Pojo integration.Pojo
type PojoArray integration.PojoArray
type HealthResult integration.HealthResult
type ErrorDTO integration.ErrorDTO
type ErrorDTO2 integration.ErrorDTO2

type ServiceConfigurationRequest struct {
	Attributes    map[string]string `json:"attributes"`
	Configuration string            `json:"configuration"`
	Description   string            `json:"description"`
	Name          string            `json:"name"`
	Notes         string            `json:"notes"`
	Service       string            `json:"service"`
	Status        string            `json:"status"`
	StatusDetails string            `json:"statusDetails"`
	SubType       string            `json:"subType"`
	Tags          []string          `json:"tags"`
	Type          string            `json:"type"`
}

type ServiceConfigurationUpdateRequest struct {
	*ServiceConfigurationRequest
	ServiceConfigId string `json:"serviceConfigId"`
}

type ServiceConfigurationResponse struct {
	Name                 string             `json:"name"`
	Description          *string            `json:"description"`
	Service              string             `json:"service"`
	Type                 string             `json:"type"`
	SubType              *string            `json:"subType"`
	Configuration        string             `json:"configuration"`
	Attributes           *map[string]string `json:"attributes"`
	Tags                 *[]string          `json:"tags"`
	Notes                *string            `json:"notes"`
	Status               string             `json:"status"`
	StatusDetails        *string            `json:"statusDetails"`
	Timestamp            time.Time          `json:"timestamp"`
	ModifiedDate         time.Time          `json:"modifiedDate"`
	EventActorID         string             `json:"eventActorId"`
	EventActorUsername   string             `json:"eventActorUsername"`
	Version              string             `json:"version"`
	ServiceConfigID      types.UUID         `json:"serviceConfigId"`
	EventActorTenantID   string             `json:"eventActorTenantId"`
	EventActorTenantName string             `json:"eventActorTenantName"`
}

type StatusUpdateRequest struct {
	Status        string `json:"status"`
	StatusDetails string `json:"statusDetails"`
}

type ServiceConfigurationAssignmentRequest struct {
	Tenants []types.UUID `json:"tenants"`
}

type ServiceConfigurationAssignmentResponse struct {
	AssignmentID       string  `json:"assignmentId"`
	ServiceConfigID    string  `json:"serviceConfigId"`
	TenantID           string  `json:"tenantId"`
	TenantName         string  `json:"tenantName"`
	Name               string  `json:"name"`
	Description        *string `json:"description"`
	Service            string  `json:"service"`
	Type               string  `json:"type"`
	SubType            *string `json:"subType"`
	Version            *string `json:"version"`
	Status             *string `json:"status"`
	StatusDetails      *string `json:"statusDetails"`
	AssignedTenantID   string  `json:"assignedTenantId"`
	AssignedTenantName string  `json:"assignedTenantName"`
}

type ServiceConfigurationApplicationRequest struct {
	Parameters       map[string]string `json:"parameters"`
	ServiceConfigID  types.UUID        `json:"serviceConfigId"`
	Status           string            `json:"status"`
	StatusDetails    string            `json:"statusDetails"`
	TargetEntityID   string            `json:"targetEntityId"`
	TargetEntityType string            `json:"targetEntityType"`
	TenantID         types.UUID        `json:"tenantId"`
}

type ServiceConfigurationApplicationStatusUpdateRequest struct {
	ApplicationID types.UUID `json:"applicationId"`
	Status        string     `json:"status"`
	StatusDetails string     `json:"statusDetails"`
}

type ServiceConfigurationApplicationResponse struct {
	ID               types.UUID        `json:"id"`
	ServiceConfigID  types.UUID        `json:"serviceConfigId"`
	TenantID         types.UUID        `json:"tenantId"`
	Timestamp        time.Time         `json:"timestamp"`
	ModifiedDate     time.Time         `json:"modifiedDate"`
	Status           string            `json:"status"`
	StatusDetails    *string           `json:"statusDetails"`
	Parameters       map[string]string `json:"parameters"`
	TargetEntityID   string            `json:"targetEntityId"`
	TargetEntityType string            `json:"targetEntityType"`
}
