package streamops

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMessageDataSink_WithChannel(t *testing.T) {
	const channelName = "channel"
	sink := new(MessageDataSink)
	got := sink.WithChannel(channelName)
	assert.True(t, got.Channel.IsPresent())
	assert.Equal(t, channelName, got.Channel.Value())
}

func TestMessageDataSink_Message(t *testing.T) {
	type fields struct {
		Channel  string
		Metadata map[string]string
		Payload  []byte
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *message.Message
	}{
		{
			name: "Simple",
			fields: fields{
				Channel: "channel",
				Metadata: map[string]string{
					"key": "value",
				},
				Payload: []byte("payload"),
			},
			args: args{
				context.Background(),
			},
			want: func() *message.Message {
				m := message.NewMessage("", []byte("payload"))
				m.Metadata["key"] = "value"
				m.SetContext(context.Background())
				return m
			}(),
		},
		{
			name: "NoMetadata",
			fields: fields{
				Channel:  "channel",
				Metadata: nil,
				Payload:  []byte("payload"),
			},
			want: func() *message.Message {
				m := message.NewMessage("", []byte("payload"))
				m.SetContext(context.Background())
				return m
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MessageDataSink{
				Channel:  types.OptionalOf(tt.fields.Channel),
				Metadata: tt.fields.Metadata,
				Payload:  tt.fields.Payload,
			}

			gotChannel, gotMessage := s.Message(tt.args.ctx)

			assert.True(t,
				reflect.DeepEqual(tt.fields.Channel, gotChannel.Value()),
				testhelpers.Diff(tt.fields.Channel, gotChannel.Value()))
			assert.True(t,
				reflect.DeepEqual(tt.want.Metadata, gotMessage.Metadata),
				testhelpers.Diff(tt.want.Metadata, gotMessage.Metadata))
			assert.True(t,
				reflect.DeepEqual(tt.want.Payload, gotMessage.Payload),
				testhelpers.Diff(tt.want.Payload, gotMessage.Payload))
			assert.True(t,
				reflect.DeepEqual(tt.want.Context(), gotMessage.Context()),
				testhelpers.Diff(tt.want.Context(), gotMessage.Context()))
		})
	}
}

func TestMessageDataSink_WithMetadataItem(t *testing.T) {
	tests := []struct {
		name     string
		metadata map[string]string
		key      string
		value    string
		want     *MessageDataSink
	}{
		{
			name:     "Add",
			metadata: map[string]string{},
			key:      "key",
			value:    "value",
			want: &MessageDataSink{
				Metadata: map[string]string{
					"key": "value",
				},
			},
		},
		{
			name: "Overwrite",
			metadata: map[string]string{
				"key": "oldvalue",
			},
			key:   "key",
			value: "value",
			want: &MessageDataSink{
				Metadata: map[string]string{
					"key": "value",
				},
			},
		},
		{
			name:     "NoMetadata",
			metadata: nil,
			key:      "key",
			value:    "value",
			want: &MessageDataSink{
				Metadata: map[string]string{
					"key": "value",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MessageDataSink{
				Metadata: tt.metadata,
			}
			got := s.WithMetadataItem(tt.key, tt.value)
			assert.True(t, reflect.DeepEqual(tt.want, got), testhelpers.Diff(tt.want, got))
		})
	}
}

func TestMessageDataSink_WithPayload(t *testing.T) {
	var payload = []byte("payload")
	sink := new(MessageDataSink)
	got := sink.WithPayload(payload)
	assert.Equal(t, payload, got.Payload)
}
