package lifecycle

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/support/log"
)

const (
	EventInit      = "initialize"
	EventConfigure = "configure"
	EventStart     = "start"
	EventReady     = "ready"
	EventStop      = "stop"
	EventFinal     = "finalize"

	PhaseBefore = ".before"
	PhaseDuring = ""
	PhaseAfter  = ".after"
)

var (
	application = newObservable()
	logger      = log.NewLogger("msx.lifecycle")
)

func OnEvent(event string, phase string, observer Observer) {
	application.On(event, phase, observer)
}

func ClearEvent(event string, phase string) {
	application.Clear(event, phase)
}

func TriggerEvent(event string) {
	logger.Infof("Event triggered: %s", event)
	application.Trigger(event)
}

func Shutdown() {
	logger.Info("Shutdown triggered")
	application.Shutdown()
}

func Context() context.Context {
	return application.ctx
}

func Run() {
	TriggerEvent(EventInit)
	TriggerEvent(EventConfigure)
	TriggerEvent(EventStart)
	TriggerEvent(EventReady)
	TriggerEvent(EventStop)
	TriggerEvent(EventFinal)
}
