// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package logtest

import (
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
