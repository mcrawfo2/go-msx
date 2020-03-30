package webservice

import (
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
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

			WriteError(req, resp, 500, e)
		}
	}()

	chain.ProcessFilter(req, resp)
}
