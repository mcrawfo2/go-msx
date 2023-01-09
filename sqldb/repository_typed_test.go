// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"modernc.org/sqlite"
	"testing"
)

func TestSuite_TypedRepository(t *testing.T) {
	table := "persons3"
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

	personsRepo := NewTypedRepository[Person](ctx, table)

	_, err = db.Exec("DROP TABLE IF EXISTS "+table, nil)
	assert.NoError(t, err)

	// sqlite does not have uuid data type
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS "+table+" (\nid VARCHAR(16) NOT NULL PRIMARY KEY,\nname TEXT\n);", nil)
	assert.NoError(t, err)

	person1 := Person{Id: uuid.MustParse("437f96b0-6722-11ed-9022-0242ac120009"), Name: "Jonee"}
	err = personsRepo.Insert(ctx, person1)
	assert.NoError(t, err)

	// no upsert for sqlite3
	/*
		person1.Name = "Jonee6"
		err = personsRepo.Upsert(ctx, person1)
		assert.NoError(t, err)
	*/

	person1.Name = "Jonee7"
	err = personsRepo.Update(ctx, goqu.Ex(map[string]interface{}{"id": person1.Id}), person1)
	assert.NoError(t, err)

	count := int64(0)
	err = personsRepo.CountAll(ctx, &count, nil)
	assert.NoError(t, err)
	logger.WithContext(ctx).Info(count)

	var destPerson Person
	err = personsRepo.FindOne(ctx, &destPerson, And(map[string]interface{}{"id": person1.Id}).Expression())
	assert.NoError(t, err)
	logger.WithContext(ctx).Info(destPerson)

	var destPersons []Person
	pagingResponse, err := personsRepo.FindAll(ctx, &destPersons,
		Where(goqu.Ex(map[string]interface{}{"name": person1.Name})),
		Keys(goqu.Ex(map[string]interface{}{"id": person1.Id})),
		// Distinct([]string{"name"}), // no distinct on for sqlite3
		Sort([]paging.SortOrder{paging.SortOrder{Property: "name", Direction: "ASC"}}),
		Paging(paging.Request{Size: 10, Page: 0}),
	)
	assert.NoError(t, err)
	logger.WithContext(ctx).Info(pagingResponse)
	logger.WithContext(ctx).Info(destPersons)

	err = personsRepo.DeleteOne(ctx, goqu.Ex(map[string]interface{}{"id": person1.Id}))
	assert.NoError(t, err)

	err = personsRepo.DeleteAll(ctx, goqu.Ex(map[string]interface{}{"name": person1.Name}))
	assert.NoError(t, err)

	// no truncate for sqlite3
	/*
		err = personsRepo.Truncate(ctx)
		assert.NoError(t, err)
	*/

	_, err = db.Exec("DROP TABLE "+table, nil)
	assert.NoError(t, err)
}

func TestTypedRepository_Insert(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mock.ExpectExec(`INSERT INTO "persons"`).WithArgs(uuid.MustParse(mockId), mockName).WillReturnResult(sqlmock.NewResult(1, 1))

	personsRepo := NewTypedRepository[Person](ctx, "persons")

	person1 := Person{Id: uuid.MustParse(mockId), Name: mockName}

	err = personsRepo.Insert(ctx, person1)
	assert.NoError(t, err)
}

func TestTypedRepository_Upsert(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mock.ExpectExec(`UPSERT INTO "persons"`).WithArgs(uuid.MustParse(mockId), mockName).WillReturnResult(sqlmock.NewResult(1, 1))

	personsRepo := NewTypedRepository[Person](ctx, "persons")

	person1 := Person{Id: uuid.MustParse(mockId), Name: mockName}

	err = personsRepo.Upsert(ctx, person1)
	assert.NoError(t, err)
}

func TestTypedRepository_Update(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mock.ExpectExec(`UPDATE "persons"`).WithArgs(uuid.MustParse(mockId), mockName, uuid.MustParse(mockId)).WillReturnResult(sqlmock.NewResult(1, 1))

	personsRepo := NewTypedRepository[Person](ctx, "persons")

	person1 := Person{Id: uuid.MustParse(mockId), Name: mockName}

	err = personsRepo.Update(ctx, goqu.Ex(map[string]interface{}{"id": person1.Id}), person1)
	assert.NoError(t, err)
}

func TestTypedRepository_CountAll(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mock.ExpectQuery(`SELECT COUNT(.+) FROM`).WillReturnRows(sqlmock.NewRows([]string{"COUNT"}).FromCSVString("5"))

	personsRepo := NewTypedRepository[Person](ctx, "persons")

	count := int64(0)
	err = personsRepo.CountAll(ctx, &count, nil)
	assert.NoError(t, err)
	logger.WithContext(ctx).Info(count)
}

func TestTypedRepository_FindOne(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	columns := []string{"id", "name"}
	mock.ExpectQuery(`SELECT (.+) FROM "persons" WHERE`).
		WithArgs(mockId).
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString(mockId + "," + mockName))

	personsRepo := NewTypedRepository[Person](ctx, "persons")

	var destPerson Person
	err = personsRepo.FindOne(ctx, &destPerson, And(map[string]interface{}{"id": uuid.MustParse(mockId)}).Expression())
	assert.NoError(t, err)
	logger.WithContext(ctx).Info(destPerson)
}

func TestTypedRepository_FindAll(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(uuid.MustParse(mockId), mockName)
	mock.ExpectQuery(`SELECT COUNT(.+) FROM`).WillReturnRows(sqlmock.NewRows([]string{"COUNT"}).FromCSVString("5"))
	mock.ExpectQuery(`SELECT (.+) FROM "persons"`).WillReturnRows(rows)

	personsRepo := NewTypedRepository[Person](ctx, "persons")

	var destPersons []Person
	pagingResponse, err := personsRepo.FindAll(ctx, &destPersons,
		Where(goqu.Ex(map[string]interface{}{"name": mockName})),
		Keys(goqu.Ex(map[string]interface{}{"id": uuid.MustParse(mockId)})),
		Distinct("name"),
		Sort([]paging.SortOrder{paging.SortOrder{Property: "name", Direction: "ASC"}}),
		Paging(paging.Request{Size: 10, Page: 0}),
	)
	assert.NoError(t, err)
	logger.WithContext(ctx).Info(pagingResponse)
	logger.WithContext(ctx).Info(destPersons)
}

func TestTypedRepository_DeleteOne(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mock.ExpectExec(`DELETE FROM "persons"`).WithArgs(uuid.MustParse(mockId)).WillReturnResult(sqlmock.NewResult(1, 1))

	personsRepo := NewTypedRepository[Person](ctx, "persons")

	err = personsRepo.DeleteOne(ctx, goqu.Ex(map[string]interface{}{"id": uuid.MustParse(mockId)}))
	assert.NoError(t, err)
}

func TestTypedRepository_DeleteAll(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mock.ExpectExec(`DELETE FROM "persons"`).WithArgs(mockName).WillReturnResult(sqlmock.NewResult(1, 1))

	personsRepo := NewTypedRepository[Person](ctx, "persons")

	err = personsRepo.DeleteAll(ctx, goqu.Ex(map[string]interface{}{"name": mockName}))
	assert.NoError(t, err)
}

func TestTypedRepository_Truncate(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mock.ExpectExec(`TRUNCATE "persons"`).WillReturnResult(sqlmock.NewResult(1, 1))

	personsRepo := NewTypedRepository[Person](ctx, "persons")

	err = personsRepo.Truncate(ctx)
	assert.NoError(t, err)
}
