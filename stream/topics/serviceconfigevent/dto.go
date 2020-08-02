package serviceconfigevent

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	MetaDataEventType    = "eventType"
	MetaDataService      = "service"
	MetaDataFrom         = "from"
	MetaDataPartitionKey = "partitionKey"

	// ServiceConfiguration
	EventTypeCreated       = "CreatedServiceConfigurationEvent"
	EventTypeUpdated       = "UpdatedServiceConfigurationEvent"
	EventTypeDeleted       = "DeletedServiceConfigurationEvent"
	EventTypeStatusUpdated = "UpdatedServiceConfigurationStatusEvent"

	// ServiceConfigurationAssignment
	EventTypeAssignmentCreated       = "CreatedServiceConfigurationAssignmentEvent"
	EventTypeAssignmentDeleted       = "DeletedServiceConfigurationAssignmentEvent"
	EventTypeAssignmentStatusUpdated = "UpdatedServiceConfigurationAssignmentStatusEvent"

	// ServiceConfigurationApplication
	EventTypeApplicationCreated       = "CreatedServiceConfigurationApplicationEvent"
	EventTypeApplicationDeleted       = "DeletedServiceConfigurationApplicationEvent"
	EventTypeApplicationStatusUpdated = "UpdatedServiceConfigurationApplicationStatusEvent"
)

type Event struct {
	Headers              map[string]string    `json:"-"`
	EventId              types.UUID           `json:"eventId"`
	Status               string               `json:"status"`
	StatusDetails        string               `json:"statusDetails"`
	ServiceConfiguration ServiceConfiguration `json:"serviceConfiguration"`
	TenantId             types.UUID           `json:"tenantId"`
	EventActorId         types.UUID           `json:"eventActorId"`
	EventActorName       string               `json:"eventActorUsername"`
	Timestamp            types.Time           `json:"timestamp"`
}

func (e Event) Validate() error {
	return types.ErrorMap{
		"eventId":              validation.Validate(&e.EventId, validation.Required, validate.Self),
		"status":               validation.Validate(&e.Status, validation.Required),
		"statusDetails":        validation.Validate(&e.StatusDetails),
		"serviceConfiguration": validation.Validate(&e.ServiceConfiguration, validate.Self),
		"tenantId":             validation.Validate(&e.TenantId, validation.Required, validate.Self),
	}.Filter()
}

func (e Event) EventType() string {
	return e.Headers[MetaDataEventType]
}

func (e Event) Service() string {
	return e.Headers[MetaDataService]
}

func (e Event) From() string {
	return e.Headers[MetaDataFrom]
}

func (e Event) PartitionKey() string {
	return e.Headers[MetaDataPartitionKey]
}

type ServiceConfiguration struct {
	ServiceConfigID types.UUID        `json:"serviceConfigId"`
	Name            string            `json:"name"`
	Description     *string           `json:"description"`
	Service         string            `json:"service"`
	Type            string            `json:"type"`
	SubType         *string           `json:"subType"`
	Configuration   string            `json:"configuration"`
	Notes           *string           `json:"notes"`
	Tags            []string          `json:"tags"`
	Attributes      map[string]string `json:"attributes"`
	TenantId        types.UUID        `json:"tenantId"`
}

func (s ServiceConfiguration) Validate() error {
	return types.ErrorMap{
		// TODO
	}.Filter()
}

type AssignmentEvent struct {
	Event
	AssignedTenantId types.UUID `json:"assignedTenantId"`
}

func (s AssignmentEvent) Validate() error {
	return types.ErrorList{
		s.Event.Validate(),
		types.ErrorMap{
			"assignedTenantId": validation.Validate(&s.AssignedTenantId, validation.Required, validate.Self),
		}.Filter(),
	}.Filter()
}

type ApplicationEvent struct {
	Event
	ApplicationId    types.UUID        `json:"applicationId"`
	Parameters       map[string]string `json:"parameters"`
	TargetEntityId   string            `json:"targetEntityId"`
	TargetEntityType string            `json:"targetEntityType"`
	Version          *int              `json:"version"`
}

func (s ApplicationEvent) Validate() error {
	return types.ErrorList{
		s.Event.Validate(),
		types.ErrorMap{
			"applicationId":    validation.Validate(&s.ApplicationId, validation.Required, validate.Self),
			"parameters":       validation.Validate(&s.Parameters),
			"targetEntityId":   validation.Validate(&s.TargetEntityId, validation.Required),
			"targetEntityType": validation.Validate(&s.TargetEntityType, validation.Required),
			"version":          validation.Validate(&s.Version),
		}.Filter(),
	}.Filter()
}
