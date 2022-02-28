package app

import (
	"cto-github.cisco.com/NFV-BU/go-msx/certificate/cacheprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/certificate/fileprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/certificate/vaultprovider"
)

func init() {
	OnEvent(EventConfigure, PhaseAfter, fileprovider.RegisterFactory)
	OnEvent(EventConfigure, PhaseAfter, vaultprovider.RegisterFactory)
	OnEvent(EventConfigure, PhaseAfter, cacheprovider.RegisterFactory)
}
