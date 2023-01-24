// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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
			TimestampFormat: "2006-01-02T15:04:05.000000Z07:00",
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
