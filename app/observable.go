package app

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Observer func(ctx context.Context) error

type MsxApplication struct {
	Callbacks  map[string][]Observer
	stage      string
	appContext context.Context    // top-level context
	cancel     context.CancelFunc // cancels Startup/Runtime or Shutdown
	refresh    chan struct{}
	sync.Mutex
}

func (o *MsxApplication) On(event string, phase string, observer Observer) {
	o.Lock()
	defer o.Unlock()

	key := event + phase
	if _, ok := o.Callbacks[event+phase]; !ok {
		o.Callbacks[key] = []Observer{observer}
	} else {
		o.Callbacks[key] = append(o.Callbacks[key], observer)
	}
}

func (o *MsxApplication) Clear(event string, phase string) {
	o.Lock()
	defer o.Unlock()
	delete(o.Callbacks, event+phase)
}

func (o *MsxApplication) callbacks(event, phase string) ([]Observer, bool) {
	o.Lock()
	defer o.Unlock()
	observers, ok := o.Callbacks[event+phase]
	return observers, ok
}

func (o *MsxApplication) triggerPhase(ctx context.Context, event, phase string) error {
	logger.Infof("Event triggered: %s%s", event, phase)
	if observers, ok := o.callbacks(event, phase); ok {
		for _, observer := range observers {
			if ctx.Err() != nil {
				break
			}
			if err := observer(ctx); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Observer for event phase %s%s returned error", event, phase))
			}
		}
	}
	return nil
}

func (o *MsxApplication) triggerEvent(ctx context.Context, event string) error {
	for _, phase := range []string{PhaseBefore, PhaseDuring, PhaseAfter} {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err := o.triggerPhase(ctx, event, phase); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Event %s failed", event))
		}
	}
	return nil
}

func (o *MsxApplication) Run() error {
	// Create a new context for startup
	var runtimeContext context.Context
	runtimeContext, o.cancel = context.WithCancel(o.appContext)
	defer o.cancel()

	// Listen for process cancellation during startup and runtime
	die := make(chan os.Signal, 1)
	signal.Notify(die, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(die)
	}()
	go func() {
		select {
		case <-die:
			o.cancel()
		}
	}()

	if err := o.startupEvents(runtimeContext); err == nil {
		// Main loop
		select {
		case <-o.refresh:
			if err = o.refreshEvents(runtimeContext); err != nil {
				logger.Error(errors.Wrap(err, "Refresh failed"))
				break
			}

		case <-runtimeContext.Done():
			break
		}
	} else {
		logger.Error(errors.Wrap(err, "Startup failed"))
	}

	// Shutdown gracefully
	return o.shutdownEvents()
}

func (o *MsxApplication) Refresh() error {
	if o.stage != EventReady {
		return errors.New("Application not ready for refresh")
	}

	o.refresh <- struct{}{}
	return nil
}

func (o *MsxApplication) Stop() error {
	if o.cancel == nil {
		return errors.New("Not currently running")
	}

	o.cancel()
	return nil
}

func (o *MsxApplication) startupEvents(runtimeContext context.Context) error {
	for _, event := range []string{EventInit, EventConfigure, EventStart, EventReady} {
		if runtimeContext.Err() != nil {
			return runtimeContext.Err()
		}

		// Set to the latest stage we even attempted
		o.stage = event

		if err := o.triggerEvent(runtimeContext, event); err != nil {
			return errors.Wrap(err, "Startup failed")
		}
	}

	return nil
}

func (o *MsxApplication) refreshEvents(runtimeContext context.Context) error {
	if runtimeContext.Err() != nil {
		return runtimeContext.Err()
	}

	if err := o.triggerEvent(runtimeContext, EventRefresh); err != nil {
		return errors.Wrap(err, "Refresh failed")
	}

	return nil
}

func (o *MsxApplication) shutdownEvents() error {
	// Create a new context for bring-up
	var shutdownContext context.Context
	shutdownContext, o.cancel = context.WithCancel(o.appContext)
	defer o.cancel()

	// Listen for process cancellation during bring-up and runtime
	die := make(chan os.Signal, 1)
	signal.Notify(die, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(die)
	}()
	go func() {
		select {
		case <-die:
			o.cancel()
		}
	}()

	var events []string
	switch o.stage {
	case EventStart, EventReady, EventRefresh:
		events = []string{EventStop, EventFinal}
	case EventConfigure, EventInit:
		events = []string{EventFinal}
	}

	for _, event := range events {
		if shutdownContext.Err() != nil {
			return shutdownContext.Err()
		}

		// Set to the latest stage we even attempted
		o.stage = event

		if err := o.triggerEvent(shutdownContext, event); err != nil {
			return errors.Wrap(err, "Shutdown failed")
		}
	}

	return nil
}

func NewMsxApplication() *MsxApplication {
	return &MsxApplication{
		Callbacks:  make(map[string][]Observer),
		Mutex:      sync.Mutex{},
		appContext: context.Background(),
		refresh:    make(chan struct{}, 1),
	}
}
