// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package streamops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewWatermillMessageEncoder(t *testing.T) {
	sink := new(MessageDataSink)
	enc := NewWatermillMessageEncoder(sink)
	assert.Equal(t, sink, enc.Sink)
}

func TestWatermillMessageEncoder_EncodeBody(t *testing.T) {
	type args struct {
		body interface{}
		mime string
		enc  ops.Encoding
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantSink *MessageDataSink
	}{
		{
			name: "JSON",
			args: args{
				body: map[string]string{
					"key": "value",
				},
				mime: "application/json",
			},
			wantErr: false,
			wantSink: &MessageDataSink{
				Metadata: map[string]string{
					"contentType": "application/json",
				},
				Payload: []byte(`{"key":"value"}` + "\n"),
			},
		},
		{
			name: "XML",
			args: args{
				body: struct {
					XMLName xml.Name `xml:"entry"`
					Key     string   `xml:"key"`
				}{
					Key: "Value",
				},
				mime: "application/xml",
			},
			wantErr: false,
			wantSink: &MessageDataSink{
				Metadata: map[string]string{
					"contentType": "application/xml",
				},
				Payload: []byte(`<entry><key>Value</key></entry>`),
			},
		},
		{
			name: "Gzip",
			args: args{
				body: map[string]string{
					"key": "value",
				},
				mime: "application/json",
				enc:  []string{"gzip"},
			},
			wantErr: false,
			wantSink: &MessageDataSink{
				Metadata: map[string]string{
					"contentType":     "application/json",
					"contentEncoding": "gzip",
				},
				Payload: []byte{
					0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xaa, 0x56, 0xca, 0x4e, 0xad, 0x54,
					0xb2, 0x52, 0x2a, 0x4b, 0xcc, 0x29, 0x4d, 0x55, 0xaa, 0xe5, 0x02, 0x04, 0x00, 0x00, 0xff, 0xff,
					0x12, 0xb0, 0x36, 0x69, 0x10, 0x00, 0x00, 0x00,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder := WatermillMessageEncoder{
				Sink: new(MessageDataSink),
			}
			err := encoder.EncodeBody(tt.args.body, tt.args.mime, tt.args.enc)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.True(t,
					reflect.DeepEqual(tt.wantSink, encoder.Sink),
					testhelpers.Diff(tt.wantSink, encoder.Sink))
			}
		})
	}
}

func TestWatermillMessageEncoder_EncodeChannelPrimitive(t *testing.T) {
	tests := []struct {
		name     string
		channel  types.Optional[string]
		wantErr  bool
		wantSink *MessageDataSink
	}{
		{
			name:    "Channel",
			channel: types.OptionalOf("channel"),
			wantSink: &MessageDataSink{
				Channel: types.OptionalOf("channel"),
			},
		},
		{
			name:     "NoChannel",
			channel:  types.OptionalEmpty[string](),
			wantSink: &MessageDataSink{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder := WatermillMessageEncoder{
				Sink: new(MessageDataSink),
			}
			err := encoder.EncodeChannelPrimitive(tt.channel)
			assert.Equal(t,
				tt.wantErr,
				err != nil)
			if !tt.wantErr {
				assert.True(t,
					reflect.DeepEqual(tt.wantSink, encoder.Sink),
					testhelpers.Diff(tt.wantSink, encoder.Sink))
			}
		})
	}
}

func TestWatermillMessageEncoder_EncodeHeaderPrimitive(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    types.Optional[string]
		wantErr  bool
		wantSink *MessageDataSink
	}{
		{
			name:  "Header",
			key:   "headerName",
			value: types.OptionalOf("headerValue"),
			wantSink: &MessageDataSink{
				Metadata: map[string]string{
					"headerName": "headerValue",
				},
			},
		},
		{
			name:     "NoHeaderValue",
			key:      "headerName",
			value:    types.OptionalEmpty[string](),
			wantSink: &MessageDataSink{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder := WatermillMessageEncoder{
				Sink: new(MessageDataSink),
			}
			err := encoder.EncodeHeaderPrimitive(tt.key, tt.value)
			assert.Equal(t,
				tt.wantErr,
				err != nil)
			if !tt.wantErr {
				assert.True(t,
					reflect.DeepEqual(tt.wantSink, encoder.Sink),
					testhelpers.Diff(tt.wantSink, encoder.Sink))
			}
		})
	}
}

func TestWatermillMessageEncoder_EncodeMessageIdPrimitive(t *testing.T) {
	tests := []struct {
		name      string
		messageId types.Optional[string]
		wantErr   bool
		wantSink  *MessageDataSink
	}{
		{
			name:      "MessageId",
			messageId: types.OptionalOf("messageId"),
			wantSink: &MessageDataSink{
				MessageId: types.OptionalOf("messageId"),
			},
		},
		{
			name:      "NoMessageId",
			messageId: types.OptionalEmpty[string](),
			wantSink:  &MessageDataSink{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder := WatermillMessageEncoder{
				Sink: new(MessageDataSink),
			}
			err := encoder.EncodeMessageIdPrimitive(tt.messageId)
			assert.Equal(t,
				tt.wantErr,
				err != nil)
			if !tt.wantErr {
				assert.True(t,
					reflect.DeepEqual(tt.wantSink, encoder.Sink),
					testhelpers.Diff(tt.wantSink, encoder.Sink))
			}
		})
	}
}
