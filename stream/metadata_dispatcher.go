// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package stream

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

type MetadataHeader string
type ListenerActionLookup func(value string) ListenerAction

type Dispatcher interface {
	Dispatch(msg *message.Message) error
}

type MetadataDispatcher struct {
	metadataHeaderName    string
	metadataHeaderActions map[MetadataHeader]ListenerAction
	actionLookup          ListenerActionLookup
}

func (d *MetadataDispatcher) OnMessage(msg *message.Message) error {
	return d.Dispatch(msg)
}

func (d *MetadataDispatcher) Dispatch(msg *message.Message) error {
	log := logger.WithContext(msg.Context()).
		WithField("messageId", msg.UUID).
		WithField("metadataHeaderName", d.metadataHeaderName)
	log.Info("Dispatching the message")

	headerValue, ok := msg.Metadata[d.metadataHeaderName]
	if !ok {
		log.Info("The message metadata does not contain the header")
		return nil
	}

	lookup := d.actionLookup
	var action ListenerAction
	if lookup == nil {
		lookup = func(value string) ListenerAction {
			return d.metadataHeaderActions[MetadataHeader(headerValue)]
		}
	}
	action = lookup(headerValue)
	if action == nil {
		log.WithField("metadataHeader", headerValue).Info("No listener action exists for the metadata header")
		return nil
	}

	log.WithField("metadataHeader", headerValue).Info("Dispatching the message to the listener action")
	if err := action(msg); err != nil {
		log.WithField("metadataHeader", headerValue).Error("Failed to process the message with the listener action")
		return err
	}

	return nil
}

func NewMetadataDispatcher(metadataHeaderName string, metadataHeaderActions map[MetadataHeader]ListenerAction) (Dispatcher, error) {
	if len(metadataHeaderName) == 0 {
		return nil, errors.New("Empty metadata header name")
	}

	if len(metadataHeaderActions) == 0 {
		return nil, errors.New("Empty metadata header actions")
	}

	return &MetadataDispatcher{
		metadataHeaderName:    metadataHeaderName,
		metadataHeaderActions: metadataHeaderActions,
	}, nil
}

func NewMetadataDispatcherIndirect(metadataHeaderName string, lookup ListenerActionLookup) (Dispatcher, error) {
	if len(metadataHeaderName) == 0 {
		return nil, errors.New("Empty metadata header name")
	}

	if lookup == nil {
		return nil, errors.New("Missing action lookup")
	}

	return &MetadataDispatcher{
		metadataHeaderName: metadataHeaderName,
		actionLookup:       lookup,
	}, nil
}
