// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/js"
	"github.com/swaggest/jsonschema-go"
)

type MessageMutator func(message *Message)

type MessagePublisherDocumentor struct {
	skip    bool
	message *Message
	mutator MessageMutator
}

func (d *MessagePublisherDocumentor) WithSkip(skip bool) *MessagePublisherDocumentor {
	d.skip = skip
	return d
}

func (d *MessagePublisherDocumentor) WithMessage(message *Message) *MessagePublisherDocumentor {
	d.message = message
	return d
}

func (d *MessagePublisherDocumentor) WithMessageMutator(fn MessageMutator) *MessagePublisherDocumentor {
	d.mutator = fn
	return d
}

func (d MessagePublisherDocumentor) DocType() string {
	return DocType
}

func (d MessagePublisherDocumentor) Document(publisher *streamops.MessagePublisher) error {
	if d.skip {
		return nil
	}

	message := d.message
	if message == nil {
		message = new(Message)
	}

	if message.ID == nil {
		message.WithID(publisher.Name())
	}

	if message.ContentType == nil {
		if publisher.ContentType() != "" {
			message.WithContentType(publisher.ContentType())
		} else {
			message.WithContentType("application/json")
		}
	}

	// Headers
	if headersSchema := d.headersSchema(publisher); headersSchema != nil {
		message.HeadersEns().WithSchema(*headersSchema)
	}

	// Body
	if bodySchema := d.bodySchema(publisher); bodySchema != nil {
		message.WithPayload(*bodySchema)
	}

	// Mutator
	if d.mutator != nil {
		d.mutator(message)
	}

	// Publish
	messageRef := createMessageReference(publisher.Name(), *message)
	addPublisherOperationMessageChoice(publisher.Channel().Name(), messageRef)

	return nil
}

func (d MessagePublisherDocumentor) headersSchema(publisher *streamops.MessagePublisher) *jsonschema.Schema {
	headerPortFields := publisher.OutputPort().Fields.All(
		ops.PortFieldHasGroup(streamops.FieldGroupStreamHeader),
	)

	if len(headerPortFields) == 0 {
		return nil
	}

	headersSchema := js.ObjectSchema()
	headersRequired := []string{}
	for _, header := range headerPortFields {
		headersSchema.WithPropertiesItem(
			header.Peer,
			jsonSchemaFromPortField(header).ToSchemaOrBool())
		if !header.Optional {
			headersRequired = append(headersRequired, header.Peer)
		}
	}
	headersSchema.Required = headersRequired

	return headersSchema
}

func (d MessagePublisherDocumentor) bodySchema(publisher *streamops.MessagePublisher) *jsonschema.Schema {
	bodyPortField := publisher.OutputPort().Fields.First(
		ops.PortFieldHasGroup(streamops.FieldGroupStreamBody),
	)

	return jsonSchemaFromPortField(bodyPortField)
}
