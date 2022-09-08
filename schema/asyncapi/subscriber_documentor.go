// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
)

type ChannelSubscriberDocumentor struct {
	skip      bool
	operation *Operation
	mutator   OperationMutator
}

func (d *ChannelSubscriberDocumentor) WithSkip(skip bool) *ChannelSubscriberDocumentor {
	d.skip = skip
	return d
}

func (d *ChannelSubscriberDocumentor) WithOperation(operation *Operation) *ChannelSubscriberDocumentor {
	d.operation = operation
	return d
}

func (d *ChannelSubscriberDocumentor) WithOperationMutator(fn OperationMutator) *ChannelSubscriberDocumentor {
	d.mutator = fn
	return d
}

func (d ChannelSubscriberDocumentor) DocType() string {
	return DocType
}

func (d ChannelSubscriberDocumentor) Document(subscriber *streamops.ChannelSubscriber) error {
	// Initialize
	operation := d.operation
	if operation == nil {
		operation = new(Operation)
	}

	// Name
	operation.WithID(subscriber.Name())

	// Bindings
	operationBindingsEns(operation, subscriber.Channel().Binder())

	// Mutator
	if d.mutator != nil {
		d.mutator(operation)
	}

	// Publish
	addSubscriberOperation(subscriber.Channel().Name(), *operation)

	// Children
	for _, messageSubscriber := range streamops.RegisteredMessageSubscribers(subscriber.Channel().Name()) {
		err := ops.DocumentorWithType[streamops.MessageSubscriber](messageSubscriber, DocType).
			OrElse(MessageSubscriberDocumentor{}).
			Document(messageSubscriber)
		if err != nil {
			return err
		}
	}

	return nil
}

func getChannelSubscriberDocumentor(subscriber *streamops.ChannelSubscriber) ops.Documentor[streamops.ChannelSubscriber] {
	documentor := subscriber.Documentor(ops.WithDocType[streamops.ChannelSubscriber](DocType))
	if documentor == nil {
		documentor = ChannelSubscriberDocumentor{}
	}
	return documentor
}
