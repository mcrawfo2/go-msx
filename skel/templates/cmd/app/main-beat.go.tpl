package main

import (
	"cto-github.cisco.com/NFV-BU/go-msx-beats/beat"
)

const (
	appName = "${app.name}"
)

func main() {
	beat.Run(appName)
}
