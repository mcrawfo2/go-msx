// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoquRepository(t *testing.T) {
	table := "persons2"
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

	rgoqu := NewGoquRepository(ctx)

	_, err = db.Exec("DROP TABLE IF EXISTS "+table, nil)
	assert.NoError(t, err)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS "+table+" (\nid VARCHAR(16) NOT NULL PRIMARY KEY,\nname TEXT\n);", nil)
	assert.NoError(t, err)

	person1 := Person{Id: uuid.MustParse("437f96b0-6722-11ed-9022-0242ac120005"), Name: "Jonee"}

	dsInsert, err := rgoqu.Insert(ctx, table)
	assert.NoError(t, err)

	err = rgoqu.ExecuteInsert(ctx, dsInsert.Rows(person1))
	assert.NoError(t, err)

	// no upsert for sqlite3
	/*
		dsUpsert, err := rgoqu.Upsert(ctx, table)
		assert.NoError(t, err)

		person1.Name = "Jonee6"
		err = rgoqu.ExecuteUpsert(ctx, dsUpsert.Rows(person1))
		assert.NoError(t, err)
	*/

	dsUpdate, err := rgoqu.Update(ctx, table)
	assert.NoError(t, err)

	person1.Name = "Jonee7"
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
