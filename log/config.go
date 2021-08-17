package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"
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

var containerOnce sync.Once
var container bool

func IsContainerized() bool {
	containerOnce.Do(func() {
		_, err := os.Stat("/.dockerenv")
		if err == nil {
			container = true
		}

		if !container {
			_, err = os.Stat("/run/secrets/kubernetes.io/serviceaccount")
			if err == nil {
				container = true
			}
		}
	})

	return container
}
