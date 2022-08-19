package streamops

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type TestChannelSubscriberDocumentor struct{}

func (t TestChannelSubscriberDocumentor) DocType() string {
	return "test"
}

func (t TestChannelSubscriberDocumentor) Document(i *ChannelSubscriber) error {
	//TODO implement me
	panic("implement me")
}

func TestChannelSubscriber_Channel(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
		"spring.application.name": "TestChannelSubscriber",
	})
	channel, err := NewChannel(ctx, "MY_TOPIC")
	assert.NoError(t, err)

	type fields struct {
		channel *Channel
		name    string
	}
	tests := []struct {
		name   string
		fields fields
		want   *Channel
	}{
		{
			name: "Success",
			fields: fields{
				channel: channel,
				name:    channel.Name(),
			},
			want: channel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewChannelSubscriber(ctx, tt.fields.channel, tt.fields.name, types.OptionalEmpty[string]())
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, p.Channel(), "Channel()")
		})
	}
}

func TestChannelSubscriber_OnMessage(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
		"spring.application.name": "TestChannelSubscriber",
	})
	channel, err := NewChannel(ctx, "MY_TOPIC")
	assert.NoError(t, err)

	var received string

	type fields struct {
		messageConsumers func(t *testing.T, cs *ChannelSubscriber)
		dispatchHeader   types.Optional[string]
	}
	type args struct {
		msg *message.Message
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantErr      bool
		wantReceived string
	}{
		{
			name: "SingleSubscriber",
			fields: fields{
				messageConsumers: func(t *testing.T, cs *ChannelSubscriber) {
					builder, err := NewMessageSubscriberBuilder(ctx, cs, "message-subscriber")
					assert.NoError(t, err)
					_, err = builder.
						WithHandler(func() { received = "ok" }).
						Build()
					assert.NoError(t, err)
				},
			},
			args: args{
				msg: message.NewMessage(
					types.MustNewUUID().String(),
					[]byte{}),
			},
			wantErr:      false,
			wantReceived: "ok",
		},
		{
			name: "SingleSubscriberFailure",
			fields: fields{
				messageConsumers: func(t *testing.T, cs *ChannelSubscriber) {
					builder, err := NewMessageSubscriberBuilder(ctx, cs, "message-subscriber")
					assert.NoError(t, err)
					_, err = builder.
						WithHandler(func() error { return errors.New("some error") }).
						Build()
					assert.NoError(t, err)
				},
			},
			args: args{
				msg: message.NewMessage(
					types.MustNewUUID().String(),
					[]byte{}),
			},
			wantErr: true,
		},
		{
			name: "SingleSubscriberRegistrationFailure",
			fields: fields{
				messageConsumers: func(t *testing.T, cs *ChannelSubscriber) {
					builder, err := NewMessageSubscriberBuilder(ctx, cs, "message-subscriber")
					assert.NoError(t, err)
					_, err = builder.
						WithHandler(func() error { return errors.New("some error") }).
						Build()
					assert.NoError(t, err)

					builder, err = NewMessageSubscriberBuilder(ctx, cs, "message-subscriber")
					assert.NoError(t, err)
					_, err = builder.
						WithHandler(func() {}).
						Build()
					assert.Error(t, err)
				},
			},
			args: args{
				msg: message.NewMessage(
					types.MustNewUUID().String(),
					[]byte{}),
			},
			wantErr: true,
		},
		{
			name: "MultiHeaderSubscriber",
			fields: fields{
				messageConsumers: func(t *testing.T, cs *ChannelSubscriber) {
					builder, err := NewMessageSubscriberBuilder(ctx, cs, "message-subscriber-1")
					assert.NoError(t, err)
					_, err = builder.
						WithMetadataFilterValues("number", "one", "three").
						WithHandler(func() { received = "odd" }).
						Build()
					assert.NoError(t, err)

					builder, err = NewMessageSubscriberBuilder(ctx, cs, "message-subscriber-2")
					assert.NoError(t, err)
					_, err = builder.
						WithMetadataFilterValues("number", "two", "four").
						WithHandler(func() { received = "even" }).
						Build()
					assert.NoError(t, err)
				},
				dispatchHeader: types.OptionalOf("number"),
			},
			args: args{
				msg: func() *message.Message {
					msg := message.NewMessage(
						types.MustNewUUID().String(),
						[]byte{})
					msg.Metadata.Set("number", "three")
					return msg
				}(),
			},
			wantErr:      false,
			wantReceived: "odd",
		},
		{
			name: "MultiHeaderSubscriberFailure",
			fields: fields{
				messageConsumers: func(t *testing.T, cs *ChannelSubscriber) {
					builder, err := NewMessageSubscriberBuilder(ctx, cs, "message-subscriber-1")
					assert.NoError(t, err)
					_, err = builder.
						WithMetadataFilterValues("number", "one", "three").
						WithHandler(func() { received = "odd" }).
						Build()
					assert.Error(t, err)
				},
				dispatchHeader: types.OptionalOf("eventType"),
			},
			args: args{
				msg: func() *message.Message {
					msg := message.NewMessage(
						types.MustNewUUID().String(),
						[]byte{})
					msg.Metadata.Set("eventType", "event")
					return msg
				}(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channelSubscriber, err := NewChannelSubscriber(ctx, channel, "channel-subscriber", tt.fields.dispatchHeader)
			assert.NoError(t, err)

			if tt.fields.messageConsumers != nil {
				tt.fields.messageConsumers(t, channelSubscriber)
			}

			received = ""
			gotErr := channelSubscriber.OnMessage(tt.args.msg)
			assert.Equal(t, tt.wantErr, gotErr != nil)
			if !tt.wantErr {
				assert.Equal(t, tt.wantReceived, received)
			}
		})
	}
}

func TestChannelSubscriber_Documentor(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
		"spring.application.name": "TestChannelSubscriber",
	})
	channel, err := NewChannel(ctx, "MY_TOPIC")
	assert.NoError(t, err)

	documentor := new(TestChannelSubscriberDocumentor)

	type fields struct {
		name        string
		documentors ops.Documentors[ChannelSubscriber]
	}
	tests := []struct {
		name    string
		fields  fields
		docType string
		want    ops.Documentor[ChannelSubscriber]
	}{
		{
			name: "Found",
			fields: fields{
				name:        "my-subscriber",
				documentors: ops.Documentors[ChannelSubscriber]{documentor},
			},
			docType: "test",
			want:    documentor,
		},
		{
			name: "NotFound",
			fields: fields{
				name:        "my-subscriber",
				documentors: ops.Documentors[ChannelSubscriber]{},
			},
			docType: "test",
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewChannelSubscriber(ctx, channel, tt.fields.name, types.OptionalEmpty[string]())
			assert.NoError(t, err)
			p.AddDocumentor(tt.fields.documentors...)

			result := ops.DocumentorWithType[ChannelSubscriber](p, tt.docType)
			assert.Equal(t, tt.want != nil, result.IsPresent())
			if tt.want != nil && result.IsPresent() {
				assert.Equalf(t, tt.want, result.Value(), "Documentor()")
			}
		})
	}
}

func TestChannelSubscriber_Name(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
		"spring.application.name": "TestChannelSubscriber",
	})
	channel, err := NewChannel(ctx, "MY_TOPIC")
	assert.NoError(t, err)

	type fields struct {
		channel *Channel
		name    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Success",
			fields: fields{
				channel: channel,
				name:    "my-subscriber-name",
			},
			want: "my-subscriber-name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewChannelSubscriber(ctx, channel, tt.fields.name, types.OptionalEmpty[string]())
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, p.Name(), "Channel()")
		})
	}
}

func TestNewChannelSubscriber(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
		"spring.application.name": "TestChannelSubscriber",
	})
	channel, err := NewChannel(ctx, "MY_TOPIC")
	assert.NoError(t, err)

	type args struct {
		ctx            context.Context
		channel        *Channel
		operationId    string
		dispatchHeader types.Optional[string]
	}
	tests := []struct {
		name    string
		args    args
		want    *ChannelSubscriber
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				ctx:            ctx,
				channel:        channel,
				operationId:    "operation-id",
				dispatchHeader: types.OptionalEmpty[string](),
			},
			want: &ChannelSubscriber{
				channel:        channel,
				name:           "operation-id",
				dispatchHeader: types.OptionalEmpty[string](),
				dispatchTable:  map[stream.MetadataHeader]stream.ListenerAction{},
			},
			wantErr: false,
		},
		{
			name: "Failure",
			args: args{
				ctx:            ctx,
				channel:        nil,
				operationId:    "operation=id",
				dispatchHeader: types.OptionalEmpty[string](),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewChannelSubscriber(tt.args.ctx, tt.args.channel, tt.args.operationId, tt.args.dispatchHeader)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr && err == nil {
				assert.True(t,
					reflect.DeepEqual(tt.want, got),
					testhelpers.Diff(tt.want, got))
			}
		})
	}
}

func TestRegisterChannelSubscriber(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
		"spring.application.name": "TestChannelSubscriber",
	})
	channel, err := NewChannel(ctx, "MY_TOPIC")
	assert.NoError(t, err)

	tests := []struct {
		name       string
		subscriber *ChannelSubscriber
	}{
		{
			name: "Success",
			subscriber: &ChannelSubscriber{
				channel:        channel,
				name:           "my-subscriber",
				dispatchHeader: types.OptionalEmpty[string](),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registeredChannelSubscribers = map[string]*ChannelSubscriber{}
			RegisterChannelSubscriber(tt.subscriber)
			assert.Equal(t, registeredChannelSubscribers[tt.subscriber.Channel().Name()], tt.subscriber)
		})
	}
}

func TestRegisteredChannelSubscriber(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
		"spring.application.name": "TestChannelSubscriber",
	})
	channel, err := NewChannel(ctx, "MY_TOPIC")
	assert.NoError(t, err)

	tests := []struct {
		name       string
		subscriber *ChannelSubscriber
	}{
		{
			name: "Success",
			subscriber: &ChannelSubscriber{
				channel:        channel,
				name:           "my-subscriber",
				dispatchHeader: types.OptionalEmpty[string](),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registeredChannelSubscribers = map[string]*ChannelSubscriber{}
			RegisterChannelSubscriber(tt.subscriber)
			got := RegisteredChannelSubscriber(tt.subscriber.Channel().Name())
			assert.Equal(t, tt.subscriber, got)
		})
	}
}
