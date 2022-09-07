// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package streamops

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/security/service"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type TestMessageSubscriberBuilderDependencies struct {
	Ctx               context.Context
	Channel           *Channel
	ChannelSubscriber *ChannelSubscriber
	Name              string
	Documentor        *TestMessageSubscriberDocumentor
	DocType           string
	Handler           func(context.Context) error
}

func NewTestMessageSubscriberBuilderDependencies(t *testing.T) TestMessageSubscriberBuilderDependencies {
	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
		"spring.application.name": "TestMessageSubscriberBuilder",
	})

	channel, err := NewChannel(ctx, "MY_TOPIC")
	assert.NoError(t, err)

	cs, err := NewChannelSubscriber(ctx, channel, "my-topic-receive", types.OptionalEmpty[string]())
	assert.NoError(t, err)

	doc := new(TestMessageSubscriberDocumentor)

	return TestMessageSubscriberBuilderDependencies{
		Ctx:               ctx,
		Channel:           channel,
		ChannelSubscriber: cs,
		Name:              "my-topic-message-subscriber",
		Documentor:        doc,
		DocType:           "test",
		Handler:           func(ctx context.Context) error { return nil },
	}
}

type TestMessageSubscriberDocumentor struct{}

func (t TestMessageSubscriberDocumentor) DocType() string {
	return "test"
}

func (t TestMessageSubscriberDocumentor) Document(i *MessageSubscriber) error {
	panic("implement me")
}

func TestMessageSubscriberBuilder_WithDecorator(t *testing.T) {
	deps := NewTestMessageSubscriberBuilderDependencies(t)

	deco := service.DefaultServiceAccount

	builder, err := NewMessageSubscriberBuilder(deps.Ctx, deps.ChannelSubscriber, deps.Name)
	assert.NoError(t, err)
	builder.WithDecorator(deco)

	assert.Len(t, builder.Filters, 1)
}

func TestMessageSubscriberBuilder_WithDocumentor(t *testing.T) {
	deps := NewTestMessageSubscriberBuilderDependencies(t)

	builder, err := NewMessageSubscriberBuilder(deps.Ctx, deps.ChannelSubscriber, deps.Name)
	assert.NoError(t, err)
	builder.WithDocumentor(deps.Documentor)

	assert.Len(t, builder.Documentors, 1)
	assert.True(t,
		reflect.DeepEqual(builder.Documentors[0], deps.Documentor),
		testhelpers.Diff(builder.Documentors[0], deps.Documentor))
}

func TestMessageSubscriberBuilder_WithFilter(t *testing.T) {
	deps := NewTestMessageSubscriberBuilderDependencies(t)

	deco := service.DefaultServiceAccount

	builder, err := NewMessageSubscriberBuilder(deps.Ctx, deps.ChannelSubscriber, deps.Name)
	assert.NoError(t, err)
	builder.WithFilter(types.NewOrderedDecorator(10, deco))

	assert.Len(t, builder.Filters, 1)
}

func TestMessageSubscriberBuilder_WithHandler(t *testing.T) {
	deps := NewTestMessageSubscriberBuilderDependencies(t)

	builder, err := NewMessageSubscriberBuilder(deps.Ctx, deps.ChannelSubscriber, deps.Name)
	assert.NoError(t, err)
	builder.WithHandler(deps.Handler)

	assert.NotNil(t, builder.Handler)
}

func TestMessageSubscriberBuilder_WithInputs(t *testing.T) {
	deps := NewTestMessageSubscriberBuilderDependencies(t)

	type inputs struct {
		EventType string `inp:"header"`
		Request   struct {
			Id   types.UUID `json:"id"`
			Name string     `json:"name"`
		} `inp:"body"`
	}

	builder, err := NewMessageSubscriberBuilder(deps.Ctx, deps.ChannelSubscriber, deps.Name)
	assert.NoError(t, err)
	builder.WithInputs(inputs{})

	assert.NotNil(t, builder.Inputs)
	assert.Equal(t, inputs{}, builder.Inputs)
}

func TestMessageSubscriberBuilder_WithMetadataFilterValues(t *testing.T) {
	deps := NewTestMessageSubscriberBuilderDependencies(t)

	builder, err := NewMessageSubscriberBuilder(deps.Ctx, deps.ChannelSubscriber, deps.Name)
	assert.NoError(t, err)
	builder.WithMetadataFilterValues("eventType", "up", "down")

	assert.NotNil(t, builder.MetadataFilterValues)
	assert.Len(t, builder.MetadataFilterValues, 1)
	assert.Len(t, builder.MetadataFilterValues["eventType"], 2)
	assert.Equal(t, builder.MetadataFilterValues["eventType"], []string{"up", "down"})
}

func TestMessageSubscriberBuilder_Build(t *testing.T) {
	deps := NewTestMessageSubscriberBuilderDependencies(t)

	type inputs struct {
		EventType string `in:"header"`
		Request   struct {
			Id   types.UUID `json:"id"`
			Name string     `json:"name"`
		} `in:"body"`
	}

	builder, err := NewMessageSubscriberBuilder(deps.Ctx, deps.ChannelSubscriber, deps.Name)
	assert.NoError(t, err)

	ms, err := builder.
		WithMetadataFilterValues("eventType", "up", "down").
		WithDocumentor(deps.Documentor).
		WithInputs(inputs{}).
		WithHandler(func(context.Context, inputs) error { return nil }).
		WithDecorator(service.DefaultServiceAccount).
		Build()
	assert.NoError(t, err)
	assert.NotNil(t, ms)
	assert.Equal(t, ms.channelSubscriber, deps.ChannelSubscriber)
	assert.Len(t, ms.documentors, 1)
	assert.Equal(t, ms.documentors[0], deps.Documentor)
	assert.Len(t, ms.filters, 1)
	assert.Equal(t, ms.name, deps.Name)
	assert.NotNil(t, ms.inputPort)
}

type TestMessageSubscriberDependencies struct {
	TestMessageSubscriberBuilderDependencies
}

func NewTestMessageSubscriberDependencies(t *testing.T) TestMessageSubscriberDependencies {
	builderDeps := NewTestMessageSubscriberBuilderDependencies(t)
	return TestMessageSubscriberDependencies{
		TestMessageSubscriberBuilderDependencies: builderDeps,
	}
}

func TestMessageSubscriber_Name(t *testing.T) {
	deps := NewTestMessageSubscriberDependencies(t)

	builder, err := NewMessageSubscriberBuilder(deps.Ctx, deps.ChannelSubscriber, deps.Name)
	assert.NoError(t, err)

	ms, err := builder.
		WithHandler(deps.Handler).
		Build()
	assert.NoError(t, err)
	assert.Equal(t, deps.Name, ms.Name())
}

func TestMessageSubscriber_Channel(t *testing.T) {
	deps := NewTestMessageSubscriberDependencies(t)

	builder, err := NewMessageSubscriberBuilder(deps.Ctx, deps.ChannelSubscriber, deps.Name)
	assert.NoError(t, err)

	ms, err := builder.
		WithHandler(deps.Handler).
		Build()
	assert.NoError(t, err)
	assert.Equal(t, deps.Channel, ms.Channel())
}

func TestMessageSubscriber_InputPort(t *testing.T) {
	deps := NewTestMessageSubscriberDependencies(t)

	builder, err := NewMessageSubscriberBuilder(deps.Ctx, deps.ChannelSubscriber, deps.Name)
	assert.NoError(t, err)

	ms, err := builder.
		WithHandler(deps.Handler).
		WithInputs(struct {
			Body []byte `in:"body"`
		}{}).
		Build()
	assert.NoError(t, err)
	assert.NotNil(t, ms.InputPort())
}

func TestMessageSubscriber_ContentType(t *testing.T) {
	deps := NewTestMessageSubscriberDependencies(t)

	builder, err := NewMessageSubscriberBuilder(deps.Ctx, deps.ChannelSubscriber, deps.Name)
	assert.NoError(t, err)

	ms, err := builder.
		WithHandler(deps.Handler).
		Build()
	assert.NoError(t, err)
	assert.Equal(t, deps.Channel.DefaultContentType(), ms.ContentType())
}

func TestMessageSubscriber_Documentor(t *testing.T) {
	deps := NewTestMessageSubscriberDependencies(t)

	builder, err := NewMessageSubscriberBuilder(deps.Ctx, deps.ChannelSubscriber, deps.Name)
	assert.NoError(t, err)

	ms, err := builder.
		WithHandler(deps.Handler).
		WithDocumentor(deps.Documentor).
		Build()
	assert.NoError(t, err)

	optionalDoc := ops.DocumentorWithType[MessageSubscriber](ms, "test")
	assert.True(t, optionalDoc.IsPresent())

	doc := optionalDoc.Value()
	assert.Equal(t, deps.Documentor, doc)
}

func TestMessageSubscriber_MetadataFilterValues(t *testing.T) {
	tests := []struct {
		name    string
		inputs  interface{}
		want    []string
		wantErr bool
	}{
		{
			name: "InputPortConstFilter",
			inputs: struct {
				Colour string `in:"header,const=blue"`
				Body   []byte `in:"body"`
			}{},
			want: []string{"blue"},
		},
		{
			name: "InputPortEnumFilter",
			inputs: struct {
				Colour string `in:"header" enum:"red,orange,green"`
				Body   []byte `in:"body"`
			}{},
			want: []string{"red", "orange", "green"},
		},
		{
			name: "InputPortFailure",
			inputs: struct {
				Colour string `in:"header"`
				Body   []byte `in:"body"`
			}{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := NewTestMessageSubscriberDependencies(t)

			builder, err := NewMessageSubscriberBuilder(deps.Ctx, deps.ChannelSubscriber, deps.Name)
			assert.NoError(t, err)

			builder.WithHandler(deps.Handler)

			ms, err := builder.
				WithInputs(tt.inputs).
				Build()
			assert.NoError(t, err)

			mfv, err := ms.MetadataFilterValues("colour")
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.True(t,
					reflect.DeepEqual(mfv, tt.want),
					testhelpers.Diff(mfv, tt.want))
			}
		})
	}
}

func TestMessageSubscriber_OnMessage(t *testing.T) {
	deps := NewTestMessageSubscriberDependencies(t)

	type inputs struct {
		Colour  string `in:"header" enum:"red,orange,green"`
		Request struct {
			A string `json:"a"`
		} `in:"body"`
	}

	builder, err := NewMessageSubscriberBuilder(deps.Ctx, deps.ChannelSubscriber, deps.Name)
	assert.NoError(t, err)

	handler := func(ctx context.Context, inp inputs, msg *message.Message, channel *Channel) error {
		assert.Equal(t, "red", inp.Colour)
		assert.Equal(t, "value", inp.Request.A)
		return errors.New("some error")
	}

	ms, err := builder.
		WithHandler(handler).
		WithInputs(inputs{}).
		Build()
	assert.NoError(t, err)

	msg := message.NewMessage(
		types.MustNewUUID().String(),
		[]byte(`{"a":"value"}`))
	msg.Metadata["colour"] = "red"
	msg.SetContext(deps.Ctx)

	err = ms.OnMessage(msg)
	assert.ErrorContains(t, err, "some error")
}

func TestRegisterMessageSubscriber(t *testing.T) {
	deps := NewTestMessageSubscriberDependencies(t)

	builder, err := NewMessageSubscriberBuilder(deps.Ctx, deps.ChannelSubscriber, deps.Name)
	assert.NoError(t, err)

	ms, err := builder.
		WithHandler(deps.Handler).
		Build()
	assert.NoError(t, err)

	registeredMessageSubscribers = []*MessageSubscriber{}
	RegisterMessageSubscriber(ms)
	assert.Equal(t, registeredMessageSubscribers[0], ms)
}

func TestRegisteredMessageSubscribers(t *testing.T) {
	deps := NewTestMessageSubscriberDependencies(t)

	builder, err := NewMessageSubscriberBuilder(deps.Ctx, deps.ChannelSubscriber, deps.Name)
	assert.NoError(t, err)

	ms, err := builder.
		WithHandler(deps.Handler).
		Build()
	assert.NoError(t, err)

	registeredMessageSubscribers = []*MessageSubscriber{}
	RegisterMessageSubscriber(ms)
	subscribers := RegisteredMessageSubscribers(deps.Channel.Name())
	assert.Len(t, subscribers, 1)
	assert.Equal(t, ms, subscribers[0])
}
