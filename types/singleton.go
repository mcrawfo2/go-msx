// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package types

import (
	"context"
	"sync"
)

type Singleton[VT any] struct {
	value     *VT
	mtx       sync.Mutex
	construct func(context.Context) (*VT, error)
}

func (s *Singleton[VT]) Factory(ctx context.Context) (*VT, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.value != nil {
		return s.value, nil
	}

	value, err := s.construct(ctx)
	if err != nil {
		return nil, err
	}

	s.value = value

	return value, nil
}

func NewSingleton[VT any](constructor func(context.Context) (*VT, error)) *Singleton[VT] {
	return &Singleton[VT]{
		value:     nil,
		mtx:       sync.Mutex{},
		construct: constructor,
	}
}
