package lucene

import (
	"fmt"
	"strconv"
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

type IndexOptions struct {
	RefreshSeconds      *int
	DirectoryPath       *string
	RamBufferMb         *string
	MaxMergeMb          *string
	MaxCachedMb         *string
	IndexingThreads     *int
	IndexingQueuesSize  *string
	ExcludedDataCenters *string
	Partitioner         *string
	Sparse              *bool
	Schema              *IndexSchema
}

func (o IndexOptions) optionStrings() []string {
	var result []string
	result = append(result, o.Option(IndexOptionRefreshSeconds, o.RefreshSeconds))
	result = append(result, o.Option(IndexOptionDirectoryPath, o.DirectoryPath))
	result = append(result, o.Option(IndexOptionRamBufferMb, o.RamBufferMb))
	result = append(result, o.Option(IndexOptionMaxMergeMb, o.MaxMergeMb))
	result = append(result, o.Option(IndexOptionMaxCachedMb, o.MaxCachedMb))
	result = append(result, o.Option(IndexOptionIndexingThreads, o.IndexingThreads))
	result = append(result, o.Option(IndexOptionIndexingQueuesSize, o.IndexingQueuesSize))
	result = append(result, o.Option(IndexOptionExcludedDataCenters, o.ExcludedDataCenters))
	result = append(result, o.Option(IndexOptionPartitioner, o.Partitioner))
	result = append(result, o.Option(IndexOptionSparse, o.Sparse))
	result = append(result, o.Option(IndexOptionSchema, o.Schema))
	return result
}

func (o IndexOptions) String() string {
	sb := new(strings.Builder)
	sb.WriteString("{")
	out := false
	for _, v := range o.optionStrings() {
		if v == "" {
			continue
		}
		if out {
			sb.WriteRune(',')
		}
		sb.WriteString(v)
		out = true
	}
	sb.WriteRune('}')

	return sb.String()
}

func (o IndexOptions) Option(key string, value interface{}) string {
	if value == nil {
		return ""
	}

	valueString := ""
	if valueStringer, ok := value.(fmt.Stringer); ok {
		valueString = valueStringer.String()
	} else if valueInt, ok := value.(*int); ok {
		valueString = strconv.Itoa(*valueInt)
	} else if valueBool, ok := value.(*bool); ok {
		if *valueBool {
			valueString = "true"
		} else {
			valueString = "false"
		}
	} else {
		valueString = fmt.Sprintf("%v", value)
	}

	return fmt.Sprintf(`'%s':'%s'`, key, valueString)
}

type Index struct {
	Name    string
	Table   string
	Column  *string
	Options IndexOptions
}

type IndexQueryBuilder struct{}

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
	sb.WriteString(i.Options.String())
	return sb.String()
}

func (b *IndexQueryBuilder) DropIndex(i Index, ifExists bool) string {
	sb := new(strings.Builder)
	sb.WriteString("DROP INDEX ")
	if ifExists {
		sb.WriteString("IF EXISTS ")
	}
	sb.WriteString(i.Name)
	return sb.String()
}
