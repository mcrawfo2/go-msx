package app

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

type errorReporter struct {
	cancel chan struct{}
	app    *MsxApplication
}

func (e errorReporter) Fatal(err error) {
	bt := types.BackTraceFromError(err)
	logger.
		WithError(err).
		WithField(log.FieldStack, bt.Stanza()).
		Error("Background task returned fatal error")
	log.Stack(logger, nil, bt)

	e.cancel <- struct{}{}
}

func (e errorReporter) NonFatal(err error) {
	bt := types.BackTraceFromError(err)
	logger.
		WithError(err).
		WithField(log.FieldStack, bt.Stanza()).
		Error("Background task returned non-fatal error")
	log.Stack(logger, nil, bt)
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
