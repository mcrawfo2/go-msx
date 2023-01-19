# TypedRepository

**TypedRepository** should be the most used going forward and it should be able to cater to most repository needs. It is most encouraged to use this repository.

If developers feel there are anything missing or things that can be improved on (specially if you feel you still need to use any repo aside from TypedRepository) please reach out to the #go-msx team.


### Insert
	err = personsRepo.Insert(ctx, person1)

### Upsert
	err = personsRepo.Upsert(ctx, person1)

### Update
	err = personsRepo.Update(ctx, goqu.Ex(map[string]interface{}{"id": person1.Id}), person1)

### CountAll
	err = personsRepo.CountAll(ctx, &count, nil)

### FindOne
	err = personsRepo.FindOne(ctx, &destPerson, sqldb.And(map[string]interface{}{"id": person1.Id}).Expression())

### FindAll
	pagingResponse, err := personsRepo.FindAll(ctx, &destPersons,
		sqldb.Where(goqu.Ex(map[string]interface{}{"name": person1.Name})),
		sqldb.Keys(goqu.Ex(map[string]interface{}{"id": person1.Id})),
		sqldb.Distinct("name"),
		sqldb.Sort([]paging.SortOrder{paging.SortOrder{Property: "name", Direction: "ASC"}}),
		sqldb.Paging(paging.Request{Size: 10, Page: 0}),
	)

### DeleteOne
	err = personsRepo.DeleteOne(ctx, goqu.Ex(map[string]interface{}{"id": person1.Id}))

### DeleteAll
	err = personsRepo.DeleteAll(ctx, goqu.Ex(map[string]interface{}{"name": person1.Name}))

### Truncate
	err = personsRepo.Truncate(ctx)

<br />

## Complete Code Examples

First create the persons table in order for the sample codes below to work.

`
CREATE TABLE persons (
    id UUID PRIMARY KEY,
    name STRING
);
`

	type Person struct {
	    Id uuid.UUID `db:"id"`
	    Name string `db:"name"`
    }

	personsRepo, err := sqldb.NewTypedRepository[Person](ctx, "persons")
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}

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
	err = personsRepo.CountAll(ctx, &count, nil)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
	logger.WithContext(ctx).Info(count)

	var destPerson Person
	err = personsRepo.FindOne(ctx, &destPerson, sqldb.And(map[string]interface{}{"id": person1.Id}).Expression())
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
	logger.WithContext(ctx).Info(destPerson)

	var destPersons []Person
	pagingResponse, err := personsRepo.FindAll(ctx, &destPersons,
		sqldb.Where(goqu.Ex(map[string]interface{}{"name": person1.Name})),
		sqldb.Keys(goqu.Ex(map[string]interface{}{"id": person1.Id})),
		sqldb.Distinct("name"),
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