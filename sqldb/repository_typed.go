// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type WhereOption interface {
	Expression() goqu.Expression
}

type WhereOptionFunc func() goqu.Expression

func (w WhereOptionFunc) Expression() goqu.Expression {
	return w()
}

type KeysOption = goqu.Ex

func All(exps ...WhereOption) WhereOption {
	var expressions []exp.Expression
	for _, expressioner := range exps {
		if expressioner == nil {
			continue
		}
		expressions = append(expressions, expressioner.Expression())
	}
	return exp.NewExpressionList(exp.AndType, expressions...)
}

func Any(exps ...WhereOption) WhereOption {
	var expressions []exp.Expression
	for _, expressioner := range exps {
		if expressioner == nil {
			continue
		}
		expressions = append(expressions, expressioner.Expression())
	}
	return exp.NewExpressionList(exp.OrType, expressions...)
}

func And(exp goqu.Ex) WhereOption {
	return exp
}

func Or(exp goqu.ExOr) WhereOption {
	return exp
}

//go:generate mockery --name=TypedRepositoryApi --inpackage --case=snake --with-expecter
type TypedRepositoryApi[I any] interface {
	CountAll(ctx context.Context, dest *int64, where WhereOption) error
	FindAll(ctx context.Context, dest *[]I, options ...FindAllOption) (pagingResponse paging.Response, err error)
	FindOne(ctx context.Context, dest *I, where WhereOption) error
	Insert(ctx context.Context, value ...I) error
	Update(ctx context.Context, where WhereOption, value I) error
	Upsert(ctx context.Context, value ...I) error
	DeleteOne(ctx context.Context, keys KeysOption) error
	DeleteAll(ctx context.Context, where WhereOption) error
	Truncate(ctx context.Context) error
}

type TypedRepository[I any] struct {
	table string
	goqu  GoquRepositoryApi
}

func NewTypedRepository[I any](ctx context.Context, table string) (TypedRepositoryApi[I], error) {
	api := ContextTypedRepository[I](table).Get(ctx)
	if api == nil {
		goquRepository, err := NewGoquRepository(ctx)
		if err != nil {
			return nil, err
		}

		api = &TypedRepository[I]{
			table: table,
			goqu:  goquRepository,
		}
	}
	return api, nil
}

func (c *TypedRepository[I]) CountAll(ctx context.Context, dest *int64, where WhereOption) error {
	ds := c.goqu.Select(c.table)

	if where != nil {
		ds = ds.Where(where.Expression())
	}

	return c.goqu.ExecuteGet(ctx, ds.Select(goqu.COUNT("*")), dest)
}

func (c *TypedRepository[I]) FindAll(ctx context.Context, dest *[]I, options ...FindAllOption) (pagingResponse paging.Response, err error) {
	// row query
	rowsQuery := c.goqu.Select(c.table)
	pgReq := paging.Request{}

	// Apply each option to the query and paging request
	for _, option := range options {
		rowsQuery, pgReq = option(rowsQuery, pgReq)
	}

	// Apply limit and offset when pagination is requested
	if pgReq.Size > 0 {
		rowsQuery = rowsQuery.
			Limit(pgReq.Size).
			Offset(pgReq.Page * pgReq.Size)
	}

	// Retrieve the matching rows
	err = c.goqu.ExecuteSelect(ctx, rowsQuery, dest)
	if err != nil {
		return
	}

	// Fill in default (uncounted) paging response
	pagingResponse = paging.Response{
		Content: dest,
		Size:    pgReq.Size,
		Number:  pgReq.Page,
		Sort:    pgReq.Sort,
	}

	if pgReq.Size > 0 {
		// Remove extraneous clauses from original query
		rowsQuery = rowsQuery.ClearLimit().ClearOffset().ClearOrder()

		// total items
		count := int64(0)
		countQuery := c.goqu.Select(c.table).From(rowsQuery).Select(goqu.COUNT("*"))

		err = c.goqu.ExecuteGet(ctx, countQuery, &count)
		if err != nil {
			return
		}

		pagingResponse.TotalItems = types.NewUintPtr(uint(count))
	}

	return
}

func (c *TypedRepository[I]) FindOne(ctx context.Context, dest *I, where WhereOption) error {
	ds := c.goqu.Get(c.table)

	if where != nil {
		ds = ds.Where(where.Expression())
	}

	return c.goqu.ExecuteGet(ctx, ds, dest)
}

func (c *TypedRepository[I]) Insert(ctx context.Context, value ...I) error {
	ds := c.goqu.Insert(c.table)
	return c.goqu.ExecuteInsert(ctx, ds.Rows(types.Slice[I](value).AnySlice()...))
}

func (c *TypedRepository[I]) Update(ctx context.Context, where WhereOption, value I) error {
	ds := c.goqu.Update(c.table)

	if where != nil {
		ds = ds.Where(where.Expression())
	}

	return c.goqu.ExecuteUpdate(ctx, ds.Set(value))
}

func (c *TypedRepository[I]) Upsert(ctx context.Context, value ...I) error {
	ds := c.goqu.Upsert(c.table)
	return c.goqu.ExecuteUpsert(ctx, ds.Rows(types.Slice[I](value).AnySlice()...))
}

func (c *TypedRepository[I]) DeleteOne(ctx context.Context, keys KeysOption) error {
	ds := c.goqu.Delete(c.table)
	return c.goqu.ExecuteDelete(ctx, ds.Where(keys))
}

func (c *TypedRepository[I]) DeleteAll(ctx context.Context, where WhereOption) error {
	ds := c.goqu.Delete(c.table)

	if where != nil {
		ds = ds.Where(where.Expression())
	}

	return c.goqu.ExecuteDelete(ctx, ds)
}

func (c *TypedRepository[I]) Truncate(ctx context.Context) error {
	ds := c.goqu.Truncate(c.table)
	return c.goqu.ExecuteTruncate(ctx, ds)
}
