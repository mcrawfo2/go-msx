package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

const (
	LogFormatJson   = "json"
	LogFormatLogFmt = "logfmt"
)

func init() {
	logrus.SetOutput(os.Stdout)
}

func SetLevel(level string) {
	logLevel := LevelFromName(strings.ToUpper(level))
	logrus.SetLevel(logLevel)
}

func SetFormat(format string) {
	switch format {
	case LogFormatJson:
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
