package stream

import (
	"errors"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
)

func TestPanicRecovererSubscriberAction(t *testing.T) {
	dummySuccessAction := func(msg *message.Message) error { return nil }
	dummyErrorAction := func(msg *message.Message) error { return errors.New("test") }
	dummyPanicAction := func(msg *message.Message) error { panic("test") }
	dummyPanicErrorAction := func(msg *message.Message) error { panic(errors.New("test")) }

	tests := []struct {
		name    string
		action  ListenerAction
		wantErr bool
	}{
		{
			name:    "dummySuccessAction",
			action:  dummySuccessAction,
			wantErr: false,
		},
		{
			name:    "dummyErrorAction",
			action:  dummyErrorAction,
			wantErr: true,
		},
		{
			name:    "dummyPanicAction",
			action:  dummyPanicAction,
			wantErr: false,
		},
		{
			name:    "dummyPanicErrorAction",
			action:  dummyPanicErrorAction,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			panicRecovererAction := PanicRecovererActionInterceptor(nil, tt.action)
			msg := message.NewMessage("test", message.Payload{})
			if err := panicRecovererAction(msg); (err != nil) != tt.wantErr {
				t.Errorf("PanicRecovererSubscriberAction.Call() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
