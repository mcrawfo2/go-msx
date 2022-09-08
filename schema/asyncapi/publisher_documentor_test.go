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

type TestChannelPublisherDocumentorDeps struct {
	Ctx              context.Context
	Channel          *streamops.Channel
	ChannelPublisher *streamops.ChannelPublisher
	MessagePublisher *streamops.MessagePublisher
}

func newTestChannelPublisherDocumentorDeps(t *testing.T) *TestChannelPublisherDocumentorDeps {
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

	return &TestChannelPublisherDocumentorDeps{
		Ctx:              ctx,
		Channel:          channel,
		ChannelPublisher: channelPublisher,
		MessagePublisher: messagePublisher,
	}
}

func TestChannelPublisherDocumentor_DocType(t *testing.T) {
	cd := new(ChannelPublisherDocumentor)
	assert.Equal(t, DocType, cd.DocType())
}

func TestChannelPublisherDocumentor_Document(t *testing.T) {
	deps := newTestChannelPublisherDocumentorDeps(t)

	type fields struct {
		skip      bool
		operation *Operation
		mutator   OperationMutator
	}
	type args struct {
		publisher *streamops.ChannelPublisher
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
				publisher: deps.ChannelPublisher,
			},
			wantErr: false,
		},
		{
			name:   "NoChannelItem",
			fields: fields{},
			args: args{
				publisher: deps.ChannelPublisher,
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
				publisher: deps.ChannelPublisher,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := ChannelPublisherDocumentor{
				skip:      tt.fields.skip,
				operation: tt.fields.operation,
				mutator:   tt.fields.mutator,
			}
			err := d.Document(tt.args.publisher)
			assert.Equal(t, tt.wantErr, err != nil)
			if err == nil {
				assert.True(t,
					reflect.DeepEqual(tt.want, d.operation),
					testhelpers.Diff(tt.want, d.operation))
			}
		})
	}
}

func TestChannelPublisherDocumentor_WithOperation(t *testing.T) {
	operation := &Operation{
		MapOfAnything: map[string]interface{}{
			"key": "value",
		},
	}

	want := &ChannelPublisherDocumentor{
		operation: operation,
	}

	got := new(ChannelPublisherDocumentor).WithOperation(operation)
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}

func TestChannelPublisherDocumentor_WithOperationMutator(t *testing.T) {
	mutator := func(operation *Operation) {}

	want := &ChannelPublisherDocumentor{
		mutator: mutator,
	}

	got := new(ChannelPublisherDocumentor).WithOperationMutator(mutator)
	assert.True(t,
		reflect.DeepEqual(
			fmt.Sprintf("%p", want.mutator),
			fmt.Sprintf("%p", got.mutator),
		),
		testhelpers.Diff(want, got))
}

func TestChannelPublisherDocumentor_WithSkip(t *testing.T) {
	skip := true

	want := &ChannelPublisherDocumentor{
		skip: skip,
	}

	got := new(ChannelPublisherDocumentor).WithSkip(skip)
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}
