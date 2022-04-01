// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"database/sql/driver"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/lib/pq"
)

var drivers = make(map[string]driver.Driver)

func init() {
	drivers["postgres"] = &pq.Driver{}
}
