package sqldbobserver

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"database/sql/driver"
	"strings"
)

// conn defines a tracing wrapper for driver.Stmt.
type stmt struct {
	stmt driver.Stmt
}

// Close implements driver.Stmt Close.
func (s *stmt) Close() error {
	return s.stmt.Close()
}

// NumInput implements driver.Stmt NumInput.
func (s *stmt) NumInput() int {
	return s.stmt.NumInput()
}

// Exec implements driver.Stmt Exec.
func (s *stmt) Exec(args []driver.Value) (result driver.Result, err error) {
	return s.stmt.Exec(args)
}

// Query implements driver.Stmt Query.
func (s *stmt) Query(args []driver.Value) (driver.Rows, error) {
	return s.stmt.Query(args)
}

// ExecContext implements driver.ExecerContext ExecContext.
func (s *stmt) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	ctx, span := trace.NewSpan(ctx, "sql.Stmt.ExecContext")
	span.SetTag(TagQuery, strings.Fields(query)[0])
	defer span.Finish()

	if execerContext, ok := s.stmt.(driver.ExecerContext); ok {
		return execerContext.ExecContext(ctx, query, args)
	}
	values, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}
	return s.Exec(values)
}

// QueryContext implements driver.QueryerContext QueryContext.
func (s *stmt) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	ctx, span := trace.NewSpan(ctx, "sql.Stmt.QueryContext")
	span.SetTag(TagQuery, strings.Fields(query)[0])
	defer span.Finish()

	if queryerContext, ok := s.stmt.(driver.QueryerContext); ok {
		rows, err := queryerContext.QueryContext(ctx, query, args)
		return rows, err
	}
	values, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}
	return s.Query(values)
}
