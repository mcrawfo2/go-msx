// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
)

type OperationMutator func(*Operation)

type ChannelPublisherDocumentor struct {
	skip      bool
	operation *Operation
	mutator   OperationMutator
}

func (d *ChannelPublisherDocumentor) WithSkip(skip bool) *ChannelPublisherDocumentor {
	d.skip = skip
	return d
}

func (d *ChannelPublisherDocumentor) WithOperation(operation *Operation) *ChannelPublisherDocumentor {
	d.operation = operation
	return d
}

func (d *ChannelPublisherDocumentor) WithOperationMutator(fn OperationMutator) *ChannelPublisherDocumentor {
	d.mutator = fn
	return d
}

func (d ChannelPublisherDocumentor) DocType() string {
	return DocType
}

func (d ChannelPublisherDocumentor) Document(publisher *streamops.ChannelPublisher) error {
	// Initialize
	operation := d.operation
	if operation == nil {
		operation = new(Operation)
	}

	// Name
	operation.WithID(publisher.Name())

	// Server

	// Bindings
	operationBindingsEns(operation, publisher.Channel().Binder())

	// Mutator
	if d.mutator != nil {
		d.mutator(operation)
	}

	// Publish
	addPublisherOperation(publisher.Channel().Name(), *operation)

	// Children
	for _, messagePublisher := range streamops.RegisteredMessagePublishers(publisher.Channel().Name()) {
		err := ops.DocumentorWithType[streamops.MessagePublisher](messagePublisher, DocType).
			OrElse(MessagePublisherDocumentor{}).
			Document(messagePublisher)
		if err != nil {
			return err
		}
	}

	return nil
}
