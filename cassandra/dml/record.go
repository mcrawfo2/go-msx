// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package dml

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra"
	"errors"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

type Record map[string]interface{}

func (r Record) Columns() []string {
	var result []string
	for k := range r {
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

func ScanTable(ctx context.Context, session *gocql.Session, table string, columns []string, record interface{}, action RecordFunc) (err error) {
	stmt, names := qb.Select(table).Columns(columns...).ToCql()
	query := gocqlx.Query(session.Query(stmt), names).WithContext(ctx)
	defer query.Release()

	iter := query.Iter()
	defer func() {
		err = iter.Close()
	}()

	for iter.StructScan(record) {
		err = action(ctx, session, record)
		if err != nil {
			return
		}
	}
	return
}

func DeleteRecord(ctx context.Context, session *gocql.Session, table string, where ...qb.Cmp) error {
	stmt, names := qb.Delete(table).Where(where...).ToCql()
	return gocqlx.Query(session.Query(stmt), names).WithContext(ctx).ExecRelease()
}

func TableExists(ctx context.Context, session *gocql.Session, table string) (bool, error) {
	if len(table) == 0 {
		return false, errors.New("table name is empty")
	}
	cassandraPool, err := cassandra.PoolFromContext(ctx)
	if err != nil {
		return false, err
	}
	schemas, err := session.KeyspaceMetadata(cassandraPool.ClusterConfig().KeyspaceName)
	if err != nil {
		return false, err
	}
	tables := schemas.Tables
	if _, ok := tables[table]; !ok {
		return false, nil
	}

	return true, nil
}
