package main

import (
	"cto-github.cisco.com/NFV-BU/go-msx/app"
//#if EXTERNAL
//#else EXTERNAL
	_ "cto-github.cisco.com/NFV-BU/go-msx-populator/populate"
//#endif EXTERNAL
)

const (
	appName = "${app.name}"
)

func main() {
	app.Run(appName)
}
