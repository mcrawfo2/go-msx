// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"database/sql/driver"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"sync"
)

const queryLogFormat = `Query: %v %v`

var logger = log.NewLogger("msx.sql")
var statementLogger = log.NewLogger("msx.sql.query")
var statements = log.NewLevelLogger(statementLogger, logrus.DebugLevel)

var drivers = make(map[string]driver.Driver)
var driverMtx sync.Mutex

func init() {
	drivers["postgres"] = &pq.Driver{}
}
