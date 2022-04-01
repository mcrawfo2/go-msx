// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package vault

import (
	"context"
)

type vaultContextKey int

const (
	contextKeyVaultPool vaultContextKey = iota
	contextKeyVaultConnectionApi
)

// Deprecated.  Use ContextWithConnection instead.
func ContextWithPool(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyVaultPool, pool)
}

// Deprecated.  Use ConnectionFromContext instead.
func PoolFromContext(ctx context.Context) *ConnectionPool {
	if raw := ctx.Value(contextKeyVaultPool); raw != nil {
		if typed, ok := raw.(*ConnectionPool); ok {
			return typed
		}

		logger.Warn("Context vault connection pool value is the wrong type")
	}
	return nil
}

func ContextWithConnection(ctx context.Context, api ConnectionApi) context.Context {
	return context.WithValue(ctx, contextKeyVaultConnectionApi, api)
}

func ConnectionFromContext(ctx context.Context) ConnectionApi {
	if raw := ctx.Value(contextKeyVaultConnectionApi); raw != nil {
		if typed, ok := raw.(ConnectionApi); ok {
			return typed
		}

		logger.Warn("Context vault connection value is the wrong type")
	}

	return nil
}
