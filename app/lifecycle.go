package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	EventInit      = "initialize"
	EventConfigure = "configure"
	EventStart     = "start"
	EventReady     = "ready"
	EventRefresh   = "refresh"
	EventStop      = "stop"
	EventFinal     = "finalize"

	PhaseBefore = ".before"
	PhaseDuring = ""
	PhaseAfter  = ".after"
)

type Observer func(ctx context.Context) error

type MsxApplication struct {
	Callbacks  map[string][]Observer
	stage      string
	appContext context.Context    // top-level context
	cancel     context.CancelFunc // cancels Startup/Runtime or Shutdown
	refresh    chan struct{}
	exitCode   int
	sync.Mutex
}

func (a *MsxApplication) On(event string, phase string, observer Observer) {
	a.Lock()
	defer a.Unlock()

	key := event + phase
	if _, ok := a.Callbacks[event+phase]; !ok {
		a.Callbacks[key] = []Observer{observer}
	} else {
		a.Callbacks[key] = append(a.Callbacks[key], observer)
	}
}

func (a *MsxApplication) Clear(event string, phase string) {
	a.Lock()
	defer a.Unlock()
	delete(a.Callbacks, event+phase)
}

func (a *MsxApplication) callbacks(event, phase string) ([]Observer, bool) {
	a.Lock()
	defer a.Unlock()
	observers, ok := a.Callbacks[event+phase]
	return observers, ok
}

func (a *MsxApplication) triggerPhase(ctx context.Context, event, phase string) error {
	logger.Debugf("Event triggered: %s%s", event, phase)
	if observers, ok := a.callbacks(event, phase); ok {
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

func (a *MsxApplication) triggerEvent(ctx context.Context, event string) error {
	for _, phase := range []string{PhaseBefore, PhaseDuring, PhaseAfter} {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// When configure.ready has completed, add the application config to the supplied context
		phaseCtx := ctx
		if isConfigured(event, phase) {
			phaseCtx = config.ContextWithConfig(ctx, applicationConfig)
		}

		if err := a.triggerPhase(phaseCtx, event, phase); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Event %s failed", event))
		}

	}
	return nil
}

func isConfigured(event, phase string) bool {
	switch event {
	case EventInit:
		return false
	case EventConfigure:
		return phase == PhaseAfter
	default:
		return true
	}
}

func (a *MsxApplication) Run() error {
	// Create a new context for startup
	var runtimeContext context.Context
	runtimeContext, a.cancel = context.WithCancel(a.appContext)
	defer a.cancel()

	// Listen for process cancellation during startup and runtime
	die := make(chan os.Signal, 1)
	signal.Notify(die, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(die)
	}()
	go func() {
		select {
		case <-die:
			a.cancel()
		}
	}()

	if err := a.startupEvents(runtimeContext); err == nil {
		// Main loop
		select {
		case <-a.refresh:
			if err = a.refreshEvents(runtimeContext); err != nil {
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
	return a.shutdownEvents()
}

func (a *MsxApplication) Refresh() error {
	if a.stage != EventReady {
		return errors.New("Application not ready for refresh")
	}

	a.refresh <- struct{}{}
	return nil
}

func (a *MsxApplication) Stop() error {
	if a.cancel == nil {
		return errors.New("Not currently running")
	}

	a.cancel()
	return nil
}

func (a *MsxApplication) startupEvents(runtimeContext context.Context) error {
	for _, event := range []string{EventInit, EventConfigure, EventStart, EventReady} {
		if runtimeContext.Err() != nil {
			return runtimeContext.Err()
		}

		// Set to the latest stage we even attempted
		a.stage = event

		if err := a.triggerEvent(runtimeContext, event); err != nil {
			return errors.Wrap(err, "Startup failed")
		}
	}

	return nil
}

func (a *MsxApplication) refreshEvents(runtimeContext context.Context) error {
	if runtimeContext.Err() != nil {
		return runtimeContext.Err()
	}

	if err := a.triggerEvent(runtimeContext, EventRefresh); err != nil {
		return errors.Wrap(err, "Refresh failed")
	}

	return nil
}

func (a *MsxApplication) shutdownEvents() error {
	// Create a new context for bring-up
	var shutdownContext context.Context
	shutdownContext, a.cancel = context.WithCancel(a.appContext)
	defer a.cancel()

	// Listen for process cancellation during bring-up and runtime
	die := make(chan os.Signal, 1)
	signal.Notify(die, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(die)
	}()
	go func() {
		select {
		case <-die:
			a.cancel()
		}
	}()

	var events []string
	switch a.stage {
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
		a.stage = event

		if err := a.triggerEvent(shutdownContext, event); err != nil {
			return errors.Wrap(err, "Shutdown failed")
		}
	}

	if a.exitCode != 0 {
		return errors.New(fmt.Sprintf("Application exited with status: %d", a.exitCode))
	} else {
		return nil
	}
}

func (a *MsxApplication) SetExitCode(exitCode int) {
	a.exitCode = exitCode
}

func NewMsxApplication() *MsxApplication {
	return &MsxApplication{
		Callbacks:  make(map[string][]Observer),
		Mutex:      sync.Mutex{},
		appContext: context.Background(),
		refresh:    make(chan struct{}, 1),
	}
}

var application = NewMsxApplication()

func OnEvent(event string, phase string, observer Observer) {
	application.On(event, phase, observer)
}
