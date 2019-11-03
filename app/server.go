package app

import (
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
)

func init() {
	OnEvent(EventStart, PhaseAfter, webservice.Start)
	OnEvent(EventStop, PhaseBefore, webservice.Stop)
}
