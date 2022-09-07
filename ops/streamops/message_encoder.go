// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package streamops

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"strings"
)

type MessageEncoder interface {
	EncodeChannelPrimitive(channel types.Optional[string]) (err error)
	EncodeHeaderPrimitive(name string, value types.Optional[string]) (err error)
	EncodeMessageIdPrimitive(messageId types.Optional[string]) (err error)
	EncodeBody(body interface{}, mime string, enc ops.Encoding) (err error)
}

type WatermillMessageEncoder struct {
	Sink *MessageDataSink
}

func (a WatermillMessageEncoder) EncodeChannelPrimitive(channel types.Optional[string]) (err error) {
	if channel.IsPresent() {
		a.Sink.WithChannel(channel.Value())
	}
	return nil
}

func (a WatermillMessageEncoder) EncodeHeaderPrimitive(name string, value types.Optional[string]) (err error) {
	if value.IsPresent() {
		a.Sink.WithMetadataItem(name, value.Value())
	}
	return nil
}

func (a WatermillMessageEncoder) EncodeMessageIdPrimitive(messageId types.Optional[string]) (err error) {
	if messageId.IsPresent() {
		a.Sink.WithMessageId(messageId.Value())
	}
	return nil
}

func (a WatermillMessageEncoder) EncodeBody(body interface{}, mime string, enc ops.Encoding) error {
	if body == nil {
		a.Sink.WithPayload(nil)
		return nil
	}

	payloadBuffer := types.CloseableByteBuffer{Buffer: new(bytes.Buffer)}

	contentOptions := ops.NewContentOptions(mime)
	contentOptions.WithEncoding(enc...)
	err := contentOptions.WriteEntity(payloadBuffer, body)
	if err != nil {
		return errors.Wrap(err, "Failed to encode body")
	}

	a.Sink.WithPayload(payloadBuffer.Bytes())

	if mime != "" {
		a.Sink.WithMetadataItem(PeerNameContentType, mime)
	}
	if len(enc) > 0 {
		a.Sink.WithMetadataItem(PeerNameContentEncoding, strings.Join(enc, ","))
	}
	return nil
}

func NewWatermillMessageEncoder(sink *MessageDataSink) WatermillMessageEncoder {
	return WatermillMessageEncoder{
		Sink: sink,
	}
}
