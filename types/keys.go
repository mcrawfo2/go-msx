package types

import (
	"fmt"
	"sort"
	"strings"
)

type Keys map[string]interface{}

func (k Keys) String() string {
	type pair struct {
		key   string
		value interface{}
	}

	var pairs []pair
	for key, value := range k {
		pairs = append(pairs, pair{key, value})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return strings.Compare(pairs[i].key, pairs[j].key) < 0
	})

	var sb = strings.Builder{}
	for _, pair := range pairs {
		if sb.Len() > 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(fmt.Sprintf("%s=%v", pair.key, pair.value))
	}
	return sb.String()
}
