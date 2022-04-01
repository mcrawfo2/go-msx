// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package types

import "context"

type ContextInjector func(ctx context.Context) context.Context
type ContextInjectors []ContextInjector

func (i *ContextInjectors) Register(injector ContextInjector) {
	*i = append(*i, injector)
}

func (i ContextInjectors) Inject(ctx context.Context) context.Context {
	for _, contextInjector := range i {
		ctx = contextInjector(ctx)
	}
	return ctx
}

func (i ContextInjectors) Clone() *ContextInjectors {
	clone := i.Slice()
	return &clone
}

func (i ContextInjectors) Slice() ContextInjectors {
	clone := ContextInjectors(append([]ContextInjector(nil), i...))
	return clone
}
