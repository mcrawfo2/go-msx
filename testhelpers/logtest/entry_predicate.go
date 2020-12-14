package logtest

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

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

