package api

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	validation "github.com/go-ozzo/ozzo-validation"
)

type UpperCamelSingularCreateRequest struct {
	//#if TENANT_DOMAIN
	TenantId types.UUID `json:"tenantId"`
	//#endif TENANT_DOMAIN
	Data string `json:"data" san:"xss"`
}

func (r *UpperCamelSingularCreateRequest) Validate() error {
	return types.ErrorMap{
		//#if TENANT_DOMAIN
		"tenantId": validation.Validate(&r.TenantId, validation.Required),
		//#endif TENANT_DOMAIN
		"data": validation.Validate(&r.Data, validation.Required),
	}
}

type UpperCamelSingularUpdateRequest struct {
	Data string `json:"data" san:"xss"`
}

func (r *UpperCamelSingularUpdateRequest) Validate() error {
	return types.ErrorMap{
		"data": validation.Validate(&r.Data, validation.Required),
	}
}

type UpperCamelSingularResponse struct {
	UpperCamelSingularId types.UUID `json:"lowerCamelSingularId"`
	//#if TENANT_DOMAIN
	TenantId types.UUID `json:"tenantId"`
	//#endif TENANT_DOMAIN
	Data string `json:"data"`
}
