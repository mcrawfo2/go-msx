# SqlRepository

**SqlRepository** is offered for those who may need to go down to the SQL level. 
This may be hard to maintain and inflexible (might be error-prone too) and thus discouraged as much as possible.

If developers feel there are anything missing or things that can be improved on (specially if you feel you still need to use any repo aside from TypedRepository) please reach out to the #go-msx team.


### SqlExecute
	err = rsql.SqlExecute(ctx, "INSERT INTO persons VALUES ($1, $2)", []interface{}{uuid.MustParse("437f96b0-6722-11ed-9022-0242ac120002"), "Jonee"})

### SqlSelect
	var destPersons3 []Person
	err = rsql.SqlSelect(ctx, "SELECT * FROM persons", nil, &destPersons3)

### SqlGet
	var destPerson3 Person
	err = rsql.SqlGet(ctx, "SELECT * FROM persons WHERE id=$1", []interface{}{uuid.MustParse("437f96b0-6722-11ed-9022-0242ac120002")}, &destPerson3)

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

	rsql, err := sqldb.NewSqlRepository(ctx)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}

	err = rsql.SqlExecute(ctx, "INSERT INTO persons VALUES ($1, $2)", []interface{}{uuid.MustParse("437f96b0-6722-11ed-9022-0242ac120002"), "Jonee"})
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}

	var destPersons3 []Person
	err = rsql.SqlSelect(ctx, "SELECT * FROM persons", nil, &destPersons3)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
	logger.WithContext(ctx).Info(destPersons3)

	var destPerson3 Person
	err = rsql.SqlGet(ctx, "SELECT * FROM persons WHERE id=$1", []interface{}{uuid.MustParse("437f96b0-6722-11ed-9022-0242ac120002")}, &destPerson3)
	if err != nil {
		logger.WithContext(ctx).Error(err)
	}
	logger.WithContext(ctx).Info(destPerson3)