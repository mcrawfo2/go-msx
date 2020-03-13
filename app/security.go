package app

import (
	"cto-github.cisco.com/NFV-BU/go-msx/security/idmdetailsprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/security/jwttokenprovider"
)

func init() {
	OnEvent(EventConfigure, PhaseAfter, jwttokenprovider.RegisterTokenProvider)
	OnEvent(EventConfigure, PhaseAfter, idmdetailsprovider.RegisterTokenDetailsProvider)
}
