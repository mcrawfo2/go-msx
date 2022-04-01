// Copyright (c) 2018 OpenTracing-SQL Authors
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldbobserver

import (
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"database/sql/driver"
)

// conn defines a tracing wrapper for driver.Tx.
type tx struct {
	tx     driver.Tx
	tracer *tracer
	span   trace.Span
}

// Commit implements driver.Tx Commit.
func (t *tx) Commit() error {
	if t.span != nil {
		defer t.span.Finish()
	}
	return t.tx.Commit()
}

// Rollback implements driver.Tx Rollback.
func (t *tx) Rollback() error {
	if t.span != nil {
		defer t.span.Finish()
	}
	return t.tx.Rollback()
}
