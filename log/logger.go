package log

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type Logger struct {
	ParentLogger
	fields logrus.Fields
}

type StdLogger logrus.StdLogger

const FieldName = "logger"

func (logger *Logger) Fields() LogContext {
	// Return all fields except `name`
	result := make(logrus.Fields)
	for k, v := range logger.fields {
		if k != FieldName {
			result[k] = v
		}
	}

	return LogContext(result)
}

func (logger *Logger) newEntry() *logrus.Entry {
	return logger.ParentLogger.WithFields(logger.fields)
}

func (logger *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return logger.newEntry().WithField(key, value)
}

func (logger *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.newEntry().WithFields(fields)
}

func (logger *Logger) WithLogContext(logCtx LogContext) *logrus.Entry {
	return logger.newEntry().WithFields(logrus.Fields(logCtx))
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func (logger *Logger) WithError(err error) *logrus.Entry {
	return logger.newEntry().WithError(err)
}

// Add a context to the log entry.
func (logger *Logger) WithContext(ctx context.Context) *logrus.Entry {
	entry := logger.newEntry().WithContext(ctx)
	if logContext, ok := LogContextFromContext(ctx); ok {
		entry = entry.WithFields(logrus.Fields(logContext))
	}
	return entry
}

// Overrides the time of the log entry.
func (logger *Logger) WithTime(t time.Time) *logrus.Entry {
	return logger.newEntry().WithTime(t)
}

func (logger *Logger) WithExtendedField(key string, value interface{}) *Logger {
	return newLogger(logger, LogContext{key: value})
}

func (logger *Logger) WithExtendedLogContext(fields ...LogContext) *Logger {
	if len(fields) == 0 {
		return logger
	}

	return newLogger(logger, fields...)
}

func (logger *Logger) Logf(level logrus.Level, format string, args ...interface{}) {
	if logger.ParentLogger.IsLevelEnabled(level) {
		logger.newEntry().Logf(level, format, args...)
	}
}

func (logger *Logger) Tracef(format string, args ...interface{}) {
	logger.newEntry().Tracef(format, args...)
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.newEntry().Debugf(format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.newEntry().Infof(format, args...)
}

func (logger *Logger) Printf(format string, args ...interface{}) {
	logger.newEntry().Printf(format, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.newEntry().Warnf(format, args...)
}

func (logger *Logger) Warningf(format string, args ...interface{}) {
	logger.newEntry().Warningf(format, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.newEntry().Errorf(format, args...)
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.newEntry().Fatalf(format, args...)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	logger.newEntry().Panicf(format, args...)
}

func (logger *Logger) Log(level logrus.Level, args ...interface{}) {
	logger.newEntry().Log(level, args...)
}

func (logger *Logger) Trace(args ...interface{}) {
	logger.newEntry().Trace(args...)
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.newEntry().Debug(args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.newEntry().Info(args...)
}

func (logger *Logger) Print(args ...interface{}) {
	logger.newEntry().Print(args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.newEntry().Warn(args...)
}

func (logger *Logger) Warning(args ...interface{}) {
	logger.newEntry().Warning(args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.newEntry().Error(args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.newEntry().Fatal(args...)
}

func (logger *Logger) Panic(args ...interface{}) {
	logger.newEntry().Panic(args...)
}

func (logger *Logger) Logln(level logrus.Level, args ...interface{}) {
	logger.newEntry().Logln(level, args...)
}

func (logger *Logger) Traceln(args ...interface{}) {
	logger.newEntry().Traceln(args...)
}

func (logger *Logger) Debugln(args ...interface{}) {
	logger.newEntry().Debugln(args...)
}

func (logger *Logger) Infoln(args ...interface{}) {
	logger.newEntry().Infoln(args...)
}

func (logger *Logger) Println(args ...interface{}) {
	logger.newEntry().Println(args...)
}

func (logger *Logger) Warnln(args ...interface{}) {
	logger.newEntry().Warnln(args...)
}

func (logger *Logger) Warningln(args ...interface{}) {
	logger.newEntry().Warningln(args...)
}

func (logger *Logger) Errorln(args ...interface{}) {
	logger.newEntry().Errorln(args...)
}

func (logger *Logger) Fatalln(args ...interface{}) {
	logger.newEntry().Fatalln(args...)
}

func (logger *Logger) Panicln(args ...interface{}) {
	logger.newEntry().Panicln(args...)
}

func (logger *Logger) Level(level logrus.Level) StdLogger {
	return NewLevelLogger(logger, level)
}

func newLogger(logger ParentLogger, fields ...LogContext) *Logger {
	allFields := make(logrus.Fields)
	for _, field := range fields {
		for k, v := range field {
			allFields[k] = v
		}
	}

	return &Logger{
		ParentLogger: logger,
		fields:       allFields,
	}
}

func NewLogger(name string, fields ...LogContext) *Logger {
	fields = append([]LogContext{{FieldName: name}}, fields...)
	return newLogger(logrus.StandardLogger(), fields...)
}
