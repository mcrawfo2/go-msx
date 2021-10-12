package app

import (
	"cto-github.cisco.com/NFV-BU/go-msx/transit"
	"cto-github.cisco.com/NFV-BU/go-msx/transit/vaultprovider"
)

func init() {
	OnEvent(EventConfigure, PhaseAfter, transit.ConfigureEncrypterFactory)
	OnEvent(EventConfigure, PhaseAfter, vaultprovider.RegisterVaultTransitProvider)
}
