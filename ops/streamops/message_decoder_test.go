package streamops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewMessageDecoder(t *testing.T) {
	type args struct {
		source             MessageDataSource
		defaultContentType string
		defaultEncoding    string
	}
	tests := []struct {
		name string
		args args
		want WatermillMessageInputDecoder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewMessageDecoder(tt.args.source, tt.args.defaultContentType, tt.args.defaultEncoding), "NewMessageDecoder(%v, %v, %v)", tt.args.source, tt.args.defaultContentType, tt.args.defaultEncoding)
		})
	}
}

func TestWatermillMessageInputDecoder_DecodeContent(t *testing.T) {
	type fields struct {
		defaultContentType string
		defaultEncoding    string
		source             MessageDataSource
	}
	type args struct {
		pf *ops.PortField
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult ops.Content
		wantErr    bool
	}{
		{
			name: "Success",
			fields: fields{
				defaultContentType: "default-content-type",
				defaultEncoding:    "",
				source: MessageDataSource{
					Channel: "channel",
					Message: func() *message.Message {
						return message.NewMessage(types.MustNewUUID().String(), []byte("data"))
					}(),
				},
			},
			args: args{
				&ops.PortField{
					Name:    "Body",
					Indices: []int{0},
					Peer:    "body",
					Group:   FieldGroupStreamBody,
					Type: ops.PortFieldType{
						Shape:        ops.FieldShapeContent,
						Type:         reflect.TypeOf(ops.Content{}),
						Indirections: 1,
						HandlerType:  reflect.TypeOf(ops.Content{}),
					},
				},
			},
			wantResult: ops.NewContentFromBytes(
				ops.ContentOptions{
					MimeType: "default-content-type",
					Encoding: nil,
				},
				[]byte("data")),
		},
		{
			name: "WrongGroup",
			fields: fields{
				defaultContentType: "default-content-type",
				defaultEncoding:    "",
				source:             MessageDataSource{},
			},
			args: args{
				&ops.PortField{
					Group: FieldGroupStreamMessageId,
				},
			},
			wantErr: true,
		},
		{
			name: "CustomContentType",
			fields: fields{
				defaultContentType: "default-content-type",
				defaultEncoding:    "",
				source: MessageDataSource{
					Channel: "channel",
					Message: func() *message.Message {
						m := message.NewMessage(types.MustNewUUID().String(), []byte("data"))
						m.Metadata["contentType"] = "custom-content-type"
						return m
					}(),
				},
			},
			args: args{
				&ops.PortField{
					Name:    "Body",
					Indices: []int{0},
					Peer:    "body",
					Group:   FieldGroupStreamBody,
					Type: ops.PortFieldType{
						Shape:        ops.FieldShapeContent,
						Type:         reflect.TypeOf(ops.Content{}),
						Indirections: 1,
						HandlerType:  reflect.TypeOf(ops.Content{}),
					},
				},
			},
			wantResult: ops.NewContentFromBytes(
				ops.ContentOptions{
					MimeType: "custom-content-type",
					Encoding: nil,
				},
				[]byte("data")),
		},
		{
			name: "CustomEncoding",
			fields: fields{
				defaultContentType: "default-content-type",
				defaultEncoding:    "",
				source: MessageDataSource{
					Channel: "channel",
					Message: func() *message.Message {
						m := message.NewMessage(types.MustNewUUID().String(), []byte("data"))
						m.Metadata["contentEncoding"] = "gzip"
						return m
					}(),
				},
			},
			args: args{
				&ops.PortField{
					Name:    "Body",
					Indices: []int{0},
					Peer:    "body",
					Group:   FieldGroupStreamBody,
					Type: ops.PortFieldType{
						Shape:        ops.FieldShapeContent,
						Type:         reflect.TypeOf(ops.Content{}),
						Indirections: 1,
						HandlerType:  reflect.TypeOf(ops.Content{}),
					},
				},
			},
			wantResult: ops.NewContentFromBytes(
				ops.ContentOptions{
					MimeType: "default-content-type",
					Encoding: []string{"gzip"},
				},
				[]byte("data")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WatermillMessageInputDecoder{
				defaultContentType: tt.fields.defaultContentType,
				defaultEncoding:    tt.fields.defaultEncoding,
				source:             tt.fields.source,
			}
			gotResult, err := w.DecodeContent(tt.args.pf)
			assert.Equal(t, tt.wantErr, err != nil, "Unexpected result")
			if !tt.wantErr {
				assert.True(t,
					reflect.DeepEqual(tt.wantResult, gotResult),
					testhelpers.Diff(tt.wantResult, gotResult))
			}
		})
	}
}

func TestWatermillMessageInputDecoder_DecodePrimitive(t *testing.T) {
	msg := message.NewMessage(types.MustNewUUID().String(), []byte("data"))
	msg.Metadata["key"] = "value"

	type fields struct {
		defaultContentType string
		defaultEncoding    string
		source             MessageDataSource
	}
	type args struct {
		pf *ops.PortField
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult types.Optional[string]
		wantErr    bool
	}{
		{
			name: "Header",
			fields: fields{
				defaultContentType: "",
				defaultEncoding:    "",
				source: MessageDataSource{
					Channel: "channel",
					Message: msg,
				},
			},
			args: args{
				&ops.PortField{
					Name:    "Key",
					Indices: []int{0},
					Peer:    "key",
					Group:   FieldGroupStreamHeader,
					Type: ops.PortFieldType{
						Shape:        ops.FieldShapePrimitive,
						Type:         reflect.TypeOf(""),
						Indirections: 0,
						HandlerType:  reflect.TypeOf(""),
					},
				},
			},
			wantResult: types.OptionalOf("value"),
		},
		{
			name: "HeaderMissing",
			fields: fields{
				defaultContentType: "",
				defaultEncoding:    "",
				source: MessageDataSource{
					Channel: "channel",
					Message: msg,
				},
			},
			args: args{
				&ops.PortField{
					Name:    "Key2",
					Indices: []int{0},
					Peer:    "key2",
					Group:   FieldGroupStreamHeader,
					Type: ops.PortFieldType{
						Shape:        ops.FieldShapePrimitive,
						Type:         reflect.TypeOf(""),
						Indirections: 0,
						HandlerType:  reflect.TypeOf(""),
					},
				},
			},
			wantResult: types.OptionalEmpty[string](),
		},
		{
			name: "HeaderDefault",
			fields: fields{
				defaultContentType: "",
				defaultEncoding:    "",
				source: MessageDataSource{
					Channel: "channel",
					Message: msg,
				},
			},
			args: args{
				&ops.PortField{
					Name:    "Key2",
					Indices: []int{0},
					Peer:    "key2",
					Group:   FieldGroupStreamHeader,
					Type: ops.PortFieldType{
						Shape:        ops.FieldShapePrimitive,
						Type:         reflect.TypeOf(""),
						Indirections: 0,
						HandlerType:  reflect.TypeOf(""),
					},
					Options: map[string]string{
						"default": "default-key2-value",
					},
				},
			},
			wantResult: types.OptionalOf("default-key2-value"),
		},
		{
			name: "HeaderDefaultMissing",
			fields: fields{
				defaultContentType: "",
				defaultEncoding:    "",
				source: MessageDataSource{
					Channel: "channel",
					Message: msg,
				},
			},
			args: args{
				&ops.PortField{
					Name:    "Key2",
					Indices: []int{0},
					Peer:    "key2",
					Group:   FieldGroupStreamHeader,
					Type: ops.PortFieldType{
						Shape:        ops.FieldShapePrimitive,
						Type:         reflect.TypeOf(""),
						Indirections: 0,
						HandlerType:  reflect.TypeOf(""),
					},
					Options: map[string]string{},
				},
			},
			wantResult: types.OptionalEmpty[string](),
		},
		{
			name: "Channel",
			fields: fields{
				defaultContentType: "",
				defaultEncoding:    "",
				source: MessageDataSource{
					Channel: "channel",
					Message: msg,
				},
			},
			args: args{
				&ops.PortField{
					Name:    "ChannelName",
					Indices: []int{0},
					Peer:    "channel",
					Group:   FieldGroupStreamChannel,
					Type: ops.PortFieldType{
						Shape:        ops.FieldShapePrimitive,
						Type:         reflect.TypeOf(""),
						Indirections: 0,
						HandlerType:  reflect.TypeOf(""),
					},
				},
			},
			wantResult: types.OptionalOf("channel"),
		},
		{
			name: "MessageId",
			fields: fields{
				defaultContentType: "",
				defaultEncoding:    "",
				source: MessageDataSource{
					Channel: "channel",
					Message: msg,
				},
			},
			args: args{
				&ops.PortField{
					Name:    "MessageId",
					Indices: []int{0},
					Peer:    "messageId",
					Group:   FieldGroupStreamMessageId,
					Type: ops.PortFieldType{
						Shape:        ops.FieldShapePrimitive,
						Type:         reflect.TypeOf(""),
						Indirections: 0,
						HandlerType:  reflect.TypeOf(""),
					},
				},
			},
			wantResult: types.OptionalOf(msg.UUID),
		},
		{
			name: "InvalidGroup)",
			fields: fields{
				defaultContentType: "",
				defaultEncoding:    "",
				source: MessageDataSource{
					Channel: "channel",
					Message: msg,
				},
			},
			args: args{
				&ops.PortField{
					Name:    "MessageId",
					Indices: []int{0},
					Peer:    "messageId",
					Group:   FieldGroupStreamBody,
					Type: ops.PortFieldType{
						Shape:        ops.FieldShapePrimitive,
						Type:         reflect.TypeOf(""),
						Indirections: 0,
						HandlerType:  reflect.TypeOf(""),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WatermillMessageInputDecoder{
				defaultContentType: tt.fields.defaultContentType,
				defaultEncoding:    tt.fields.defaultEncoding,
				source:             tt.fields.source,
			}
			gotResult, err := w.DecodePrimitive(tt.args.pf)
			assert.Equal(t, tt.wantErr, err != nil, "Unexpected result")
			if !tt.wantErr {
				assert.True(t,
					reflect.DeepEqual(tt.wantResult, gotResult),
					testhelpers.Diff(tt.wantResult, gotResult))
			}
		})
	}
}
