// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"modernc.org/sqlite"
	"testing"
)

func TestSuite_GoquRepository(t *testing.T) {
	table := "persons2"
	ctx := sqlite3DBContext()

	err := ConfigurePool(ctx)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
	ctx = ContextWithPool(ctx)

	myPool, err := PoolFromContext(ctx)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}

	drivers["sqlite3"] = &sqlite.Driver{}

	db, err := myPool.NewSqlConnection()
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}

	rgoqu := NewGoquRepository(ctx)

	_, err = db.Exec("DROP TABLE IF EXISTS "+table, nil)
	assert.NoError(t, err)

	// sqlite does not have uuid data type
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS "+table+" (\nid VARCHAR(16) NOT NULL PRIMARY KEY,\nname TEXT\n);", nil)
	assert.NoError(t, err)

	person1 := Person{Id: uuid.MustParse("437f96b0-6722-11ed-9022-0242ac120005"), Name: "Jonee"}

	dsInsert, err := rgoqu.Insert(ctx, table)
	assert.NoError(t, err)

	err = rgoqu.ExecuteInsert(ctx, dsInsert.Rows(person1))
	assert.NoError(t, err)

	// no upsert for sqlite3
	/*
		person1.Name = "Jonee6"
		dsUpsert, err := rgoqu.Upsert(ctx, table)
		assert.NoError(t, err)

		err = rgoqu.ExecuteUpsert(ctx, dsUpsert.Rows(person1))
		assert.NoError(t, err)
	*/

	person1.Name = "Jonee7"
	dsUpdate, err := rgoqu.Update(ctx, table)
	assert.NoError(t, err)

	err = rgoqu.ExecuteUpdate(ctx, dsUpdate.Where(goqu.Ex(map[string]interface{}{"id": person1.Id})).Set(person1))
	assert.NoError(t, err)

	var destPerson Person
	dsGet, err := rgoqu.Get(ctx, table)
	assert.NoError(t, err)

	err = rgoqu.ExecuteGet(ctx, dsGet.Where(goqu.Ex(map[string]interface{}{"id": person1.Id})), &destPerson)
	assert.NoError(t, err)
	logger.WithContext(ctx).Info(destPerson)

	var destPersons []Person
	dsSelect, err := rgoqu.Select(ctx, table)
	assert.NoError(t, err)

	err = rgoqu.ExecuteSelect(ctx, dsSelect.Where(goqu.Ex(map[string]interface{}{"name": person1.Name})), &destPersons)
	assert.NoError(t, err)
	logger.WithContext(ctx).Info(destPersons)

	dsDelete, err := rgoqu.Delete(ctx, table)
	assert.NoError(t, err)

	err = rgoqu.ExecuteDelete(ctx, dsDelete.Where(goqu.Ex(map[string]interface{}{"id": person1.Id})))
	assert.NoError(t, err)

	// no truncate for sqlite3
	/*
		dsTruncate, err := rgoqu.Truncate(ctx, table)
		assert.NoError(t, err)

		err = rgoqu.ExecuteTruncate(ctx, dsTruncate)
		assert.NoError(t, err)
	*/

	_, err = db.Exec("DROP TABLE "+table, nil)
	assert.NoError(t, err)
}

func TestGoquRepository_ExecuteInsert(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mock.ExpectExec(`INSERT INTO "persons"`).WithArgs(uuid.MustParse(mockId), mockName).WillReturnResult(sqlmock.NewResult(1, 1))

	rgoqu := NewGoquRepository(ctx)

	person1 := Person{Id: uuid.MustParse(mockId), Name: mockName}

	dsInsert, err := rgoqu.Insert(ctx, "persons")
	assert.NoError(t, err)

	err = rgoqu.ExecuteInsert(ctx, dsInsert.Rows(person1))
	assert.NoError(t, err)
}

func TestGoquRepository_ExecuteUpsert(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mock.ExpectExec(`UPSERT INTO "persons"`).WithArgs(uuid.MustParse(mockId), mockName).WillReturnResult(sqlmock.NewResult(1, 1))

	rgoqu := NewGoquRepository(ctx)

	person1 := Person{Id: uuid.MustParse(mockId), Name: mockName}

	dsUpsert, err := rgoqu.Upsert(ctx, "persons")
	assert.NoError(t, err)

	err = rgoqu.ExecuteUpsert(ctx, dsUpsert.Rows(person1))
	assert.NoError(t, err)
}

func TestGoquRepository_ExecuteUpdate(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mock.ExpectExec(`UPDATE "persons"`).WithArgs(uuid.MustParse(mockId), mockName, uuid.MustParse(mockId)).WillReturnResult(sqlmock.NewResult(1, 1))

	rgoqu := NewGoquRepository(ctx)

	person1 := Person{Id: uuid.MustParse(mockId), Name: mockName}

	dsUpdate, err := rgoqu.Update(ctx, "persons")
	assert.NoError(t, err)

	err = rgoqu.ExecuteUpdate(ctx, dsUpdate.Where(goqu.Ex(map[string]interface{}{"id": person1.Id})).Set(person1))
	assert.NoError(t, err)
}

func TestGoquRepository_ExecuteGet(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	columns := []string{"id", "name"}
	mock.ExpectQuery(`SELECT (.+) FROM "persons" WHERE`).
		WithArgs(mockId).
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString(mockId + "," + mockName))

	rgoqu := NewGoquRepository(ctx)

	var destPerson Person
	dsGet, err := rgoqu.Get(ctx, "persons")
	assert.NoError(t, err)

	err = rgoqu.ExecuteGet(ctx, dsGet.Where(goqu.Ex(map[string]interface{}{"id": uuid.MustParse(mockId)})), &destPerson)
	assert.NoError(t, err)
	logger.WithContext(ctx).Info(destPerson)
}

func TestGoquRepository_ExecuteSelect(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(uuid.MustParse(mockId), mockName)
	mock.ExpectQuery(`SELECT (.+) FROM "persons"`).WillReturnRows(rows)

	rgoqu := NewGoquRepository(ctx)

	var destPersons []Person
	dsSelect, err := rgoqu.Select(ctx, "persons")
	assert.NoError(t, err)

	err = rgoqu.ExecuteSelect(ctx, dsSelect.Where(goqu.Ex(map[string]interface{}{"name": mockName})), &destPersons)
	assert.NoError(t, err)
	logger.WithContext(ctx).Info(destPersons)
}

func TestGoquRepository_ExecuteDelete(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mock.ExpectExec(`DELETE FROM "persons"`).WithArgs(uuid.MustParse(mockId)).WillReturnResult(sqlmock.NewResult(1, 1))

	rgoqu := NewGoquRepository(ctx)

	dsDelete, err := rgoqu.Delete(ctx, "persons")
	assert.NoError(t, err)

	err = rgoqu.ExecuteDelete(ctx, dsDelete.Where(goqu.Ex(map[string]interface{}{"id": uuid.MustParse(mockId)})))
	assert.NoError(t, err)
}

func TestGoquRepository_ExecuteTruncate(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mock.ExpectExec(`TRUNCATE "persons"`).WillReturnResult(sqlmock.NewResult(1, 1))

	rgoqu := NewGoquRepository(ctx)

	dsTruncate, err := rgoqu.Truncate(ctx, "persons")
	assert.NoError(t, err)

	err = rgoqu.ExecuteTruncate(ctx, dsTruncate)
	assert.NoError(t, err)
}
