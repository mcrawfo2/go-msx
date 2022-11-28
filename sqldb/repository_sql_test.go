// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/fs"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Person struct {
	Id   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

func testMock() context.Context {
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

func TestSqlRepository(t *testing.T) {
	table := "persons1"
	ctx := testMock()

	err := ConfigurePool(ctx)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
	ctx = ContextWithPool(ctx)

	myPool, err := PoolFromContext(ctx)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}

	drivers["sqlite3"] = &sqlite3.SQLiteDriver{}

	db, err := myPool.NewSqlConnection()
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}

	rsql := NewSqlRepository(ctx)

	_, err = db.Exec("DROP TABLE IF EXISTS "+table, nil)
	assert.NoError(t, err)

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
