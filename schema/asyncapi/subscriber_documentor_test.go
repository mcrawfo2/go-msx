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

type TestChannelSubscriberDocumentorDeps struct {
	Ctx               context.Context
	Channel           *streamops.Channel
	ChannelSubscriber *streamops.ChannelSubscriber
	MessageSubscriber *streamops.MessageSubscriber
}

func newTestChannelSubscriberDocumentorDeps(t *testing.T) *TestChannelSubscriberDocumentorDeps {
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

	return &TestChannelSubscriberDocumentorDeps{
		Ctx:               ctx,
		Channel:           channel,
		ChannelSubscriber: channelSubscriber,
		MessageSubscriber: messageSubscriber,
	}
}

func TestChannelSubscriberDocumentor_DocType(t *testing.T) {
	cd := new(ChannelSubscriberDocumentor)
	assert.Equal(t, DocType, cd.DocType())
}

func TestChannelSubscriberDocumentor_Document(t *testing.T) {
	deps := newTestChannelSubscriberDocumentorDeps(t)

	type fields struct {
		skip      bool
		operation *Operation
		mutator   OperationMutator
	}
	type args struct {
		subscriber *streamops.ChannelSubscriber
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Operation
		wantErr bool
	}{

		{
			name: "Skip",
			fields: fields{
				skip: true,
			},
			args: args{
				subscriber: deps.ChannelSubscriber,
			},
			wantErr: false,
		},
		{
			name:   "NoChannelItem",
			fields: fields{},
			args: args{
				subscriber: deps.ChannelSubscriber,
			},
			wantErr: false,
		},
		{
			name: "Mutator",
			fields: fields{
				mutator: func(operation *Operation) {
					operation.Description = types.NewStringPtr("mutated")
				},
			},
			args: args{
				subscriber: deps.ChannelSubscriber,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := ChannelSubscriberDocumentor{
				skip:      tt.fields.skip,
				operation: tt.fields.operation,
				mutator:   tt.fields.mutator,
			}
			err := d.Document(tt.args.subscriber)
			assert.Equal(t, tt.wantErr, err != nil)
			if err == nil {
				assert.True(t,
					reflect.DeepEqual(tt.want, d.operation),
					testhelpers.Diff(tt.want, d.operation))
			}
		})
	}
}

func TestChannelSubscriberDocumentor_WithOperation(t *testing.T) {
	operation := &Operation{
		MapOfAnything: map[string]interface{}{
			"key": "value",
		},
	}

	want := &ChannelSubscriberDocumentor{
		operation: operation,
	}

	got := new(ChannelSubscriberDocumentor).WithOperation(operation)
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}

func TestChannelSubscriberDocumentor_WithOperationMutator(t *testing.T) {
	mutator := func(operation *Operation) {}

	want := &ChannelSubscriberDocumentor{
		mutator: mutator,
	}

	got := new(ChannelSubscriberDocumentor).WithOperationMutator(mutator)
	assert.True(t,
		reflect.DeepEqual(
			fmt.Sprintf("%p", want.mutator),
			fmt.Sprintf("%p", got.mutator),
		),
		testhelpers.Diff(want, got))
}

func TestChannelSubscriberDocumentor_WithSkip(t *testing.T) {
	skip := true

	want := &ChannelSubscriberDocumentor{
		skip: skip,
	}

	got := new(ChannelSubscriberDocumentor).WithSkip(skip)
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}

func Test_getChannelSubscriberDocumentor(t *testing.T) {
	deps := newTestChannelSubscriberDocumentorDeps(t)

	want := &ChannelSubscriberDocumentor{}
	deps.ChannelSubscriber.AddDocumentor(want)

	got := getChannelSubscriberDocumentor(deps.ChannelSubscriber)
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}
