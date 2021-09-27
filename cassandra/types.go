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
