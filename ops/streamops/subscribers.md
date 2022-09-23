# Subscribers

Stream Operations Subscribers are used to publish messages on streams.
They consist of a number of components:

| Component                         | AsyncApi Documentor                  | Documentation Model |
|-----------------------------------|--------------------------------------|---------------------|
| Your Message Subscriber (service) | -                                    | -                   |
| Your Input Port (struct)          | asyncapi.MessageSubscriberDocumentor | jsonschema.Schema   |
| Your Input DTO (struct)           | asyncapi.MessageSubscriberDocumentor | jsonschema.Schema   |
| streamops.MessageSubscriber       | asyncapi.MessageSubscriberDocumentor | asyncapi.Message    |
| streamops.ChannelSubscriber       | asyncapi.ChannelSubscriberDocumentor | asyncapi.Operation  |
| streamops.Channel                 | asyncapi.ChannelDocumentor           | asyncapi.Channel    |

## Generation

All go-msx code generation is done using the `skel` tool.

* To generate a channel supporting a single message subscriber:
    ```bash
    skel generate-channel-subscriber "COMPLIANCE_EVENT_TOPIC"
    ```

* To generate a channel supporting multiple message subscribers,
  or add another message subscriber to an existing multi-message subscriber
  channel:
    ```bash
    skel generate-channel-subscriber "COMPLIANCE_EVENT_TOPIC" --message "DriftCheck"
    ```

* To generate a consumer for channels from an existing AsyncApi specification via url:
    ```bash
    export ASYNCAPI_SPEC_URL="https://cto-github.cisco.com/raw/NFV-BU/merakigoservice/develop/api/asyncapi.yaml?token=..."
    skel generate-channel-asyncapi "$ASYNCAPI_SPEC_URL" COMPLIANCE_EVENT_TOPIC 
    ```

* To generate a consumer for channels from an existing AsyncApi specification from a local
  specification:
    ```bash
    skel generate-channel-asyncapi "$ASYNCAPI_SPEC_URL" COMPLIANCE_EVENT_TOPIC 
    ```

The above examples will generate these components to `/internal/stream/complianceevent`:

* `pkg.go`
    - Package-wide logger
    - Context Key type definition
    - Channel for `COMPLIANCE_EVENT_TOPIC`
    - Channel documentation (`asyncapi.Channel`)
* `channel_subscriber.go`
    * Channel subscriber for the package channel
    * Channel subscriber documentation (`asyncapi.Operation`)
* `message_subscriber_*.go`
    * Message subscriber for incoming messages
    * Message subscriber documentation (`asyncapi.Message`)
* `*.go`
    * DTOs for subscribed messages (eg `DriftCheckResponse`)

## Components

### Channel

The channel component represents the stream itself (SQS or Kafka topic, Redis stream, Go channel, SQLDB table, etc).
It is implemented as a singleton that should be created after configuration but before start-up.

#### Example

`pkg.go`
```go
package complianceevent

import (
  "context"
  "cto-github.cisco.com/NFV-BU/go-msx/app"
  "cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
  "cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
)

type contextKeyNamed string

const ChannelName = "COMPLIANCE_EVENT_TOPIC"

var channel *streamops.Channel

// Initialize the channel definition during the Configure phase
func init() {
	app.OnEvent(app.EventConfigure, app.PhaseAfter, func(ctx context.Context) (err error) {
		// Create a channel instance
		channel, err = streamops.NewChannel(ctx, ChannelName)
		if err != nil {
			return
		}
		
		// Declare some documentation
		doc := new(asyncapi.ChannelDocumentor).
			WithChannelItem(new(asyncapi.ChannelItem).
				WithDescription(
					"Commands originating from the Compliance service.  Compliance implementors " +
					"should subscribe to the topic and perform the specified action on the enclosed " +
					"entity.  Responses should be published to COMPLIANCE_UPDATE_TOPIC."))

		// Attach the documentation to the channel
		channel.WithDocumentor(doc)

		return nil
	})
}

```

### Channel Subscriber

The channel subscriber component represents the set of subscribable messages for a given stream.
It is implemented as a service, and should have one of your application services as a dependency.


#### Example

`channel_subscriber.go`
```go
package complianceevent

import (
  "context"
  "cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
  "cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
  "cto-github.cisco.com/NFV-BU/go-msx/types"
)

// Context Key for accessing context-injected instances
const contextKeyChannelSubscriber = contextKey("ChannelSubscriber")

// ContextChannelSubscriber returns an accessor for context-injected instances
func ContextChannelSubscriber() types.ContextKeyAccessor[*streamops.ChannelSubscriber] {
	return types.NewContextKeyAccessor[*streamops.ChannelSubscriber](contextKeyChannelSubscriber)
}

// Constructor
func newChannelSubscriber(ctx context.Context) (channelSubscriber *streamops.ChannelSubscriber, err error) {
	channelSubscriber = ContextChannelSubscriber().Get(ctx)
	if channelSubscriber == nil {
		// Declare some documentation
		doc := new(asyncapi.ChannelSubscriberDocumentor).
			WithOperation(new(asyncapi.Operation).
				WithID("sendComplianceUpdate").
				WithSummary("Send compliance action results.").
				WithTags(*asyncapi.NewTag("fromSouth")))

		// Create a channel subscriber
		channelSubscriber, err = streamops.NewChannelSubscriber(ctx,
			channel,
			"sendComplianceUpdate",
			types.OptionalOf("eventType")) // empty for no header-based dispatch
		if err != nil {
			return nil, err
		}

		// Attach the documentation to the subscriber
		channelSubscriber.AddDocumentor(doc)
	}

	return channelSubscriber, nil
}

```


### Message Subscriber

The message subscriber component represents one of the publishable messages for a given stream.
It is implemented as a service created after configuration but before start-up.  
Notice that it has a defined API interface for mocking, and should be mocked by dependent services
during testing.

`message_subscriber_drift_check.go`
```go
package complianceevent

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
	"cto-github.cisco.com/NFV-BU/go-msx/security/service"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/testservice/internal/compliance"
	"cto-github.cisco.com/NFV-BU/testservice/pkg/api"
)

// Dependency Management

// ToggleComplianceResponseHandler is the handler interface used for mocking during tests
type ToggleComplianceResponseHandler interface {
	OnToggleComplianceUpdate(ctx context.Context, payload api.ToggleComplianceResponse) error
}

// Context Key for accessing context-injected instances
const contextKeyToggleComplianceResponseSubscriber = contextKey("ToggleComplianceResponseSubscriber")

// ContextToggleComplianceResponseSubscriber returns an accessor for context-injected instances
func ContextToggleComplianceResponseSubscriber() types.ContextKeyAccessor[*streamops.MessageSubscriber] {
	return types.NewContextKeyAccessor[*streamops.MessageSubscriber](contextKeyToggleComplianceResponseSubscriber)
}

// Constructor
func newToggleComplianceResponseSubscriber(ctx context.Context) (*streamops.MessageSubscriber, error) {
	svc := ContextToggleComplianceResponseSubscriber().Get(ctx)
	if svc == nil {
		// Declare the inputs for your subscriber
		type toggleComplianceResponseInput struct {
			// EventType contains the eventType header value.   Only when it matches "ComplianceUpdate"
			// will this message subscriber will be executed.
			EventType string                       `in:"header" const:"ComplianceUpdate"`
			// Payload contains the message body that was received.
			Payload   api.ToggleComplianceResponse `in:"body"`
		}

		// Declare some documentation for your subscriber
		doc := new(asyncapi.MessageSubscriberDocumentor).
			WithMessage(new(asyncapi.Message).
				WithTitle("Toggle Compliance Response").
				WithSummary("Inform about enable/add entity to compliance monitoring.").
				WithTags(
					*asyncapi.NewTag("toggleCompliance"),
					*asyncapi.NewTag("fromSouth"),
				))

		// Obtain the channel subscriber
		cs, err := newChannelSubscriber(ctx)
		if err != nil {
			return nil, err
		}

		// Create an instance of your application service to call
		handler, err := compliance.NewComplianceService(ctx)
		if err != nil {
			return nil, err
		}

		// Create a builder
		sb, err := streamops.NewMessageSubscriberBuilder(ctx, cs, "ToggleComplianceResponse")
		if err != nil {
			return nil, err
		}

		// Configure and build the subscriber
		svc, err = sb.
			WithInputs(toggleComplianceResponseInput{}).
			//WithMetadataFilterValues("eventType", "ComplianceUpdate"). // not required, already specified in inputs
			WithDecorator(service.DefaultServiceAccount).
			WithHandler(func(ctx context.Context, in *toggleComplianceResponseInput) error {
				return handler.OnToggleComplianceUpdate(ctx, in.Payload)
			}).
			WithDocumentor(doc).
			Build()
		if err != nil {
			return nil, err
		}
	}

	return svc, nil
}

func init() {
	// Register the subscriber during module initialization.
	app.OnCommandsEvent(
		[]string{app.CommandRoot, app.CommandAsyncApi},
		app.EventStart,
		app.PhaseBefore,
		func(ctx context.Context) error {
			_, err := newToggleComplianceResponseSubscriber(ctx)
			return err
		})
}
```

### Payload DTO

The payload DTO will contain the parsed message that is subscribed.
Before dispatch to your subscriber, the message will be validated using
the JSON-schema annotations and any `Validatable` interface implementation
on your DTO.

`toggle_compliance_response.go`
```go
package complianceupdate

import "cto-github.cisco.com/NFV-BU/go-msx/types"

type ToggleComplianceResponse struct {
	// Action corresponds to the JSON schema field "action".
	Action ToggleComplianceResponseAction `json:"action" enum:"enableCompliance,disableCompliance"`

	// Domain corresponds to the JSON schema field "domain".
	Domain string `json:"domain"`

	// EntityId corresponds to the JSON schema field "entityId".
	EntityId string `json:"entityId"`

	// EntityType corresponds to the JSON schema field "entityType".
	EntityType string `json:"entityType"`

	// GroupId corresponds to the JSON schema field "groupId".
	GroupId *types.UUID `json:"groupId,omitempty"`

	// Message corresponds to the JSON schema field "message".
	Message string `json:"message"`

	// Status corresponds to the JSON schema field "status".
	Status string `json:"status"`

	// Timestamp corresponds to the JSON schema field "timestamp".
	Timestamp types.Time `json:"timestamp"`

	// _ corresponds to the JSON schema of the parent structure
	_ struct {} `json:"-" title:"ToggleComplianceResponse"`
}

type ToggleComplianceResponseAction string

const ToggleComplianceResponseActionDisableCompliance ToggleComplianceResponseAction = "disableCompliance"
const ToggleComplianceResponseActionEnableCompliance ToggleComplianceResponseAction = "enableCompliance"
```