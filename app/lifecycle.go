package app

import (
	"context"
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
)

func OnEvent(event string, phase string, observer Observer) {
	application.On(event, phase, observer)
}

func ClearEvent(event string, phase string) {
	application.Clear(event, phase)
}

func TriggerEvent(event string) {
	application.Trigger(event)
}

func Shutdown() {
	application.Shutdown()
}

func Context() context.Context {
	return application.ctx
}

func Lifecycle() {
	TriggerEvent(EventInit)
	TriggerEvent(EventConfigure)
	TriggerEvent(EventStart)
	TriggerEvent(EventReady)
	TriggerEvent(EventStop)
	TriggerEvent(EventFinal)
}
