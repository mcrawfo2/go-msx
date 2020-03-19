package app

import "cto-github.cisco.com/NFV-BU/go-msx/fs"

func init() {
	OnEvent(EventConfigure, PhaseAfter, withConfig(fs.ConfigureFileSystem))
}
