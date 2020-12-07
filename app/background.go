package app

type errorReporter struct {
	cancel chan struct{}
	app    *MsxApplication
}

func (e errorReporter) Fatal(err error) {
	logger.WithError(err).Error("Background task returned fatal error")
	e.cancel <- struct{}{}
}

func (e errorReporter) NonFatal(err error) {
	logger.WithError(err).Error("Background task returned non-fatal error")
}

func (e errorReporter) C() <-chan struct{} {
	return e.cancel
}

func newErrorReporter(app *MsxApplication) errorReporter {
	return errorReporter{
		cancel: make(chan struct{}),
		app:    app,
	}
}
