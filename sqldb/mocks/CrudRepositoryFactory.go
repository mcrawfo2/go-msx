package mocks

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/sqldb"
)

type CrudRepositoryFactory struct {
	Repository *CrudRepositoryApi
}

func (c CrudRepositoryFactory) NewCrudRepository(tableName string) sqldb.CrudRepositoryApi {
	return c.Repository
}

func NewMockCrudRepositoryFactory() sqldb.CrudRepositoryFactoryApi {
	return &CrudRepositoryFactory{
		Repository: &CrudRepositoryApi{},
	}
}

func ContextWithMockCrudRepositoryFactory(ctx context.Context) context.Context {
	factory := NewMockCrudRepositoryFactory()
	return sqldb.ContextWithCrudRepositoryFactory(ctx, factory)
}
