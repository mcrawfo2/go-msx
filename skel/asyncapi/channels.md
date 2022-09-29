# Channels

* To interactively generate a channel publisher or subscriber, for one or
  more messages:
    ```bash
    skel generate-channel
    ```

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
  
## Files

From the above examples, the following files may be generated:

* `pkg.go`
    - Package-wide logger
    - Context Key type definition
    - Channel for `COMPLIANCE_EVENT_TOPIC`
    - Channel documentation (`asyncapi.Channel`)
* `publisher_channel.go`
    * Channel publisher for the package channel
    * Channel publisher documentation (`asyncapi.Operation`)
* `subscriber_channel.go`
    * Channel subscriber for the package channel
    * Channel subscriber documentation (`asyncapi.Operation`)
* `publisher_*.go`
    * Message publisher for individual outgoing messages
    * Message publisher documentation (`asyncapi.Message`)
* `subscriber_*.go`
    * Message subscriber for individual incoming messages
    * Message subscriber documentation (`asyncapi.Message`)
* `api/*.go`
    * DTOs for published messages (eg `DriftCheckRequest`)
