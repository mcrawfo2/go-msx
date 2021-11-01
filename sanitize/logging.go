package sanitize

import (
	"github.com/sirupsen/logrus"
)

type LoggingFormatter struct {
	Base      logrus.Formatter
}

func (l *LoggingFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Message = secretSanitizer.Secrets(entry.Message)
	return l.Base.Format(entry)
}

var loggingFormatter = &LoggingFormatter{
	Base: logrus.StandardLogger().Formatter,
}

func init() {
	// Wrap the base formatter with our sanitizer
	logrus.StandardLogger().Formatter = loggingFormatter
}
