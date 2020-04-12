package log

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type Logger struct {
	ParentLogger
	fields        logrus.Fields
	levelListener func()
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
	if logger.IsLevelEnabled(level) {
		logger.newEntry().Logf(level, format, args...)
	}
}

func (logger *Logger) Tracef(format string, args ...interface{}) {
	if logger.IsLevelEnabled(logrus.TraceLevel) {
		logger.newEntry().Tracef(format, args...)
	}
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	if logger.IsLevelEnabled(logrus.DebugLevel) {
		logger.newEntry().Debugf(format, args...)
	}
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	if logger.IsLevelEnabled(logrus.InfoLevel) {
		logger.newEntry().Infof(format, args...)
	}
}

func (logger *Logger) Printf(format string, args ...interface{}) {
	if logger.IsLevelEnabled(logrus.InfoLevel) {
		logger.newEntry().Printf(format, args...)
	}
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	if logger.IsLevelEnabled(logrus.WarnLevel) {
		logger.newEntry().Warnf(format, args...)
	}
}

func (logger *Logger) Warningf(format string, args ...interface{}) {
	if logger.IsLevelEnabled(logrus.WarnLevel) {
		logger.newEntry().Warningf(format, args...)
	}
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	if logger.IsLevelEnabled(logrus.ErrorLevel) {
		logger.newEntry().Errorf(format, args...)
	}
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	if logger.IsLevelEnabled(logrus.FatalLevel) {
		logger.newEntry().Fatalf(format, args...)
	}
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	if logger.IsLevelEnabled(logrus.PanicLevel) {
		logger.newEntry().Panicf(format, args...)
	}
}

func (logger *Logger) Log(level logrus.Level, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		logger.newEntry().Log(level, args...)
	}
}

func (logger *Logger) Trace(args ...interface{}) {
	if logger.IsLevelEnabled(logrus.TraceLevel) {
		logger.newEntry().Trace(args...)
	}
}

func (logger *Logger) Debug(args ...interface{}) {
	if logger.IsLevelEnabled(logrus.DebugLevel) {
		logger.newEntry().Debug(args...)
	}
}

func (logger *Logger) Info(args ...interface{}) {
	if logger.IsLevelEnabled(logrus.InfoLevel) {
		logger.newEntry().Info(args...)
	}
}

func (logger *Logger) Print(args ...interface{}) {
	if logger.IsLevelEnabled(logrus.InfoLevel) {
		logger.newEntry().Print(args...)
	}
}

func (logger *Logger) Warn(args ...interface{}) {
	if logger.IsLevelEnabled(logrus.WarnLevel) {
		logger.newEntry().Warn(args...)
	}
}

func (logger *Logger) Warning(args ...interface{}) {
	if logger.IsLevelEnabled(logrus.WarnLevel) {
		logger.newEntry().Warning(args...)
	}
}

func (logger *Logger) Error(args ...interface{}) {
	if logger.IsLevelEnabled(logrus.ErrorLevel) {
		logger.newEntry().Error(args...)
	}
}

func (logger *Logger) Fatal(args ...interface{}) {
	if logger.IsLevelEnabled(logrus.FatalLevel) {
		logger.newEntry().Fatal(args...)
	}
}

func (logger *Logger) Panic(args ...interface{}) {
	if logger.IsLevelEnabled(logrus.PanicLevel) {
		logger.newEntry().Panic(args...)
	}
}

func (logger *Logger) Logln(level logrus.Level, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		logger.newEntry().Logln(level, args...)
	}
}

func (logger *Logger) Traceln(args ...interface{}) {
	if logger.IsLevelEnabled(logrus.TraceLevel) {
		logger.newEntry().Traceln(args...)
	}
}

func (logger *Logger) Debugln(args ...interface{}) {
	if logger.IsLevelEnabled(logrus.DebugLevel) {
		logger.newEntry().Debugln(args...)
	}
}

func (logger *Logger) Infoln(args ...interface{}) {
	if logger.IsLevelEnabled(logrus.InfoLevel) {
		logger.newEntry().Infoln(args...)
	}
}

func (logger *Logger) Println(args ...interface{}) {
	if logger.IsLevelEnabled(logrus.InfoLevel) {
		logger.newEntry().Println(args...)
	}
}

func (logger *Logger) Warnln(args ...interface{}) {
	if logger.IsLevelEnabled(logrus.WarnLevel) {
		logger.newEntry().Warnln(args...)
	}
}

func (logger *Logger) Warningln(args ...interface{}) {
	if logger.IsLevelEnabled(logrus.WarnLevel) {
		logger.newEntry().Warningln(args...)
	}
}

func (logger *Logger) Errorln(args ...interface{}) {
	if logger.IsLevelEnabled(logrus.ErrorLevel) {
		logger.newEntry().Errorln(args...)
	}
}

func (logger *Logger) Fatalln(args ...interface{}) {
	if logger.IsLevelEnabled(logrus.FatalLevel) {
		logger.newEntry().Fatalln(args...)
	}
}

func (logger *Logger) Panicln(args ...interface{}) {
	if logger.IsLevelEnabled(logrus.PanicLevel) {
		logger.newEntry().Panicln(args...)
	}
}

func (logger *Logger) Level(level logrus.Level) StdLogger {
	return NewLevelLogger(logger, level)
}

func (logger *Logger) SetLevel(level logrus.Level) {
	logger.ParentLogger.SetLevel(level)
	name := logger.fields[FieldName].(string)
	levels[name] = level
	if logger.levelListener != nil {
		logger.levelListener()
	}
}

func (logger *Logger) IsLevelEnabled(level logrus.Level) bool {
	return logger.ParentLogger.GetLevel() >= level
}

func (logger *Logger) OnLevelChange(fn func()) {
	logger.levelListener = fn
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

var loggers = make(map[string]*Logger)
var levels = make(map[string]logrus.Level)

func NewLogger(name string, fields ...LogContext) *Logger {
	var level logrus.Level
	var ok bool
	if level, ok = levels[name]; !ok {
		level = InfoLevel
	}

	fields = append([]LogContext{{FieldName: name}}, fields...)
	logger := newLogger(&logrus.Logger{
		Out:       logrus.StandardLogger().Out,
		Formatter: logrus.StandardLogger().Formatter,
		Level:     level,
	}, fields...)

	logger.SetLevel(level)
	loggers[name] = logger

	return logger
}

func SetLoggerLevel(name string, level logrus.Level) {
	logger, ok := loggers[name]
	if ok {
		logger.SetLevel(level)
	} else {
		levels[name] = level
	}
}

func GetLoggerLevels() map[string]logrus.Level {
	return levels
}
