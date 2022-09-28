package ${async.channel.package}

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

const contextKeyChannelSubscriber = contextKeyNamed("ChannelSubscriber")

func ContextChannelSubscriber() types.ContextKeyAccessor[*streamops.ChannelSubscriber] {
	return types.NewContextKeyAccessor[*streamops.ChannelSubscriber](contextKeyChannelSubscriber)
}

func newChannelSubscriber(ctx context.Context) (channelSubscriber *streamops.ChannelSubscriber, err error) {
	channelSubscriber = ContextChannelSubscriber().Get(ctx)
	if channelSubscriber == nil {
		doc := new(asyncapi.ChannelSubscriberDocumentor).
			WithOperation(new(asyncapi.Operation).
				WithID("${async.operation.id}").
				WithSummary("${async.operation.id}"))

		channelSubscriber, err = streamops.NewChannelSubscriber(ctx,
			channel,
			"${async.operation.id}",
			types.OptionalOf("eventType"))
		if err != nil {
			return nil, err
		}

		channelSubscriber.AddDocumentor(doc)
	}

	return channelSubscriber, nil
}
