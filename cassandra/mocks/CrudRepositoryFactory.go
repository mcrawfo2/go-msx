// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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
