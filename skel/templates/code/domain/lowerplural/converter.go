package lowerplural

import "cto-github.cisco.com/NFV-BU/go-msx/skel/templates/code/domain/api"

type lowerCamelSingularConverter struct{}

func (c *lowerCamelSingularConverter) FromCreateRequest(request api.UpperCamelSingularCreateRequest) lowerCamelSingular {
	return lowerCamelSingular{
		Name: request.Name,
		//#if TENANT_DOMAIN
		TenantId: request.TenantId.ToByteArray(),
		//#endif TENANT_DOMAIN
		Data: request.Data,
	}
}

func (c *lowerCamelSingularConverter) FromUpdateRequest(target lowerCamelSingular, request api.UpperCamelSingularUpdateRequest) lowerCamelSingular {
	target.Data = request.Data
	return target
}

func (c *lowerCamelSingularConverter) ToUpperCamelSingularListResponse(sources []lowerCamelSingular) (results []api.UpperCamelSingularResponse) {
	results = []api.UpperCamelSingularResponse{}
	for _, source := range sources {
		results = append(results, c.ToUpperCamelSingularResponse(source))
	}
	return
}

func (c *lowerCamelSingularConverter) ToUpperCamelSingularResponse(source lowerCamelSingular) api.UpperCamelSingularResponse {
	return api.UpperCamelSingularResponse{
		Name: source.Name,
		//#if TENANT_DOMAIN
		TenantId: source.TenantId[:],
		//#endif TENANT_DOMAIN
		Data: source.Data,
	}
}
