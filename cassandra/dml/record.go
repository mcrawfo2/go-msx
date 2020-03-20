package dml

import (
	"context"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

type Record map[string]interface{}

func (r Record) Columns() []string {
	var result []string
	for k, _ := range r {
		result = append(result, k)
	}
	return result
}

func SeedRecords(ctx context.Context, session *gocql.Session, table string, records []Record) error {
	for _, record := range records {
		stmt, names := qb.
			Insert(table).
			Columns(record.Columns()...).
			ToCql()

		err := gocqlx.Query(session.Query(stmt), names).
			WithContext(ctx).
			BindMap(record).
			ExecRelease()

		if err != nil {
			return err
		}
	}

	return nil
}

type RecordFunc func(ctx context.Context, session *gocql.Session, record interface{}) error

func ScanTable(ctx context.Context, session *gocql.Session, table string, columns []string, record interface{}, action RecordFunc) error {
	stmt, names := qb.Select(table).Columns(columns...).ToCql()
	query := gocqlx.Query(session.Query(stmt), names)
	defer query.Release()

	iter := query.Iter()
	for iter.Scan(record) {
		err := action(ctx, session, record)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteRecord(ctx context.Context, session *gocql.Session, table string, where ...qb.Cmp) error {
	stmt, names := qb.Delete(table).Where(where...).ToCql()
	return gocqlx.Query(session.Query(stmt), names).WithContext(ctx).ExecRelease()
}
