package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"runtime/debug"
)

func recoveryFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	defer func() {
		if r := recover(); r != nil {
			var e error
			if err, ok := r.(error); ok {
				e = err
			} else {
				e = errors.Errorf("Exception: %v", r)
			}

			logger.WithContext(req.Request.Context()).WithError(e).Error("Recovered from panic")
			log.ErrorMessage(logger, req.Request.Context(), e)
			bt := types.BackTraceFromDebugStackTrace(debug.Stack())
			log.Stack(logger, req.Request.Context(), bt)

			WriteError(req, resp, 500, e)
		}
	}()

	chain.ProcessFilter(req, resp)
}
