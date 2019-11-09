package ddl

import (
	"strings"
)

const (
	ReplicationOptionsKeyClass   = "class"
	ClassSimpleStrategy          = "SimpleStrategy"
	ClassNetworkTopologyStrategy = "NetworkTopologyStrategy"
)

type Keyspace struct {
	Name               string
	ReplicationOptions map[string]string
	DurableWrites      bool
}

type KeyspaceQueryBuilder struct{}

func (b *KeyspaceQueryBuilder) CreateKeyspace(keyspace Keyspace, ifNotExists bool) string {
	sb := new(strings.Builder)
	sb.WriteString("CREATE KEYSPACE ")
	if ifNotExists {
		sb.WriteString("IF NOT EXISTS ")
	}
	sb.WriteString(keyspace.Name)

	sb.WriteString(" WITH replication = {")
	n := 0
	for k, v := range keyspace.ReplicationOptions {
		if n > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("'")
		sb.WriteString(k)
		sb.WriteString("': '")
		sb.WriteString(v)
		sb.WriteString("'")
		n++
	}
	sb.WriteString("}")

	if keyspace.DurableWrites {
		sb.WriteString(" AND durable_writes = true")
	}

	return sb.String()
}

func (b *KeyspaceQueryBuilder) DropKeyspace(keyspace Keyspace, ifExists bool) string {
	sb := new(strings.Builder)
	sb.WriteString("DROP KEYSPACE ")
	if ifExists {
		sb.WriteString("IF EXISTS ")
	}
	sb.WriteString(keyspace.Name)
	return sb.String()
}
