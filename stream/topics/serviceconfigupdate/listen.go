package serviceconfigupdate

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

type ApplicationStatusUpdateRequestHandler func(ctx context.Context, request ApplicationStatusUpdateRequest) error

func NewApplicationStatusUpdateRequestListener(fn ApplicationStatusUpdateRequestHandler, filters []stream.MessageFilter) stream.ListenerAction {
	return func(msg *message.Message) error {
		for _, filter := range filters {
			if !filter(msg.Context(), msg.Metadata) {
				return nil
			}
		}

		var request actionRequest
		var event ApplicationStatusUpdateRequest
		request.Request = &event
		err := json.Unmarshal(msg.Payload, &request)
		if err != nil {
			return errors.Wrap(err, "Failed to unmarshal message payload to ApplicationStatusUpdateRequest")
		}

		err = event.Validate()
		if err != nil {
			return errors.Wrap(err, "Failed to validate message payload to ApplicationStatusUpdateRequest")
		}

		return fn(msg.Context(), event)
	}
}

func AddApplicationStatusUpdateListener(fn ApplicationStatusUpdateRequestHandler) error {
	listener := NewApplicationStatusUpdateRequestListener(fn, []stream.MessageFilter{
		stream.FilterByMetaData(MetaDataEventType, EventTypeApplicationStatusUpdate),
	})
	return stream.AddListener(TopicServiceConfigUpdateTopic, listener)
}

func AddApplicationStatusUpdateServiceListener(fn ApplicationStatusUpdateRequestHandler, service string) error {
	listener := NewApplicationStatusUpdateRequestListener(fn, []stream.MessageFilter{
		stream.FilterByMetaData(MetaDataEventType, EventTypeApplicationStatusUpdate),
		stream.FilterByMetaData(MetaDataService, service),
	})
	return stream.AddListener(TopicServiceConfigUpdateTopic, listener)
}

type UpdateRequestHandler func(ctx context.Context, request UpdateRequest) error

func NewUpdateRequestListener(fn UpdateRequestHandler, filters []stream.MessageFilter) stream.ListenerAction {
	return func(msg *message.Message) error {
		for _, filter := range filters {
			if !filter(msg.Context(), msg.Metadata) {
				return nil
			}
		}

		var request actionRequest
		var event UpdateRequest
		request.Request = &event
		err := json.Unmarshal(msg.Payload, &request)
		if err != nil {
			return errors.Wrap(err, "Failed to unmarshal message payload to UpdateRequest")
		}

		err = event.Validate()
		if err != nil {
			return errors.Wrap(err, "Failed to validate message payload to UpdateRequest")
		}

		return fn(msg.Context(), event)
	}
}

func AddUpdateListener(fn UpdateRequestHandler) error {
	listener := NewUpdateRequestListener(fn, []stream.MessageFilter{
		stream.FilterByMetaData(MetaDataEventType, EventTypeUpdate),
	})
	return stream.AddListener(TopicServiceConfigUpdateTopic, listener)
}
