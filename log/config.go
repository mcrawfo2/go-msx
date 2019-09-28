package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func init() {
	logrus.SetOutput(os.Stdout)
}

func SetLevel(level string) {
	switch level {
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	}
}

func SetFormat(format string) {
	switch format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint:     true,
			TimestampFormat: time.RFC3339Nano,
		})

	default:
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339Nano,
		})
	}
}

