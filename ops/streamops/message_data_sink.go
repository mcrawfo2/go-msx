// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package streamops

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/ThreeDotsLabs/watermill/message"
)

type MessageDataSink struct {
	Channel   types.Optional[string]
	MessageId types.Optional[string]
	Metadata  map[string]string
	Payload   []byte
}

func (s *MessageDataSink) WithChannel(channel string) *MessageDataSink {
	s.Channel = types.OptionalOf(channel)
	return s
}

func (s *MessageDataSink) WithMetadataItem(key, value string) *MessageDataSink {
	if s.Metadata == nil {
		s.Metadata = make(map[string]string)
	}
	s.Metadata[key] = value
	return s
}

func (s *MessageDataSink) WithPayload(payload []byte) *MessageDataSink {
	s.Payload = payload
	return s
}

func (s *MessageDataSink) WithMessageId(messageId string) *MessageDataSink {
	s.MessageId = types.OptionalOf(messageId)
	return s
}

func (s *MessageDataSink) Message(ctx context.Context) (types.Optional[string], *message.Message) {
	msg := message.NewMessage(
		s.MessageId.OrElse(types.MustNewUUID().String()),
		s.Payload)

	if s.Metadata != nil {
		msg.Metadata = s.Metadata
	} else {
		msg.Metadata = make(map[string]string)
	}

	msg.SetContext(ctx)

	return s.Channel, msg
}
