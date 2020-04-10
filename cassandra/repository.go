package cassandra

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra/ddl"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"encoding/base64"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/scylladb/go-reflectx"
	"github.com/scylladb/gocqlx"
	gocqlxqb "github.com/scylladb/gocqlx/qb"
	"reflect"
	"time"
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
	CountAll(ctx context.Context, dest interface{}) error
	CountAllBy(ctx context.Context, where map[string]interface{}, dest interface{}) error
	FindAll(ctx context.Context, dest interface{}) (err error)
	FindAllPagedBy(ctx context.Context, where map[string]interface{}, preq paging.Request, dest interface{}) (presp paging.Response, err error)
	FindAllBy(ctx context.Context, where map[string]interface{}, dest interface{}) (err error)
	FindAllCql(ctx context.Context, stmt string, names []string, where map[string]interface{}, dest interface{}) (err error)
	FindAllByLuceneSearch(ctx context.Context, index, search string, dest interface{}) (err error)
	FindAllByPagedLuceneSearch(ctx context.Context, index, search string, dest interface{}, request paging.Request) (response paging.Response, err error)
	FindOneBy(ctx context.Context, where map[string]interface{}, dest interface{}) (err error)
	FindPartitionKeys(ctx context.Context, dest interface{}) (err error)
	Save(ctx context.Context, value interface{}) (err error)
	SaveWithTtl(ctx context.Context, value interface{}, ttl time.Duration) (err error)
	SaveAll(ctx context.Context, values []interface{}) (err error)
	UpdateBy(ctx context.Context, where map[string]interface{}, values map[string]interface{}) (err error)
	DeleteBy(ctx context.Context, where map[string]interface{}) (err error)
}

type CrudRepository struct {
	Table ddl.Table
}

func (r *CrudRepository) CountAll(ctx context.Context, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return
	}

	err = pool.WithSession(func(session *gocql.Session) error {
		stmt, names := gocqlxqb.
			Select(r.Table.Name).
			CountAll().
			ToCql()

		return gocqlx.
			Query(session.Query(stmt), names).
			WithContext(ctx).
			GetRelease(dest)
	})

	return
}

func (r *CrudRepository) CountAllBy(ctx context.Context, where map[string]interface{}, dest interface{}) (err error) {
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
			CountAll().
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

func (r *CrudRepository) FindAllPagedBy(ctx context.Context, where map[string]interface{}, request paging.Request, dest interface{}) (response paging.Response, err error) {
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

		pageState, err := r.getPageState(ctx, session, stmt, names, where, request)
		if err != nil {
			return err
		} else if len(pageState) == 0 && request.Page > 0 {
			response = paging.Response{
				Content: dest,
				Size:    request.Size,
				Number:  request.Page,
				Sort:    request.Sort,
				State:   nil,
			}
			return nil
		}

		qx := gocqlx.
			Query(session.Query(stmt), names).
			WithContext(ctx).
			BindMap(where).
			PageState(pageState).
			PageSize(int(request.Size))

		iter := gocqlx.Iter(qx.Query)
		defer qx.Release()
		err = iter.Select(dest)
		if err != nil {
			return err
		}

		response = paging.Response{
			Content: dest,
			Size:    request.Size,
			Number:  request.Page,
			Sort:    request.Sort,
			State:   base64.StdEncoding.EncodeToString(iter.PageState()),
		}

		return nil
	})

	return
}

// Find the paging state for the requested page by executing the query for all prior records
func (r *CrudRepository) getPageState(ctx context.Context, session *gocql.Session, stmt string, names []string, where map[string]interface{}, request paging.Request) (pageState []byte, err error) {
	if request.State != nil {
		switch request.State.(type) {
		case *string:
			pagingStateStringPtr := request.State.(*string)
			if pagingStateStringPtr != nil {
				return base64.StdEncoding.DecodeString(*pagingStateStringPtr)
			}
		}
	} else if request.Page == 0 {
		return nil, nil
	}

	qx := gocqlx.
		Query(session.Query(stmt), names).
		WithContext(ctx).
		PageSize(int(request.Page) * int(request.Size)).
		BindMap(where)
	defer qx.Release()

	iter := qx.Query.Iter()
	defer iter.Close()
	return iter.PageState(), nil
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

func (r *CrudRepository) FindAllByLuceneSearch(ctx context.Context, index, search string, dest interface{}) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return
	}

	err = pool.WithSession(func(session *gocql.Session) error {
		const luceneQuery = "luceneQuery"

		stmt, names := gocqlxqb.
			Select(r.Table.Name).
			Columns(r.Table.ColumnNames()...).
			ToCql()

		stmt += fmt.Sprintf(" where expr(%s,?)", index)
		names = append(names, luceneQuery)
		where := map[string]interface{}{
			luceneQuery: search,
		}

		return gocqlx.
			Query(session.Query(stmt), names).
			WithContext(ctx).
			BindMap(where).
			SelectRelease(dest)
	})

	return
}

func (r *CrudRepository) FindAllByPagedLuceneSearch(ctx context.Context, index, search string, dest interface{}, request paging.Request) (response paging.Response, err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return
	}

	err = pool.WithSession(func(session *gocql.Session) error {
		const luceneQuery = "luceneQuery"

		stmt, names := gocqlxqb.
			Select(r.Table.Name).
			Columns(r.Table.ColumnNames()...).
			ToCql()

		stmt += fmt.Sprintf(" where expr(%s,?)", index)
		names = append(names, luceneQuery)
		where := map[string]interface{}{
			luceneQuery: search,
		}

		pageState, err := r.getPageState(ctx, session, stmt, names, where, request)
		if err != nil {
			return err
		} else if len(pageState) == 0 && request.Page > 0 {
			response = paging.Response{
				Content: dest,
				Size:    request.Size,
				Number:  request.Page,
				Sort:    request.Sort,
				State:   nil,
			}

			return nil
		}

		return gocqlx.
			Query(session.Query(stmt), names).
			WithContext(ctx).
			BindMap(where).
			PageState(pageState).
			PageSize(int(request.Size)).
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

func (r *CrudRepository) SaveWithTtl(ctx context.Context, value interface{}, ttl time.Duration) (err error) {
	pool, err := PoolFromContext(ctx)
	if err != nil {
		return
	}

	err = pool.WithSessionRetry(ctx, func(session *gocql.Session) error {
		stmt, names := gocqlxqb.
			Insert(r.Table.Name).
			Columns(r.Table.ColumnNames()...).
			TTL(ttl).
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
					batchMap[batchPrefix+"."+itemNames[i]] = val.Interface()
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
