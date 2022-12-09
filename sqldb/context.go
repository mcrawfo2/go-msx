// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
)

type contextKey string

const (
	contextKeySqlPool               = contextKey("SqlPool")
	contextKeyCrudRepositoryFactory = contextKey("CrudRepositoryFactory")
	contextKeySqlExecutor           = contextKey("SqlExecutor")
	contextKeyTransactionManager    = contextKey("TransactionManager")
	contextKeyRepositorySql         = contextKey("RepositorySql")
	contextKeyRepositoryGoqu        = contextKey("RepositoryGoqu")
	contextKeyRepositoryTypedPrefix = "RepositoryTyped"
)

var ErrDisabled = errors.New("Sql connection disabled")

func ContextWithPool(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeySqlPool, pool)
}

func PoolFromContext(ctx context.Context) (*ConnectionPool, error) {
	connectionPoolInterface := ctx.Value(contextKeySqlPool)
	if connectionPoolInterface == nil {
		return nil, ErrDisabled
	}
	return connectionPoolInterface.(*ConnectionPool), nil
}

func ContextWithCrudRepositoryFactory(ctx context.Context, factory CrudRepositoryFactoryApi) context.Context {
	return context.WithValue(ctx, contextKeyCrudRepositoryFactory, factory)
}

func CrudRepositoryFactoryFromContext(ctx context.Context) CrudRepositoryFactoryApi {
	api, _ := ctx.Value(contextKeyCrudRepositoryFactory).(CrudRepositoryFactoryApi)
	return api
}

func ContextSqlRepository() types.ContextKeyAccessor[SqlRepositoryApi] {
	return types.NewContextKeyAccessor[SqlRepositoryApi](contextKeyRepositorySql)
}

func ContextGoquRepository() types.ContextKeyAccessor[GoquRepositoryApi] {
	return types.NewContextKeyAccessor[GoquRepositoryApi](contextKeyRepositoryGoqu)
}

func ContextTypedRepository[I any](table string) types.ContextKeyAccessor[TypedRepositoryApi[I]] {
	contextKey := contextKey(contextKeyRepositoryTypedPrefix + table)
	return types.NewContextKeyAccessor[TypedRepositoryApi[I]](contextKey)
}
