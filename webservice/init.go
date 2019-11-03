package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	restfulLog "github.com/emicklei/go-restful/log"
)

var restfulLogger = log.NewLogger("restful")

func init() {
	// Reconfigure the restful logging
	restfulLog.SetLogger(restfulLogger)
}
