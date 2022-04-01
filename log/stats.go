// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package log

import "github.com/sirupsen/logrus"

type LoggingStatsCollector interface {
	Inc(loggerName string, level logrus.Level)
}

func RegisterStatsCollector(l LoggingStatsCollector) {
	logrus.AddHook(newStatsHook(l))
}

type statsHook struct {
	c LoggingStatsCollector
}

func (s statsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (s statsHook) Fire(entry *logrus.Entry) error {
	loggerName, _ := entry.Data[FieldLogger].(string)
	s.c.Inc(loggerName, entry.Level)
	return nil
}

func newStatsHook(c LoggingStatsCollector) logrus.Hook {
	return statsHook{c: c}
}
