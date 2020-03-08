package mocks

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra/ddl"
)

type CrudRepositoryFactory struct {
	Repository *CrudRepositoryApi
}

func (c CrudRepositoryFactory) NewCrudRepository(table ddl.Table) cassandra.CrudRepositoryApi {
	return c.Repository
}

func NewMockCrudRepositoryFactory() cassandra.CrudRepositoryFactoryApi {
	return &CrudRepositoryFactory{
		Repository: &CrudRepositoryApi{},
	}
}

func ContextWithMockCrudRepositoryFactory(ctx context.Context) context.Context {
	factory := NewMockCrudRepositoryFactory()
	return cassandra.ContextWithCrudRepositoryFactory(ctx, factory)
}
