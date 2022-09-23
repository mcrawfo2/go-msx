# Publishers

Stream Operations Publishers are used to publish messages on streams. 
They consist of a number of components:

| Component                        | AsyncApi Documentor                 | Documentation Model |
|----------------------------------|-------------------------------------|---------------------|
| Your Message Publisher (service) | -                                   | -                   |
| Your Output Port (struct)        | asyncapi.MessagePublisherDocumentor | jsonschema.Schema   |
| Your Output DTO (struct)         | asyncapi.MessagePublisherDocumentor | jsonschema.Schema   |
| streamops.MessagePublisher       | asyncapi.MessagePublisherDocumentor | asyncapi.Message    |
| streamops.ChannelPublisher       | asyncapi.ChannelPublisherDocumentor | asyncapi.Operation  |
| streamops.Channel                | asyncapi.ChannelDocumentor          | asyncapi.Channel    |

## Generation

All go-msx code generation is done using the `skel` tool.

* To generate a channel supporting a single message publisher:
    ```bash
    skel generate-channel-publisher "COMPLIANCE_EVENT_TOPIC"
    ```

* To generate a channel supporting multiple message publishers,
  or add another message publisher to an existing multi-message publisher
  channel:
    ```bash
    skel generate-channel-publisher "COMPLIANCE_EVENT_TOPIC" --message "DriftCheck"
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
* `channel_publisher.go`
    * Channel publisher for the package channel
    * Channel publisher documentation (`asyncapi.Operation`)
* `message_publisher_*.go`
    * Message publisher for outgoing messages
    * Message publisher documentation (`asyncapi.Message`)
* `*.go`
    * DTOs for published messages (eg `DriftCheckRequest`)

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

		// Attach documentation to the channel
		channel.WithDocumentor(new(asyncapi.ChannelDocumentor).
			WithChannelItem(new(asyncapi.ChannelItem).
				WithDescription(
					"Commands originating from the Compliance service.  Compliance implementors " +
						"should subscribe to the topic and perform the specified action on the enclosed " +
						"entity.  Responses should be published to COMPLIANCE_UPDATE_TOPIC.")))

		return nil
	})
}
```

### Channel Publisher

The channel publisher component represents the set of publishable messages for a given stream.
It is implemented as a service that should be created as a dependency of your message publisher.


#### Example

`channel_publisher.go`
```go
package complianceevent

import (
  "context"
  "cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
  "cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
  "cto-github.cisco.com/NFV-BU/go-msx/types"
)

// Context Key for accessing context-injected instances
const contextKeyChannelPublisher = contextKeyNamed("ChannelPublisher")

// ContextChannelPublisher returns an accessor for context-injected instances
func ContextChannelPublisher() types.ContextKeyAccessor[*streamops.ChannelPublisher] {
  return types.NewContextKeyAccessor[*streamops.ChannelPublisher](contextKeyChannelPublisher)
}

// Constructor
func newChannelPublisher(ctx context.Context) (svc *streamops.ChannelPublisher, err error) {
    svc = contextChannelPublisher().Get(ctx)
    if svc == nil {
        // Create the publisher  
        svc, err = streamops.NewChannelPublisher(ctx, channel, "onComplianceEvent")
        if err != nil {
          return nil, err
        }
    
		// Declare some documentation
        doc := new(asyncapi.ChannelPublisherDocumentor).
                WithOperation(new(asyncapi.Operation).
                  WithID("onComplianceEvent").
                  WithSummary("Receive compliance commands.").
                  WithTags(*asyncapi.NewTag("toSouth")))
		
		// Attach the documentation to the publisher
        svc.AddDocumentor(doc)
    }
  
    return svc, nil
}

```


### Message Publisher

The message publisher component represents one of the publishable messages for a given stream.
It is implemented as a service created after configuration but before start-up.  
Notice that it has a defined API interface for mocking, and should be mocked by dependent services
during testing.

`message_publisher_drift_check.go`
```go
package complianceevent

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/streamops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/testservice/pkg/api"
)

// Dependencies

//go:generate mockery --inpackage --name=DriftCheckRequestPublisher --structname=MockDriftCheckRequestPublisher --filename mock_DriftCheckRequestPublisher.go

// DriftCheckRequestPublisher is the publisher interface used for mocking during tests
type DriftCheckRequestPublisher interface {
	PublishDriftCheckRequest(ctx context.Context, payload api.DriftCheckRequest) error
}

// Context Key for accessing context-injected instances
const contextKeyDriftCheckRequestPublisher = contextKeyNamed("DriftCheckRequestPublisher")

// ContextDriftCheckRequestPublisher returns an accessor for context-injected instances
func ContextDriftCheckRequestPublisher() types.ContextKeyAccessor[DriftCheckRequestPublisher] {
	return types.NewContextKeyAccessor[DriftCheckRequestPublisher](contextKeyDriftCheckRequestPublisher)
}

// Implementation

// Concrete implementation of interface to be used by default.
type driftCheckRequestPublisher struct {
	messagePublisher *streamops.MessagePublisher
}

// Declare the outputs for your publisher.
type driftCheckRequestOutput struct {
	EventType string                `out:"header" const:"DriftCheck"`
	Payload   api.DriftCheckRequest `out:"body"`
}

// PublishDriftCheckRequest publishes a DriftCheckRequest message to the channel.
func (p driftCheckRequestPublisher) PublishDriftCheckRequest(ctx context.Context, payload DriftCheckRequest) error {
	return p.messagePublisher.Publish(ctx, driftCheckRequestOutput{
		Payload: payload,
	})
}

// NewDriftCheckRequestPublisher constructs a new instance of DriftCheckRequestPublisher.
// If an instance exists in the context, it is returned instead.
func NewDriftCheckRequestPublisher(ctx context.Context) (DriftCheckRequestPublisher, error) {
	svc := ContextDriftCheckRequestPublisher().Get(ctx)
	if svc == nil {
		// Declare some documentation for your publisher
		doc := new(asyncapi.MessagePublisherDocumentor).
			WithMessage(new(asyncapi.Message).
				WithTitle("Drift Check Request").
				WithSummary("Request consumer to check drift of entity configuration.").
				WithTags(
					*asyncapi.NewTag("driftCheck"),
					*asyncapi.NewTag("toSouth"),
				))

		// Obtain the channel publisher
		cp, err := newChannelPublisher(ctx)
		if err != nil {
			return nil, err
		}

		// Create a message publisher builder
		mpb, err := streamops.NewMessagePublisherBuilder(ctx, cp, "DriftCheckRequest", driftCheckRequestOutput{})
		if err != nil {
			return nil, err
		}

		// Configure and build the publisher
		mp, err := mpb.WithDocumentor(doc).Build()
		if err != nil {
			return nil, err
		}

		svc = &driftCheckRequestPublisher{
			messagePublisher: mp,
		}
	}

	return svc, nil
}
```

### Payload DTO

The payload DTO will contain the body of message that is to be published.
Before dispatch to the underlying stream, the message will be validated using
the JSON-schema annotations on your DTO.

`drift_check_request.go`
```go
package complianceevent

import "cto-github.cisco.com/NFV-BU/go-msx/types"

type DriftCheckRequest struct {
	// Action corresponds to the JSON schema field "action".
	Action DriftCheckRequestAction `json:"action" const:"checkDrift"`

	// Domain corresponds to the JSON schema field "domain".
	Domain string `json:"domain"`

	// EntityId corresponds to the JSON schema field "entityId".
	EntityId string `json:"entityId"`

	// EntityLevelCompliance corresponds to the JSON schema field
	// "entityLevelCompliance".
	EntityLevelCompliance DriftCheckRequestEntityLevelCompliance `json:"entityLevelCompliance" enum:"full,partial"`

	// EntityType corresponds to the JSON schema field "entityType".
	EntityType string `json:"entityType"`

	// GroupId corresponds to the JSON schema field "groupId".
	GroupId *types.UUID `json:"groupId,omitempty"`

	// MessageId corresponds to the JSON schema field "messageId".
	MessageId *types.UUID `json:"messageId,omitempty"`

	// Standards corresponds to the JSON schema field "standards".
	Standards []ConfigPayload `json:"standards" required:"true" minItems:"1"`

	// Timestamp corresponds to the JSON schema field "timestamp".
	Timestamp types.Time `json:"timestamp"`

	// _ corresponds to the JSON schema of the parent structure
	_ struct {} `json:"-" title:"DriftCheckRequest"`
}

type DriftCheckRequestAction string

const DriftCheckRequestActionCheckDrift DriftCheckRequestAction = "checkDrift"
type DriftCheckRequestEntityLevelCompliance string

const DriftCheckRequestEntityLevelComplianceFull DriftCheckRequestEntityLevelCompliance = "full"
const DriftCheckRequestEntityLevelCompliancePartial DriftCheckRequestEntityLevelCompliance = "partial"
```