package monitor

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
)

const (
	endpointNameGetDevicesHealth = "getDevicesHealth"
	serviceName                  = integration.ServiceNameMonitor
)

var (
	logger    = log.NewLogger("msx.integration.monitor")
	endpoints = map[string]integration.MsxServiceEndpoint{
		endpointNameGetDevicesHealth: {Method: "GET", Path: "/api/v2/health/devices"},
	}
)

type Integration struct {
	integration.MsxServiceExecutor
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

func (i Integration) GetDeviceHealth(deviceIds string) (*integration.MsxResponse, error) {
	return i.Execute(&integration.MsxEndpointRequest{
		EndpointName: endpointNameGetDevicesHealth,
		QueryParameters: map[string][]string{
			"deviceIds": {deviceIds},
		},
		Payload:        new(Pojo),
		ExpectEnvelope: true,
	})
}
