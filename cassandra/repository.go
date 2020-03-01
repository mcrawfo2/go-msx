package cassandra

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra/ddl"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/scylladb/go-reflectx"
	"github.com/scylladb/gocqlx"
	gocqlxqb "github.com/scylladb/gocqlx/qb"
	"reflect"
)

type CrudRepositoryFactoryApi interface {
	NewCrudRepository(table ddl.Table) CrudRepositoryApi
}

type ProductionCrudRepositoryFactory struct{}

func (f *ProductionCrudRepositoryFactory) NewCrudRepository(table ddl.Table) CrudRepositoryApi {
	return &CrudRepository{Table: table}
}

func NewProductionCrudRepositoryFactory() CrudRepositoryFactoryApi {
	return new(ProductionCrudRepositoryFactory)
}

type CrudRepositoryApi interface {
	FindAll(ctx context.Context, dest interface{}) (err error)
	FindAllBy(ctx context.Context, where map[string]interface{}, dest interface{}) (err error)
	FindAllCql(ctx context.Context, stmt string, names []string, where map[string]interface{}, dest interface{}) (err error)
	FindOneBy(ctx context.Context, where map[string]interface{}, dest interface{}) (err error)
	FindPartitionKeys(ctx context.Context, dest interface{}) (err error)
	Save(ctx context.Context, value interface{}) (err error)
	SaveAll(ctx context.Context, values []interface{}) (err error)
	UpdateBy(ctx context.Context, where map[string]interface{}, values map[string]interface{}) (err error)
	DeleteBy(ctx context.Context, where map[string]interface{}) (err error)
}

type CrudRepository struct {
	Table ddl.Table
}

func (r *CrudRepository) FindAll(ctx context.Context, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return
	}

	err = pool.WithSession(func(session *gocql.Session) error {
		stmt, names := gocqlxqb.
			Select(r.Table.Name).
			Columns(r.Table.ColumnNames()...).
			ToCql()

		return gocqlx.
			Query(session.Query(stmt), names).
			WithContext(ctx).
			SelectRelease(dest)
	})

	return
}

func (r *CrudRepository) FindAllBy(ctx context.Context, where map[string]interface{}, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return
	}

	err = pool.WithSession(func(session *gocql.Session) error {
		var cmps []gocqlxqb.Cmp
		for k, v := range where {
			if reflect.TypeOf(v).Kind() == reflect.Slice {
				cmps = append(cmps, gocqlxqb.InNamed(k, k))
			} else {
				cmps = append(cmps, gocqlxqb.EqNamed(k, k))
			}
		}

		stmt, names := gocqlxqb.
			Select(r.Table.Name).
			Columns(r.Table.ColumnNames()...).
			Where(cmps...).
			ToCql()

		return gocqlx.
			Query(session.Query(stmt), names).
			WithContext(ctx).
			BindMap(where).
			SelectRelease(dest)
	})

	return
}

func (r *CrudRepository) FindAllCql(ctx context.Context, stmt string, names []string, where map[string]interface{}, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return
	}

	err = pool.WithSession(func(session *gocql.Session) error {
		return gocqlx.
			Query(session.Query(stmt), names).
			WithContext(ctx).
			BindMap(where).
			SelectRelease(dest)
	})

	return
}

func (r *CrudRepository) FindPartitionKeys(ctx context.Context, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return err
	}

	err = pool.WithSession(func(session *gocql.Session) error {
		stmt, names := gocqlxqb.
			Select(r.Table.Name).
			Distinct(r.Table.PartitionKeys...).
			ToCql()

		return gocqlx.
			Query(session.Query(stmt), names).
			WithContext(ctx).
			SelectRelease(dest)
	})

	return
}

func (r *CrudRepository) FindOneBy(ctx context.Context, where map[string]interface{}, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return
	}

	err = pool.WithSession(func(session *gocql.Session) error {
		var cmps []gocqlxqb.Cmp
		for k, _ := range where {
			cmps = append(cmps, gocqlxqb.EqNamed(k, k))
		}

		stmt, names := gocqlxqb.
			Select(r.Table.Name).
			Columns(r.Table.ColumnNames()...).
			Where(cmps...).
			ToCql()

		return gocqlx.
			Query(session.Query(stmt), names).
			WithContext(ctx).
			BindMap(where).
			GetRelease(dest)
	})

	return
}

func (r *CrudRepository) Save(ctx context.Context, value interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return
	}

	err = pool.WithSessionRetry(ctx, func(session *gocql.Session) error {
		stmt, names := gocqlxqb.
			Insert(r.Table.Name).
			Columns(r.Table.ColumnNames()...).
			ToCql()

		return gocqlx.Query(session.Query(stmt), names).
			WithContext(ctx).
			BindStruct(value).
			ExecRelease()
	})

	return
}

func (r *CrudRepository) SaveAll(ctx context.Context, values []interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return
	}

	err = pool.WithSessionRetry(ctx, func(session *gocql.Session) error {

		insertBuilder := gocqlxqb.
			Insert(r.Table.Name).
			Columns(r.Table.ColumnNames()...)

		itemNames := r.Table.ColumnNames()
		batchBuilder := gocqlxqb.Batch()
		batchMap := make(map[string]interface{})
		for i, batchValue := range values {
			batchPrefix := createBatchKey(i)
			batchBuilder.AddWithPrefix(batchPrefix, insertBuilder)

			// Flatten the batchValue object into the map (one entry per field)
			mapper := gocqlx.DefaultMapper
			v := reflect.ValueOf(batchValue)
			err := mapper.TraversalsByNameFunc(v.Type(), itemNames, func(i int, t []int) error {
				if len(t) != 0 {
					val := reflectx.FieldByIndexesReadOnly(v, t)
					batchMap[batchPrefix + "." + itemNames[i]] = val.Interface()
					return nil
				} else {
					return fmt.Errorf("could not find name %q in %#v", itemNames[i], batchValue)
				}
			})
			if err != nil {
				return err
			}
		}

		stmt, names := batchBuilder.ToCql()

		return gocqlx.Query(session.Query(stmt), names).
			WithContext(ctx).
			BindMap(batchMap).
			ExecRelease()
	})

	return
}

func createBatchKey(i int) string {
	var key []byte
	scale := 16
	for {
		remainder := i % scale
		base := 97
		if i < scale && len(key) > 0 {
			base = 96
		}
		key = append([]byte{byte(base + remainder)}, key...)
		i = i / scale
		if i == 0 {
			break
		}
	}
	return string(key)
}

func (r *CrudRepository) UpdateBy(ctx context.Context, where map[string]interface{}, values map[string]interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return
	}

	err = pool.WithSessionRetry(ctx, func(session *gocql.Session) error {
		var bind = make(map[string]interface{})
		var cmps []gocqlxqb.Cmp
		for k, v := range where {
			cmps = append(cmps, gocqlxqb.EqNamed(k, k+"Where"))
			bind[k+"Where"] = v
		}

		var sets []string
		for k, v := range values {
			sets = append(sets, k)
			bind[k] = v
		}

		stmt, names := gocqlxqb.
			Update(r.Table.Name).
			Set(sets...).
			Where(cmps...).
			ToCql()

		return gocqlx.Query(session.Query(stmt), names).
			WithContext(ctx).
			BindMap(bind).
			ExecRelease()
	})

	return
}

func (r *CrudRepository) DeleteBy(ctx context.Context, where map[string]interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return
	}

	err = pool.WithSessionRetry(ctx, func(session *gocql.Session) error {
		var cmps []gocqlxqb.Cmp
		for k, _ := range where {
			cmps = append(cmps, gocqlxqb.EqNamed(k, k))
		}

		stmt, names := gocqlxqb.
			Delete(r.Table.Name).
			Where(cmps...).
			ToCql()

		return gocqlx.Query(session.Query(stmt), names).
			WithContext(ctx).
			BindMap(where).
			ExecRelease()
	})

	return
}
