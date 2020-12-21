package ipam

import "cto-github.cisco.com/NFV-BU/go-msx/types"

type IpamCIDRRequest struct {
	CIDR     string     `json:"cidr"`
	TenantID types.UUID `json:"tenantId"`
}

type IpamCIDRResponse struct {
	CIDR     string     `json:"cidr"`
	TenantID types.UUID `json:"tenantId"`
}

type IpamCIDRListResponse []IpamCIDRResponse

type IpamIPResponse struct {
	CIDR     string     `json:"cidr"`
	TenantID types.UUID `json:"tenantId"`
	IP       string     `json:"ip"`
}
