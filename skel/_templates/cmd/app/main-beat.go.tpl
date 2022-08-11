package main

import (
	"cto-github.cisco.com/NFV-BU/go-msx-beats/beat"
	_ "cto-github.cisco.com/NFV-BU/go-msx-populator/populate"
)

const (
	appName = "${app.name}"
)

func main() {
	beat.Run(appName)
}
