// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package stream

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddListener(t *testing.T) {
	dummyAction := func(msg *message.Message) error { return nil }

	type args struct {
		topic  string
		action ListenerAction
	}
	tests := []struct {
		name      string
		args      args
		listeners map[string][]ListenerAction
		wantErr   bool
	}{
		{
			name: "Success",
			args: args{
				topic:  "mock",
				action: dummyAction,
			},
			listeners: make(map[string][]ListenerAction),
			wantErr:   false,
		},
		{
			name: "NoTopic",
			args: args{
				topic:  "",
				action: dummyAction,
			},
			listeners: make(map[string][]ListenerAction),
			wantErr:   true,
		},
		{
			name: "NoAction",
			args: args{
				topic:  "mock",
				action: nil,
			},
			listeners: make(map[string][]ListenerAction),
			wantErr:   true,
		},
		{
			name: "Running",
			args: args{
				topic:  "mock",
				action: dummyAction,
			},
			listeners: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listeners = tt.listeners
			listenerCount := len(tt.listeners[tt.args.topic])
			if err := AddListener(tt.args.topic, tt.args.action); (err != nil) != tt.wantErr {
				t.Errorf("AddListener() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				assert.Len(t, listeners[tt.args.topic], listenerCount+1)
			}
		})
	}
}
