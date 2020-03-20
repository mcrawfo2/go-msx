package app

import "cto-github.cisco.com/NFV-BU/go-msx/transit/vaultprovider"

func init() {
	OnEvent(EventConfigure, PhaseAfter, vaultprovider.RegisterVaultTransitProvider)
}
