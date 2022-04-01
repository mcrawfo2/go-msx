// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package stream

import (
	"errors"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
)

func TestNewMetadataDispatcher_Success(t *testing.T) {
	mockListenerAction := func(msg *message.Message) error {
		return nil
	}
	metadataHeaderActions := map[MetadataHeader]ListenerAction{
		"mock event type": mockListenerAction,
	}
	dispatcher, err := NewMetadataDispatcher("eventType", metadataHeaderActions)
	assert.NoError(t, err)
	assert.NotNil(t, dispatcher)
}
func TestNewMetadataDispatcher_EmptyMetadataHeaderName(t *testing.T) {
	mockListenerAction := func(msg *message.Message) error {
		return nil
	}
	metadataHeaderActions := map[MetadataHeader]ListenerAction{
		"mock event type": mockListenerAction,
	}
	dispatcher, err := NewMetadataDispatcher("", metadataHeaderActions)
	assert.Error(t, err)
	assert.Nil(t, dispatcher)
}

func TestNewMetadataDispatcher_EmptyMetadataHeaderActions(t *testing.T) {
	metadataHeaderActions := map[MetadataHeader]ListenerAction{}
	dispatcher, err := NewMetadataDispatcher("eventType", metadataHeaderActions)
	assert.Error(t, err)
	assert.Nil(t, dispatcher)
}

func TestDispatch(t *testing.T) {
	testMetadataHeaderName := "eventType"
	testMetadataHeader := MetadataHeader("test event type")
	tests := []struct {
		name                  string
		getMsg                func() *message.Message
		metadataHeaderActions map[MetadataHeader]ListenerAction
		wantErr               bool
	}{
		{
			name: "MetadataNotContainingHeader",
			getMsg: func() *message.Message {
				msg := message.NewMessage("", []byte{})
				return msg

			},
			metadataHeaderActions: map[MetadataHeader]ListenerAction{
				testMetadataHeader: func(msg *message.Message) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "IgnoredMetadataHeader",
			getMsg: func() *message.Message {
				msg := message.NewMessage("", []byte{})
				msg.Metadata[testMetadataHeaderName] = "ignored event type"
				return msg

			},
			metadataHeaderActions: map[MetadataHeader]ListenerAction{
				testMetadataHeader: func(msg *message.Message) error {
					return nil
				},
			},
			wantErr: false,
		},

		{
			name: "MessageProcessed_Success",
			getMsg: func() *message.Message {
				msg := message.NewMessage("", []byte{})
				msg.Metadata[testMetadataHeaderName] = string(testMetadataHeader)
				return msg

			},
			metadataHeaderActions: map[MetadataHeader]ListenerAction{
				testMetadataHeader: func(msg *message.Message) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "MessageProcessed_Failure",
			getMsg: func() *message.Message {
				msg := message.NewMessage("", []byte{})
				msg.Metadata[testMetadataHeaderName] = string(testMetadataHeader)
				return msg

			},
			metadataHeaderActions: map[MetadataHeader]ListenerAction{
				testMetadataHeader: func(msg *message.Message) error {
					return errors.New("Failed to process the message")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dispatcher, err := NewMetadataDispatcher(testMetadataHeaderName, tt.metadataHeaderActions)
			err = dispatcher.Dispatch(tt.getMsg())
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Dispatch() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if tt.wantErr {
					t.Errorf("Dispatch() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			return
		})
	}
}
