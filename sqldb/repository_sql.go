// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"context"
	"github.com/jmoiron/sqlx"
)

//go:generate mockery --inpackage --name=SqlRepositoryApi --structname=MockSqlRepositoryApi
type SqlRepositoryApi interface {
	// SqlGet retrieves one row using raw parameterized sql
	SqlGet(ctx context.Context, stmt string, args []interface{}, dest interface{}) error
	// SqlSelect retrieves multiple rows using raw parameterized sql
	SqlSelect(ctx context.Context, stmt string, args []interface{}, dest interface{}) error
	// SqlExecute executes a DML statement using raw parameterized sql
	SqlExecute(ctx context.Context, stmt string, args []interface{}) error
}

type SqlRepository struct {
}

func NewSqlRepository(ctx context.Context) SqlRepositoryApi {
	api := ContextSqlRepository().Get(ctx)
	if api == nil {
		api = &SqlRepository{}
	}

	return api
}

func (c *SqlRepository) Rebind(conn SqlExecutor, stmt string) string {
	driver := conn.DriverName()
	baseDriver := baseDriverName(driver)
	bindType := sqlx.BindType(baseDriver)
	return sqlx.Rebind(bindType, stmt)
}

func (c *SqlRepository) SqlGet(ctx context.Context, stmt string, args []interface{}, dest interface{}) error {
	return WithSqlExecutor(ctx, func(ctx context.Context, conn SqlExecutor) error {
		stmt = c.Rebind(conn, stmt)
		statements.Printf(queryLogFormat, stmt, args)
		return conn.GetContext(ctx, dest, stmt, args...)
	})
}

func (c *SqlRepository) SqlSelect(ctx context.Context, stmt string, args []interface{}, dest interface{}) error {
	return WithSqlExecutor(ctx, func(ctx context.Context, conn SqlExecutor) error {
		stmt = c.Rebind(conn, stmt)
		statements.Printf(queryLogFormat, stmt, args)
		return conn.SelectContext(ctx, dest, stmt, args...)
	})
}

func (c *SqlRepository) SqlExecute(ctx context.Context, stmt string, args []interface{}) error {
	return WithSqlExecutor(ctx, func(ctx context.Context, conn SqlExecutor) error {
		stmt = c.Rebind(conn, stmt)
		statements.Printf(queryLogFormat, stmt, args)
		_, err := conn.ExecContext(ctx, stmt, args...)
		return err
	})
}
