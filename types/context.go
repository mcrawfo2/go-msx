// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package types

import "context"

type ContextKeyAccessor[T any] struct {
	key interface{}
}

func (a ContextKeyAccessor[T]) Get(ctx context.Context) T {
	value, ok := ctx.Value(a.key).(T)
	if !ok {
		var def T
		return def
	}
	return value
}

func (a ContextKeyAccessor[T]) Set(ctx context.Context, value T) context.Context {
	return context.WithValue(ctx, a.key, value)
}

func NewContextKeyAccessor[T any](key interface{}) ContextKeyAccessor[T] {
	return ContextKeyAccessor[T]{
		key: key,
	}
}

type ContextKeyGetter[T any] struct {
	accessor ContextKeyAccessor[T]
}

func (a ContextKeyGetter[T]) Get(ctx context.Context) T {
	return a.accessor.Get(ctx)
}

func NewContextKeyGetter[T any](key interface{}) ContextKeyGetter[T] {
	return ContextKeyGetter[T]{
		accessor: NewContextKeyAccessor[T](key),
	}
}

type ContextKeySetter[T any] struct {
	accessor ContextKeyAccessor[T]
}

func (a ContextKeySetter[T]) Set(ctx context.Context, value T) context.Context {
	return a.accessor.Set(ctx, value)
}

func NewContextKeySetter[T any](key interface{}) ContextKeySetter[T] {
	return ContextKeySetter[T]{
		accessor: NewContextKeyAccessor[T](key),
	}
}
