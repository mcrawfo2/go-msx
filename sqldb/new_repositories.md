Example codes for the new Repositories

First create the persons table in order for the sample codes below to work. 

CREATE TABLE persons (
    id UUID PRIMARY KEY,
    name STRING
);

## Top Level (TypedRepository)

	type Person struct {
	    Id uuid.UUID `db:"id"`
	    Name string `db:"name"`
    }

    personsRepo := sqldb.NewTypedRepository[Person](ctx, "persons")

	person1 := Person{Id: uuid.MustParse("437f96b0-6722-11ed-9022-0242ac120009"), Name: "Jonee"}
	err = personsRepo.Insert(ctx, person1)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}

	person1.Name = "Jonee6"
	err = personsRepo.Upsert(ctx, person1)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}

	person1.Name = "Jonee7"
	err = personsRepo.Update(ctx, goqu.Ex(map[string]interface{}{"id": person1.Id}), person1)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}

	count := int64(0)
	err = personsRepo.CountAll(ctx, &count, types.Optional[sqldb.WhereOption]{})
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
	logger.WithContext(ctx).Info(count)

	var destPerson Person
	err = personsRepo.FindOne(ctx, &destPerson, types.OptionalOf(goqu.Ex(map[string]interface{}{"id": person1.Id}).Expression()))
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
	logger.WithContext(ctx).Info(destPerson)

	var destPersons []Person
	pagingResponse, err := personsRepo.FindAll(ctx, &destPersons,
		sqldb.Where(goqu.Ex(map[string]interface{}{"name": person1.Name})),
		sqldb.Keys(goqu.Ex(map[string]interface{}{"id": person1.Id})),
		sqldb.Distinct([]string{"name"}),
		sqldb.Sort([]paging.SortOrder{paging.SortOrder{Property: "name", Direction: "ASC"}}),
		sqldb.Paging(paging.Request{Size: 10, Page: 0}),
	)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
	logger.WithContext(ctx).Info(pagingResponse)
	logger.WithContext(ctx).Info(destPersons)

	err = personsRepo.DeleteOne(ctx, goqu.Ex(map[string]interface{}{"id": person1.Id}))
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}

	err = personsRepo.DeleteAll(ctx, goqu.Ex(map[string]interface{}{"name": person1.Name}))
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}

	err = personsRepo.Truncate(ctx)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}


## Middle Level (GoquRepository)

TypedRepository should already be complete and there is no need for middle or low level repositories. 
Do reach out if you still find a need to use GoquRepository or SqlRepository. 

	rgoqu := sqldb.NewGoquRepository(ctx)

	person1 := Person{Id: uuid.MustParse("437f96b0-6722-11ed-9022-0242ac120005"), Name: "Jonee"}

	dsInsert, err := rgoqu.Insert(ctx, "persons")
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
	err = rgoqu.ExecuteInsert(ctx, dsInsert.Rows(person1))
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}


	dsUpsert, err := rgoqu.Upsert(ctx, "persons")
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}

	person1.Name = "Jonee6"
	err = rgoqu.ExecuteUpsert(ctx, dsUpsert.Rows(person1))
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}


	dsUpdate, err := rgoqu.Update(ctx, "persons")
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}

	person1.Name = "Jonee7"
	err = rgoqu.ExecuteUpdate(ctx, dsUpdate.Where(goqu.Ex(map[string]interface{}{"id": person1.Id})).Set(person1))
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}


	var destPerson Person
	dsGet, err := rgoqu.Get(ctx, "persons")
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
	err = rgoqu.ExecuteGet(ctx, dsGet.Where(goqu.Ex(map[string]interface{}{"id": person1.Id})), &destPerson)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
	logger.WithContext(ctx).Info(destPerson)


	var destPersons []Person
	dsSelect, err := rgoqu.Select(ctx, "persons")
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
	err = rgoqu.ExecuteSelect(ctx, dsSelect.Where(goqu.Ex(map[string]interface{}{"name": person1.Name})), &destPersons)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
	logger.WithContext(ctx).Info(destPersons)


	dsDelete, err := rgoqu.Delete(ctx, "persons")
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
	err = rgoqu.ExecuteDelete(ctx, dsDelete.Where(goqu.Ex(map[string]interface{}{"id": person1.Id})))
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}


	dsTruncate, err := rgoqu.Truncate(ctx, "persons")
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
	err = rgoqu.ExecuteTruncate(ctx, dsTruncate)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}


## Low Level (SqlRepository)

Use of this is discouraged. 

	rsql := sqldb.NewSqlRepository(ctx)

	err = rsql.SqlExecute(ctx, "INSERT INTO persons VALUES ($1, $2)", []interface{}{uuid.MustParse("437f96b0-6722-11ed-9022-0242ac120002"), "Jonee"})
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}

	var dest1 []Person
	err = rsql.SqlSelect(ctx, "SELECT * FROM persons", nil, &dest1)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
	logger.WithContext(ctx).Info(dest1)

	var dest Person
	err = rsql.SqlGet(ctx, "SELECT * FROM persons WHERE id=$1", []interface{}{uuid.MustParse("437f96b0-6722-11ed-9022-0242ac120002")}, &dest)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
	logger.WithContext(ctx).Info(dest)