package main

import (
	"cto-github.cisco.com/NFV-BU/go-msx/app"
)

const (
	appName = "${app.name}"
)

func main() {
	app.Run(appName)
}
