package config

import (
	"context"
	"sync"
)

type OnceLoader struct {
	once     sync.Once
	provider Provider
	settings map[string]string
}

func (l *OnceLoader) Cache() map[string]string {
	return l.settings
}

func (l *OnceLoader) Description() string {
	return l.provider.Description()
}

func (l *OnceLoader) Load(ctx context.Context) (map[string]string, error) {
	var err error

	l.once.Do(
		func() {
			l.settings, err = l.provider.Load(ctx)
		},
	)

	if err != nil {
		return nil, err
	}

	settings := map[string]string{}

	for key, value := range l.settings {
		settings[key] = value
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return settings, err
}

func NewOnceLoader(provider Provider) *OnceLoader {
	return &OnceLoader{
		once:     sync.Once{},
		provider: provider,
		settings: map[string]string{},
	}
}
