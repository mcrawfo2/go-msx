package types

import "context"

type ContextInjector func(ctx context.Context) context.Context
type ContextInjectors []ContextInjector

func (i *ContextInjectors) Register(injector ContextInjector) {
	*i = append(*i, injector)
}

func (i *ContextInjectors) Inject(ctx context.Context) context.Context {
	for _, contextInjector := range *i {
		ctx = contextInjector(ctx)
	}
	return ctx
}

func (i *ContextInjectors) Clone() *ContextInjectors {
	slice := *i
	clone := slice[:]
	return &clone
}
