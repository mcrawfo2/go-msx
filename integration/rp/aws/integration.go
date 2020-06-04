package aws

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
)

const (
	endpointNameConnect              = "connect"
	endpointNameGetRegions           = "getRegions"
	endpointNameGetAvailabilityZones = "getAvailabilityZones"
	endpointNameGetResources         = "getResources"

	serviceName = integration.ResourceProviderNameAws
)

var (
	logger    = log.NewLogger("msx.integration.rp.aws")
	endpoints = map[string]integration.MsxServiceEndpoint{
		endpointNameConnect:              {Method: "POST", Path: "/api/v1/connect"},
		endpointNameGetRegions:           {Method: "GET", Path: "/api/v1/regions"},
		endpointNameGetAvailabilityZones: {Method: "GET", Path: "/api/v1/availabilityzones"},
		endpointNameGetResources:         {Method: "GET", Path: "/api/v1/resources"},
	}
)

func NewIntegration(ctx context.Context) (Api, error) {
	return &Integration{
		MsxService: integration.NewMsxServiceResourceProvider(ctx, serviceName, endpoints),
	}, nil
}

type Integration struct {
	*integration.MsxService
}

func (i *Integration) Connect(request AwsConnectRequest) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	var payload = ""

	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:   endpointNameConnect,
		Body:           bodyBytes,
		ExpectEnvelope: true,
		Payload:        &payload,
	})
}

func (i *Integration) GetRegions(controlPlaneId types.UUID) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetRegions,
		QueryParameters: map[string][]string{
			"controlPlaneId": {controlPlaneId.String()},
		},
		ExpectEnvelope: true,
		Payload:        &[]Region{},
	})
}

func (i *Integration) GetAvailabilityZones(controlPlaneId types.UUID, region string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetAvailabilityZones,
		QueryParameters: map[string][]string{
			"controlPlaneId": {controlPlaneId.String()},
			"region":         {region},
		},
		ExpectEnvelope: true,
		Payload:        &[]AvailabilityZone{},
	})
}

func (i *Integration) GetResources(serviceConfigurationApplicationId types.UUID) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetResources,
		QueryParameters: map[string][]string{
			"serviceConfigurationApplicationId": {serviceConfigurationApplicationId.String()},
		},
		ExpectEnvelope: true,
	})
}
