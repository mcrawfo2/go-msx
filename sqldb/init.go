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
