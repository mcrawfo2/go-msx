// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"mime/multipart"
)

//go:generate mockery --name=InputDecoder --inpackage ==case=snake --testonly

type InputDecoder interface {
	DecodePrimitive(pf *PortField) (result types.Optional[string], err error)
	DecodeContent(pf *PortField) (content Content, err error)
	DecodeArray(pf *PortField) (result []string, err error)
	DecodeObject(pf *PortField) (result types.Pojo, err error)
	DecodeFile(pf *PortField) (result *multipart.FileHeader, err error)
	DecodeFileArray(pf *PortField) (result []*multipart.FileHeader, err error)
	DecodeAny(pf *PortField) (result types.Optional[any], err error)
}
