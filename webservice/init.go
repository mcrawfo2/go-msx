// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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
