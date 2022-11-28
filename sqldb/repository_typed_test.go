// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypedRepository(t *testing.T) {
	table := "persons3"
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

	personsRepo := NewTypedRepository[Person](ctx, table)

	_, err = db.Exec("DROP TABLE IF EXISTS "+table, nil)
	assert.NoError(t, err)

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
	err = personsRepo.CountAll(ctx, &count, types.Optional[WhereOption]{})
	assert.NoError(t, err)
	logger.WithContext(ctx).Info(count)

	var destPerson Person
	err = personsRepo.FindOne(ctx, &destPerson, types.OptionalOf(goqu.Ex(map[string]interface{}{"id": person1.Id}).Expression()))
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
