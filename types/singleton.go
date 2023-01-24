// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package types

import (
	"context"
	"sync"
)

type Constructor[I any] func(context.Context) (I, error)
type Accessor[I any] func() ContextKeyAccessor[I]

type Singleton[I any] struct {
	mtx         sync.Mutex
	value       Optional[I]
	constructor Constructor[I]
	accessor    Accessor[I]
}

func (s *Singleton[I]) Factory(ctx context.Context) (value I, err error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.value.IsPresent() {
		value = s.value.Value()
		return
	}

	value, ok := s.accessor().TryGet(ctx)
	if ok {
		return
	}

	value, err = s.constructor(ctx)
	if err != nil {
		return
	}

	s.value = OptionalOf(value)

	return
}

func NewSingleton[I any](constructor Constructor[I], accessor Accessor[I]) *Singleton[I] {
	return &Singleton[I]{
		constructor: constructor,
		accessor:    accessor,
	}
}
