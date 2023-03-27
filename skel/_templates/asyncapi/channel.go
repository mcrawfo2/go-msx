//#id channelpackage ${async.channel.package}
package channelpackage

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

// Context

const contextKeyChannel = contextKeyNamed("Channel")

func ContextChannel() types.ContextKeyAccessor[*streamops.Channel] {
	return types.NewContextKeyAccessor[*streamops.Channel](contextKeyChannel)
}

// Constructor

func newChannel(ctx context.Context) (*streamops.Channel, error) {
	doc := new(asyncapi.ChannelDocumentor).
		WithChannelItem(new(asyncapi.ChannelItem).
			WithDescription("Description of the ${async.channel} channel."))

	ch, err := streamops.NewChannel(ctx, "${async.channel}")
	if err != nil {
		return nil, err
	}

	ch.WithDocumentor(doc)

	return ch, nil
}

// Singleton

var channel = types.NewSingleton(
	newChannel,
	ContextChannel)
