// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

type CrudRepositoryFactoryApi interface {
	// NewCrudRepository Deprecated
	NewCrudRepository(tableName string) CrudRepositoryApi
	NewCrudPreparedRepository(tableName string) CrudRepositoryApi
}

type ProductionCrudRepositoryFactory struct{}

// NewCrudRepository Deprecated
func (f *ProductionCrudRepositoryFactory) NewCrudRepository(tableName string) CrudRepositoryApi {
	return newCrudRepository(tableName)
}

func (f *ProductionCrudRepositoryFactory) NewCrudPreparedRepository(tableName string) CrudRepositoryApi {
	return newCrudPreparedRepository(tableName)
}

func NewProductionCrudRepositoryFactory() CrudRepositoryFactoryApi {
	return new(ProductionCrudRepositoryFactory)
}
