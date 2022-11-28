// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"strings"
)

//go:generate mockery --inpackage --name=GoquRepositoryApi --structname=MockGoquRepositoryApi
type GoquRepositoryApi interface {
	Get(ctx context.Context, table string) (ds *goqu.SelectDataset, err error)
	Select(ctx context.Context, table string) (ds *goqu.SelectDataset, err error)
	Insert(ctx context.Context, table string) (ds *goqu.InsertDataset, err error)
	Update(ctx context.Context, table string) (ds *goqu.UpdateDataset, err error)
	Upsert(ctx context.Context, table string) (ds *goqu.InsertDataset, err error)
	Delete(ctx context.Context, table string) (ds *goqu.DeleteDataset, err error)
	Truncate(ctx context.Context, table string) (ds *goqu.TruncateDataset, err error)

	ExecuteGet(ctx context.Context, ds *goqu.SelectDataset, dest interface{}) error
	ExecuteSelect(ctx context.Context, ds *goqu.SelectDataset, dest interface{}) error
	ExecuteInsert(ctx context.Context, ds *goqu.InsertDataset) error
	ExecuteUpdate(ctx context.Context, ds *goqu.UpdateDataset) error
	ExecuteUpsert(ctx context.Context, ds *goqu.InsertDataset) error
	ExecuteDelete(ctx context.Context, ds *goqu.DeleteDataset) error
	ExecuteTruncate(ctx context.Context, ds *goqu.TruncateDataset) error
}

type GoquRepository struct {
	sql SqlRepositoryApi
}

func NewGoquRepository(ctx context.Context) GoquRepositoryApi {
	api := ContextGoquRepository().Get(ctx)
	if api == nil {
		api = &GoquRepository{
			sql: NewSqlRepository(ctx),
		}
	}

	return api
}

func (c *GoquRepository) dialect(conn SqlExecutor) goqu.DialectWrapper {
	return goqu.Dialect(conn.DriverName())
}

func (c *GoquRepository) Rebind(conn SqlExecutor, stmt string) string {
	driver := conn.DriverName()
	baseDriver := baseDriverName(driver)
	bindType := sqlx.BindType(baseDriver)
	return sqlx.Rebind(bindType, stmt)
}

func (c *GoquRepository) Get(ctx context.Context, table string) (ds *goqu.SelectDataset, err error) {
	return c.Select(ctx, table)
}

func (c *GoquRepository) Select(ctx context.Context, table string) (ds *goqu.SelectDataset, err error) {
	err = WithSqlExecutor(ctx, func(ctx context.Context, conn SqlExecutor) error {
		ds = c.dialect(conn).From(table).Prepared(true)
		return nil
	})
	return
}

func (c *GoquRepository) Insert(ctx context.Context, table string) (ds *goqu.InsertDataset, err error) {
	err = WithSqlExecutor(ctx, func(ctx context.Context, conn SqlExecutor) error {
		ds = c.dialect(conn).Insert(table).Prepared(true)
		return nil
	})
	return
}

func (c *GoquRepository) Update(ctx context.Context, table string) (ds *goqu.UpdateDataset, err error) {
	err = WithSqlExecutor(ctx, func(ctx context.Context, conn SqlExecutor) error {
		ds = c.dialect(conn).Update(table).Prepared(true)
		return nil
	})
	return
}

func (c *GoquRepository) Upsert(ctx context.Context, table string) (ds *goqu.InsertDataset, err error) {
	return c.Insert(ctx, table)
}

func (c *GoquRepository) Delete(ctx context.Context, table string) (ds *goqu.DeleteDataset, err error) {
	err = WithSqlExecutor(ctx, func(ctx context.Context, conn SqlExecutor) error {
		ds = c.dialect(conn).Delete(table).Prepared(true)
		return nil
	})
	return
}

func (c *GoquRepository) Truncate(ctx context.Context, table string) (ds *goqu.TruncateDataset, err error) {
	err = WithSqlExecutor(ctx, func(ctx context.Context, conn SqlExecutor) error {
		ds = c.dialect(conn).Truncate(table).Prepared(true)
		return nil
	})
	return
}

func (c *GoquRepository) ExecuteGet(ctx context.Context, ds *goqu.SelectDataset, dest interface{}) error {
	stmt, args, err := ds.ToSQL()
	if err != nil {
		return err
	}
	return c.sql.SqlGet(ctx, stmt, args, dest)
}

func (c *GoquRepository) ExecuteSelect(ctx context.Context, ds *goqu.SelectDataset, dest interface{}) error {
	stmt, args, err := ds.ToSQL()
	if err != nil {
		return err
	}
	return c.sql.SqlSelect(ctx, stmt, args, dest)
}

func (c *GoquRepository) ExecuteInsert(ctx context.Context, ds *goqu.InsertDataset) error {
	stmt, args, err := ds.ToSQL()
	if err != nil {
		return err
	}
	return c.sql.SqlExecute(ctx, stmt, args)
}

func (c *GoquRepository) ExecuteUpdate(ctx context.Context, ds *goqu.UpdateDataset) error {
	stmt, args, err := ds.ToSQL()
	if err != nil {
		return err
	}
	return c.sql.SqlExecute(ctx, stmt, args)
}

func (c *GoquRepository) ExecuteUpsert(ctx context.Context, ds *goqu.InsertDataset) error {
	stmt, args, err := ds.ToSQL()
	if err != nil {
		return err
	}
	stmt = "UPSERT" + strings.TrimPrefix(stmt, "INSERT")
	return c.sql.SqlExecute(ctx, stmt, args)
}

func (c *GoquRepository) ExecuteDelete(ctx context.Context, ds *goqu.DeleteDataset) error {
	stmt, args, err := ds.ToSQL()
	if err != nil {
		return err
	}
	return c.sql.SqlExecute(ctx, stmt, args)
}

func (c *GoquRepository) ExecuteTruncate(ctx context.Context, ds *goqu.TruncateDataset) error {
	stmt, args, err := ds.ToSQL()
	if err != nil {
		return err
	}
	return c.sql.SqlExecute(ctx, stmt, args)
}
