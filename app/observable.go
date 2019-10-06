package app

import (
	"context"
	"sync"
)

type Observer func()

type observable struct {
	Callbacks map[string][]Observer
	cancel    func()
	ctx       context.Context
	sync.Mutex
}

func (o *observable) isCancelled() bool {
	return o.ctx.Err() != nil
}

func (o *observable) On(event string, phase string, observer Observer) {
	o.Lock()
	defer o.Unlock()

	key := event + phase
	if _, ok := o.Callbacks[event+phase]; !ok {
		o.Callbacks[key] = []Observer{observer}
	} else {
		o.Callbacks[key] = append(o.Callbacks[key], observer)
	}
}

func (o *observable) Clear(event string, phase string) {
	o.Lock()
	defer o.Unlock()
	delete(o.Callbacks, event+phase)
}

func (o *observable) isIgnored(event, phase string) bool {
	return o.isCancelled() && event != EventStop && event != EventFinal
}

func (o *observable) callbacks(event, phase string) ([]Observer, bool) {
	o.Lock()
	defer o.Unlock()
	observers, ok := o.Callbacks[event+phase]
	return observers, ok
}

func (o *observable) trigger(event, phase string) {
	logger.Infof("Event triggered: %s%s", event, phase)
	if observers, ok := o.callbacks(event, phase); ok {
		for _, observer := range observers {
			if o.isIgnored(event, phase) {
				break
			}
			observer()
		}
	}
}

func (o *observable) Trigger(event string) {
	for _, phase := range []string{PhaseBefore, PhaseDuring, PhaseAfter} {
		if o.isIgnored(event, phase) {
			break
		}
		o.trigger(event, phase)
	}
}

func (o *observable) Shutdown() {
	if o.isCancelled() {
		return
	}

	logger.Info("Shutdown requested")

	o.Lock()
	defer o.Unlock()
	o.cancel()
}

func newObservable() *observable {
	ctx, cancel := context.WithCancel(context.Background())
	return &observable{
		Callbacks: make(map[string][]Observer),
		Mutex:     sync.Mutex{},
		ctx:       ctx,
		cancel:    cancel,
	}
}
