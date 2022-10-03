// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/sirupsen/logrus"
)

const queryLogFormat = `Query: %v %v`

var logger = log.NewLogger("msx.sql")
var statementLogger = log.NewLogger("msx.sql.query")
var statements = log.NewLevelLogger(statementLogger, logrus.DebugLevel)
