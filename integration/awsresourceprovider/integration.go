package awsresourceprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
)

const (
	endpointNameConnect                           = "connect"
	endpointNameGetRegions                        = "getRegions"
	endpointNameGetAvailabilityZones              = "getAvailabilityZones"
	endpointNameGetResources                      = "getResources"
	endpointNameGetVpnConnections                 = "getVpnConnections"
	endpointNameGetEc2InstanceStatus              = "getEc2InstanceStatus"
	endpointNameGetTransitGatewayStatus           = "getTransitGatewayStatus"
	endpointNameGetTransitGatewayAttachmentStatus = "getTransitGatewayAttachmentStatus"
	endpointNameGetStackOutput                    = "getStackOutput"
	serviceName                                   = integration.ResourceProviderNameAws
)

var (
	logger    = log.NewLogger("msx.integration.rp.aws")
	endpoints = map[string]integration.MsxServiceEndpoint{
		endpointNameConnect:                           {Method: "POST", Path: "/api/v1/connect"},
		endpointNameGetRegions:                        {Method: "GET", Path: "/api/v1/regions"},
		endpointNameGetAvailabilityZones:              {Method: "GET", Path: "/api/v1/availabilityzones"},
		endpointNameGetResources:                      {Method: "GET", Path: "/api/v1/resources"},
		endpointNameGetVpnConnections:                 {Method: "GET", Path: "/api/v1/vpnconnection"},
		endpointNameGetTransitGatewayStatus:           {Method: "GET", Path: "/api/v1/transitgateway/status"},
		endpointNameGetTransitGatewayAttachmentStatus: {Method: "GET", Path: "/api/v1/transitgatewayattachment/status"},
		endpointNameGetEc2InstanceStatus:              {Method: "GET", Path: "/api/v1/ec2instance/status"},
		endpointNameGetStackOutput:                    {Method: "GET", Path: "/api/v1/stack/outputs"},
	}
)

func NewIntegration(ctx context.Context) (Api, error) {
	integrationInstance := IntegrationFromContext(ctx)
	if integrationInstance == nil {
		integrationInstance = &Integration{
			MsxService: integration.NewMsxServiceResourceProvider(ctx, serviceName, endpoints),
		}
	}
	return integrationInstance, nil
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
		Payload:        &[]Resource{},
	})
}

func (i *Integration) GetVpnConnectionDetails(controlPlaneId types.UUID, vpnConnectionIds []string, region string) (*integration.MsxResponse, error) {
	queryParams := map[string][]string{
		"controlPlaneId":   {controlPlaneId.String()},
		"region":           {region},
		"vpnConnectionIds": vpnConnectionIds,
	}
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName:    endpointNameGetVpnConnections,
		QueryParameters: queryParams,
		ExpectEnvelope:  true,
		Payload:         &[]VpnConnection{},
	})
}

func (i *Integration) GetEc2InstanceStatus(controlPlaneId types.UUID, region string, instanceId string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetEc2InstanceStatus,
		QueryParameters: map[string][]string{
			"controlPlaneId": {controlPlaneId.String()},
			"region":         {region},
			"instanceId":     {instanceId},
		},
		ExpectEnvelope: true,
		Payload:        &AwsEc2InstanceStatuses{},
	})
}

func (i *Integration) GetTransitGatewayStatus(controlPlaneId types.UUID, region string, transitGatewayIds []string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetTransitGatewayStatus,
		QueryParameters: map[string][]string{
			"controlPlaneId":   {controlPlaneId.String()},
			"region":           {region},
			"transitGatewayId": transitGatewayIds,
		},
		ExpectEnvelope: true,
		Payload:        &[]AwsTransitGatewayStatus{},
	})
}

func (i *Integration) GetTransitGatewayAttachmentStatus(controlPlaneId types.UUID, region string, transitGatewayAttachmentIds []string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetTransitGatewayAttachmentStatus,
		QueryParameters: map[string][]string{
			"controlPlaneId":             {controlPlaneId.String()},
			"region":                     {region},
			"transitGatewayAttachmentId": transitGatewayAttachmentIds,
		},
		ExpectEnvelope: true,
		Payload:        &[]AwsTransitGatewayAttachmentStatus{},
	})
}

func (i *Integration) GetStackOutputs(serviceConfigurationApplicationId types.UUID) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetStackOutput,
		QueryParameters: map[string][]string{
			"serviceConfigurationApplicationId": {serviceConfigurationApplicationId.String()},
		},
		ExpectEnvelope: true,
		Payload:        &[]StackOutput{},
	})
}
