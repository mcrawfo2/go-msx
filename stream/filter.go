package stream

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/ThreeDotsLabs/watermill/message"
)

type MessageFilter func(ctx context.Context, metadata message.Metadata) bool

func FilterByMetaData(key string, values ...string) MessageFilter {
	return func(ctx context.Context, metadata message.Metadata) bool {
		metaDataValue := metadata[key]
		if !types.StringStack(values).Contains(metaDataValue) {
			logger.WithContext(ctx).Infof("Filtering out message with %s %q", key, metaDataValue)
			return false
		}

		return true
	}
}

func FilterMessage(msg *message.Message, filters []MessageFilter) bool {
	for _, filter := range filters {
		if !filter(msg.Context(), msg.Metadata) {
			return false
		}
	}
	return true
}
