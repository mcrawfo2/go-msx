package lucene

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Search(t *testing.T) {
	var s = new(SearchBuilder).
		WithRefresh(true).
		WithQuery(Must(Match("createdby", "system"))).
		WithSort(Field("createdby")).
		Build()
	assert.Equal(t, `{`+
		`"query":[{"type":"boolean","must":[{"type":"match","field":"createdby","value":"system"}]}],`+
		`"sort":[{"type":"simple","field":"createdby"}],`+
		`"refresh":true`+
		`}`, s)
}

func Test_Search_NoRefresh(t *testing.T) {
	var s = new(SearchBuilder).
		WithQuery(Must(Match("createdby", "system"))).
		WithSort(Field("createdby")).
		Build()
	assert.Equal(t, `{`+
		`"query":[{"type":"boolean","must":[{"type":"match","field":"createdby","value":"system"}]}],`+
		`"sort":[{"type":"simple","field":"createdby"}]`+
		`}`, s)
}
