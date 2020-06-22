package sqldb

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

var logger = log.NewLogger("msx.sql")
var ErrNotFound = sql.ErrNoRows

type CrudRepositoryApi interface {
	//CountAll(ctx context.Context, dest interface{}) error
	//CountAllBy(ctx context.Context, where map[string]interface{}, dest interface{}) error
	FindAll(ctx context.Context, dest interface{}) (err error)
	//FindAllPagedBy(ctx context.Context, where map[string]interface{}, preq paging.Request, dest interface{}) (presp paging.Response, err error)
	FindAllBy(ctx context.Context, where map[string]interface{}, dest interface{}) (err error)
	//FindAllDataSet(ctx context.Context, ds *goqu.SelectDataset, where map[string]interface{}, dest interface{}) (err error)
	FindOneBy(ctx context.Context, where map[string]interface{}, dest interface{}) (err error)
	//FindPartitionKeys(ctx context.Context, dest interface{}) (err error)
	Insert(ctx context.Context, value interface{}) (err error)
	//SaveWithTtl(ctx context.Context, value interface{}, ttl time.Duration) (err error)
	//SaveAll(ctx context.Context, values []interface{}) (err error)
	//UpdateBy(ctx context.Context, where map[string]interface{}, values map[string]interface{}) (err error)
	DeleteBy(ctx context.Context, where map[string]interface{}) (err error)
	Truncate(ctx context.Context) error
}

type CrudRepository struct {
	tableName string
	dialect   goqu.DialectWrapper
}

func (c *CrudRepository) FindAll(ctx context.Context, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect.From(c.tableName).ToSQL()
		if err != nil {
			return err
		}
		return conn.SelectContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudRepository) FindAllBy(ctx context.Context, where map[string]interface{}, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect.From(c.tableName).Where(goqu.Ex(where)).ToSQL()
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
		stmt, args, err := c.dialect.From(c.tableName).Where(goqu.Ex(where)).ToSQL()
		if err != nil {
			return err
		}
		return conn.GetContext(ctx, dest, stmt, args...)
	})
}

func (c *CrudRepository) Insert(ctx context.Context, value interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return pool.WithSqlxConnection(ctx, func(ctx context.Context, conn *sqlx.DB) error {
		stmt, args, err := c.dialect.Insert(c.tableName).Rows(value).ToSQL()
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
		stmt, args, err := c.dialect.Update(c.tableName).Where(goqu.Ex(where)).Set(value).ToSQL()
		if err != nil {
			return err
		}
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
		stmt, args, err := c.dialect.Delete(c.tableName).Where(goqu.Ex(where)).ToSQL()
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
		stmt, args, err := c.dialect.Truncate(c.tableName).ToSQL()
		if err != nil {
			return err
		}
		_, err = conn.ExecContext(ctx, stmt, args...)
		return err
	})
}

func NewCrudRepository(driver, tableName string) CrudRepositoryApi {
	return &CrudRepository{
		tableName: tableName,
		dialect:   goqu.Dialect(driver),
	}
}
