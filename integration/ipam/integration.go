package ipam

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"net"
	"net/url"
	"strconv"
)

const (
	serviceName              = integration.ServiceNameIpam
	endpointNameRegisterCIDR = "registerCIDR"
	endpointNameGetCIDRs     = "getCIDRs"
	endpointNameAcquireIP    = "acquireIP"
	endpointNameReleaseIP    = "releaseIP"
)

var (
	logger    = log.NewLogger("msx.integration.ipamservice")
	endpoints = map[string]integration.MsxServiceEndpoint{
		endpointNameRegisterCIDR: {Method: "POST", Path: "/api/v1/cidrs"},
		endpointNameGetCIDRs:     {Method: "GET", Path: "/api/v1/cidrs"},
		endpointNameAcquireIP:    {Method: "PUT", Path: "/api/v1/ips"},
		endpointNameReleaseIP:    {Method: "DELETE", Path: "/api/v1/ips"},
	}
)

type Integration struct {
	integration.MsxServiceExecutor
}

func (i Integration) AquireIP(requestCIDR IpamCIDRRequest) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(requestCIDR)
	if err != nil {
		return nil, err
	}
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameAcquireIP,
		Body:         bodyBytes,
		Payload:      new(IpamIPResponse),
	})
}

func (i Integration) ReleaseIP(cidr net.IPNet, ipAddress net.IP, tenantId types.UUID) (*integration.MsxResponse, error) {
	queryParameters := make(url.Values)
	queryParameters["cidr"] = []string{cidr.String()}
	queryParameters["ip"] = []string{ipAddress.String()}
	queryParameters["tenantId"] = []string{tenantId.String()}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:    endpointNameReleaseIP,
		QueryParameters: queryParameters,
	})

}

func (i Integration) GetCIDRs(page, pageSize int, tenantId types.UUID) (*integration.MsxResponse, error) {
	queryParameters := make(url.Values)
	queryParameters["page"] = []string{strconv.Itoa(page)}
	queryParameters["pageSize"] = []string{strconv.Itoa(pageSize)}
	queryParameters["tenantId"] = []string{tenantId.String()}

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:    endpointNameGetCIDRs,
		QueryParameters: queryParameters,
		Payload: &paging.PaginatedResponseV8{
			Contents: new(IpamCIDRListResponse),
		},
	})
}

func (i Integration) RegisterCIDR(requestCIDR IpamCIDRRequest) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(requestCIDR)
	if err != nil {
		return nil, err
	}
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameRegisterCIDR,
		Body:         bodyBytes,
	})

}

func NewIntegration(ctx context.Context) (Api, error) {
	integrationInstance := IntegrationFromContext(ctx)
	if integrationInstance == nil {
		integrationInstance = &Integration{
			MsxServiceExecutor: integration.NewMsxService(ctx, serviceName, endpoints),
		}
	}
	return integrationInstance, nil
}

func NewIntegrationWithExecutor(executor integration.MsxServiceExecutor) Api {
	return &Integration{
		MsxServiceExecutor: executor,
	}
}
