// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package stream

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/ThreeDotsLabs/watermill/message"
)

// MessageFilter returns true if the message should be processed
type MessageFilter func(ctx context.Context, metadata message.Metadata) bool

// FilterByMetaData returns a new MessageFilter that allows messages
// which have any of the specified values in the specified Metadata key
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

// FilterMessage returns true if the message is allowed by all of the specified filters
func FilterMessage(msg *message.Message, filters []MessageFilter) bool {
	for _, filter := range filters {
		if !filter(msg.Context(), msg.Metadata) {
			return false
		}
	}
	return true
}
