package awsresourceprovider

import (
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

type AwsEc2InstanceStatuses struct {
	AvailabilityZone string `json:"availabilityzone"`
	InstanceId       string `json:"instanceid"`
	InstanceState    struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
	} `json:"instancestate"`
	InstanceStatus struct {
		Details struct {
			Name   string `json:"name"`
			Status string `json:"status"`
		} `json:"details"`
		Status string `json:"status"`
	} `json:"instancestatus"`
	SystemStatus struct {
		Details struct {
			Name   string `json:"name"`
			Status string `json:"status"`
		} `json:"details"`
		Status string `json:"status"`
	} `json:"systemstatus"`
}
