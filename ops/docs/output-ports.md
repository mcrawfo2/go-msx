# Output Ports

An Output Port is a go structure used to describe a target payload to be generated,
such as an HTTP Response or an outgoing Stream Message.  

The port structure must be populated by you before being either:
- passed into the stream Message Publisher; or
- returned from your HTTP Endpoint Controller

go-msx will validate the contents of your structure and if it is valid,
publish the message or response.

An example stream message output port:
```go
type driftCheckResponseOutputs struct {
    EventType string                 `out:"header" const:"DriftCheck"`
    Payload   api.DriftCheckResponse `out:"body"`
}
```

An example HTTP response input port:
```go
type createEntityOutputs struct {
    Payload api.CreateEntityResponse `resp:"body"`	
}
```

### Struct Tags

Each field with an output struct tag will be applied to the outgoing payload.

Note that the struct tag prefix depends on the protocol being described:
- For HTTP Requests, the output struct tag must be `resp` (for backwards compatibility)
- For Stream Messages, the output struct tag must be `out`.

The full syntax of the input struct tag is one of the following, appropriate for
the handling protocol:

    output:"<fieldGroup>[=<peerName>]"
    resp:"<fieldGroup>[=<peerName>]"

The output struct tag contains the following subcomponents:

- `<fieldGroup>` (_Required_):
The name of the message/response part to which the value will be injected.

- `[=<peerName>]` (_Optional_):
A _peer_ is a field or property in the target payload. For example an HTTP response may have
a header with the name date, which can be requested using the following output struct tag: `resp:"header=Date"` 

See [HTTP Response Ports]() and [Stream Ports](../streamops/ports.md) for available field groups and peer
name conventions for your specific protocol.

