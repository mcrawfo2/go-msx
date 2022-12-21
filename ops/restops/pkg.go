// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/pkg/errors"
)

var logger = log.NewPackageLogger()

const (
	PortTypeRequest  = webservice.StructTagRequest
	PortTypeResponse = webservice.StructTagResponse

	PathApiRoot = "/api"
)

var ErrNotImplemented = errors.New("Not implemented")
