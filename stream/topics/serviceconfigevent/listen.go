package serviceconfigevent

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

const TopicServiceConfigEventTopic = "SERVICECONFIG_EVENT_TOPIC"

type EventHandler func(ctx context.Context, event Event) error
type AssignmentEventHandler func(ctx context.Context, event AssignmentEvent) error
type ApplicationEventHandler func(ctx context.Context, event ApplicationEvent) error

func NewEventListener(fn EventHandler, filters []stream.MessageFilter) stream.ListenerAction {
	return func(msg *message.Message) error {
		if !stream.FilterMessage(msg, filters) {
			return nil
		}

		var event Event
		err := json.Unmarshal(msg.Payload, &event)
		if err != nil {
			return errors.Wrap(err, "Failed to unmarshal message payload to Event")
		}

		event.Headers = msg.Metadata

		err = event.Validate()
		if err != nil {
			return errors.Wrap(err, "Failed to validate message payload to Event")
		}

		return fn(msg.Context(), event)
	}
}

// AddEventListener adds a listener for ServiceConfig events
func AddEventListener(fn EventHandler, eventTypes ...string) error {
	for _, eventType := range eventTypes {
		switch eventType {
		case EventTypeCreated, EventTypeUpdated, EventTypeDeleted, EventTypeStatusUpdated:
		default:
			return errors.Errorf("Invalid event type for event listener: %q", eventType)
		}
	}

	listener := NewEventListener(fn, []stream.MessageFilter{
		stream.FilterByMetaData(MetaDataEventType, eventTypes...),
	})
	return stream.AddListener(TopicServiceConfigEventTopic, listener)
}

func NewAssignmentEventListener(fn AssignmentEventHandler, filters []stream.MessageFilter) stream.ListenerAction {
	return func(msg *message.Message) error {
		if !stream.FilterMessage(msg, filters) {
			return nil
		}

		var event AssignmentEvent
		err := json.Unmarshal(msg.Payload, &event)
		if err != nil {
			return errors.Wrap(err, "Failed to unmarshal message payload to AssignmentEvent")
		}

		event.Headers = msg.Metadata

		err = event.Validate()
		if err != nil {
			return errors.Wrap(err, "Failed to validate message payload to Event")
		}

		return fn(msg.Context(), event)
	}
}

// AddAssignmentEventListener adds a listener for ServiceConfigAssignment events
func AddAssignmentEventListener(fn AssignmentEventHandler, eventTypes ...string) error {
	for _, eventType := range eventTypes {
		switch eventType {
		case EventTypeAssignmentCreated, EventTypeAssignmentDeleted, EventTypeAssignmentStatusUpdated:
		default:
			return errors.Errorf("Invalid event type for assignment event listener: %q", eventType)
		}
	}

	listener := NewAssignmentEventListener(fn, []stream.MessageFilter{
		stream.FilterByMetaData(MetaDataEventType, eventTypes...),
	})
	return stream.AddListener(TopicServiceConfigEventTopic, listener)
}

func NewApplicationEventListener(fn ApplicationEventHandler, filters []stream.MessageFilter) stream.ListenerAction {
	return func(msg *message.Message) error {
		if !stream.FilterMessage(msg, filters) {
			return nil
		}

		var event ApplicationEvent
		err := json.Unmarshal(msg.Payload, &event)
		if err != nil {
			return errors.Wrap(err, "Failed to unmarshal message payload")
		}

		event.Headers = msg.Metadata

		err = event.Validate()
		if err != nil {
			return errors.Wrap(err, "Failed to validate message payload")
		}

		return fn(msg.Context(), event)
	}
}

// AddApplicationEventListener adds a listener for ServiceConfigApplication events
func AddApplicationEventListener(fn ApplicationEventHandler, filters []stream.MessageFilter, eventTypes ...string) error {
	for _, eventType := range eventTypes {
		switch eventType {
		case EventTypeApplicationCreated, EventTypeApplicationDeleted, EventTypeApplicationStatusUpdated:
		default:
			return errors.Errorf("Invalid event type for application event listener: %q", eventType)
		}
	}

	listener := NewApplicationEventListener(fn, append([]stream.MessageFilter{
		stream.FilterByMetaData(MetaDataEventType, eventTypes...),
	}, filters...))
	return stream.AddListener(TopicServiceConfigEventTopic, listener)
}
