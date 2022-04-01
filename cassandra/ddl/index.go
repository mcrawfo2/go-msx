// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ddl

import (
	"strings"
)

type Index struct {
	Name   string
	Table  string
	Column string
}

type CreateIndexQueryBuilder struct{}

func (b *CreateIndexQueryBuilder) CreateIndex(index Index, ifNotExists bool) string {
	sb := new(strings.Builder)
	sb.WriteString("CREATE INDEX ")
	if ifNotExists {
		sb.WriteString("IF NOT EXISTS ")
	}
	sb.WriteString(index.Name)
	sb.WriteString(" ON ")
	sb.WriteString(index.Table)
	sb.WriteString(" (")
	sb.WriteString(index.Column)
	sb.WriteString(")")
	return sb.String()
}
