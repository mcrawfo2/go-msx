package integration

import (
	"net/http"
	"net/url"
)

type MsxRequest struct {
	EndpointName       string
	EndpointParameters map[string]string
	Headers            http.Header
	QueryParameters    url.Values
	Body               []byte
	ExpectEnvelope     bool
	NoToken            bool
	Payload            interface{}
	ErrorPayload       interface{}
}
