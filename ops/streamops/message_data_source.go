// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package streamops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/ThreeDotsLabs/watermill/message"
)

type MessageDataSource struct {
	Channel string
	Message *message.Message
}

func (s MessageDataSource) ChannelName() string {
	return s.Channel
}

func (s MessageDataSource) MessageId() string {
	return s.Message.UUID
}

func (s MessageDataSource) MetadataItem(key string) types.Optional[string] {
	if val, ok := s.Message.Metadata[key]; ok {
		return types.OptionalOf(val)
	}
	return types.OptionalEmpty[string]()
}

func (s MessageDataSource) Payload() []byte {
	return s.Message.Payload
}

func NewMessageDataSource(channel string, msg *message.Message) MessageDataSource {
	return MessageDataSource{
		Channel: channel,
		Message: msg,
	}
}
