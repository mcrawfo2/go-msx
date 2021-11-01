package sanitize

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
)

type LoggingFormatter struct {
	Base logrus.Formatter
}

func (l *LoggingFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Sanitize the fields to a new set
	newData := make(logrus.Fields)
	for i, field := range entry.Data {
		newData[i] = LogField(field)
	}

	// Swap in the sanitized versions until we are done formatting
	newData, entry.Data = entry.Data, newData
	defer func() {
		// Restore the original entry.Data
		entry.Data = newData
	}()

	// Render the message
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

type ErrorSanitizer struct {
	err error
}

func (e ErrorSanitizer) Error() string {
	return secretSanitizer.Secrets(e.err.Error())
}

type StringerSanitizer struct {
	value fmt.Stringer
}

func (e StringerSanitizer) String() string {
	return secretSanitizer.Secrets(e.value.String())
}

func LogField(i interface{}) (result interface{}) {
	defer func() {
		if err := recover(); err != nil {
			errString := fmt.Sprintf("%+v", err)
			result = "%!v(PANIC=LogField method: " + errString + ")"
		}
	}()

	// Check for self-rendering
	if err, ok := i.(error); ok {
		return ErrorSanitizer{err:err}
	} else if stringer, ok := i.(fmt.Stringer); ok {
		return StringerSanitizer{value:stringer}
	}

	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
		i = v.Interface()
		b, err := json.Marshal(i)
		if err != nil {
			panic(err)
		}
		return secretSanitizer.Secrets(string(b))
	case reflect.Interface:
		return LogField(v.Interface())
	case reflect.String:
		return secretSanitizer.Secrets(v.Interface().(string))
	default:
		return v.Interface()
	}
}
