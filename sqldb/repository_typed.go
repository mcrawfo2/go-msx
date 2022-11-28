// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

type WhereOption = goqu.Expression
type KeysOption = goqu.Ex

//go:generate mockery --inpackage --name=TypedRepositoryApi --structname=MockTypedRepositoryApi
type TypedRepositoryApi[I any] interface {
	CountAll(ctx context.Context, dest *int64, where types.Optional[WhereOption]) error
	FindAll(ctx context.Context, dest *[]I, options ...FindAllOption) (pagingResponse paging.Response, err error)
	FindOne(ctx context.Context, dest *I, where types.Optional[WhereOption]) error
	Insert(ctx context.Context, value I) error
	Update(ctx context.Context, where WhereOption, value I) error
	Upsert(ctx context.Context, value I) error
	DeleteOne(ctx context.Context, keys KeysOption) error
	DeleteAll(ctx context.Context, where WhereOption) error
	Truncate(ctx context.Context) error
}

type TypedRepository[I any] struct {
	table string
	goqu  GoquRepositoryApi
}

func NewTypedRepository[I any](ctx context.Context, table string) TypedRepositoryApi[I] {
	api := ContextTypedRepository[I](table).Get(ctx)
	if api == nil {
		api = &TypedRepository[I]{
			table: table,
			goqu:  NewGoquRepository(ctx),
		}
	}
	return api
}

func (c *TypedRepository[I]) dialect(conn SqlExecutor) goqu.DialectWrapper {
	return goqu.Dialect(conn.DriverName())
}

func (c *TypedRepository[I]) Rebind(conn SqlExecutor, stmt string) string {
	driver := conn.DriverName()
	baseDriver := baseDriverName(driver)
	bindType := sqlx.BindType(baseDriver)
	return sqlx.Rebind(bindType, stmt)
}

func (c *TypedRepository[I]) CountAll(ctx context.Context, dest *int64, where types.Optional[WhereOption]) error {
	ds, err := c.goqu.Select(ctx, c.table)
	if err != nil {
		return err
	}

	if where.IsPresent() {
		ds = ds.Where(where.Value())
	}

	return c.goqu.ExecuteGet(ctx, ds.Select(goqu.COUNT("*")), dest)
}

func (c *TypedRepository[I]) FindAll(ctx context.Context, dest *[]I, options ...FindAllOption) (pagingResponse paging.Response, err error) {
	// sub query
	sub, err := c.goqu.Select(ctx, c.table)
	if err != nil {
		return
	}

	pgReq := &paging.Request{}

	for _, option := range options {
		sub, pgReq = option(sub, pgReq)
	}

	sub = sub.ClearLimit().ClearOffset()

	// total items
	count := int64(0)
	countDs, err := c.goqu.Select(ctx, c.table)
	if err != nil {
		return
	}
	countDs = countDs.From(sub).Select(goqu.COUNT("*"))

	err = c.goqu.ExecuteGet(ctx, countDs, &count)
	if err != nil {
		return
	}

	totalItems := uint(count)

	// with paging
	pagedDs, err := c.goqu.Select(ctx, c.table)
	if err != nil {
		return
	}
	pagedDs = pagedDs.From(sub)

	if pgReq.Size > 0 {
		pagedDs = pagedDs.
			Limit(pgReq.Size).
			Offset(pgReq.Page * pgReq.Size)
	}

	err = c.goqu.ExecuteSelect(ctx, pagedDs, dest)
	if err != nil {
		return
	}

	pagingResponse = paging.Response{
		Content:    dest,
		Size:       pgReq.Size,
		Number:     pgReq.Page,
		Sort:       pgReq.Sort,
		TotalItems: &totalItems,
	}

	return
}

func (c *TypedRepository[I]) FindOne(ctx context.Context, dest *I, where types.Optional[WhereOption]) error {
	ds, err := c.goqu.Get(ctx, c.table)
	if err != nil {
		return err
	}

	if where.IsPresent() {
		ds = ds.Where(where.Value())
	}

	return c.goqu.ExecuteGet(ctx, ds, dest)
}

func (c *TypedRepository[I]) Insert(ctx context.Context, value I) error {
	ds, err := c.goqu.Insert(ctx, c.table)
	if err != nil {
		return err
	}
	return c.goqu.ExecuteInsert(ctx, ds.Rows(value))
}

func (c *TypedRepository[I]) Update(ctx context.Context, where WhereOption, value I) error {
	ds, err := c.goqu.Update(ctx, c.table)
	if err != nil {
		return err
	}
	return c.goqu.ExecuteUpdate(ctx, ds.Where(where).Set(value))
}

func (c *TypedRepository[I]) Upsert(ctx context.Context, value I) error {
	ds, err := c.goqu.Upsert(ctx, c.table)
	if err != nil {
		return err
	}
	return c.goqu.ExecuteUpsert(ctx, ds.Rows(value))
}

func (c *TypedRepository[I]) DeleteOne(ctx context.Context, keys KeysOption) error {
	ds, err := c.goqu.Delete(ctx, c.table)
	if err != nil {
		return err
	}
	return c.goqu.ExecuteDelete(ctx, ds.Where(keys))
}

func (c *TypedRepository[I]) DeleteAll(ctx context.Context, where WhereOption) error {
	ds, err := c.goqu.Delete(ctx, c.table)
	if err != nil {
		return err
	}
	return c.goqu.ExecuteDelete(ctx, ds.Where(where))
}

func (c *TypedRepository[I]) Truncate(ctx context.Context) error {
	ds, err := c.goqu.Truncate(ctx, c.table)
	if err != nil {
		return err
	}
	return c.goqu.ExecuteTruncate(ctx, ds)
}
