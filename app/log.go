package app

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/sanitize"
)

var logger = log.NewLogger("msx.app")

func init() {
	OnEvent(EventConfigure, PhaseAfter, sanitize.ConfigureSecretSanitizer)
}
