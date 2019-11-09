package ddl

import "strings"

const (
	LuceneIndexOptionRefreshSeconds      = "refresh_seconds"
	LuceneIndexOptionDirectoryPath       = "directory_path"
	LuceneIndexOptionRamBufferMb         = "ram_buffer_mb"
	LuceneIndexOptionMaxMergeMb          = "max_merge_mb"
	LuceneIndexOptionMaxCachedMb         = "max_cached_mb"
	LuceneIndexOptionIndexingThreads     = "indexing_threads"
	LuceneIndexOptionIndexingQueuesSize  = "indexing_queues_size"
	LuceneIndexOptionExcludedDataCenters = "excluded_data_centers"
	LuceneIndexOptionPartitioner         = "partitioner"
	LuceneIndexOptionSparse              = "sparse"
	LuceneIndexOptionSchema              = "schema"
)

type LuceneIndex struct {
	Name    string
	Table   string
	Column  *string
	Options map[string]string
}

type LuceneIndexQueryBuilder struct {
	o OptionsQueryPartBuilder
}

func (b *LuceneIndexQueryBuilder) CreateIndex(i LuceneIndex) string {
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
