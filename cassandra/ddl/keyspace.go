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

type KeyspaceQueryBuilder struct {
	options OptionsQueryPartBuilder
}

func (b *KeyspaceQueryBuilder) CreateKeyspace(keyspace Keyspace, ifNotExists bool) string {
	sb := new(strings.Builder)
	sb.WriteString("CREATE KEYSPACE ")
	if ifNotExists {
		sb.WriteString("IF NOT EXISTS ")
	}
	sb.WriteString(keyspace.Name)

	sb.WriteString(" WITH replication = ")
	sb.WriteString(b.options.Options(keyspace.ReplicationOptions))
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
