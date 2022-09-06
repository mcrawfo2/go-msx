// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
	"cto-github.cisco.com/NFV-BU/go-msx/redis"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type TestMessageSubscriberDocumentorDeps struct {
	Ctx               context.Context
	Channel           *streamops.Channel
	ChannelSubscriber *streamops.ChannelSubscriber
	MessageSubscriber *streamops.MessageSubscriber
}

func newTestMessageSubscriberDocumentorDeps(t *testing.T) *TestMessageSubscriberDocumentorDeps {
	ctx := context.Background()
	ctx = configtest.ContextWithNewInMemoryConfig(ctx,
		map[string]string{
			"cli.flag.disconnected":                          "true",
			"spring.redis.enable":                            "true",
			"spring.cloud.stream.bindings.my-channel.binder": "redis",
			"spring.application.name":                        "TestChannelDocumentor",
		})
	_ = redis.ConfigurePool(ctx)

	channel, err := streamops.NewChannel(ctx, "my-channel")
	if err != nil {
		assert.Failf(t, "unexpected error", "%+v", err)
	}

	channelSubscriber, err := streamops.NewChannelSubscriber(ctx, channel, "my-channel-subscriber", types.OptionalEmpty[string]())
	if err != nil {
		assert.Failf(t, "unexpected error", "%+v", err)
	}

	type inputs struct {
		Body []byte `in:"body"`
	}

	mpb, err := streamops.NewMessageSubscriberBuilder(ctx, channelSubscriber, "my-message-subscriber")
	if err != nil {
		assert.Failf(t, "Test preconditions not met.", "%+v", err)
	}

	messageSubscriber, err := mpb.
		WithInputs(inputs{}).
		WithHandler(func(context.Context) {}).
		Build()
	if err != nil {
		assert.Failf(t, "Test preconditions not met.", "%+v", err)
	}

	return &TestMessageSubscriberDocumentorDeps{
		Ctx:               ctx,
		Channel:           channel,
		ChannelSubscriber: channelSubscriber,
		MessageSubscriber: messageSubscriber,
	}
}

func TestMessageSubscriberDocumentor_DocType(t *testing.T) {
	cd := new(MessageSubscriberDocumentor)
	assert.Equal(t, DocType, cd.DocType())
}

func TestMessageSubscriberDocumentor_Document(t *testing.T) {
	deps := newTestMessageSubscriberDocumentorDeps(t)

	type fields struct {
		skip    bool
		message *Message
		mutator MessageMutator
	}
	type args struct {
		subscriber *streamops.MessageSubscriber
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Message
		wantErr bool
	}{

		{
			name: "Skip",
			fields: fields{
				skip: true,
			},
			args: args{
				subscriber: deps.MessageSubscriber,
			},
			wantErr: false,
		},
		{
			name:   "NoMessage",
			fields: fields{},
			args: args{
				subscriber: deps.MessageSubscriber,
			},
			wantErr: false,
		},
		{
			name: "Mutator",
			fields: fields{
				mutator: func(message *Message) {
					message.Description = types.NewStringPtr("mutated")
				},
			},
			args: args{
				subscriber: deps.MessageSubscriber,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := MessageSubscriberDocumentor{
				skip:    tt.fields.skip,
				message: tt.fields.message,
				mutator: tt.fields.mutator,
			}
			err := d.Document(tt.args.subscriber)
			assert.Equal(t, tt.wantErr, err != nil)
			if err == nil {
				assert.True(t,
					reflect.DeepEqual(tt.want, d.message),
					testhelpers.Diff(tt.want, d.message))
			}
		})
	}
}

func TestMessageSubscriberDocumentor_WithMessage(t *testing.T) {
	message := &Message{
		MapOfAnything: map[string]interface{}{
			"key": "value",
		},
	}

	want := &MessageSubscriberDocumentor{
		message: message,
	}

	got := new(MessageSubscriberDocumentor).WithMessage(message)
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}

func TestMessageSubscriberDocumentor_WithMessageMutator(t *testing.T) {
	mutator := func(message *Message) {}

	want := &MessageSubscriberDocumentor{
		mutator: mutator,
	}

	got := new(MessageSubscriberDocumentor).WithMessageMutator(mutator)
	assert.True(t,
		reflect.DeepEqual(
			fmt.Sprintf("%p", want.mutator),
			fmt.Sprintf("%p", got.mutator),
		),
		testhelpers.Diff(want, got))
}

func TestMessageSubscriberDocumentor_WithSkip(t *testing.T) {
	skip := true

	want := &MessageSubscriberDocumentor{
		skip: skip,
	}

	got := new(MessageSubscriberDocumentor).WithSkip(skip)
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}
