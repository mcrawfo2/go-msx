package sqldb

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"strings"
)

var logger = log.NewLogger("msx.sql")
var ErrNotFound = sql.ErrNoRows
var ErrNotImplemented = errors.New("Feature not implemented")

type CrudRepositoryApi interface {
	CountAll(ctx context.Context, dest *int64) error
	CountAllBy(ctx context.Context, where map[string]interface{}, dest *int64) error
	CountAllByExpression(ctx context.Context, where goqu.Expression, dest *int64) error
	FindAll(ctx context.Context, dest interface{}) (err error)
	FindAllPagedBy(ctx context.Context, where map[string]interface{}, preq paging.Request, dest interface{}) (presp paging.Response, err error)
	FindAllBy(ctx context.Context, where map[string]interface{}, dest interface{}) (err error)
	FindAllByExpression(ctx context.Context, where goqu.Expression, dest interface{}) (err error)
	FindAllPagedByExpression(ctx context.Context, where goqu.Expression, preq paging.Request, dest interface{}) (presp paging.Response, err error)
	FindAllDistinctBy(ctx context.Context, distinct []string, where map[string]interface{}, dest interface{}) (err error)
	FindAllSortedBy(ctx context.Context, where map[string]interface{}, sortOrder paging.SortOrder, dest interface{}) (err error)
	FindAllSortedByExpression(ctx context.Context, where goqu.Expression, sortOrder paging.SortOrder, dest interface{}) (err error)
	FindOneBy(ctx context.Context, where map[string]interface{}, dest interface{}) (err error)
	FindOneSortedBy(ctx context.Context, where map[string]interface{}, sortOrder paging.SortOrder, dest interface{}) (err error)
	Insert(ctx context.Context, value interface{}) (err error)
	Update(ctx context.Context, where map[string]interface{}, value interface{}) (err error)
	Save(ctx context.Context, value interface{}) (err error)
	SaveAll(ctx context.Context, values []interface{}) (err error)
	DeleteBy(ctx context.Context, where map[string]interface{}) (err error)
	Truncate(ctx context.Context) error
}

// Deprecated.  Use CrudPreparedRepository instead.
type CrudRepository struct {
	tableName string
}

func (c *CrudRepository) CountAll(ctx context.Context, dest *int64) error {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).From(c.tableName).Select(goqu.COUNT("*")).ToSQL()
		if err != nil {
			return err
		}
		return conn.GetContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudRepository) CountAllBy(ctx context.Context, where map[string]interface{}, dest *int64) error {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).From(c.tableName).Select(goqu.COUNT("*")).Where(goqu.Ex(where)).ToSQL()
		if err != nil {
			return err
		}
		return conn.GetContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudRepository) CountAllByExpression(ctx context.Context, where goqu.Expression, dest *int64) error {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		if where == nil {
			where = goqu.Literal("true")
		}
		stmt, args, err := c.dialect(conn).From(c.tableName).Select(goqu.COUNT("*")).Where(where).ToSQL()
		if err != nil {
			return err
		}
		return conn.GetContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudRepository) dialect(conn *sqlx.DB) goqu.DialectWrapper {
	return goqu.Dialect(conn.DriverName())
}

func (c *CrudRepository) FindAll(ctx context.Context, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).From(c.tableName).ToSQL()
		if err != nil {
			return err
		}
		return conn.SelectContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudRepository) FindAllPagedBy(ctx context.Context, where map[string]interface{}, pagingRequest paging.Request, dest interface{}) (pagingResponse paging.Response, err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return pagingResponse, err
	}

	err = pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		selectDataSet := c.dialect(conn).
			From(c.tableName).
			Where(goqu.Ex(where)).
			Limit(pagingRequest.Size).
			Offset(pagingRequest.Page * pagingRequest.Size)

		for _, sortOrder := range pagingRequest.Sort {
			ident := goqu.I(sortOrder.Property)
			switch sortOrder.Direction {
			case paging.SortDirectionDesc:
				selectDataSet = selectDataSet.OrderAppend(ident.Desc())
			default:
				selectDataSet = selectDataSet.OrderAppend(ident.Asc())
			}
		}

		stmt, args, err := selectDataSet.ToSQL()
		if err != nil {
			return err
		}

		err = conn.SelectContext(ctx, dest, stmt, args...)
		if err != nil {
			return err
		}

		pagingResponse = paging.Response{
			Content: dest,
			Size:    pagingRequest.Size,
			Number:  pagingRequest.Page,
			Sort:    pagingRequest.Sort,
		}

		return nil
	})
	return pagingResponse, err
}

func (c *CrudRepository) FindAllBy(ctx context.Context, where map[string]interface{}, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).From(c.tableName).Where(goqu.Ex(where)).ToSQL()
		if err != nil {
			return err
		}
		return conn.SelectContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudRepository) FindAllByExpression(ctx context.Context, where goqu.Expression, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		if where == nil {
			where = goqu.Literal("true")
		}
		stmt, args, err := c.dialect(conn).From(c.tableName).Where(where).ToSQL()
		if err != nil {
			return err
		}
		return conn.SelectContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudRepository) FindAllPagedByExpression(ctx context.Context, where goqu.Expression, pagingRequest paging.Request, dest interface{}) (pagingResponse paging.Response, err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return
	}

	err = pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		if where == nil {
			where = goqu.Literal("true")
		}

		ds := c.dialect(conn).
			From(c.tableName).
			Where(where).
			Limit(pagingRequest.Size).
			Offset(pagingRequest.Page * pagingRequest.Size)

		for _, sortOrder := range pagingRequest.Sort {
			ident := goqu.I(sortOrder.Property)
			switch sortOrder.Direction {
			case paging.SortDirectionDesc:
				ds = ds.OrderAppend(ident.Desc())
			default:
				ds = ds.OrderAppend(ident.Asc())
			}
		}

		stmt, args, err := ds.ToSQL()
		if err != nil {
			return err
		}

		err = conn.SelectContext(ctx, dest, stmt, args...)
		if err != nil {
			return err
		}

		pagingResponse = paging.Response{
			Content: dest,
			Size:    pagingRequest.Size,
			Number:  pagingRequest.Page,
			Sort:    pagingRequest.Sort,
		}

		return nil

	})

	return
}

func (c *CrudRepository) FindAllDistinctBy(ctx context.Context, distinct []string, where map[string]interface{}, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).From(c.tableName).Distinct(distinct).Where(goqu.Ex(where)).ToSQL()
		if err != nil {
			return err
		}
		return conn.SelectContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudRepository) FindAllSortedBy(ctx context.Context, where map[string]interface{}, sortOrder paging.SortOrder, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.constructSortedQueryWithArgs(conn, where, sortOrder)
		if err != nil {
			return err
		}
		return conn.SelectContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudRepository) FindAllSortedByExpression(ctx context.Context, where goqu.Expression, sortOrder paging.SortOrder, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		if where == nil {
			where = goqu.Literal("true")
		}
		stmt, args, err := c.constructSortedQueryWithArgsByExpression(conn, where, sortOrder)
		if err != nil {
			return err
		}
		return conn.SelectContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudRepository) FindOneBy(ctx context.Context, where map[string]interface{}, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).From(c.tableName).Where(goqu.Ex(where)).ToSQL()
		if err != nil {
			return err
		}
		return conn.GetContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudRepository) FindOneSortedBy(ctx context.Context, where map[string]interface{}, sortOrder paging.SortOrder, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.constructSortedQueryWithArgs(conn, where, sortOrder)
		if err != nil {
			return err
		}

		return conn.GetContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudRepository) constructSortedQueryWithArgs(conn *sqlx.DB, where map[string]interface{}, sortOrder paging.SortOrder) (sql string, args []interface{}, err error) {
	selectDataSet := c.dialect(conn).
		From(c.tableName).
		Where(goqu.Ex(where))

	ident := goqu.I(sortOrder.Property)
	switch sortOrder.Direction {
	case paging.SortDirectionDesc:
		selectDataSet = selectDataSet.OrderAppend(ident.Desc())
	default:
		selectDataSet = selectDataSet.OrderAppend(ident.Asc())
	}

	sql, args, err = selectDataSet.ToSQL()

	return
}

func (c *CrudRepository) constructSortedQueryWithArgsByExpression(conn *sqlx.DB, where goqu.Expression, sortOrder paging.SortOrder) (sql string, args []interface{}, err error) {
	selectDataSet := c.dialect(conn).
		From(c.tableName).
		Where(where)

	ident := goqu.I(sortOrder.Property)
	switch sortOrder.Direction {
	case paging.SortDirectionDesc:
		selectDataSet = selectDataSet.OrderAppend(ident.Desc())
	default:
		selectDataSet = selectDataSet.OrderAppend(ident.Asc())
	}

	sql, args, err = selectDataSet.ToSQL()

	return
}

func (c *CrudRepository) Insert(ctx context.Context, value interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).Insert(c.tableName).Rows(value).ToSQL()
		if err != nil {
			return err
		}
		_, err = conn.ExecContext(ctx, stmt, args...)
		return err
	})
}

func (c *CrudRepository) Update(ctx context.Context, where map[string]interface{}, value interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).Update(c.tableName).Where(goqu.Ex(where)).Set(value).ToSQL()
		if err != nil {
			return err
		}
		_, err = conn.ExecContext(ctx, stmt, args...)
		return err
	})
}

func (c *CrudRepository) Save(ctx context.Context, value interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	// Cockroach only for now
	if pool.Config().Driver != "postgres" {
		return ErrNotImplemented
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).Insert(c.tableName).Rows(value).ToSQL()
		if err != nil {
			return err
		}

		stmt = "UPSERT" + strings.TrimPrefix(stmt, "INSERT")

		_, err = conn.ExecContext(ctx, stmt, args...)
		return err
	})
}

func (c *CrudRepository) SaveAll(ctx context.Context, values []interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	// Cockroach only for now
	if pool.Config().Driver != "postgres" {
		return ErrNotImplemented
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).Insert(c.tableName).Rows(values...).ToSQL()
		if err != nil {
			return err
		}

		stmt = "UPSERT" + strings.TrimPrefix(stmt, "INSERT")

		_, err = conn.ExecContext(ctx, stmt, args...)
		return err
	})
}

func (c *CrudRepository) DeleteBy(ctx context.Context, where map[string]interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).Delete(c.tableName).Where(goqu.Ex(where)).ToSQL()
		if err != nil {
			return err
		}
		_, err = conn.ExecContext(ctx, stmt, args...)
		return err
	})
}

func (c *CrudRepository) Truncate(ctx context.Context) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).Truncate(c.tableName).ToSQL()
		if err != nil {
			return err
		}
		_, err = conn.ExecContext(ctx, stmt, args...)
		return err
	})
}

// Deprecated.  Use newCrudPreparedRepository instead.
func newCrudRepository(tableName string) CrudRepositoryApi {
	return &CrudRepository{
		tableName: tableName,
	}
}
