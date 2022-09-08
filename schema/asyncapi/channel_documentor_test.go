// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
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

type TestChannelDocumentorDeps struct {
	Ctx               context.Context
	Channel           *streamops.Channel
	ChannelPublisher  *streamops.ChannelPublisher
	ChannelSubscriber *streamops.ChannelSubscriber
}

func newTestChannelDocumentorDeps(t *testing.T) *TestChannelDocumentorDeps {
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

	channelSubscriber, err := streamops.NewChannelSubscriber(ctx, channel, "my-channel-publisher", types.OptionalEmpty[string]())
	if err != nil {
		assert.Failf(t, "unexpected error", "%+v", err)
	}

	return &TestChannelDocumentorDeps{
		Ctx:               ctx,
		Channel:           channel,
		ChannelPublisher:  channelPublisher,
		ChannelSubscriber: channelSubscriber,
	}
}

func TestChannelDocumentor_DocType(t *testing.T) {
	cd := new(ChannelDocumentor)
	assert.Equal(t, DocType, cd.DocType())
}

func TestChannelDocumentor_Document(t *testing.T) {
	deps := newTestChannelDocumentorDeps(t)

	type fields struct {
		skip        bool
		channelItem *ChannelItem
		mutator     ops.DocumentElementMutator[ChannelItem]
	}
	type args struct {
		c *streamops.Channel
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ChannelItem
		wantErr bool
	}{
		{
			name: "Skip",
			fields: fields{
				skip: true,
			},
			args: args{
				c: deps.Channel,
			},
			wantErr: false,
		},
		{
			name:   "NoChannelItem",
			fields: fields{},
			args: args{
				c: deps.Channel,
			},
			wantErr: false,
		},
		{
			name: "Mutator",
			fields: fields{
				mutator: func(channelItem *ChannelItem) {
					channelItem.Description = types.NewStringPtr("mutated")
				},
			},
			args: args{
				c: deps.Channel,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := ChannelDocumentor{
				skip:        tt.fields.skip,
				channelItem: tt.fields.channelItem,
				mutator:     tt.fields.mutator,
			}
			err := d.Document(tt.args.c)
			assert.Equal(t, tt.wantErr, err != nil)
			if err == nil {
				assert.True(t,
					reflect.DeepEqual(tt.want, d.channelItem),
					testhelpers.Diff(tt.want, d.channelItem))
			}
		})
	}
}

func TestChannelDocumentor_WithChannelItem(t *testing.T) {
	channelItem := &ChannelItem{
		MapOfAnything: map[string]interface{}{
			"key": "value",
		},
	}

	want := &ChannelDocumentor{
		channelItem: channelItem,
	}

	got := new(ChannelDocumentor).WithChannelItem(channelItem)
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}

func TestChannelDocumentor_WithChannelItemMutator(t *testing.T) {
	mutator := func(channelItem *ChannelItem) {}

	want := &ChannelDocumentor{
		mutator: mutator,
	}

	got := new(ChannelDocumentor).WithChannelItemMutator(mutator)
	assert.True(t,
		reflect.DeepEqual(
			fmt.Sprintf("%p", want.mutator),
			fmt.Sprintf("%p", got.mutator),
		),
		testhelpers.Diff(want, got))
}

func TestChannelDocumentor_WithSkip(t *testing.T) {
	skip := true

	want := &ChannelDocumentor{
		skip: skip,
	}

	got := new(ChannelDocumentor).WithSkip(skip)
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}
