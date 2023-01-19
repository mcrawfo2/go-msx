# GoquRepository

**GoquRepository** is for those who may want more flexibility (that otherwise TypedRepository cannot provide) and would like to work on the Goqu level.

If developers feel there are anything missing or things that can be improved on (specially if you feel you still need to use any repo aside from TypedRepository) please reach out to the #go-msx team.


### Insert
	dsInsert := rgoqu.Insert("persons")

	err = rgoqu.ExecuteInsert(ctx, dsInsert.Rows(person2))

### Upsert
	dsUpsert := rgoqu.Upsert("persons")

	err = rgoqu.ExecuteUpsert(ctx, dsUpsert.Rows(person2))

### Update
	dsUpdate := rgoqu.Update("persons")

	err = rgoqu.ExecuteUpdate(ctx, dsUpdate.Where(goqu.Ex(map[string]interface{}{"id": person2.Id})).Set(person2))

### Get
	dsGet := rgoqu.Get("persons")

	err = rgoqu.ExecuteGet(ctx, dsGet.Where(goqu.Ex(map[string]interface{}{"id": person2.Id})), &destPerson2)

### Select
	dsSelect := rgoqu.Select("persons")

	err = rgoqu.ExecuteSelect(ctx, dsSelect.Where(goqu.Ex(map[string]interface{}{"name": person2.Name})), &destPersons2)

### Delete
	dsDelete := rgoqu.Delete("persons")

	err = rgoqu.ExecuteDelete(ctx, dsDelete.Where(goqu.Ex(map[string]interface{}{"id": person2.Id})))

### Truncate
	dsTruncate := rgoqu.Truncate("persons")

	err = rgoqu.ExecuteTruncate(ctx, dsTruncate)

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

	rgoqu, err := sqldb.NewGoquRepository(ctx)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}

	person2 := Person{Id: uuid.MustParse("437f96b0-6722-11ed-9022-0242ac120005"), Name: "Jonee"}

	dsInsert := rgoqu.Insert("persons")

	err = rgoqu.ExecuteInsert(ctx, dsInsert.Rows(person2))
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}


	dsUpsert := rgoqu.Upsert("persons")

	person2.Name = "Jonee6"
	err = rgoqu.ExecuteUpsert(ctx, dsUpsert.Rows(person2))
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}


	dsUpdate := rgoqu.Update("persons")

	person2.Name = "Jonee7"
	err = rgoqu.ExecuteUpdate(ctx, dsUpdate.Where(goqu.Ex(map[string]interface{}{"id": person2.Id})).Set(person2))
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}


	var destPerson2 Person
	dsGet := rgoqu.Get("persons")

	err = rgoqu.ExecuteGet(ctx, dsGet.Where(goqu.Ex(map[string]interface{}{"id": person2.Id})), &destPerson2)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
	logger.WithContext(ctx).Info(destPerson2)


	var destPersons2 []Person
	dsSelect := rgoqu.Select("persons")

	err = rgoqu.ExecuteSelect(ctx, dsSelect.Where(goqu.Ex(map[string]interface{}{"name": person2.Name})), &destPersons2)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
	logger.WithContext(ctx).Info(destPersons2)


	dsDelete := rgoqu.Delete("persons")

	err = rgoqu.ExecuteDelete(ctx, dsDelete.Where(goqu.Ex(map[string]interface{}{"id": person2.Id})))
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}


	dsTruncate := rgoqu.Truncate("persons")

	err = rgoqu.ExecuteTruncate(ctx, dsTruncate)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
