package streamops

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type TestChannelPublisherDocumentor struct{}

func (t TestChannelPublisherDocumentor) DocType() string {
	return "test"
}

func (t TestChannelPublisherDocumentor) Document(i *ChannelPublisher) error {
	//TODO implement me
	panic("implement me")
}

func TestChannelPublisher_Channel(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
		"spring.application.name": "TestChannelSubscriber_Channel",
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
			p, err := NewChannelPublisher(ctx, tt.fields.channel, tt.fields.name)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, p.Channel(), "Channel()")
		})
	}
}

func TestChannelPublisher_Documentor(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
		"spring.application.name": "TestChannelSubscriber_Channel",
	})
	channel, err := NewChannel(ctx, "MY_TOPIC")
	assert.NoError(t, err)

	documentor := new(TestChannelPublisherDocumentor)

	type fields struct {
		name        string
		documentors ops.Documentors[ChannelPublisher]
	}
	tests := []struct {
		name    string
		fields  fields
		docType string
		want    ops.Documentor[ChannelPublisher]
	}{
		{
			name: "Found",
			fields: fields{
				name:        "my-publisher",
				documentors: ops.Documentors[ChannelPublisher]{documentor},
			},
			docType: "test",
			want:    documentor,
		},
		{
			name: "NotFound",
			fields: fields{
				name:        "my-publisher",
				documentors: ops.Documentors[ChannelPublisher]{},
			},
			docType: "test",
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewChannelPublisher(ctx, channel, tt.fields.name)
			assert.NoError(t, err)
			p.AddDocumentor(tt.fields.documentors...)

			result := ops.DocumentorWithType[ChannelPublisher](p, tt.docType)
			assert.Equal(t, tt.want != nil, result.IsPresent())
			if tt.want != nil && result.IsPresent() {
				assert.Equalf(t, tt.want, result.Value(), "Documentor()")
			}
		})
	}
}

func TestChannelPublisher_Name(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
		"spring.application.name": "TestChannelPublisher_Channel",
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
			p, err := NewChannelPublisher(ctx, channel, tt.fields.name)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, p.Name(), "Channel()")
		})
	}
}

func TestChannelPublisher_Publish(t *testing.T) {
	type fields struct {
		channel          *Channel
		name             string
		publisherService stream.PublisherService
		documentors      ops.Documentors[ChannelPublisher]
	}
	type args struct {
		ctx      context.Context
		payload  []byte
		metadata map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ChannelPublisher{
				channel:          tt.fields.channel,
				name:             tt.fields.name,
				publisherService: tt.fields.publisherService,
				documentors:      tt.fields.documentors,
			}
			tt.wantErr(t, p.Publish(tt.args.ctx, tt.args.payload, tt.args.metadata), fmt.Sprintf("Publish(%v, %v, %v)", tt.args.ctx, tt.args.payload, tt.args.metadata))
		})
	}
}

func TestNewChannelPublisher(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
		"spring.application.name": "TestChannelSubscriber_Channel",
	})
	channel, err := NewChannel(ctx, "MY_TOPIC")
	assert.NoError(t, err)

	type args struct {
		ctx         context.Context
		channel     *Channel
		operationId string
	}
	tests := []struct {
		name    string
		args    args
		want    *ChannelPublisher
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				ctx:         ctx,
				channel:     channel,
				operationId: "operation-id",
			},
			want: &ChannelPublisher{
				channel: channel,
				name:    "operation-id",
			},
			wantErr: false,
		},
		{
			name: "Failure",
			args: args{
				ctx:         ctx,
				channel:     nil,
				operationId: "operation=id",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewChannelPublisher(tt.args.ctx, tt.args.channel, tt.args.operationId)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr && err == nil {
				assert.True(t,
					reflect.DeepEqual(tt.want.channel, got.channel),
					testhelpers.Diff(tt.want.channel, got.channel))
				assert.True(t,
					reflect.DeepEqual(tt.want.name, got.name),
					testhelpers.Diff(tt.want.name, got.name))
				assert.True(t,
					reflect.DeepEqual(tt.want.documentors, got.documentors),
					testhelpers.Diff(tt.want.documentors, got.documentors))
			}
		})
	}
}

func TestRegisterChannelPublisher(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
		"spring.application.name": "TestChannelPublisher_Channel",
	})
	channel, err := NewChannel(ctx, "MY_TOPIC")
	assert.NoError(t, err)

	tests := []struct {
		name      string
		publisher *ChannelPublisher
	}{
		{
			name: "Success",
			publisher: &ChannelPublisher{
				channel: channel,
				name:    "my-publisher",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registeredChannelPublishers = map[string]*ChannelPublisher{}
			RegisterChannelPublisher(tt.publisher)
			assert.Equal(t, registeredChannelPublishers[tt.publisher.Channel().Name()], tt.publisher)
		})
	}
}

func TestRegisteredChannelPublisher(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
		"spring.application.name": "TestChannelPublisher_Channel",
	})
	channel, err := NewChannel(ctx, "MY_TOPIC")
	assert.NoError(t, err)

	tests := []struct {
		name      string
		publisher *ChannelPublisher
	}{
		{
			name: "Success",
			publisher: &ChannelPublisher{
				channel: channel,
				name:    "my-publisher",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registeredChannelPublishers = map[string]*ChannelPublisher{}
			RegisterChannelPublisher(tt.publisher)
			got := RegisteredChannelPublisher(tt.publisher.Channel().Name())
			assert.Equal(t, tt.publisher, got)
		})
	}
}
