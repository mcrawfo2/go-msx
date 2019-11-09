package ddl

import (
	"github.com/scylladb/gocqlx/table"
	"strings"
)

type Column struct {
	Name     string
	DataType string
}

type Table struct {
	Name          string
	Columns       []Column
	PartitionKeys []string
	ClusterKeys   []string
}

func (t Table) Metadata() table.Metadata {
	return table.Metadata{
		Name:    t.Name,
		Columns: t.ColumnNames(),
		PartKey: t.PartitionKeys,
		SortKey: t.ClusterKeys,
	}
}

func (t Table) ColumnNames() []string {
	var result []string
	for _, column := range t.Columns {
		result = append(result, column.Name)
	}
	return result
}

type CreateTableQueryBuilder struct{}

func (b *CreateTableQueryBuilder) CreateTable(table Table, ifNotExists bool) string {
	sb := new(strings.Builder)
	sb.WriteString("CREATE TABLE ")
	if ifNotExists {
		sb.WriteString("IF NOT EXISTS ")
	}
	sb.WriteString(table.Name)
	sb.WriteString(" (")

	n := 0
	for _, column := range table.Columns {
		if n > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(column.Name)
		sb.WriteRune(' ')
		sb.WriteString(column.DataType)
		n++
	}

	// Keys
	if n > 0 {
		sb.WriteString(", ")
	}
	sb.WriteString("PRIMARY KEY ((")
	sb.WriteString(strings.Join(table.PartitionKeys, ","))
	sb.WriteString(")")
	if len(table.ClusterKeys) > 0 {
		sb.WriteString(", ")
		sb.WriteString(strings.Join(table.ClusterKeys, ","))
	}
	sb.WriteString("))")

	return sb.String()
}
