// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package log

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/go-stack/stack"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type Logger struct {
	ParentLogger
	fields        logrus.Fields
	levelListener func()
}

type StdLogger logrus.StdLogger

const (
	FieldLogger = "logger"
	FieldStack  = "stack"
)

func (l *Logger) Fields() LogContext {
	// Return all fields except `name`
	result := make(logrus.Fields)
	for k, v := range l.fields {
		if k != FieldLogger {
			result[k] = v
		}
	}

	return LogContext(result)
}

func (l *Logger) newEntry() *logrus.Entry {
	return l.ParentLogger.WithFields(l.fields)
}

func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.newEntry().WithField(key, value)
}

func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.newEntry().WithFields(fields)
}

func (l *Logger) WithLogContext(logCtx LogContext) *logrus.Entry {
	return l.newEntry().WithFields(logrus.Fields(logCtx))
}

// WithError adds an error as single field to a new log entry.
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.newEntry().WithError(err)
}

// WithContext adds a Context to the log entry.
func (l *Logger) WithContext(ctx context.Context) *logrus.Entry {
	if ctx == nil {
		return l.newEntry()
	}

	entry := l.newEntry().WithContext(ctx)
	if logContext, ok := LogContextFromContext(ctx); ok {
		entry = entry.WithFields(logrus.Fields(logContext))
	}
	return entry
}

// WithTime overrides the time of the log entry.
func (l *Logger) WithTime(t time.Time) *logrus.Entry {
	return l.newEntry().WithTime(t)
}

func (l *Logger) WithExtendedField(key string, value interface{}) *Logger {
	return newLogger(l, LogContext{key: value})
}

func (l *Logger) WithExtendedLogContext(fields ...LogContext) *Logger {
	if len(fields) == 0 {
		return l
	}

	return newLogger(l, fields...)
}

func (l *Logger) Logf(level logrus.Level, format string, args ...interface{}) {
	if l.IsLevelEnabled(level) {
		l.newEntry().Logf(level, format, args...)
	}
}

func (l *Logger) Tracef(format string, args ...interface{}) {
	if l.IsLevelEnabled(logrus.TraceLevel) {
		l.newEntry().Tracef(format, args...)
	}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.IsLevelEnabled(logrus.DebugLevel) {
		l.newEntry().Debugf(format, args...)
	}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if l.IsLevelEnabled(logrus.InfoLevel) {
		l.newEntry().Infof(format, args...)
	}
}

func (l *Logger) Printf(format string, args ...interface{}) {
	if l.IsLevelEnabled(logrus.InfoLevel) {
		l.newEntry().Printf(format, args...)
	}
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.IsLevelEnabled(logrus.WarnLevel) {
		l.newEntry().Warnf(format, args...)
	}
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	if l.IsLevelEnabled(logrus.WarnLevel) {
		l.newEntry().Warningf(format, args...)
	}
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.IsLevelEnabled(logrus.ErrorLevel) {
		l.newEntry().Errorf(format, args...)
	}
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	if l.IsLevelEnabled(logrus.FatalLevel) {
		l.newEntry().Fatalf(format, args...)
	}
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	if l.IsLevelEnabled(logrus.PanicLevel) {
		l.newEntry().Panicf(format, args...)
	}
}

func (l *Logger) Log(level logrus.Level, args ...interface{}) {
	if l.IsLevelEnabled(level) {
		l.newEntry().Log(level, args...)
	}
}

func (l *Logger) Trace(args ...interface{}) {
	if l.IsLevelEnabled(logrus.TraceLevel) {
		l.newEntry().Trace(args...)
	}
}

func (l *Logger) Debug(args ...interface{}) {
	if l.IsLevelEnabled(logrus.DebugLevel) {
		l.newEntry().Debug(args...)
	}
}

func (l *Logger) Info(args ...interface{}) {
	if l.IsLevelEnabled(logrus.InfoLevel) {
		l.newEntry().Info(args...)
	}
}

func (l *Logger) Print(args ...interface{}) {
	if l.IsLevelEnabled(logrus.InfoLevel) {
		l.newEntry().Print(args...)
	}
}

func (l *Logger) Warn(args ...interface{}) {
	if l.IsLevelEnabled(logrus.WarnLevel) {
		l.newEntry().Warn(args...)
	}
}

func (l *Logger) Warning(args ...interface{}) {
	if l.IsLevelEnabled(logrus.WarnLevel) {
		l.newEntry().Warning(args...)
	}
}

func (l *Logger) Error(args ...interface{}) {
	if l.IsLevelEnabled(logrus.ErrorLevel) {
		l.newEntry().Error(args...)
	}
}

func (l *Logger) Fatal(args ...interface{}) {
	if l.IsLevelEnabled(logrus.FatalLevel) {
		l.newEntry().Fatal(args...)
	}
}

func (l *Logger) Panic(args ...interface{}) {
	if l.IsLevelEnabled(logrus.PanicLevel) {
		l.newEntry().Panic(args...)
	}
}

func (l *Logger) Logln(level logrus.Level, args ...interface{}) {
	if l.IsLevelEnabled(level) {
		l.newEntry().Logln(level, args...)
	}
}

func (l *Logger) Traceln(args ...interface{}) {
	if l.IsLevelEnabled(logrus.TraceLevel) {
		l.newEntry().Traceln(args...)
	}
}

func (l *Logger) Debugln(args ...interface{}) {
	if l.IsLevelEnabled(logrus.DebugLevel) {
		l.newEntry().Debugln(args...)
	}
}

func (l *Logger) Infoln(args ...interface{}) {
	if l.IsLevelEnabled(logrus.InfoLevel) {
		l.newEntry().Infoln(args...)
	}
}

func (l *Logger) Println(args ...interface{}) {
	if l.IsLevelEnabled(logrus.InfoLevel) {
		l.newEntry().Println(args...)
	}
}

func (l *Logger) Warnln(args ...interface{}) {
	if l.IsLevelEnabled(logrus.WarnLevel) {
		l.newEntry().Warnln(args...)
	}
}

func (l *Logger) Warningln(args ...interface{}) {
	if l.IsLevelEnabled(logrus.WarnLevel) {
		l.newEntry().Warningln(args...)
	}
}

func (l *Logger) Errorln(args ...interface{}) {
	if l.IsLevelEnabled(logrus.ErrorLevel) {
		l.newEntry().Errorln(args...)
	}
}

func (l *Logger) Fatalln(args ...interface{}) {
	if l.IsLevelEnabled(logrus.FatalLevel) {
		l.newEntry().Fatalln(args...)
	}
}

func (l *Logger) Panicln(args ...interface{}) {
	if l.IsLevelEnabled(logrus.PanicLevel) {
		l.newEntry().Panicln(args...)
	}
}

func (l *Logger) Level(level logrus.Level) StdLogger {
	return NewLevelLogger(l, level)
}

func (l *Logger) SetLevel(level logrus.Level) {
	l.ParentLogger.SetLevel(level)
	name := l.fields[FieldLogger].(string)
	levels[name] = level
	if l.levelListener != nil {
		l.levelListener()
	}
}

func (l *Logger) IsLevelEnabled(level logrus.Level) bool {
	return l.ParentLogger.GetLevel() >= level
}

func (l *Logger) OnLevelChange(fn func()) {
	l.levelListener = fn
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

type GlobalFormatter struct{}

func (g GlobalFormatter) Format(e *logrus.Entry) ([]byte, error) {
	return logrus.StandardLogger().Formatter.Format(e)
}

var loggers = make(map[string]*Logger)
var levels = make(map[string]logrus.Level)

func NewLogger(name string, fields ...LogContext) *Logger {
	var level logrus.Level
	var ok bool
	if level, ok = levels[name]; !ok {
		level = InfoLevel
	}

	stdLogger := logrus.StandardLogger()
	fields = append([]LogContext{{FieldLogger: name}}, fields...)
	logger := newLogger(&logrus.Logger{
		Out:       stdLogger.Out,
		Formatter: GlobalFormatter{},
		Level:     level,
		Hooks:     stdLogger.Hooks,
	}, fields...)

	logger.SetLevel(level)
	loggers[name] = logger

	return logger
}

func NewPackageLogger() *Logger {
	// Find caller directory name
	st := stack.Caller(1)
	tf := st.Frame()
	fp := strings.Split(tf.Function, "/")

	// Strip root package from caller package (minus 1)
	fp = fp[2:]
	fp[0] = strings.TrimPrefix(fp[0], "go-")
	fp[0] = strings.ReplaceAll(fp[0], "-", ".")
	fp[0] = strings.TrimSuffix(fp[0], "service")

	// Strip function name from caller function
	lp := strings.SplitN(fp[len(fp)-1], ".", 2)
	fp[len(fp)-1] = lp[0]

	// Convert stripped path to dotted notation
	pn := strings.Join(fp, ".")

	// Return a logger
	return NewLogger(pn)
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

func ErrorMessage(logger *Logger, ctx context.Context, err error) {
	if IsContainerized() {
		return
	}

	message := err.Error()
	if !strings.Contains(message, "\n") {
		return
	}

	for _, line := range strings.Split(message, "\n") {
		logger.WithContext(ctx).Error(line)
	}
}

func Stack(logger *Logger, ctx context.Context, bt types.BackTrace) {
	if IsContainerized() {
		return
	}

	for _, line := range bt.Lines() {
		logger.WithContext(ctx).Error(line)
	}
}
