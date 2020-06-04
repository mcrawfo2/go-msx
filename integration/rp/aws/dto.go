package aws

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
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

type Pojo integration.Pojo
