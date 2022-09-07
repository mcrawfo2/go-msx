// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package streamops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMessageDataSource_ChannelName(t *testing.T) {
	const channelName = "channel"
	source := MessageDataSource{
		Channel: channelName,
	}
	got := source.ChannelName()
	assert.Equal(t, channelName, got)
}

func TestMessageDataSource_MessageId(t *testing.T) {
	const messageId = "message-id"
	source := MessageDataSource{
		Message: &message.Message{
			UUID: messageId,
		},
	}
	got := source.MessageId()
	assert.Equal(t, messageId, got)
}

func TestMessageDataSource_MetadataItem(t *testing.T) {
	tests := []struct {
		name    string
		message *message.Message
		key     string
		want    types.Optional[string]
	}{
		{
			name: "Exists",
			message: &message.Message{
				Metadata: map[string]string{
					"key": "value",
				},
			},
			key:  "key",
			want: types.OptionalOf("value"),
		},
		{
			name:    "NotExists",
			message: &message.Message{},
			key:     "key",
			want:    types.OptionalEmpty[string](),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MessageDataSource{
				Message: tt.message,
			}
			assert.Equalf(t,
				tt.want,
				s.MetadataItem(tt.key),
				"MetadataItem(%v)", tt.key)
		})
	}
}

func TestMessageDataSource_Payload(t *testing.T) {
	const payload = "payload"
	source := MessageDataSource{
		Message: &message.Message{
			Payload: []byte(payload),
		},
	}
	got := source.Payload()
	assert.NotNil(t, got)
	assert.Equal(t, payload, string(got))
}

func TestNewMessageDataSource(t *testing.T) {
	const channel = "channel"
	msg := &message.Message{}
	got := NewMessageDataSource(channel, msg)
	want := MessageDataSource{
		Channel: channel,
		Message: msg,
	}

	assert.True(t,
		reflect.DeepEqual(got, want),
		testhelpers.Diff(got, want))
}
