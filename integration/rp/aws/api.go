package aws

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

type Api interface {
	Connect(request AwsConnectRequest) (*integration.MsxResponse, error)
	GetRegions(controlPlaneId types.UUID) (*integration.MsxResponse, error)
	GetAvailabilityZones(controlPlaneId types.UUID, region string) (*integration.MsxResponse, error)
	GetResources(serviceConfigurationApplicationId types.UUID) (*integration.MsxResponse, error)
}
