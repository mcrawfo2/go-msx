package ${async.channel.package}

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

// Context

const contextKeyChannelSubscriber = contextKeyNamed("ChannelSubscriber")

func ContextChannelSubscriber() types.ContextKeyAccessor[*streamops.ChannelSubscriber] {
	return types.NewContextKeyAccessor[*streamops.ChannelSubscriber](contextKeyChannelSubscriber)
}

// Constructor

func newChannelSubscriber(ctx context.Context) (channelSubscriber *streamops.ChannelSubscriber, err error) {
	doc := new(asyncapi.ChannelSubscriberDocumentor).
		WithOperation(new(asyncapi.Operation).
			WithID("${async.operation.id}").
			WithSummary("${async.operation.id}"))

	ch, err := channel.Factory(ctx)
	if err != nil {
		return nil, err
	}

	channelSubscriber, err = streamops.NewChannelSubscriber(ctx,
		ch,
		"${async.operation.id}",
		types.OptionalOf("eventType"))
	if err != nil {
		return nil, err
	}

	channelSubscriber.AddDocumentor(doc)

	return channelSubscriber, nil
}

// Singleton

var channelSubscriber = types.NewSingleton(
	newChannelSubscriber,
	ContextChannelSubscriber,
)
