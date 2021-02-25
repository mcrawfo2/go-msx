package stream

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/logtest"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
)

func TestNewWatermillLoggerAdapter(t *testing.T) {
	testLogger := log.NewLogger("TestNewWatermillLoggerAdapter")
	expected := &WatermillLoggerAdapter{logger: testLogger}
	if got := NewWatermillLoggerAdapter(testLogger); !reflect.DeepEqual(got, expected) {
		t.Errorf("NewWatermillLoggerAdapter() = %v, want %v", got, expected)
	}
}

func TestWatermillLoggerAdapter_Debug(t *testing.T) {
	tests := []struct {
		name    string
		fields  watermill.LogFields
		wantLog []logtest.Check
	}{
		{
			name: "Message",
			wantLog: []logtest.Check{
				{
					Validators: []logtest.EntryPredicate{
						logtest.HasLevel(logrus.DebugLevel),
						logtest.HasMessage("Message"),
					},
				},
			},
		},
		{
			name: "MessageWithFields",
			fields: watermill.LogFields{
				"Key": "Value",
			},
			wantLog: []logtest.Check{
				{
					Validators: []logtest.EntryPredicate{
						logtest.HasLevel(logrus.DebugLevel),
						logtest.HasMessage("MessageWithFields"),
						logtest.HasFieldValue("Key", "Value"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := logtest.RecordLogging()

			// Create logger
			l := log.NewLogger(tt.name)
			w := &WatermillLoggerAdapter{
				logger: l,
			}

			// Send message to logger
			w.Debug(tt.name, tt.fields)

			// Validate the message
			for _, check := range tt.wantLog {
				testhelpers.ReportErrors(t, "Log", check.Check(r))
			}
		})
	}
}

func TestWatermillLoggerAdapter_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		fields  watermill.LogFields
		wantLog []logtest.Check
	}{
		{
			name: "Message",
			err: errors.New("Message"),
			wantLog: []logtest.Check{
				{
					Validators: []logtest.EntryPredicate{
						logtest.HasLevel(logrus.ErrorLevel),
						logtest.HasMessage("Message"),
						logtest.HasError("Message"),
					},
				},
			},
		},
		{
			name: "MessageWithFields",
			err: errors.New("MessageWithFields"),
			fields: watermill.LogFields{
				"Key": "Value",
			},
			wantLog: []logtest.Check{
				{
					Validators: []logtest.EntryPredicate{
						logtest.HasLevel(logrus.ErrorLevel),
						logtest.HasMessage("MessageWithFields"),
						logtest.HasFieldValue("Key", "Value"),
						logtest.HasError("MessageWithFields"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := logtest.RecordLogging()

			// Create logger
			l := log.NewLogger(tt.name)
			w := &WatermillLoggerAdapter{
				logger: l,
			}

			// Send message to logger
			w.Error(tt.name, tt.err, tt.fields)

			// Validate the message
			for _, check := range tt.wantLog {
				testhelpers.ReportErrors(t, "Log", check.Check(r))
			}
		})
	}}

func TestWatermillLoggerAdapter_Trace(t *testing.T) {
	tests := []struct {
		name    string
		fields  watermill.LogFields
		wantLog []logtest.Check
	}{
		{
			name: "Message",
			wantLog: []logtest.Check{
				{
					Validators: []logtest.EntryPredicate{
						logtest.HasLevel(logrus.TraceLevel),
						logtest.HasMessage("Message"),
					},
				},
			},
		},
		{
			name: "MessageWithFields",
			fields: watermill.LogFields{
				"Key": "Value",
			},
			wantLog: []logtest.Check{
				{
					Validators: []logtest.EntryPredicate{
						logtest.HasLevel(logrus.TraceLevel),
						logtest.HasMessage("MessageWithFields"),
						logtest.HasFieldValue("Key", "Value"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := logtest.RecordLogging()

			// Create logger
			l := log.NewLogger(tt.name)
			w := &WatermillLoggerAdapter{
				logger: l,
			}

			// Send message to logger
			w.Trace(tt.name, tt.fields)

			// Validate the message
			for _, check := range tt.wantLog {
				testhelpers.ReportErrors(t, "Log", check.Check(r))
			}
		})
	}
}

func TestWatermillLoggerAdapter_Info(t *testing.T) {
	tests := []struct {
		name    string
		fields  watermill.LogFields
		wantLog []logtest.Check
	}{
		{
			name: "Message",
			wantLog: []logtest.Check{
				{
					Validators: []logtest.EntryPredicate{
						logtest.HasLevel(logrus.InfoLevel),
						logtest.HasMessage("Message"),
					},
				},
			},
		},
		{
			name: "MessageWithFields",
			fields: watermill.LogFields{
				"Key": "Value",
			},
			wantLog: []logtest.Check{
				{
					Validators: []logtest.EntryPredicate{
						logtest.HasLevel(logrus.InfoLevel),
						logtest.HasMessage("MessageWithFields"),
						logtest.HasFieldValue("Key", "Value"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := logtest.RecordLogging()

			// Create logger
			l := log.NewLogger(tt.name)
			w := &WatermillLoggerAdapter{
				logger: l,
			}

			// Send message to logger
			w.Info(tt.name, tt.fields)

			// Validate the message
			for _, check := range tt.wantLog {
				testhelpers.ReportErrors(t, "Log", check.Check(r))
			}
		})
	}
}

func TestWatermillLoggerAdapter_With(t *testing.T) {
	tests := []struct {
		name    string
		with    watermill.LogFields
		fields  watermill.LogFields
		wantLog []logtest.Check
	}{
		{
			name: "Message",
			with: watermill.LogFields{
				"KeyWith": "ValueWith",
			},
			wantLog: []logtest.Check{
				{
					Validators: []logtest.EntryPredicate{
						logtest.HasLevel(logrus.InfoLevel),
						logtest.HasMessage("Message"),
						logtest.HasFieldValue("KeyWith", "ValueWith"),
					},
				},
			},
		},
		{
			name: "MessageWithFields",
			with: watermill.LogFields{
				"KeyWith": "ValueWith",
			},
			fields: watermill.LogFields{
				"Key": "Value",
			},
			wantLog: []logtest.Check{
				{
					Validators: []logtest.EntryPredicate{
						logtest.HasLevel(logrus.InfoLevel),
						logtest.HasMessage("MessageWithFields"),
						logtest.HasFieldValue("KeyWith", "ValueWith"),
						logtest.HasFieldValue("Key", "Value"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := logtest.RecordLogging()

			// Create logger
			l := log.NewLogger(tt.name)
			w := &WatermillLoggerAdapter{
				logger: l,
			}

			// Send message to logger
			w.With(tt.with).Info(tt.name, tt.fields)

			// Validate the message
			for _, check := range tt.wantLog {
				testhelpers.ReportErrors(t, "Log", check.Check(r))
			}
		})
	}}
