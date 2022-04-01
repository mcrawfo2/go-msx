// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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
		Description: fmt.Sprintf("entry.Data['%s'] = %v", name, value),
		Matches: func(entry logrus.Entry) bool {
			return entry.Data[name] == value
		},
	}
}

func HasError(message string) EntryPredicate {
	return EntryPredicate{
		Description: fmt.Sprintf("entry.Data['error'].Error() = %q", message),
		Matches: func(entry logrus.Entry) bool {
			return entry.Data["error"].(error).Error() == message
		},
	}
}

func Index(skip int, keep int) EntryPredicate {
	return EntryPredicate{
		Description: fmt.Sprintf("entry index %d < n < %d", skip, keep),
		Matches: func(entry logrus.Entry) bool {
			if skip > 0 {
				skip--
				return false
			}
			if keep > 0 {
				keep--
				return true
			}
			return false
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
