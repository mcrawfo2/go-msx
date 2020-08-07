package serviceconfigupdate

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	TopicServiceConfigUpdateTopic = "SERVICECONFIG_UPDATE_TOPIC"

	MetaDataEventType = "eventType"
	MetaDataService   = "service"

	EventTypeApplicationStatusUpdate = "UpdateServiceConfigurationApplicationStatusEvent"
	EventTypeUpdate                  = "UpdateServiceConfigurationEvent"

	ApplicationStatusSuccess      = "applied"
	ApplicationStatusFailed       = "apply failed"
	ApplicationStatusDeleting     = "deleting"
	ApplicationStatusDeleteFailed = "delete failed"
	ApplicationStatusDeleted      = "deleted"
)

type actionRequest struct {
	Request interface{} `json:"request"`
}

type ApplicationStatusUpdateRequest struct {
	ApplicationId   types.UUID `json:"applicationId"`
	ServiceConfigId types.UUID `json:"serviceConfigId"`
	Service         string     `json:"service"`
	Status          string     `json:"status"`
	StatusDetails   string     `json:"statusDetails"`
}

func (a ApplicationStatusUpdateRequest) Validate() error {
	return types.ErrorMap{
		"applicationId":   validation.Validate(&a.ApplicationId, validation.Required, validate.Self),
		"serviceConfigId": validation.Validate(&a.ServiceConfigId, validation.Required, validate.Self),
		"status":          validation.Validate(&a.Status, validation.Required),
	}.Filter()
}

type UpdateRequest struct {
	ServiceConfigId types.UUID        `json:"serviceConfigId"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Service         string            `json:"service"`
	Type            string            `json:"type"`
	SubType         string            `json:"subType"`
	Configuration   string            `json:"configuration"`
	Notes           string            `json:"notes"`
	Tags            []string          `json:"tags"`
	Attributes      map[string]string `json:"attributes"`
	Status          string            `json:"status"`
	StatusDetails   string            `json:"statusDetails"`
}

func (u UpdateRequest) Validate() error {
	return types.ErrorMap{
		"serviceConfigId": validation.Validate(&u.ServiceConfigId, validation.Required, validate.Self),
		"name":            validation.Validate(&u.Name, validation.Required),
		"description":     validation.Validate(&u.Description, validation.Length(0, 500)),
		"service":         validation.Validate(&u.Service, validation.Required),
		"type":            validation.Validate(&u.Type, validation.Required),
		"subType":         validation.Validate(&u.SubType),
		"configuration":   validation.Validate(&u.Configuration, validation.Required, validation.Length(0, 33554432)),
		"notes":           validation.Validate(&u.Notes),
		"tags":            validation.Validate(&u.Tags),
		"attributes":      validation.Validate(&u.Attributes),
		"status":          validation.Validate(&u.Status, validation.Required),
		"statusDetails":   validation.Validate(&u.StatusDetails),
	}.Filter()
}
