// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	personsRepo, _ := NewTypedRepository[Person](ctx, table)

	_, err = db.Exec("DROP TABLE IF EXISTS "+table, nil)
	assert.NoError(t, err)

	// sqlite does not have uuid data type
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS "+table+" (\nid VARCHAR(16) NOT NULL PRIMARY KEY,\nname TEXT\n);", nil)
	assert.NoError(t, err)

	person1 := Person{Id: uuid.MustParse("437f96b0-6722-11ed-9022-0242ac120009"), Name: "Jonee"}
	err = personsRepo.Insert(ctx, person1)
	assert.NoError(t, err)

	person1.Name = "Jonee6"
	err = personsRepo.Upsert(ctx, person1)
	assert.NoError(t, err)

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
		Where(And(map[string]interface{}{"name": person1.Name})),
		Keys(map[string]interface{}{"id": person1.Id}),
		//Distinct("name"), // no distinct on for sqlite3
		Sort([]paging.SortOrder{{Property: "name", Direction: "ASC"}}),
		Paging(paging.Request{Size: 10, Page: 0}),
	)
	assert.NoError(t, err)
	logger.WithContext(ctx).Info(pagingResponse)
	logger.WithContext(ctx).Info(destPersons)

	err = personsRepo.DeleteOne(ctx, goqu.Ex(map[string]interface{}{"id": person1.Id}))
	assert.NoError(t, err)

	err = personsRepo.DeleteAll(ctx, goqu.Ex(map[string]interface{}{"name": person1.Name}))
	assert.NoError(t, err)

	err = personsRepo.Truncate(ctx)
	assert.NoError(t, err)

	_, err = db.Exec("DROP TABLE "+table, nil)
	assert.NoError(t, err)
}

func TestTypedRepository_Insert(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mock.ExpectExec("INSERT INTO `persons`").
		WithArgs(uuid.MustParse(mockId), mockName).
		WillReturnResult(sqlmock.NewResult(1, 1))

	personsRepo, _ := NewTypedRepository[Person](ctx, "persons")

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

	mock.ExpectExec("REPLACE INTO `persons`").
		WithArgs(uuid.MustParse(mockId), mockName).
		WillReturnResult(sqlmock.NewResult(1, 1))

	personsRepo, _ := NewTypedRepository[Person](ctx, "persons")

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

	mock.ExpectExec("UPDATE `persons`").
		WithArgs(
			uuid.MustParse(mockId),
			mockName,
			uuid.MustParse(mockId)).
		WillReturnResult(
			sqlmock.NewResult(1, 1))

	personsRepo, _ := NewTypedRepository[Person](ctx, "persons")

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

	personsRepo, _ := NewTypedRepository[Person](ctx, "persons")

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
	mock.ExpectQuery("SELECT (.+) FROM `persons` WHERE").
		WithArgs(mockId).
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString(mockId + "," + mockName))

	personsRepo, _ := NewTypedRepository[Person](ctx, "persons")

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
	mock.ExpectQuery("SELECT \\* FROM `persons` WHERE \\(\\(`name` = \\?\\) AND \\(`id` = \\?\\)\\) ORDER BY `name` ASC LIMIT \\?").
		WillReturnRows(rows)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM \\(SELECT \\* FROM `persons` WHERE \\(\\(`name` = \\?\\) AND \\(`id` = \\?\\)\\)\\) AS `t1`").
		WillReturnRows(sqlmock.NewRows([]string{"COUNT"}).FromCSVString("5"))

	personsRepo, _ := NewTypedRepository[Person](ctx, "persons")

	var destPersons []Person
	pagingResponse, err := personsRepo.FindAll(ctx, &destPersons,
		Where(goqu.Ex(map[string]interface{}{"name": mockName})),
		Keys(goqu.Ex(map[string]interface{}{"id": uuid.MustParse(mockId)})),
		//Distinct("name"),
		Sort([]paging.SortOrder{{Property: "name", Direction: "ASC"}}),
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

	mock.ExpectExec("DELETE FROM `persons` WHERE").WithArgs(uuid.MustParse(mockId)).WillReturnResult(sqlmock.NewResult(1, 1))

	personsRepo, _ := NewTypedRepository[Person](ctx, "persons")

	err = personsRepo.DeleteOne(ctx, goqu.Ex(map[string]interface{}{"id": uuid.MustParse(mockId)}))
	assert.NoError(t, err)
}

func TestTypedRepository_DeleteAll(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mock.ExpectExec("DELETE FROM `persons` WHERE").
		WithArgs(mockName).
		WillReturnResult(
			sqlmock.NewResult(1, 1))

	personsRepo, _ := NewTypedRepository[Person](ctx, "persons")

	err = personsRepo.DeleteAll(ctx, goqu.Ex(map[string]interface{}{"name": mockName}))
	assert.NoError(t, err)
}

func TestTypedRepository_Truncate(t *testing.T) {
	ctx, mockDB, mock, err := newReposSqlMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mock.ExpectExec("DELETE FROM `persons`").WillReturnResult(sqlmock.NewResult(1, 1))

	personsRepo, _ := NewTypedRepository[Person](ctx, "persons")

	err = personsRepo.Truncate(ctx)
	assert.NoError(t, err)
}

func TestQueryGeneration(t *testing.T) {
	tests := []struct {
		name    string
		driver  string
		setup   func(mock *MockSqlRepositoryApi)
		call    func(ctx context.Context, repo TypedRepositoryApi[Person]) error
		wantErr bool
	}{
		{
			name: "FindByKey",
			setup: func(mockApi *MockSqlRepositoryApi) {
				mockApi.EXPECT().
					SqlGet(
						mock.Anything,
						"SELECT * FROM `person` WHERE (`id` = ?)",
						[]any{mockId},
						mock.Anything).
					Return(nil)
			},
			call: func(ctx context.Context, repo TypedRepositoryApi[Person]) error {
				var person Person
				return repo.FindOne(ctx, &person, And(map[string]any{
					columnId: mockId,
				}))
			},
			wantErr: false,
		},
		{
			name: "FindByAnd",
			setup: func(mockApi *MockSqlRepositoryApi) {
				mockApi.EXPECT().
					SqlGet(
						mock.Anything,
						"SELECT * FROM `person` WHERE ((`id` = ?) AND (`name` = ?))",
						[]any{mockId, mockName},
						mock.Anything).
					Return(nil)
			},
			call: func(ctx context.Context, repo TypedRepositoryApi[Person]) error {
				var person Person
				return repo.FindOne(ctx, &person, And(map[string]any{
					columnId:   mockId,
					columnName: mockName,
				}))
			},
			wantErr: false,
		},
		{
			name: "FindByOr",
			setup: func(mockApi *MockSqlRepositoryApi) {
				mockApi.EXPECT().
					SqlGet(
						mock.Anything,
						"SELECT * FROM `person` WHERE ((`id` = ?) OR (`name` = ?))",
						[]any{mockId, mockName},
						mock.Anything).
					Return(nil)
			},
			call: func(ctx context.Context, repo TypedRepositoryApi[Person]) error {
				var person Person
				return repo.FindOne(ctx, &person, Or(map[string]any{
					columnId:   mockId,
					columnName: mockName,
				}))
			},
			wantErr: false,
		},
		{
			name: "FindByAny",
			setup: func(mockApi *MockSqlRepositoryApi) {
				mockApi.EXPECT().
					SqlGet(
						mock.Anything,
						"SELECT * FROM `person` WHERE ((`id` = ?) OR (`name` = ?))",
						[]any{mockId, mockName},
						mock.Anything).
					Return(nil)
			},
			call: func(ctx context.Context, repo TypedRepositoryApi[Person]) error {
				var person Person
				return repo.FindOne(ctx, &person, Any(
					And(map[string]any{columnId: mockId}),
					And(map[string]any{columnName: mockName}),
				))
			},
			wantErr: false,
		},
		{
			name: "FindByAll",
			setup: func(mockApi *MockSqlRepositoryApi) {
				mockApi.EXPECT().
					SqlGet(
						mock.Anything,
						"SELECT * FROM `person` WHERE ((`id` = ?) AND (`name` = ?))",
						[]any{mockId, mockName},
						mock.Anything).
					Return(nil)
			},
			call: func(ctx context.Context, repo TypedRepositoryApi[Person]) error {
				var person Person
				return repo.FindOne(ctx, &person, All(
					And(map[string]any{columnId: mockId}),
					And(map[string]any{columnName: mockName}),
				))
			},
			wantErr: false,
		},
		{
			name: "FindAll",
			setup: func(mockApi *MockSqlRepositoryApi) {
				mockApi.EXPECT().
					SqlSelect(
						mock.Anything,
						"SELECT * FROM `person` WHERE ((`id` = ?) AND (`name` = ?)) LIMIT ? OFFSET ?",
						[]any{mockId, mockName, int64(10), int64(100)},
						mock.Anything).
					Return(nil)

				mockApi.EXPECT().
					SqlGet(
						mock.Anything,
						"SELECT COUNT(*) FROM (SELECT * FROM `person` WHERE ((`id` = ?) AND (`name` = ?))) AS `t1`",
						[]any{mockId, mockName},
						mock.Anything).
					Return(nil)

			},
			call: func(ctx context.Context, repo TypedRepositoryApi[Person]) error {
				var persons []Person
				_, err := repo.FindAll(ctx, &persons,
					Where(And(map[string]any{
						columnId:   mockId,
						columnName: mockName,
					})),
					Paging(paging.Request{
						Size: 10,
						Page: 10,
					}))
				return err
			},
			wantErr: false,
		},
		{
			name:   "FindAllOptions",
			driver: "postgres",
			setup: func(mockApi *MockSqlRepositoryApi) {
				mockApi.EXPECT().
					SqlSelect(
						mock.AnythingOfType("*context.valueCtx"),
						"SELECT DISTINCT ON (\"name\") * FROM \"person\" WHERE ((\"name\" = $1) AND (\"id\" = $2)) ORDER BY \"name\" ASC LIMIT $3 OFFSET $4",
						[]any{mockName, mockId, int64(2), int64(2)},
						mock.AnythingOfType("*[]sqldb.Person")).
					Return(nil)

				mockApi.EXPECT().
					SqlGet(
						mock.AnythingOfType("*context.valueCtx"),
						"SELECT COUNT(*) FROM (SELECT DISTINCT ON (\"name\") * FROM \"person\" WHERE ((\"name\" = $1) AND (\"id\" = $2))) AS \"t1\"",
						[]any{mockName, mockId},
						mock.Anything).
					Return(nil)
			},
			call: func(ctx context.Context, repo TypedRepositoryApi[Person]) error {
				var destPersons []Person
				_, err := repo.FindAll(ctx, &destPersons,
					Where(goqu.Ex(map[string]any{"name": mockName})),
					Keys(map[string]any{"id": uuid.MustParse(mockId)}),
					Distinct("name"),
					Sort([]paging.SortOrder{{Property: "name", Direction: "ASC"}}),
					Paging(paging.Request{Size: 2, Page: 1}),
				)
				return err
			},
		},
		{
			name: "Upsert",
			setup: func(mockApi *MockSqlRepositoryApi) {
				mockApi.EXPECT().
					SqlExecute(
						mock.Anything,
						"REPLACE INTO `person` (`id`, `name`) VALUES (?, ?)",
						[]any{mockId, mockName}).
					Return(nil)
			},
			call: func(ctx context.Context, repo TypedRepositoryApi[Person]) error {
				var person = Person{
					Id:   uuid.MustParse(mockId),
					Name: mockName,
				}
				return repo.Upsert(ctx, person)
			},
			wantErr: false,
		},
		{
			name:   "UpsertPg",
			driver: "postgres",
			setup: func(mockApi *MockSqlRepositoryApi) {
				mockApi.EXPECT().
					SqlExecute(
						mock.Anything,
						`UPSERT INTO "person" ("id", "name") VALUES ($1, $2)`,
						[]any{mockId, mockName}).
					Return(nil)
			},
			call: func(ctx context.Context, repo TypedRepositoryApi[Person]) error {
				var person = Person{
					Id:   uuid.MustParse(mockId),
					Name: mockName,
				}
				return repo.Upsert(ctx, person)
			},
			wantErr: false,
		},
		{
			name: "Truncate",
			setup: func(mockApi *MockSqlRepositoryApi) {
				mockApi.EXPECT().
					SqlExecute(
						mock.Anything,
						"DELETE FROM `person`",
						[]any{}).
					Return(nil)
			},
			call: func(ctx context.Context, repo TypedRepositoryApi[Person]) error {
				return repo.Truncate(ctx)
			},
			wantErr: false,
		},
		{
			name:   "TruncatePg",
			driver: "postgres",
			setup: func(mockApi *MockSqlRepositoryApi) {
				mockApi.EXPECT().
					SqlExecute(
						mock.Anything,
						`TRUNCATE "person"`,
						[]any{}).
					Return(nil)
			},
			call: func(ctx context.Context, repo TypedRepositoryApi[Person]) error {
				return repo.Truncate(ctx)
			},
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			driver := test.driver
			if driver == "" {
				driver = "sqlite3"
			}

			ctx, _, _ := newSqlMockDependencies(driver)

			mockSqlRepositoryApi := NewMockSqlRepositoryApi(t)
			ctx = ContextSqlRepository().Set(ctx, mockSqlRepositoryApi)

			repo, err := NewTypedRepository[Person](ctx, tableNamePerson)
			assert.NoError(t, err)

			if test.setup != nil {
				test.setup(mockSqlRepositoryApi)
			}

			// when
			err = test.call(ctx, repo)

			// then
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
