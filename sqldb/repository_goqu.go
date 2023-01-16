// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"
	"strings"
)

//go:generate mockery --name=GoquRepositoryApi --inpackage --case=snake --with-expecter
type GoquRepositoryApi interface {
	Get(table string) *goqu.SelectDataset
	Select(table string) *goqu.SelectDataset
	Insert(table string) *goqu.InsertDataset
	Update(table string) *goqu.UpdateDataset
	Upsert(table string) *goqu.InsertDataset
	Delete(table string) *goqu.DeleteDataset
	Truncate(table string) *goqu.TruncateDataset

	ExecuteGet(ctx context.Context, ds *goqu.SelectDataset, dest interface{}) error
	ExecuteSelect(ctx context.Context, ds *goqu.SelectDataset, dest interface{}) error
	ExecuteInsert(ctx context.Context, ds *goqu.InsertDataset) error
	ExecuteUpdate(ctx context.Context, ds *goqu.UpdateDataset) error
	ExecuteUpsert(ctx context.Context, ds *goqu.InsertDataset) error
	ExecuteDelete(ctx context.Context, ds *goqu.DeleteDataset) error
	ExecuteTruncate(ctx context.Context, ds *goqu.TruncateDataset) error
}

type GoquRepository struct {
	sql        SqlRepositoryApi
	driverName string
}

func NewGoquRepository(ctx context.Context) (GoquRepositoryApi, error) {
	api := ContextGoquRepository().Get(ctx)
	if api == nil {
		driverName, err := SqlDriverName(ctx)
		if err != nil {
			return nil, err
		}

		sqlRepository, err := NewSqlRepository(ctx)
		if err != nil {
			return nil, err
		}

		api = &GoquRepository{
			sql:        sqlRepository,
			driverName: BaseDriverName(driverName),
		}
	}

	return api, nil
}

func (c *GoquRepository) Get(table string) *goqu.SelectDataset {
	return c.Select(table)
}

func (c *GoquRepository) Select(table string) *goqu.SelectDataset {
	return goqu.Dialect(c.driverName).From(table).Prepared(true)
}

func (c *GoquRepository) Insert(table string) *goqu.InsertDataset {
	return goqu.Dialect(c.driverName).Insert(table).Prepared(true)
}

func (c *GoquRepository) Update(table string) *goqu.UpdateDataset {
	return goqu.Dialect(c.driverName).Update(table).Prepared(true)
}

func (c *GoquRepository) Upsert(table string) *goqu.InsertDataset {
	return c.Insert(table)
}

func (c *GoquRepository) Delete(table string) *goqu.DeleteDataset {
	return goqu.Dialect(c.driverName).Delete(table).Prepared(true)
}

func (c *GoquRepository) Truncate(table string) *goqu.TruncateDataset {
	return goqu.Dialect(c.driverName).Truncate(table).Prepared(true)
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

	switch c.driverName {
	case "postgres":
		stmt = "UPSERT" + strings.TrimPrefix(stmt, "INSERT")
	case "sqlite", "sqlite3":
		stmt = "REPLACE" + strings.TrimPrefix(stmt, "INSERT")
	default:
		return errors.Wrap(ErrNotImplemented, "Upsert not supported")
	}

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

	switch c.driverName {
	case "sqlite", "sqlite3":
		stmt = "DELETE FROM" + strings.TrimPrefix(stmt, "TRUNCATE")
	}

	return c.sql.SqlExecute(ctx, stmt, args)
}
