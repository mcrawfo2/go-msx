package config

import (
	"context"
)

type Describer interface {
	Description() string
}

type Named struct {
	name string
}

func (n Named) Description() string {
	return n.name
}

func NewNamed(name string) Named {
	return Named {
		name: name,
	}
}

type Loader interface {
	Load(ctx context.Context) (ProviderEntries, error)
}

type Notifier interface {
	Run(ctx context.Context)
	Notify() <-chan struct{}
}

type SilentNotifier struct{}

func (p SilentNotifier) Run(_ context.Context)   {}
func (p SilentNotifier) Notify() <-chan struct{} { return nil }

type Provider interface {
	Describer
	Loader
	Notifier
}
