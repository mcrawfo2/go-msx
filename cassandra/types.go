// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package cassandra

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/gocql/gocql"
)

func ToModelUuid(value types.UUID) gocql.UUID {
	return value.ToByteArray()
}

func ToApiUuid(value gocql.UUID) types.UUID {
	return value[:]
}
