package dnac

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"encoding/json"
)

const (
	endpointNameConnect              = "connect"
	endpointNameRetrieveExtendedData = "retrieveExtendedData"

	serviceName        = integration.ResourceProviderNameDnac
)

var (
	logger    = log.NewLogger("msx.integration.rp.dnac")
	endpoints = map[string]integration.MsxServiceEndpoint{
		endpointNameConnect:              {Method: "POST", Path: "/api/v1/connect"},
		endpointNameRetrieveExtendedData: {Method: "POST", Path: "/api/v1/retrieveExtendedData"},
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

func (i *Integration) Connect(request DnacConnectRequest) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	var payload = ""

	return i.Execute(&integration.MsxRequest{
		EndpointName:       endpointNameConnect,
		Body:               bodyBytes,
		ExpectEnvelope:     true,
		Payload:            &payload,
	})
}

func (i *Integration) RetrieveExtendedData(request DnacExtendedRequest) (*integration.MsxResponse, error) {
	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return i.Execute(&integration.MsxRequest{
		EndpointName:       endpointNameRetrieveExtendedData,
		Body:               bodyBytes,
		ExpectEnvelope:     true,
		Payload:            &Pojo{},
	})
}

