# AsyncAPI

* To interactively generate a channel publisher or subscriber, for one or
  more channels from an existing AsyncApi specification via url or local path:
    ```bash
    skel generate-channel-asyncapi
    ```

* To generate a consumer for channels from an existing AsyncApi specification via url:
    ```bash
    export ASYNCAPI_SPEC_URL="https://cto-github.cisco.com/raw/NFV-BU/merakigoservice/develop/api/asyncapi.yaml?token=..."
    skel generate-channel-asyncapi "$ASYNCAPI_SPEC_URL" COMPLIANCE_EVENT_TOPIC 
    ```

* To generate a consumer for channels from an existing AsyncApi specification from a local
  specification:
    ```bash
    skel generate-channel-asyncapi "api/asyncapi.yaml" COMPLIANCE_EVENT_TOPIC 
    ```

* To generate a consumer for channels from an existing AsyncApi specification via url:
    ```bash
    export ASYNCAPI_SPEC_URL="https://cto-github.cisco.com/raw/NFV-BU/merakigoservice/develop/api/asyncapi.yaml?token=..."
    skel generate-channel-asyncapi "$ASYNCAPI_SPEC_URL" COMPLIANCE_EVENT_TOPIC 
    ```

* To generate a consumer for channels from an existing AsyncApi specification from a local
  specification:
    ```bash
    skel generate-channel-asyncapi "api/asyncapi.yaml" COMPLIANCE_EVENT_TOPIC 
    ```
