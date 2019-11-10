package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/emicklei/go-restful"
)

var restfulLogger = log.NewLogger("restful")

func init() {
	// Reconfigure the restful logging
	restful.TraceLogger(restfulLogger)
	restful.SetLogger(restfulLogger)
}
