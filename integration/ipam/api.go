//go:generate mockery --inpackage --name=Api --structname=MockApi
package ipam

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"net"
)

type Api interface {
	RegisterCIDR(requestCIDR IpamCIDRRequest) (*integration.MsxResponse, error)
	GetCIDRs(page, pageSize int, tenantId types.UUID) (*integration.MsxResponse, error)
	AquireIP(requestCIDR IpamCIDRRequest) (*integration.MsxResponse, error)
	ReleaseIP(cidr net.IPNet, ipAddress net.IP, tenantId types.UUID) (*integration.MsxResponse, error)
}
