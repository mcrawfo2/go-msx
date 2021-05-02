package sqldb

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"strings"
)

type CrudPreparedRepository struct {
	tableName string
}

func (c *CrudPreparedRepository) Rebind(conn *sqlx.DB, stmt string) string {
	driver := conn.DriverName()
	baseDriver := baseDriverName(driver)
	bindType := sqlx.BindType(baseDriver)
	return sqlx.Rebind(bindType, stmt)
}

func (c *CrudPreparedRepository) CountAll(ctx context.Context, dest *int64) error {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).From(c.tableName).Select(goqu.COUNT("*")).Prepared(true).ToSQL()
		if err != nil {
			return err
		}
		stmt = c.Rebind(conn, stmt)
		return conn.GetContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudPreparedRepository) CountAllBy(ctx context.Context, where map[string]interface{}, dest *int64) error {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).From(c.tableName).Select(goqu.COUNT("*")).Where(goqu.Ex(where)).Prepared(true).ToSQL()
		if err != nil {
			return err
		}
		stmt = c.Rebind(conn, stmt)
		return conn.GetContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudPreparedRepository) dialect(conn *sqlx.DB) goqu.DialectWrapper {
	return goqu.Dialect(conn.DriverName())
}

func (c *CrudPreparedRepository) FindAll(ctx context.Context, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).From(c.tableName).Prepared(true).ToSQL()
		if err != nil {
			return err
		}
		stmt = c.Rebind(conn, stmt)
		return conn.SelectContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudPreparedRepository) FindAllPagedBy(ctx context.Context, where map[string]interface{}, pagingRequest paging.Request, dest interface{}) (pagingResponse paging.Response, err error) {
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

		stmt, args, err := selectDataSet.Prepared(true).ToSQL()
		if err != nil {
			return err
		}

		stmt = c.Rebind(conn, stmt)
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

func (c *CrudPreparedRepository) FindAllBy(ctx context.Context, where map[string]interface{}, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).From(c.tableName).Where(goqu.Ex(where)).Prepared(true).ToSQL()
		if err != nil {
			return err
		}
		stmt = c.Rebind(conn, stmt)
		return conn.SelectContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudPreparedRepository) FindAllByExpression(ctx context.Context, where goqu.Expression, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		if where == nil {
			where = goqu.Literal("true")
		}
		stmt, args, err := c.dialect(conn).From(c.tableName).Where(where).Prepared(true).ToSQL()
		if err != nil {
			return err
		}
		stmt = c.Rebind(conn, stmt)
		return conn.SelectContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudPreparedRepository) FindAllDistinctBy(ctx context.Context, distinct []string, where map[string]interface{}, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).From(c.tableName).Distinct(distinct).Where(goqu.Ex(where)).Prepared(true).ToSQL()
		if err != nil {
			return err
		}
		stmt = c.Rebind(conn, stmt)
		return conn.SelectContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudPreparedRepository) FindAllSortedBy(ctx context.Context, where map[string]interface{}, sortOrder paging.SortOrder, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.constructSortedQueryWithArgs(conn, where, sortOrder)
		if err != nil {
			return err
		}
		stmt = c.Rebind(conn, stmt)
		return conn.SelectContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudPreparedRepository) FindOneBy(ctx context.Context, where map[string]interface{}, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).From(c.tableName).Where(goqu.Ex(where)).Prepared(true).ToSQL()
		if err != nil {
			return err
		}
		stmt = c.Rebind(conn, stmt)
		return conn.GetContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudPreparedRepository) FindOneSortedBy(ctx context.Context, where map[string]interface{}, sortOrder paging.SortOrder, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.constructSortedQueryWithArgs(conn, where, sortOrder)
		if err != nil {
			return err
		}
		stmt = c.Rebind(conn, stmt)
		return conn.GetContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudPreparedRepository) constructSortedQueryWithArgs(conn *sqlx.DB, where map[string]interface{}, sortOrder paging.SortOrder) (sql string, args []interface{}, err error) {
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

	sql, args, err = selectDataSet.Prepared(true).ToSQL()

	return
}

func (c *CrudPreparedRepository) Insert(ctx context.Context, value interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).Insert(c.tableName).Rows(value).Prepared(true).ToSQL()
		if err != nil {
			return err
		}
		stmt = c.Rebind(conn, stmt)
		_, err = conn.ExecContext(ctx, stmt, args...)
		return err
	})
}

func (c *CrudPreparedRepository) Update(ctx context.Context, where map[string]interface{}, value interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).Update(c.tableName).Where(goqu.Ex(where)).Set(value).Prepared(true).ToSQL()
		if err != nil {
			return err
		}
		stmt = c.Rebind(conn, stmt)
		_, err = conn.ExecContext(ctx, stmt, args...)
		return err
	})
}

func (c *CrudPreparedRepository) Save(ctx context.Context, value interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	// Cockroach only for now
	if pool.Config().Driver != "postgres" {
		return ErrNotImplemented
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).Insert(c.tableName).Rows(value).Prepared(true).ToSQL()
		if err != nil {
			return err
		}

		stmt = "UPSERT" + strings.TrimPrefix(stmt, "INSERT")
		stmt = c.Rebind(conn, stmt)
		_, err = conn.ExecContext(ctx, stmt, args...)
		return err
	})
}

func (c *CrudPreparedRepository) SaveAll(ctx context.Context, values []interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	// Cockroach only for now
	if pool.Config().Driver != "postgres" {
		return ErrNotImplemented
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).Insert(c.tableName).Rows(values...).Prepared(true).ToSQL()
		if err != nil {
			return err
		}

		stmt = "UPSERT" + strings.TrimPrefix(stmt, "INSERT")
		stmt = c.Rebind(conn, stmt)
		_, err = conn.ExecContext(ctx, stmt, args...)
		return err
	})
}

func (c *CrudPreparedRepository) DeleteBy(ctx context.Context, where map[string]interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).Delete(c.tableName).Where(goqu.Ex(where)).Prepared(true).ToSQL()
		if err != nil {
			return err
		}
		stmt = c.Rebind(conn, stmt)
		_, err = conn.ExecContext(ctx, stmt, args...)
		return err
	})
}

func (c *CrudPreparedRepository) Truncate(ctx context.Context) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect(conn).Truncate(c.tableName).Prepared(true).ToSQL()
		if err != nil {
			return err
		}
		stmt = c.Rebind(conn, stmt)
		_, err = conn.ExecContext(ctx, stmt, args...)
		return err
	})
}

func newCrudPreparedRepository(tableName string) CrudRepositoryApi {
	return &CrudPreparedRepository{
		tableName: tableName,
	}
}
