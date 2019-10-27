package stream

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/ThreeDotsLabs/watermill"
)

type WatermillLoggerAdapter struct {
	logger *log.Logger
}

func (w *WatermillLoggerAdapter) Error(msg string, err error, fields watermill.LogFields) {
	w.logger.WithLogContext(log.LogContext(fields)).WithError(err).Error(msg)
}

func (w *WatermillLoggerAdapter) Info(msg string, fields watermill.LogFields) {
	w.logger.WithLogContext(log.LogContext(fields)).Info(msg)
}

func (w *WatermillLoggerAdapter) Debug(msg string, fields watermill.LogFields) {
	w.logger.WithLogContext(log.LogContext(fields)).Debug(msg)
}

func (w *WatermillLoggerAdapter) Trace(msg string, fields watermill.LogFields) {
	w.logger.WithLogContext(log.LogContext(fields)).Trace(msg)
}

func (w *WatermillLoggerAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &WatermillLoggerAdapter{
		logger: w.logger.WithExtendedLogContext(log.LogContext(fields)),
	}
}

func NewWatermillLoggerAdapter(logger *log.Logger) watermill.LoggerAdapter {
	return &WatermillLoggerAdapter{logger: logger}
}
