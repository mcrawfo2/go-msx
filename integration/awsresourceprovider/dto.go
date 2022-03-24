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

type CheckStatusRequest struct {
	CallbackId string `json:"CallbackId"`
}

type Region struct {
	Regionname string `json:"regionname"`
	Endpoint   string `json:"endpoint"`
	AmiId      string `json:"amiId"`
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
	CustomerGatewayConfiguration *string      `json:"customerGatewayConfiguration"`
	CustomerGatewayId            *string      `json:"customerGatewayId"`
	VpnConnectionId              *string      `json:"vpnConnectionId"`
	VpnGatewayId                 *string      `json:"vpnGatewayId"`
	VpnTunnels                   []*VpnTunnel `json:"vpnTunnels"`
}

type VpnTunnel struct {
	AcceptedRouteCount *int64  `json:"acceptedRouteCount"`
	OutsideIpAddress   *string `json:"outsideIpAddress"`
	InsideIpv4Cidr     *string `json:"insideIpv4Cidr"`
	Status             *string `json:"status"`
	StatusMessage      *string `json:"statusMessage"`
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
	OwnerID           string `json:"OwnerId"`
	State             string `json:"State"`
	Tags              []Tag  `json:"Tags"`
	TransitGatewayArn string `json:"TransitGatewayArn"`
	TransitGatewayID  string `json:"TransitGatewayId"`
}

type AwsTransitGatewayAttachmentStatus struct {
	Association struct {
		State                      string `json:"State"`
		TransitGatewayRouteTableID string `json:"TransitGatewayRouteTableId"`
	} `json:"Association"`
	CreationTime               time.Time `json:"CreationTime"`
	ResourceID                 string    `json:"ResourceId"`
	ResourceOwnerID            string    `json:"ResourceOwnerId"`
	ResourceType               string    `json:"ResourceType"`
	State                      string    `json:"State"`
	Tags                       []Tag     `json:"Tags"`
	TransitGatewayAttachmentID string    `json:"TransitGatewayAttachmentId"`
	TransitGatewayID           string    `json:"TransitGatewayId"`
	TransitGatewayOwnerID      string    `json:"TransitGatewayOwnerId"`
}

type AwsTransitVPCStatus struct {
	CidrBlock               string `json:"CidrBlock"`
	CidrBlockAssociationSet []struct {
		AssociationID  string `json:"AssociationId"`
		CidrBlock      string `json:"CidrBlock"`
		CidrBlockState struct {
			State         string      `json:"State"`
			StatusMessage interface{} `json:"StatusMessage"`
		} `json:"CidrBlockState"`
	} `json:"CidrBlockAssociationSet"`
	DhcpOptionsID               string      `json:"DhcpOptionsId"`
	InstanceTenancy             string      `json:"InstanceTenancy"`
	Ipv6CidrBlockAssociationSet interface{} `json:"Ipv6CidrBlockAssociationSet"`
	IsDefault                   bool        `json:"IsDefault"`
	OwnerID                     string      `json:"OwnerId"`
	State                       string      `json:"State"`
	Tags                        []Tag       `json:"Tags"`
	VpcID                       string      `json:"VpcId"`
}

type Tag struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type StackOutput struct {
	Description *string `json:"Description"`
	OutputKey   *string `json:"OutputKey"`
	OutputValue *string `json:"OutputValue"`
}

type AwsAmiRegion struct {
	AmiName string      `json:"amiName"`
	Regions []AwsRegion `json:"regions"`
}

type AwsRegion struct {
	RegionName string `json:"regionname"`
	EndPoint   string `json:"endpoint"`
	AMIId      string `json:"amiId"`
}

type VpcRouteTable struct {
	Associations []struct {
		AssociationState struct {
			State         string  `json:"State"`
			StatusMessage *string `json:"StatusMessage"`
		} `json:"AssociationState"`
		GatewayID               *string `json:"GatewayId"`
		Main                    bool    `json:"Main"`
		RouteTableAssociationID string  `json:"RouteTableAssociationId"`
		RouteTableID            string  `json:"RouteTableId"`
		SubnetID                string  `json:"SubnetId"`
	} `json:"Associations"`
	OwnerID         string  `json:"OwnerId"`
	PropagatingVgws *string `json:"PropagatingVgws"`
	RouteTableID    string  `json:"RouteTableId"`
	Routes          []struct {
		DestinationCidrBlock        string  `json:"DestinationCidrBlock"`
		DestinationIpv6CidrBlock    *string `json:"DestinationIpv6CidrBlock"`
		DestinationPrefixListID     *string `json:"DestinationPrefixListId"`
		EgressOnlyInternetGatewayID *string `json:"EgressOnlyInternetGatewayId"`
		GatewayID                   string  `json:"GatewayId"`
		InstanceID                  *string `json:"InstanceId"`
		InstanceOwnerID             *string `json:"InstanceOwnerId"`
		LocalGatewayID              *string `json:"LocalGatewayId"`
		NatGatewayID                *string `json:"NatGatewayId"`
		NetworkInterfaceID          *string `json:"NetworkInterfaceId"`
		Origin                      string  `json:"Origin"`
		State                       string  `json:"State"`
		TransitGatewayID            *string `json:"TransitGatewayId"`
		VpcPeeringConnectionID      *string `json:"VpcPeeringConnectionId"`
	} `json:"Routes"`
	Tags []struct {
		Key   string `json:"Key"`
		Value string `json:"Value"`
	} `json:"Tags"`
	VpcID string `json:"VpcId"`
}

type StackOutputList []StackOutput

func (s StackOutputList) Map() map[string]string {
	var result = make(map[string]string)
	for _, output := range s {
		pk, pv := output.OutputKey, output.OutputValue
		if pk != nil && pv != nil {
			result[*pk] = *pv
		}
	}
	return result
}

type Secrets struct {
	Name         string `json:"name"`
	SecretString string `json:"secretString"`
}
