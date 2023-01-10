// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/fs"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"modernc.org/sqlite"
	"testing"
)

type Person struct {
	Id   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

const (
	mockName = "Jonee"
	mockId   = "437f96b0-6722-11ed-9022-0242ac120002"
)

func sqlite3DBContext() context.Context {
	// configure the file system
	cfg := configtest.NewInMemoryConfig(map[string]string{
		"spring.datasource.driver":           "sqlite3",
		"spring.datasource.enabled":          "true",
		"spring.datasource.name":             "TestUnitTest.db",
		"spring.datasource.data-source-name": "file:${spring.datasource.name}?cache=shared&mode=memory",
	})
	fs.ConfigureFileSystem(cfg)

	mockedCtx := context.Background()
	mockedCtx = config.ContextWithConfig(mockedCtx, cfg)

	return mockedCtx
}

func TestSuite_SqlRepository(t *testing.T) {
	table := "persons1"
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

	rsql := NewSqlRepository(ctx)

	_, err = db.Exec("DROP TABLE IF EXISTS "+table, nil)
	assert.NoError(t, err)

	// sqlite does not have uuid data type
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS "+table+" (\nid VARCHAR(16) NOT NULL PRIMARY KEY,\nname TEXT\n);", nil)
	assert.NoError(t, err)

	err = rsql.SqlExecute(ctx, "INSERT INTO "+table+" VALUES ($1, $2)", []interface{}{uuid.MustParse("437f96b0-6722-11ed-9022-0242ac120002"), "Jonee"})
	assert.NoError(t, err)

	var destPersons []Person
	err = rsql.SqlSelect(ctx, "SELECT * FROM "+table, nil, &destPersons)
	assert.NoError(t, err)
	logger.WithContext(ctx).Info(destPersons)

	var destPerson Person
	err = rsql.SqlGet(ctx, "SELECT * FROM "+table+" WHERE id=$1", []interface{}{uuid.MustParse("437f96b0-6722-11ed-9022-0242ac120002")}, &destPerson)
	assert.NoError(t, err)
	logger.WithContext(ctx).Info(destPerson)

	_, err = db.Exec("DROP TABLE "+table, nil)
	assert.NoError(t, err)
}

func newReposSqlMock() (context.Context, *sql.DB, sqlmock.Sqlmock, error) {
	ctx := sqlite3DBContext()

	mockDB, mock, err := sqlmock.New()

	mockSqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	ctx = ContextSqlExecutor().Set(ctx, mockSqlxDB)

	return ctx, mockDB, mock, err
}

func TestSqlRepository_SqlExecute(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mock.ExpectExec("INSERT INTO persons").WithArgs(uuid.MustParse(mockId), mockName).WillReturnResult(sqlmock.NewResult(1, 1))

	rsql := NewSqlRepository(ctx)
	err = rsql.SqlExecute(ctx, "INSERT INTO persons VALUES ($1, $2)", []interface{}{uuid.MustParse(mockId), mockName})
	assert.NoError(t, err)
}

func TestSqlRepository_SqlSelect(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(uuid.MustParse(mockId), mockName)
	mock.ExpectQuery("^SELECT (.+) FROM persons$").WillReturnRows(rows)

	rsql := NewSqlRepository(ctx)
	var destPersons []Person
	err = rsql.SqlSelect(ctx, "SELECT * FROM persons", nil, &destPersons)
	assert.NoError(t, err)
	logger.WithContext(ctx).Info(destPersons)
}

func TestSqlRepository_SqlGet(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	columns := []string{"id", "name"}
	mock.ExpectQuery("^SELECT (.+) FROM persons WHERE id=").
		WithArgs(mockId).
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString(mockId + "," + mockName))

	rsql := NewSqlRepository(ctx)
	var destPerson Person
	err = rsql.SqlGet(ctx, "SELECT * FROM persons WHERE id=$1", []interface{}{uuid.MustParse(mockId)}, &destPerson)
	assert.NoError(t, err)
	logger.WithContext(ctx).Info(destPerson)
}
