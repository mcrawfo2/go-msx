// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package mocks

import (
	"cto-github.cisco.com/NFV-BU/go-msx/sqldb"
	"testing"
)

func TestImplementations(t *testing.T) {
	var _ sqldb.CrudRepositoryApi = &CrudRepositoryApi{}
	var _ sqldb.CrudRepositoryFactoryApi = &CrudRepositoryFactory{}
}
