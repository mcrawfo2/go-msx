// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/schema/js"
)

type PortFieldValidationSchemaFunc func(field *PortField) (schema js.ValidationSchema, err error)
