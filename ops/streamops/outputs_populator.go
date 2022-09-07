// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package streamops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type OutputsPopulator struct {
	Outputs         *interface{}
	OutputPort      *ops.Port
	Channel         types.Optional[string]
	ContentType     string
	ContentEncoding string
	Encoder         MessageEncoder
}

func (p *OutputsPopulator) PopulateOutputs() error {
	if p.Outputs == nil {
		return errors.New("Outputs field not set")
	}

	if err := p.populateOutputMessageId(); err != nil {
		return errors.Wrap(err, "Failed to set message id")
	}

	if err := p.populateOutputChannel(); err != nil {
		return errors.Wrap(err, "Failed to set channel name")
	}

	if err := p.populateOutputMetadata(); err != nil {
		return errors.Wrap(err, "Failed to set message metadata")
	}

	if err := p.populateOutputPayload(); err != nil {
		return errors.Wrap(err, "Failed to set message payload")
	}

	return nil
}

func (p *OutputsPopulator) extractor(pf *ops.PortField) ops.PortFieldExtractor {
	return ops.NewPortFieldExtractor(pf, *p.Outputs)
}

func (p *OutputsPopulator) populateOutputChannel() error {
	channelPortField := p.OutputPort.Fields.First(
		ops.PortFieldHasGroup(FieldGroupStreamChannel))

	if channelPortField == nil {
		return nil
	}

	channel := p.Channel

	value, err := p.extractor(channelPortField).ExtractPrimitive()
	if err != nil {
		return err
	}

	if value.IsPresent() {
		channel = value
	}

	return p.Encoder.EncodeChannelPrimitive(channel)
}

func (p *OutputsPopulator) populateOutputMessageId() error {
	messageIdPortField := p.OutputPort.Fields.First(
		ops.PortFieldHasGroup(FieldGroupStreamMessageId))

	if messageIdPortField == nil {
		return nil
	}

	value, err := p.extractor(messageIdPortField).ExtractPrimitive()
	if err != nil {
		return err
	}

	if !value.IsPresent() {
		return nil
	}

	return p.Encoder.EncodeMessageIdPrimitive(value)
}

func (p *OutputsPopulator) populateOutputMetadata() error {
	headerPortFields := p.OutputPort.Fields.All(
		ops.PortFieldHasGroup(FieldGroupStreamHeader))

	for _, headerPortField := range headerPortFields {
		// Calculate header value
		value, err := p.extractor(headerPortField).ExtractPrimitive()
		if err != nil {
			return err
		}

		// Apply header value
		if err = p.Encoder.EncodeHeaderPrimitive(headerPortField.Peer, value); err != nil {
			return err
		}

		if value.IsPresent() {
			// Store contentType and contentEncoding
			switch headerPortField.Peer {
			case PeerNameContentType:
				p.ContentType = value.Value()

			case PeerNameContentEncoding:
				p.ContentEncoding = value.Value()
			}
		}
	}

	return nil
}

func (p *OutputsPopulator) populateOutputPayload() (err error) {
	bodyPortField := p.OutputPort.Fields.First(
		ops.PortFieldHasGroup(FieldGroupStreamBody))

	bodyValue, err := p.extractor(bodyPortField).ExtractValue()
	if err != nil {
		return errors.Wrap(err, "Failed to extract body from outputs")
	}

	if bodyValue.Kind() == reflect.Invalid {
		// Optional body
		return nil
	}

	contentType := p.ContentType
	if contentType == "" {
		return errors.New("Content-Type not specified for output population")
	}

	encoding := ops.Encoding{}
	if p.ContentEncoding != "" {
		encoding = append(encoding, strings.Split(p.ContentEncoding, ",")...)
	}

	if err = p.Encoder.EncodeBody(bodyValue.Interface(), contentType, encoding); err != nil {
		return errors.Wrap(err, "Failed to encode/marshal body")
	}

	return nil
}
