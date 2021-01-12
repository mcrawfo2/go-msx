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

var dockerOnce sync.Once
var docker bool

func IsDocker() bool {
	dockerOnce.Do(func() {
		_, err := os.Stat("/.dockerenv")
		if err == nil {
			docker = true
		} else {
			docker = !os.IsNotExist(err)
		}
	})

	return docker
}
