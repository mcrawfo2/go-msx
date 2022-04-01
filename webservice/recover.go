// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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

			bt := types.BackTraceFromDebugStackTrace(debug.Stack())
			logger.WithContext(req.Request.Context()).WithError(e).WithField(log.FieldStack, bt.Stanza()).Error("Recovered from panic")
			log.ErrorMessage(logger, req.Request.Context(), e)
			log.Stack(logger, req.Request.Context(), bt)

			WriteError(req, resp, 500, e)
		}
	}()

	chain.ProcessFilter(req, resp)
}
