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

type MessageSubscriberDocumentor struct {
	skip    bool
	message *Message
	mutator MessageMutator
}

func (d *MessageSubscriberDocumentor) WithSkip(skip bool) *MessageSubscriberDocumentor {
	d.skip = skip
	return d
}

func (d *MessageSubscriberDocumentor) WithMessage(message *Message) *MessageSubscriberDocumentor {
	d.message = message
	return d
}

func (d *MessageSubscriberDocumentor) WithMessageMutator(fn MessageMutator) *MessageSubscriberDocumentor {
	d.mutator = fn
	return d
}

func (d MessageSubscriberDocumentor) DocType() string {
	return DocType
}

func (d MessageSubscriberDocumentor) Document(subscriber *streamops.MessageSubscriber) error {
	if d.skip {
		return nil
	}

	message := d.message
	if message == nil {
		message = new(Message)
	}

	if message.ContentType == nil {
		if subscriber.ContentType() != "" {
			message.WithContentType(subscriber.ContentType())
		} else {
			message.WithContentType("application/json")
		}
	}

	// Headers
	if headersSchema := d.headersSchema(subscriber); headersSchema != nil {
		message.HeadersEns().WithSchema(*headersSchema)
	}

	// Body
	if bodySchema := d.bodySchema(subscriber); bodySchema != nil {
		message.WithPayload(*bodySchema)
	}

	// Mutator
	if d.mutator != nil {
		d.mutator(message)
	}

	// Publish
	messageRef := createMessageReference(subscriber.Name(), *message)
	addSubscriberOperationMessageChoice(subscriber.Channel().Name(), messageRef)

	return nil
}

func (d MessageSubscriberDocumentor) headersSchema(subscriber *streamops.MessageSubscriber) *jsonschema.Schema {
	headerPortFields := subscriber.InputPort().Fields.All(
		ops.PortFieldHasGroup(streamops.FieldGroupStreamHeader),
	)

	if len(headerPortFields) == 0 {
		return nil
	}

	headersSchema := js.ObjectSchema()
	for _, header := range headerPortFields {
		headersSchema.WithPropertiesItem(
			header.Peer,
			schemaFromPortField(header).ToSchemaOrBool())
	}

	return headersSchema
}

func (d MessageSubscriberDocumentor) bodySchema(subscriber *streamops.MessageSubscriber) *jsonschema.Schema {
	bodyPortField := subscriber.InputPort().Fields.First(
		ops.PortFieldHasGroup(streamops.FieldGroupStreamBody),
	)

	return schemaFromPortField(bodyPortField)
}
