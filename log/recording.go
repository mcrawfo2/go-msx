package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type Recording struct {
	Entries   []logrus.Entry
	formatter logrus.Formatter
}

func (m *Recording) Format(entry *logrus.Entry) ([]byte, error) {
	m.Entries = append(m.Entries, *entry)
	return m.formatter.Format(entry)
}

func (m *Recording) Reset() {
	m.Entries = nil
}

func RecordLogging() *Recording {
	recording := &Recording{
		Entries:   make([]logrus.Entry, 0),
		formatter: logrus.StandardLogger().Formatter,
	}
	logrus.StandardLogger().SetFormatter(recording)
	return recording
}

type EntryPredicate struct {
	Description string
	Matches     func(logrus.Entry) bool
}

func HasLevel(level logrus.Level) EntryPredicate {
	return EntryPredicate{
		Description: fmt.Sprintf("entry.level == %d", level),
		Matches: func(entry logrus.Entry) bool {
			return entry.Level == level
		},
	}
}

func HasMessage(msg string) EntryPredicate {
	return EntryPredicate{
		Description: fmt.Sprintf("entry.message == %q", msg),
		Matches: func(entry logrus.Entry) bool {
			return entry.Message == msg
		},
	}
}

func HasFieldValue(name string, value interface{}) EntryPredicate {
	return EntryPredicate{
		Description: fmt.Sprintf("entry.Data['%s'] = %q", name, value),
		Matches: func(entry logrus.Entry) bool {
			return entry.Data[name] == value
		},
	}
}

func FieldValue(name string) EntryPredicate {
	return EntryPredicate{
		Description: fmt.Sprintf("entry.Data['%s'] set", name),
		Matches: func(entry logrus.Entry) bool {
			_, ok := entry.Data[name]
			return ok
		},
	}
}


func NoFieldValue(name string) EntryPredicate {
	return EntryPredicate{
		Description: fmt.Sprintf("entry.Data['%s'] unset", name),
		Matches: func(entry logrus.Entry) bool {
			_, ok := entry.Data[name]
			return !ok
		},
	}
}

type Matcher struct {
	Filters []EntryPredicate
}

func (m Matcher) MatchEntries(r *Recording) []logrus.Entry {
	var results []logrus.Entry
	for _, entry := range r.Entries {
		if m.isMatch(entry) {
			results = append(results, entry)
		}
	}

	return results
}

func (m Matcher) isMatch(entry logrus.Entry) bool {
	for _, filter := range m.Filters {
		if !filter.Matches(entry) {
			return false
		}
	}
	return true
}

type CheckError struct {
	Entry     logrus.Entry
	Validator EntryPredicate
}

func (c CheckError) Error() string {
	return fmt.Sprintf("Failed validator: %s - %+v", c.Validator.Description, c.Entry)
}

type Check struct {
	Filters    []EntryPredicate
	Validators []EntryPredicate
}

func (c Check) Check(r *Recording) []error {
	matcher := Matcher{
		Filters: c.Filters,
	}

	var results []error

	for _, entry := range matcher.MatchEntries(r) {
		for _, predicate := range c.Validators {
			if !predicate.Matches(entry) {
				results = append(results, CheckError{
					Entry:     entry,
					Validator: predicate,
				})
			}
		}
	}

	return results
}
