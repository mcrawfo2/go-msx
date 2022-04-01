// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var standardLogger *Logger

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		QuoteEmptyFields:       true,
		TimestampFormat:        time.RFC3339Nano,
		DisableLevelTruncation: true,
	})

	logrus.SetOutput(os.Stdout)

	standardLogger = NewLogger("msx")
}

func StandardLogger() *Logger {
	return standardLogger
}
