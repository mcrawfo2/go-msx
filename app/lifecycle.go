package app

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

var (
	application = NewMsxApplication()
)

func OnEvent(event string, phase string, observer Observer) {
	application.On(event, phase, observer)
}

func ClearEvent(event string, phase string) {
	application.Clear(event, phase)
}

// Startup: cancels further startup, moves to shutdown
// Run: cancels running, moves to shutdown
// Shutdown: cancels further shutdown
func Stop() error {
	return application.Stop()
}

// Run: executes refresh events
func Refresh() error {
	return application.Refresh()
}

func Lifecycle() error {
	return application.Run()
}
