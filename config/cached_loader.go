package config

import "context"

type CachedLoader struct {
	valid    bool
	provider Provider
	settings map[string]string
}

func (f *CachedLoader) Cache() map[string]string {
	return f.settings
}

func (f *CachedLoader) Description() string {
	if f.provider == nil {
		return ""
	}
	return f.provider.Description()
}

func (f *CachedLoader) Load(ctx context.Context) (map[string]string, error) {
	if f.provider == nil {
		return nil, nil
	}

	if !f.valid {
		if err := f.reload(ctx); err != nil {
			return nil, err
		}
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return f.clone(), nil
}

func (f *CachedLoader) reload(ctx context.Context) error {
	if settings, err := f.provider.Load(ctx); err != nil {
		return err
	} else {
		f.settings = settings
		f.valid = true
		return nil
	}
}

func (f *CachedLoader) clone() map[string]string {
	settings := map[string]string{}

	for key, value := range f.settings {
		settings[key] = value
	}

	return settings
}

func (f *CachedLoader) Invalidate() {
	f.valid = false
}

func NewCachedLoader(provider Provider) *CachedLoader {
	return &CachedLoader{
		provider: provider,
		settings: map[string]string{},
	}
}
