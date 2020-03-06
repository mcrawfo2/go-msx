package lucene

import (
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra/ddl"
	"strings"
)

const (
	IndexOptionRefreshSeconds      = "refresh_seconds"
	IndexOptionDirectoryPath       = "directory_path"
	IndexOptionRamBufferMb         = "ram_buffer_mb"
	IndexOptionMaxMergeMb          = "max_merge_mb"
	IndexOptionMaxCachedMb         = "max_cached_mb"
	IndexOptionIndexingThreads     = "indexing_threads"
	IndexOptionIndexingQueuesSize  = "indexing_queues_size"
	IndexOptionExcludedDataCenters = "excluded_data_centers"
	IndexOptionPartitioner         = "partitioner"
	IndexOptionSparse              = "sparse"
	IndexOptionSchema              = "schema"
)

type Index struct {
	Name    string
	Table   string
	Column  *string
	Options map[string]string
}

type IndexQueryBuilder struct {
	o ddl.OptionsQueryPartBuilder
}

func (b *IndexQueryBuilder) CreateIndex(i Index) string {
	sb := new(strings.Builder)
	sb.WriteString("CREATE CUSTOM INDEX ")
	sb.WriteString(i.Name)
	sb.WriteString(" ON ")
	sb.WriteString(i.Table)
	sb.WriteRune('(')
	if i.Column != nil {
		sb.WriteString(*i.Column)
	}
	sb.WriteRune(')')
	sb.WriteString("USING 'com.stratio.cassandra.lucene.index' WITH OPTIONS = ")
	sb.WriteString(b.o.Options(i.Options))
	return sb.String()
}
