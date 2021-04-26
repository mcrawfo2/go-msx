package sqldb

type CrudRepositoryFactoryApi interface {
	// Deprecated
	NewCrudRepository(tableName string) CrudRepositoryApi
	NewCrudPreparedRepository(tableName string) CrudRepositoryApi
}

type ProductionCrudRepositoryFactory struct{}

// Deprecated
func (f *ProductionCrudRepositoryFactory) NewCrudRepository(tableName string) CrudRepositoryApi {
	return newCrudRepository(tableName)
}

func (f *ProductionCrudRepositoryFactory) NewCrudPreparedRepository(tableName string) CrudRepositoryApi {
	return newCrudPreparedRepository(tableName)
}

func NewProductionCrudRepositoryFactory() CrudRepositoryFactoryApi {
	return new(ProductionCrudRepositoryFactory)
}
