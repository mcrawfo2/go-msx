package ${async.channel.package}

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

const contextKeyChannelPublisher = contextKeyNamed("ChannelPublisher")

func contextChannelPublisher() types.ContextKeyAccessor[*streamops.ChannelPublisher] {
	return types.NewContextKeyAccessor[*streamops.ChannelPublisher](contextKeyChannelPublisher)
}

func newChannelPublisher(ctx context.Context) (svc *streamops.ChannelPublisher, err error) {
	svc = contextChannelPublisher().Get(ctx)
	if svc == nil {
		var ch *streamops.Channel
		ch, err = newChannel(ctx)
		if err != nil {
			return nil, err
		}

		svc, err = streamops.NewChannelPublisher(ctx, ch, "${async.operation.id}")
		if err != nil {
			return nil, err
		}

		doc := new(asyncapi.ChannelPublisherDocumentor).
			WithOperation(new(asyncapi.Operation).
				WithID("${async.operation.id}").
				WithSummary("${async.operation.id}"))
		svc.AddDocumentor(doc)
	}

	return svc, nil
}
