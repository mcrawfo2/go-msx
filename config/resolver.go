package config

import (
	"context"
)

type Resolver struct {
	provider Provider
	mappings map[string]string
}

func (f *Resolver) Description() string {
	return f.provider.Description()
}

func (f *Resolver) Load(ctx context.Context) (map[string]string, error) {
	settings, err := f.provider.Load(ctx)
	if err != nil {
		return nil, err
	}

	resolved := map[string]string{}
	for key, val := range settings {
		if dest, ok := f.mappings[key]; ok {
			resolved[dest] = val
		} else {
			resolved[key] = val
		}
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return resolved, nil
}

func NewResolver(provider Provider, mappings map[string]string) *Resolver {
	return &Resolver{
		provider: provider,
		mappings: mappings,
	}
}
