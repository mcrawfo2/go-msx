// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package streamops

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/js"
	"cto-github.cisco.com/NFV-BU/go-msx/security/service"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

type TestMessagePublisherBuilderDependencies struct {
	Ctx              context.Context
	Channel          *Channel
	ChannelPublisher *ChannelPublisher
	PublisherService *stream.MockPublisherService
	Name             string
	Documentor       *TestMessagePublisherDocumentor
	DocType          string
	Outputs          interface{}
}

func NewTestMessagePublisherBuilderDependencies(t *testing.T) TestMessagePublisherBuilderDependencies {
	RegisterPortFieldValidationSchemaFunc(func(field *ops.PortField) (schema js.ValidationSchema, err error) {
		return js.ValidationSchema{}, nil
	})

	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
		"spring.application.name": "TestMessagePublisherBuilder",
	})

	channel, err := NewChannel(ctx, "MY_TOPIC")
	assert.NoError(t, err)

	mockPublisherService := new(stream.MockPublisherService)
	ctx = stream.ContextWithPublisherService(ctx, mockPublisherService)

	cs, err := NewChannelPublisher(ctx, channel, "my-topic-receive")
	assert.NoError(t, err)

	doc := new(TestMessagePublisherDocumentor)

	outputs := struct {
		SomeHeader string `out:"header"`
		Body       struct {
			Key string `json:"key"`
		} `out:"body"`
	}{}

	return TestMessagePublisherBuilderDependencies{
		Ctx:              ctx,
		Channel:          channel,
		ChannelPublisher: cs,
		PublisherService: mockPublisherService,
		Name:             "my-topic-message-publisher",
		Documentor:       doc,
		DocType:          "test",
		Outputs:          outputs,
	}
}

type TestMessagePublisherDocumentor struct{}

func (t TestMessagePublisherDocumentor) DocType() string {
	return "test"
}

func (t TestMessagePublisherDocumentor) Document(source *MessagePublisher) error {
	panic("implement me")
}

func TestMessagePublisherBuilder_WithDecorator(t *testing.T) {
	deps := NewTestMessagePublisherBuilderDependencies(t)

	deco := service.DefaultServiceAccount

	builder, err := NewMessagePublisherBuilder(deps.Ctx, deps.ChannelPublisher, deps.Name, deps.Outputs)
	assert.NoError(t, err)
	builder.WithDecorator(deco)

	assert.Len(t, builder.Filters, 1)
}

func TestMessagePublisherBuilder_WithDocumentor(t *testing.T) {
	deps := NewTestMessagePublisherBuilderDependencies(t)

	builder, err := NewMessagePublisherBuilder(deps.Ctx, deps.ChannelPublisher, deps.Name, deps.Outputs)
	assert.NoError(t, err)
	builder.WithDocumentor(deps.Documentor)

	assert.Len(t, builder.Documentors, 1)
	assert.True(t,
		reflect.DeepEqual(builder.Documentors[0], deps.Documentor),
		testhelpers.Diff(builder.Documentors[0], deps.Documentor))
}

func TestMessagePublisherBuilder_WithFilter(t *testing.T) {
	deps := NewTestMessagePublisherBuilderDependencies(t)

	deco := service.DefaultServiceAccount

	builder, err := NewMessagePublisherBuilder(deps.Ctx, deps.ChannelPublisher, deps.Name, deps.Outputs)
	assert.NoError(t, err)
	builder.WithFilter(types.NewOrderedDecorator(10, deco))

	assert.Len(t, builder.Filters, 1)
}

func TestMessagePublisherBuilder_Build(t *testing.T) {
	deps := NewTestMessagePublisherBuilderDependencies(t)

	type outputs struct {
		EventType string `out:"header"`
		Request   struct {
			Id   types.UUID `json:"id"`
			Name string     `json:"name"`
		} `out:"body"`
	}

	builder, err := NewMessagePublisherBuilder(deps.Ctx, deps.ChannelPublisher, deps.Name, outputs{})
	assert.NoError(t, err)

	mp, err := builder.
		WithDocumentor(deps.Documentor).
		WithDecorator(service.DefaultServiceAccount).
		Build()
	assert.NoError(t, err)
	assert.NotNil(t, mp)
	assert.Equal(t, mp.channelPublisher, deps.ChannelPublisher)
	assert.Len(t, mp.documentors, 1)
	assert.Equal(t, mp.documentors[0], deps.Documentor)
	assert.Len(t, mp.filters, 1)
	assert.Equal(t, mp.name, deps.Name)
	assert.NotNil(t, mp.outputPort)
}

type TestMessagePublisherDependencies struct {
	TestMessagePublisherBuilderDependencies
}

func NewTestMessagePublisherDependencies(t *testing.T) TestMessagePublisherDependencies {
	builderDeps := NewTestMessagePublisherBuilderDependencies(t)
	return TestMessagePublisherDependencies{
		TestMessagePublisherBuilderDependencies: builderDeps,
	}
}

func TestMessagePublisher_Name(t *testing.T) {
	deps := NewTestMessagePublisherDependencies(t)

	builder, err := NewMessagePublisherBuilder(deps.Ctx, deps.ChannelPublisher, deps.Name, deps.Outputs)
	assert.NoError(t, err)

	ms, err := builder.Build()
	assert.NoError(t, err)
	assert.Equal(t, deps.Name, ms.Name())
}

func TestMessagePublisher_Channel(t *testing.T) {
	deps := NewTestMessagePublisherDependencies(t)

	builder, err := NewMessagePublisherBuilder(deps.Ctx, deps.ChannelPublisher, deps.Name, deps.Outputs)
	assert.NoError(t, err)

	ms, err := builder.Build()
	assert.NoError(t, err)
	assert.Equal(t, deps.Channel, ms.Channel())
}

func TestMessagePublisher_OutputPort(t *testing.T) {
	deps := NewTestMessagePublisherDependencies(t)

	builder, err := NewMessagePublisherBuilder(deps.Ctx, deps.ChannelPublisher, deps.Name, deps.Outputs)
	assert.NoError(t, err)

	ms, err := builder.
		Build()
	assert.NoError(t, err)
	assert.NotNil(t, ms.OutputPort())
}

func TestMessagePublisher_ContentType(t *testing.T) {
	deps := NewTestMessagePublisherDependencies(t)

	builder, err := NewMessagePublisherBuilder(deps.Ctx, deps.ChannelPublisher, deps.Name, deps.Outputs)
	assert.NoError(t, err)

	ms, err := builder.Build()
	assert.NoError(t, err)
	assert.Equal(t, deps.Channel.DefaultContentType(), ms.ContentType())
}

func TestMessagePublisher_Documentor(t *testing.T) {
	deps := NewTestMessagePublisherDependencies(t)

	builder, err := NewMessagePublisherBuilder(deps.Ctx, deps.ChannelPublisher, deps.Name, deps.Outputs)
	assert.NoError(t, err)

	ms, err := builder.
		WithDocumentor(deps.Documentor).
		Build()
	assert.NoError(t, err)

	optionalDoc := ops.DocumentorWithType[MessagePublisher](ms, "test")
	assert.True(t, optionalDoc.IsPresent())

	doc := optionalDoc.Value()
	assert.Equal(t, deps.Documentor, doc)
}

func TestMessagePublisher_Publish(t *testing.T) {
	deps := NewTestMessagePublisherDependencies(t)
	outputs := struct {
		Payload struct {
			Key string `json:"key"`
		} `out:"body"`
	}{}
	outputs.Payload.Key = "Value"

	builder, err := NewMessagePublisherBuilder(deps.Ctx, deps.ChannelPublisher, deps.Name, outputs)
	assert.NoError(t, err)

	ms, err := builder.
		WithDocumentor(deps.Documentor).
		Build()
	assert.NoError(t, err)

	deps.PublisherService.
		On("Publish",
			mock.AnythingOfType("*context.valueCtx"),
			deps.Channel.Name(),
			mock.AnythingOfType("[]uint8"),
			mock.AnythingOfType("map[string]string")).
		Run(func(args mock.Arguments) {
			assert.Contains(t, args[3], "contentType")
			assert.Equal(t, []byte(`{"key":"Value"}`+"\n"), args[2])
		}).
		Return(errors.New("some error"))

	err = ms.Publish(deps.Ctx, outputs)
	assert.ErrorContains(t, err, "some error")
}

func TestRegisterMessagePublisher(t *testing.T) {
	deps := NewTestMessagePublisherDependencies(t)

	builder, err := NewMessagePublisherBuilder(deps.Ctx, deps.ChannelPublisher, deps.Name, deps.Outputs)
	assert.NoError(t, err)
	if err != nil {
		assert.FailNow(t, "Test preconditions not met")
	}

	ms, err := builder.Build()
	assert.NoError(t, err)
	if err != nil {
		assert.FailNow(t, "Test preconditions not met")
	}

	registeredMessagePublishers = []*MessagePublisher{}
	RegisterMessagePublisher(ms)
	assert.Equal(t, registeredMessagePublishers[0], ms)
}

func TestRegisteredMessagePublishers(t *testing.T) {
	deps := NewTestMessagePublisherDependencies(t)

	builder, err := NewMessagePublisherBuilder(deps.Ctx, deps.ChannelPublisher, deps.Name, deps.Outputs)
	assert.NoError(t, err)

	ms, err := builder.Build()
	assert.NoError(t, err)

	registeredMessagePublishers = []*MessagePublisher{}
	RegisterMessagePublisher(ms)
	publishers := RegisteredMessagePublishers(deps.Channel.Name())
	assert.Len(t, publishers, 1)
	assert.Equal(t, ms, publishers[0])
}
