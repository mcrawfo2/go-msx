package stream

import (
	"errors"

	"github.com/ThreeDotsLabs/watermill/message"
)

type MetadataHeader string

type Dispatcher interface {
	Dispatch(msg *message.Message) error
}

type MetadataDispatcher struct {
	metadataHeaderName    string
	metadataHeaderActions map[MetadataHeader]ListenerAction
}

func (d *MetadataDispatcher) Dispatch(msg *message.Message) error {
	log := logger.WithContext(msg.Context()).
		WithField("messageId", msg.UUID).
		WithField("metadataHeaderName", d.metadataHeaderName)
	log.Info("Dispatching the message")

	metadataHeader, ok := msg.Metadata[d.metadataHeaderName]
	if !ok {
		log.Info("The message metadata does not contain the header")
		return nil
	}

	action, ok := d.metadataHeaderActions[MetadataHeader(metadataHeader)]
	if !ok {
		log.WithField("metadataHeader", metadataHeader).Info("No listener action exists for the metadata header")
		return nil
	}

	log.WithField("metadataHeader", metadataHeader).Info("Dispatching the message to the listener action")
	if err := action(msg); err != nil {
		log.WithField("metadataHeader", metadataHeader).Error("Failed to process the message with the listener action")
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
