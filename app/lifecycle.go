package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	EventCommand   = "command"
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

type Observer types.ActionFunc
type CommandObserver func(context.Context, []string) error

type MsxApplication struct {
	Callbacks map[string][]Observer
	stage     string
	ctx       context.Context    // background context
	cancel    context.CancelFunc // cancels Startup/Runtime or Shutdown
	refresh   chan struct{}
	exitCode  int
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

func (a *MsxApplication) triggerObserver(event, phase string, observer Observer) error {
	// Inject all of the registered values into the context
	untracedContext := contextInjectors.Inject(a.ctx)

	// Store the context without the trace so observers can
	// have a context for asynchronous operations
	ctx := trace.ContextWithUntracedContext(untracedContext)

	observerName := types.ShortFunctionName(observer)
	ctx = log.NewContextWithLogContext(ctx, log.LogContext{
		"event":    event,
		"phase":    phase,
		"observer": observerName,
	})

	// Start a new trace span
	ctx, span := trace.NewSpan(ctx, event+phase+"."+observerName)
	defer span.Finish()
	span.SetTag(trace.FieldOperation, observerName)

	err := observer(ctx)
	if err != nil {
		span.LogFields(trace.Error(err))
		logger.WithContext(ctx).Error(err)
	}

	return err
}

func (a *MsxApplication) triggerPhase(ctx context.Context, event, phase string) error {
	logger.Debugf("Event triggered: %s%s", event, phase)
	if observers, ok := a.callbacks(event, phase); ok {
		for _, observer := range observers {
			if ctx.Err() != nil {
				break
			}

			if err := a.triggerObserver(event, phase, observer); err != nil {
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

		if err := a.triggerPhase(ctx, event, phase); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Event %s failed", event))
		}

	}
	return nil
}

func (a *MsxApplication) Run(command string) error {
	// Create a new context for startup
	a.ctx, a.cancel = context.WithCancel(context.Background())
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

	logger.WithContext(a.ctx).Infof("Command selected: %s", command)
	if err := a.triggerPhase(a.ctx, EventCommand, command); err != nil {
		logger.WithContext(a.ctx).WithError(err).Error("Command initialization failed")
		logger.WithContext(a.ctx).Errorf("%+v", err)
	}

	if err := a.startupEvents(a.ctx); err == nil {
		// Main loop
		for a.ctx.Err() == nil {
			select {
			case <-a.refresh:
				if err = a.refreshEvents(a.ctx); err != nil {
					logger.Error(errors.Wrap(err, "Refresh failed"))
					logger.WithContext(a.ctx).Errorf("%+v", err)
					break
				}

			case <-a.ctx.Done():
				break
			}
		}
	} else {
		logger.WithContext(a.ctx).WithError(err).Error("Startup failed")
		logger.WithContext(a.ctx).Errorf("%+v", err)
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
	// Create a new context for shutdown
	a.ctx, a.cancel = context.WithCancel(context.Background())
	defer a.cancel()

	// Listen for process cancellation during shutdown
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
		if a.ctx.Err() != nil {
			return a.ctx.Err()
		}

		// Set to the latest stage we even attempted
		a.stage = event

		if err := a.triggerEvent(a.ctx, event); err != nil {
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
		Callbacks: make(map[string][]Observer),
		Mutex:     sync.Mutex{},
		refresh:   make(chan struct{}, 1),
	}
}

var application = NewMsxApplication()

func OnEvent(event string, phase string, observer Observer) {
	application.On(event, phase, observer)
}
