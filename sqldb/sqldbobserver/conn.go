package sqldbobserver

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"database/sql/driver"
	"strings"
)

// conn defines a tracing wrapper for driver.conn.
type conn struct {
	conn driver.Conn
}

// Prepare implements driver.Conn Prepare.
func (c *conn) Prepare(query string) (driver.Stmt, error) {
	s, err := c.conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &stmt{stmt: s}, nil
}

// Close implements driver.Conn Close.
func (c *conn) Close() error {
	return c.conn.Close()
}

// Prepare implements driver.conn Begin.
func (c *conn) Begin() (driver.Tx, error) {
	t, err := c.conn.Begin()
	if err != nil {
		return nil, err
	}
	return &tx{tx: t}, nil
}

// BeginTx implements driver.Conn BeginTx.
func (c *conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	ctx, s := trace.NewSpan(ctx, "sql.DB.BeginTx")
	if connBeginTx, ok := c.conn.(driver.ConnBeginTx); ok {
		t, err := connBeginTx.BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		return &tx{tx: t, span: s}, nil
	}
	return c.conn.Begin()
}

// PrepareContext implements driver.Conn PrepareContext.
func (c *conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	if connPrepareContext, ok := c.conn.(driver.ConnPrepareContext); ok {
		s, err := connPrepareContext.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}
		return &stmt{stmt: s}, nil
	}
	return c.conn.Prepare(query)
}

// Exec implements driver.Execer Exec.
func (c *conn) Exec(query string, args []driver.Value) (result driver.Result, err error) {
	if execer, ok := c.conn.(driver.Execer); ok {
		err = observeStats(query, func() error {
			result, err = execer.Exec(query, args)
			return err
		})
		return
	}
	return nil, ErrUnsupported
}

// Exec implements driver.StmtExecContext ExecContext.
func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (result driver.Result, err error) {
	ctx, s := trace.NewSpan(ctx, "sql.DB.ExecContext")
	s.SetTag(TagQuery, strings.Fields(query)[0])
	defer s.Finish()

	if execerContext, ok := c.conn.(driver.ExecerContext); ok {
		err = observeStats(query, func() error {
			result, err = execerContext.ExecContext(ctx, query, args)
			return err
		})
		return
	}

	values, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}
	return c.Exec(query, values)
}

// Ping implements driver.Pinger Ping.
func (c *conn) Ping(ctx context.Context) error {
	if pinger, ok := c.conn.(driver.Pinger); ok {
		ctx, s := trace.NewSpan(ctx, "sql.DB.Ping")
		defer s.Finish()
		return observeStats("PING", func() error {
			return pinger.Ping(ctx)
		})
	}
	return ErrUnsupported
}

// Query implements driver.Queryer Query.
func (c *conn) Query(query string, args []driver.Value) (rows driver.Rows, err error) {
	if queryer, ok := c.conn.(driver.Queryer); ok {
		err = observeStats(query, func() error {
			rows, err = queryer.Query(query, args)
			return err
		})
		return
	}
	return nil, ErrUnsupported
}

// QueryContext implements driver.QueryerContext QueryContext.
func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	ctx, s := trace.NewSpan(ctx, "sql.DB.QueryContext")
	s.SetTag(TagQuery, strings.Fields(query)[0])
	defer s.Finish()

	if queryerContext, ok := c.conn.(driver.QueryerContext); ok {
		err = observeStats(query, func() error {
			rows, err = queryerContext.QueryContext(ctx, query, args)
			return err
		})
		return
	}

	values, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}
	return c.Query(query, values)
}
