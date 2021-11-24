package main

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"strconv"
)

var BuildNumber = "0"

func main() {
	buildNumber, _ := strconv.ParseInt(BuildNumber, 10, 64)
	skel.Run(int(buildNumber))
}
