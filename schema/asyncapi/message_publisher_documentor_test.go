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

type TestMessagePublisherDocumentorDeps struct {
	Ctx              context.Context
	Channel          *streamops.Channel
	ChannelPublisher *streamops.ChannelPublisher
	MessagePublisher *streamops.MessagePublisher
}

func newTestMessagePublisherDocumentorDeps(t *testing.T) *TestMessagePublisherDocumentorDeps {
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

	channelPublisher, err := streamops.NewChannelPublisher(ctx, channel, "my-channel-publisher")
	if err != nil {
		assert.Failf(t, "unexpected error", "%+v", err)
	}

	type outputs struct {
		Body []byte `out:"body"`
	}

	mpb, err := streamops.NewMessagePublisherBuilder(ctx, channelPublisher, "my-message-publisher", outputs{})
	if err != nil {
		assert.Failf(t, "Test preconditions not met.", "%+v", err)
	}

	messagePublisher, err := mpb.Build()
	if err != nil {
		assert.Failf(t, "Test preconditions not met.", "%+v", err)
	}

	return &TestMessagePublisherDocumentorDeps{
		Ctx:              ctx,
		Channel:          channel,
		ChannelPublisher: channelPublisher,
		MessagePublisher: messagePublisher,
	}
}

func TestMessagePublisherDocumentor_DocType(t *testing.T) {
	cd := new(MessagePublisherDocumentor)
	assert.Equal(t, DocType, cd.DocType())
}

func TestMessagePublisherDocumentor_Document(t *testing.T) {
	deps := newTestMessagePublisherDocumentorDeps(t)

	type fields struct {
		skip    bool
		message *Message
		mutator MessageMutator
	}
	type args struct {
		publisher *streamops.MessagePublisher
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
				publisher: deps.MessagePublisher,
			},
			wantErr: false,
		},
		{
			name:   "NoMessage",
			fields: fields{},
			args: args{
				publisher: deps.MessagePublisher,
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
				publisher: deps.MessagePublisher,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := MessagePublisherDocumentor{
				skip:    tt.fields.skip,
				message: tt.fields.message,
				mutator: tt.fields.mutator,
			}
			err := d.Document(tt.args.publisher)
			assert.Equal(t, tt.wantErr, err != nil)
			if err == nil {
				assert.True(t,
					reflect.DeepEqual(tt.want, d.message),
					testhelpers.Diff(tt.want, d.message))
			}
		})
	}
}

func TestMessagePublisherDocumentor_WithMessage(t *testing.T) {
	message := &Message{
		MapOfAnything: map[string]interface{}{
			"key": "value",
		},
	}

	want := &MessagePublisherDocumentor{
		message: message,
	}

	got := new(MessagePublisherDocumentor).WithMessage(message)
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}

func TestMessagePublisherDocumentor_WithMessageMutator(t *testing.T) {
	mutator := func(message *Message) {}

	want := &MessagePublisherDocumentor{
		mutator: mutator,
	}

	got := new(MessagePublisherDocumentor).WithMessageMutator(mutator)
	assert.True(t,
		reflect.DeepEqual(
			fmt.Sprintf("%p", want.mutator),
			fmt.Sprintf("%p", got.mutator),
		),
		testhelpers.Diff(want, got))
}

func TestMessagePublisherDocumentor_WithSkip(t *testing.T) {
	skip := true

	want := &MessagePublisherDocumentor{
		skip: skip,
	}

	got := new(MessagePublisherDocumentor).WithSkip(skip)
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}
