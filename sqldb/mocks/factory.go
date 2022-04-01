// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package mocks

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/sqldb"
)

type CrudRepositoryFactory struct {
	Repository *CrudRepositoryApi
}

func (c CrudRepositoryFactory) NewCrudRepository(_ string) sqldb.CrudRepositoryApi {
	return c.Repository
}

func (c CrudRepositoryFactory) NewCrudPreparedRepository(_ string) sqldb.CrudRepositoryApi {
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
