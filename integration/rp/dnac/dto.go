package dnac

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

type ControlPlaneAuthenticationType string

const (
	Local ControlPlaneAuthenticationType = "local"
	Sso   ControlPlaneAuthenticationType = "sso"
)

type DnacConnectRequest struct {
	Url         string `json:"url"`
	TlsInsecure bool   `json:"tlsInsecure"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type DnacExtendedRequest struct {
	ControlPlaneId     types.UUID                     `json:"controlPlaneId"`
	Url                string                         `json:"url"`
	TlsInsecure        bool                           `json:"tlsInsecure"`
	AuthenticationType ControlPlaneAuthenticationType `json:"controlPlaneAuthenticationEnum"`
}

type Pojo integration.Pojo
