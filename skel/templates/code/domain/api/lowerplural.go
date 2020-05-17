package api

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	validation "github.com/go-ozzo/ozzo-validation"
)

type UpperCamelSingularCreateRequest struct {
	Name string `json:"name"`
	//#if TENANT_DOMAIN
	TenantId types.UUID `json:"tenantId"`
	//#endif TENANT_DOMAIN
	Data string `json:"data"`
}

func (r *UpperCamelSingularCreateRequest) Validate() error {
	return types.ErrorMap{
		"name": validation.Validate(&r.Name, validation.Required),
		//#if TENANT_DOMAIN
		"tenantId": validation.Validate(&r.TenantId, validation.Required),
		//#endif TENANT_DOMAIN
		"data": validation.Validate(&r.Data, validation.Required),
	}
}

type UpperCamelSingularUpdateRequest struct {
	Data string `json:"data"`
}

func (r *UpperCamelSingularUpdateRequest) Validate() error {
	return types.ErrorMap{
		"data": validation.Validate(&r.Data, validation.Required),
	}
}

type UpperCamelSingularResponse struct {
	Name string `json:"name"`
	//#if TENANT_DOMAIN
	TenantId types.UUID `json:"tenantId"`
	//#endif TENANT_DOMAIN
	Data string `json:"data"`
}
