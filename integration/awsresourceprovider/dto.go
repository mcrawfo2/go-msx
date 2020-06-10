package awsresourceprovider

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"time"
)

type ControlPlaneAuthenticationType string

const (
	Local ControlPlaneAuthenticationType = "local"
	Sso   ControlPlaneAuthenticationType = "sso"
)

type AwsConnectRequest struct {
	Url         string `json:"url"`
	TlsInsecure bool   `json:"tlsInsecure"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type Region struct {
	Regionname string `json:"regionname"`
	Endpoint   string `json:"endpoint"`
}

type AvailabilityZone struct {
	Regionname string `json:"regionname"`
	State      string `json:"state"`
	Zonename   string `json:"zonename"`
}

type Resource struct {
	LastUpdatedTimestamp time.Time `json:"LastUpdatedTimestamp"`
	LogicalResourceId    *string   `json:"LogicalResourceId"`
	PhysicalResourceId   *string   `json:"PhysicalResourceId"`
	ResourceStatus       *string   `json:"ResourceStatus"`
	ResourceStatusReason *string   `json:"ResourceStatusReason"`
	ResourceType         *string   `json:"ResourceType"`
}

type VpnConnection struct {
	CustomerGatewayConfiguration *string `json:"customerGatewayConfiguration"`
	CustomerGatewayId            *string `json:"customerGatewayId"`
	VpnConnectionId              *string `json:"vpnConnectionId"`
	VpnGatewayId                 *string `json:"vpnGatewayId"`
}

type Pojo integration.Pojo
