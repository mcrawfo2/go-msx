package dml

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

type Record goqu.Record

func (r Record) Columns() []string {
	return goqu.Record(r).Cols()
}

type RecordSet []goqu.Record

func (rs RecordSet) ToInterfaceSlice() []interface{} {
	var result []interface{}
	for _, r := range rs {
		result = append(result, r)
	}
	return result
}

func SeedRecords(ctx context.Context, db *sqlx.DB, table string, records []interface{}) error {
	dialect := goqu.Dialect(db.DriverName())

	stmt, args, err := dialect.Insert(table).Rows(records).ToSQL()
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, stmt, args...)
	return err
}

func UpdateRecord(ctx context.Context, db *sqlx.DB, table string, where goqu.Ex, record interface{}) error {

	dialect := goqu.Dialect(db.DriverName())

	stmt, _, err := dialect.Update(table).Set(record).Where(where).ToSQL()
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, stmt)
	return err
}

type RecordFunc func(ctx context.Context, db *sqlx.DB, record interface{}) error

func ScanTable(ctx context.Context, db *sqlx.DB, table string, columns []string, record interface{}, action RecordFunc) error {
	dialect := goqu.Dialect(db.DriverName())

	stmt, _, err := dialect.Select(types.StringSlice(columns).ToInterfaceSlice()...).From(table).ToSQL()
	if err != nil {
		return err
	}

	rows, err := db.QueryxContext(ctx, stmt)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(record)
		if err != nil {
			return err
		}

		err = action(ctx, db, record)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteRecord(ctx context.Context, db *sqlx.DB, table string, where goqu.Ex) error {
	dialect := goqu.Dialect(db.DriverName())

	stmt, _, err := dialect.Delete(table).Where(where).ToSQL()
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, stmt)
	return err
}

func TruncateTable(ctx context.Context, db *sqlx.DB, table string, cascade bool) error {
	dialect := goqu.Dialect(db.DriverName())

	td := dialect.Truncate(table)
	if cascade {
		td = td.Cascade()
	}

	stmt, _, err := td.ToSQL()
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, stmt)
	return err
}
