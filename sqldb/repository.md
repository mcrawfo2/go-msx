# SQL Database Repository

MSX promotes the usage of the common Controller > Service > Repository layered architecture within microservices.

The role of the Repository is to query and mutate the persistent storage of Models.

## Defining the Repository

To define a repository, create a standard Go structure with an anonymous field for the CrudRepository:

```go
type deviceSqlRepository struct {
    sqldb.CrudRepositoryApi
}
```

The `CrudRepositoryApi` object provides access to the SQL database underneath using generic objects 
and slices.

## Writing a Constructor

A standard repository constructor allows for dependency injection (during testing) and normal creation (during runtme):

```go
func newDeviceRepository(ctx context.Context) deviceRepositoryApi {
	repo := deviceRepositoryFromContext(ctx)
	if repo == nil {
		repo = &deviceSqlRepository{
			CrudRepositoryApi: sqldb.
				CrudRepositoryFactoryFromContext(ctx).
				NewCrudRepository("device"),
		}
	}
	return repo
}
```

- The `CrudRepositoryFactory` allows us to test the repository without requiring an actual database implementation.
- The `deviceRepositoryFromContext` allows us to test this repositories reverse-dependencies without requiring
an actual `deviceSqlRepository`.

## Implementing Common Access Methods

A basic repository will likely have the following common methods:
- `FindAll`
  - Retrieve all models
- `FindByKey`
  - Retrieve a single model by its primary key
- `Save`
  - Store a single model
- `Delete`
  - Remove a single model

More advanced repositories may have some less-common methods:
- `FindAllByIndexXXX`
  - Retrieve all models matching the specified criteria using an index
- `FindAllPagedBy`
  - Retrieve a subset of models matching the specified criteria, using the specified sorting and pagination
- `Truncate`
  - Remove all models

### FindAll

```go
func (r *deviceSqlRepository) FindAll(ctx context.Context) (results []device, err error) {
	logger.WithContext(ctx).Info("Retrieving all Device records")
	err = r.CrudRepositoryApi.FindAll(ctx, &results)
	return
}
```

1. Log method intention
2. Delegate to our internal CrudRepository to perform the record retrieval and struct mapping.

### FindByKey

```go
func (r *deviceSqlRepository) FindByKey(ctx context.Context, name string) (result *device, err error) {
	logger.WithContext(ctx).Infof("Retrieving Device by key %q", name)
	var res device
	err = r.CrudRepositoryApi.FindOneBy(ctx, map[string]interface{}{
		"name": name,
	}, &res)
	if err == sqldb.ErrNotFound {
		err = repository.ErrNotFound
	} else if err == nil {
		result = &res
	}
	return
}
```

1. Log method intention, including the primary key
2. Delegate to our internal CrudRepository to perform the record retrieval and struct mapping.
   - The `CrudRepositoryApi.FindOneBy` method accepts a map of criteria to search by -- 
     in this case, the primary key.
3. Normalize the `sqldb` error code to use `repository` error codes.

### Save

```go
func (r *deviceSqlRepository) Save(ctx context.Context, device device) (err error) {
	logger.WithContext(ctx).Infof("Storing Device with key %q", device.Name)
	return r.CrudRepositoryApi.Save(ctx, device)
}
```

1. Log method intention, including the primary key
2. Delegate to our internal CrudRepository to perform the record storage and struct mapping.
   - The `CrudRepositoryApi.Save` method performs an `UPSERT` query in Cockroach so it behaves
     in much the same way as the `Save` method from a KV repository.

### Delete

```go
func (r *deviceSqlRepository) Delete(ctx context.Context, name string) (err error) {
	logger.WithContext(ctx).Infof("Deleting Device by key %q", name)
	return r.CrudRepositoryApi.DeleteBy(ctx, map[string]interface{}{
		columnDeviceName: name,
	})
}
```

1. Log method intention, including the primary key
2. Delegate to our internal CrudRepository to perform the record retrieval and struct mapping.
   - The `CrudRepositoryApi.DeleteBy` method accepts a map of criteria to delete by -- 
     in this case, the primary key.
