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
	AvailabilityZone string `json:"availabilityZone"`
	InstanceId       string `json:"instanceId"`
	InstanceState    struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"instanceState"`
	InstanceStatus AwsEc2InstanceSubStatus `json:"instanceStatus"`
	SystemStatus   AwsEc2InstanceSubStatus `json:"systemStatus"`
}

type AwsEc2InstanceSubStatus struct {
	Details AwsEc2InstanceSubStatusDetails `json:"details"`
	Status  string                         `json:"status"`
}

type AwsEc2InstanceSubStatusDetails struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type AwsTransitGatewayStatus struct {
	CreationTime time.Time `json:"CreationTime"`
	Description  string    `json:"Description"`
	Options      struct {
		AmazonSideAsn                  int         `json:"AmazonSideAsn"`
		AssociationDefaultRouteTableID string      `json:"AssociationDefaultRouteTableId"`
		AutoAcceptSharedAttachments    string      `json:"AutoAcceptSharedAttachments"`
		DefaultRouteTableAssociation   string      `json:"DefaultRouteTableAssociation"`
		DefaultRouteTablePropagation   string      `json:"DefaultRouteTablePropagation"`
		DNSSupport                     string      `json:"DnsSupport"`
		MulticastSupport               interface{} `json:"MulticastSupport"`
		PropagationDefaultRouteTableID string      `json:"PropagationDefaultRouteTableId"`
		VpnEcmpSupport                 string      `json:"VpnEcmpSupport"`
	} `json:"Options"`
	OwnerID string `json:"OwnerId"`
	State   string `json:"State"`
	Tags    []struct {
		Key   string `json:"Key"`
		Value string `json:"Value"`
	} `json:"Tags"`
	TransitGatewayArn string `json:"TransitGatewayArn"`
	TransitGatewayID  string `json:"TransitGatewayId"`
}

type AwsTransitGatewayAttachmentStatus struct {
	Association struct {
		State                      string `json:"State"`
		TransitGatewayRouteTableID string `json:"TransitGatewayRouteTableId"`
	} `json:"Association"`
	CreationTime    time.Time `json:"CreationTime"`
	ResourceID      string    `json:"ResourceId"`
	ResourceOwnerID string    `json:"ResourceOwnerId"`
	ResourceType    string    `json:"ResourceType"`
	State           string    `json:"State"`
	Tags            []struct {
		Key   string `json:"Key"`
		Value string `json:"Value"`
	} `json:"Tags"`
	TransitGatewayAttachmentID string `json:"TransitGatewayAttachmentId"`
	TransitGatewayID           string `json:"TransitGatewayId"`
	TransitGatewayOwnerID      string `json:"TransitGatewayOwnerId"`
}

type StackOutput struct {
	Description string `json:"Description"`
	ExportName  string `json:"ExportName"`
	OutputKey   string `json:"OutputKey"`
	OutputValue string `json:"OutputValue"`
}
