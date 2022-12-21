// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import "cto-github.cisco.com/NFV-BU/go-msx/ops"

// Predicates

func PortFieldIsBody(pf *ops.PortField) bool {
	return pf.Group == FieldGroupHttpBody
}

func PortFieldIsPaging(pf *ops.PortField) bool {
	return pf.Group == FieldGroupHttpPaging
}

func PortFieldIsSuccessBody(pf *ops.PortField) bool {
	isError, _ := pf.BoolOption("error")
	return pf.Group == FieldGroupHttpBody && !isError
}

func PortFieldIsSuccessHeader(pf *ops.PortField) bool {
	isError, _ := pf.BoolOption("error")
	return pf.Group == FieldGroupHttpHeader && !isError
}

func PortFieldIsErrorHeader(pf *ops.PortField) bool {
	isError, _ := pf.BoolOption("error")
	return pf.Group == FieldGroupHttpHeader && isError
}

func PortFieldIsHeader(pf *ops.PortField) bool {
	return pf.Group == FieldGroupHttpHeader
}

func PortFieldIsCode(pf *ops.PortField) bool {
	return pf.Group == FieldGroupHttpCode
}

func PortFieldIsForm(pf *ops.PortField) bool {
	return pf.Group == FieldGroupHttpForm
}

func PortFieldIsError(pf *ops.PortField) bool {
	e, _ := pf.BoolOption("error")
	return e
}
