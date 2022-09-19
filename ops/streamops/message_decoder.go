// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package streamops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"strings"
)

type MessageDecoder struct {
	defaultContentType string
	defaultEncoding    string
	source             MessageDataSource
}

func (w MessageDecoder) getBodyContentOptions() ops.ContentOptions {
	contentType := w.source.MetadataItem(PeerNameContentType).OrElse(w.defaultContentType)

	content := ops.NewContentOptions(contentType)

	if enc := w.source.MetadataItem(PeerNameContentEncoding).OrElse(w.defaultEncoding); enc != "" {
		content.WithEncoding(strings.Split(enc, ",")...)
	}

	return content
}

func (w MessageDecoder) DecodeContent(pf *ops.PortField) (result ops.Content, err error) {
	switch pf.Group {
	case FieldGroupStreamBody:
		bodyContentOptions := w.getBodyContentOptions()
		payload := w.source.Payload()
		result = ops.NewContentFromBytes(bodyContentOptions, payload)

	default:
		err = errors.Errorf("Cannot retrieve %q value from field group %q", pf.Type.Shape, pf.Group)
	}

	return
}

func (w MessageDecoder) decodePrimitiveDefault(pf *ops.PortField) (result types.Optional[string]) {
	defaultValue := pf.Default()
	if defaultValue != nil {
		return types.OptionalOf(defaultValue.(string))
	}

	return types.OptionalEmpty[string]()
}

func (w MessageDecoder) DecodePrimitive(pf *ops.PortField) (result types.Optional[string], err error) {
	switch pf.Group {
	case FieldGroupStreamHeader:
		value := w.source.MetadataItem(pf.Peer)
		if !value.IsPresent() {
			value = w.decodePrimitiveDefault(pf)
		}
		result = value

	case FieldGroupStreamMessageId:
		result = types.OptionalOf(w.source.MessageId())

	case FieldGroupStreamChannel:
		result = types.OptionalOf(w.source.ChannelName())

	default:
		err = errors.Errorf("Cannot retrieve %q value from field group %q", pf.Type.Shape, pf.Group)
	}

	return
}

func NewMessageDecoder(source MessageDataSource, defaultContentType, defaultEncoding string) MessageDecoder {
	return MessageDecoder{
		defaultContentType: defaultContentType,
		defaultEncoding:    defaultEncoding,
		source:             source,
	}
}
