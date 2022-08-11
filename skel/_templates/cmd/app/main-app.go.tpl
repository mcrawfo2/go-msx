package main

import (
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	_ "cto-github.cisco.com/NFV-BU/go-msx-populator/populate"
)

const (
	appName = "${app.name}"
)

func main() {
	app.Run(appName)
}
