//go:generate mockery --inpackage --name=Api --structname=MockAwsResourceProvider

package awsresourceprovider

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

type Api interface {
	Connect(request AwsConnectRequest) (*integration.MsxResponse, error)
	GetRegions(controlPlaneId types.UUID) (*integration.MsxResponse, error)
	GetAvailabilityZones(controlPlaneId types.UUID, region string) (*integration.MsxResponse, error)
	GetResources(serviceConfigurationApplicationId types.UUID) (*integration.MsxResponse, error)
	GetVpnConnectionDetails(controlPlaneId types.UUID, vpnConnectionIds []string, region string) (*integration.MsxResponse, error)
	GetEc2InstanceStatus(controlPlaneId types.UUID, region string, instanceId string) (*integration.MsxResponse, error)
	GetTransitGatewayStatus(controlPlaneId types.UUID, region string, transitGatewayIds []string) (*integration.MsxResponse, error)
	GetTransitGatewayAttachmentStatus(controlPlaneId types.UUID, region string, transitGatewayAttachmentIds []string) (*integration.MsxResponse, error)
	GetStackOutputs(serviceConfigurationApplicationId types.UUID) (*integration.MsxResponse, error)
}
