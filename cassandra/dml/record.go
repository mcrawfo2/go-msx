package dml

import (
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

func SeedRecords(session *gocql.Session, table string, records []Record) error {
	for _, record := range records {
		stmt, names := qb.
			Insert(table).
			Columns(record.Columns()...).
			ToCql()

		err := gocqlx.Query(session.Query(stmt), names).
			BindMap(record).
			ExecRelease()

		if err != nil {
			return err
		}
	}

	return nil
}

type RecordFunc func(session *gocql.Session, record interface{}) error

func ScanTable(session *gocql.Session, table string, columns []string, record interface{}, action RecordFunc) error {
	stmt, names := qb.Select(table).Columns(columns...).ToCql()
	query := gocqlx.Query(session.Query(stmt), names)
	defer query.Release()

	iter := query.Iter()
	for iter.Scan(record) {
		err := action(session, record)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteRecord(session *gocql.Session, table string, where ...qb.Cmp) error {
	stmt, names := qb.Delete(table).Where(where...).ToCql()
	return gocqlx.Query(session.Query(stmt), names).ExecRelease()
}
