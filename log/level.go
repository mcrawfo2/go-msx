package log

import (
	"github.com/sirupsen/logrus"
)

const (
	TraceLevel = logrus.TraceLevel
	DebugLevel = logrus.DebugLevel
	InfoLevel  = logrus.InfoLevel
	WarnLevel  = logrus.WarnLevel
	ErrorLevel = logrus.ErrorLevel
	FatalLevel = logrus.FatalLevel
	PanicLevel = logrus.PanicLevel
)

type LevelLogger struct {
	parent ParentLogger
	level  logrus.Level

	printf  func(string, ...interface{})
	print   func(...interface{})
	println func(...interface{})
}

func NewLevelLogger(logger ParentLogger, level logrus.Level) logrus.StdLogger {
	var fnPrintf func(string, ...interface{})
	var fnPrint func(...interface{})
	var fnPrintln func(...interface{})

	switch level {
	case TraceLevel:
		fnPrintf = logger.Tracef
		fnPrint = logger.Trace
		fnPrintln = logger.Traceln
	case DebugLevel:
		fnPrintf = logger.Debugf
		fnPrint = logger.Debug
		fnPrintln = logger.Debugln
	case InfoLevel:
		fnPrintf = logger.Infof
		fnPrint = logger.Info
		fnPrintln = logger.Infoln
	case WarnLevel:
		fnPrintf = logger.Warnf
		fnPrint = logger.Warn
		fnPrintln = logger.Warnln
	case ErrorLevel:
		fnPrintf = logger.Errorf
		fnPrint = logger.Error
		fnPrintln = logger.Errorln
	case FatalLevel:
		fnPrintf = logger.Fatalf
		fnPrint = logger.Fatal
		fnPrintln = logger.Fatalln
	case PanicLevel:
		fnPrintf = logger.Panicf
		fnPrint = logger.Panic
		fnPrintln = logger.Panicln
	}

	return &LevelLogger{
		parent: logger,
		level:  level,

		printf:  fnPrintf,
		print:   fnPrint,
		println: fnPrintln,
	}
}

func (l *LevelLogger) Printf(template string, values ...interface{}) {
	l.printf(template, values...)
}

func (l *LevelLogger) Print(values ...interface{}) {
	l.print(values...)
}

func (l *LevelLogger) Println(values ...interface{}) {
	l.println(values...)
}

func (l *LevelLogger) Fatal(values ...interface{}) {
	l.parent.Fatal(values...)
}

func (l *LevelLogger) Fatalf(template string, values ...interface{}) {
	l.parent.Fatalf(template, values...)
}

func (l *LevelLogger) Fatalln(values ...interface{}) {
	l.parent.Fatalln(values...)
}

func (l *LevelLogger) Panic(values ...interface{}) {
	l.parent.Panic(values...)
}

func (l *LevelLogger) Panicf(template string, values ...interface{}) {
	l.parent.Panicf(template, values...)
}

func (l *LevelLogger) Panicln(values ...interface{}) {
	l.parent.Panicln(values...)
}
