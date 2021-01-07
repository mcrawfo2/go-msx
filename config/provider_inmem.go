package config

import (
	"context"
	"fmt"
	"sync"
)

type ActionFunc func()

type InMemoryProvider struct {
	Describer
	settings map[string]string
	settingsMtx sync.Mutex
	work     chan ActionFunc
	notify   chan struct{}
}

func (p *InMemoryProvider) Load(ctx context.Context) (ProviderEntries, error) {
	p.settingsMtx.Lock()
	defer p.settingsMtx.Unlock()

	var results = make(ProviderEntries, 0, len(p.settings))
	for name, value := range p.settings {
		results = append(results, NewEntry(p, name, value))
	}
	return results, nil
}

func (p *InMemoryProvider) Notify() <-chan struct{} {
	return p.notify
}

func (p *InMemoryProvider) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Infof("In-Memory provider context cancelled.  Exiting work queue.")
			return

		case fn, ok := <-p.work:
			if ok {
				fn()
			}
		}
	}
}

func (p *InMemoryProvider) SetValue(name, value string) error {
	if name == "" {
		return ErrEmptyKey
	}

	p.work <- func() {

		logger.Infof("Setting key %q to %q in %q", name, value, p.Description())

		p.settingsMtx.Lock()
		p.settings[name] = value
		p.settingsMtx.Unlock()

		p.notify <- struct{}{}
	}

	return nil
}

func (p *InMemoryProvider) UnsetValue(name string) error {
	if name == "" {
		return ErrEmptyKey
	}

	p.work <- func() {
		logger.Infof("Unsetting key %q in %q", name, p.Description())

		p.settingsMtx.Lock()
		delete(p.settings, name)
		p.settingsMtx.Unlock()

		p.notify <- struct{}{}
	}

	return nil
}

func NewInMemoryProvider(name string, settings map[string]string) *InMemoryProvider {
	return &InMemoryProvider{
		Describer: Named{
			name: fmt.Sprintf("InMemory: %s", name),
		},
		settings: settings,
		work:     make(chan ActionFunc),
		notify:   make(chan struct{}),
	}
}
