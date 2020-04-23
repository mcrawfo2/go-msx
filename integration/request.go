package integration

import (
	"bytes"
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient/discoveryinterceptor"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient/loginterceptor"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient/rpinterceptor"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient/traceinterceptor"
	"cto-github.cisco.com/NFV-BU/go-msx/security/httprequest"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	SchemeMsxService = "http"
)

type Target struct {
	ServiceName  string
	ServiceType  ServiceType
	EndpointName string
	Method       string
	Path         string
}

func (t Target) Endpoint() Endpoint {
	return Endpoint{
		Name:   t.EndpointName,
		Method: t.Method,
		Path:   t.Path,
	}
}

func (t Target) Url(endpointParameters map[string]string, queryParameters url.Values) (string, error) {
	return t.Endpoint().Url(
		SchemeMsxService,
		t.ServiceName,
		endpointParameters,
		queryParameters)
}

type MsxRequest struct {
	Target             Target
	EndpointParameters map[string]string
	Headers            http.Header
	QueryParameters    url.Values
	Body               []byte
	ExpectEnvelope     bool
	NoToken            bool
	Payload            interface{}
	ErrorPayload       interface{}
}

func (v *MsxRequest) newHttpRequest(ctx context.Context) (*http.Request, error) {
	var req *http.Request
	var err error
	var buf io.Reader

	fullUrl, err := v.Target.Url(v.EndpointParameters, v.QueryParameters)
	if err != nil {
		return nil, err
	}

	if v.Body != nil {
		buf = bytes.NewBuffer(v.Body)
	} else {
		buf = http.NoBody
	}

	req, err = http.NewRequest(v.Target.Method, fullUrl, buf)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	req.Header.Set("Accept", httpclient.MimeTypeApplicationJson)

	if v.Body != nil {
		if v.Headers.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", httpclient.MimeTypeApplicationJson)
		}
	}

	if !v.NoToken {
		httprequest.InjectToken(req)
	}

	for k, vs := range v.Headers {
		for _, v := range vs {
			req.Header.Set(k, v)
		}
	}

	return req, nil
}

func (v *MsxRequest) newHttpClientDo(ctx context.Context) (httpclient.DoFunc, error) {
	factory := httpclient.FactoryFromContext(ctx)
	if factory == nil {
		return nil, errors.New("Failed to retrieve http client factory from context")
	}

	httpClient := factory.NewHttpClient()
	httpClientDo := loginterceptor.NewInterceptor(httpClient.Do)
	switch v.Target.ServiceType {
	case ServiceTypeMicroservice:
		httpClientDo = discoveryinterceptor.NewInterceptor(httpClientDo)
	case ServiceTypeResourceProvider:
		httpClientDo = rpinterceptor.NewInterceptor(httpClientDo)
	case ServiceTypeProbe:
		// Integration is tied to the particular discovery.ServiceInstance
	default:
		return nil, errors.Errorf("Unknown service type %q", v.Target.ServiceType)
	}
	httpClientDo = traceinterceptor.NewInterceptor(httpClientDo)

	return httpClientDo, nil
}

func (v *MsxRequest) Execute(ctx context.Context) (response *MsxResponse, err error) {
	httpRequest, err := v.newHttpRequest(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create request")
	}

	httpRequest = httpRequest.WithContext(httpclient.ContextWithOperationName(
		httpRequest.Context(),
		v.Target.ServiceName+"."+v.Target.EndpointName))

	httpClientDo, err := v.newHttpClientDo(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create http client")
	}

	var resp *http.Response
	if resp, err = httpClientDo(httpRequest); err != nil {
		return nil, errors.Wrap(err, "Failed to execute request")
	}
	defer resp.Body.Close()

	ctx = httpRequest.Context()

	response = &MsxResponse{}
	if response.Body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, errors.Wrap(err, "Failed to read response body")
	}

	response.BodyString = string(response.Body)
	response.StatusCode = resp.StatusCode
	response.Status = resp.Status
	response.Headers = resp.Header

	if response.StatusCode > 399 {
		err = v.UnmarshalError(ctx, response)
	} else {
		err = v.UnmarshalSuccess(ctx, response)
	}
	return
}

func (v *MsxRequest) UnmarshalError(ctx context.Context, response *MsxResponse) (err error) {
	switch {
	case response.Body == nil:
		logger.WithContext(ctx).Debug("No body to unmarshal")

	case v.ExpectEnvelope:
		response.Envelope = &MsxEnvelope{}
		if err = json.Unmarshal(response.Body, response.Envelope); err != nil {
			return errors.Wrap(err, "Failed to unmarshal envelope")
		}

		// Extract the envelope error if possible
		if response.Envelope.IsError() {
			return response.Envelope.Error()
		}

	case v.ErrorPayload == nil && v.Payload == nil:
		logger.WithContext(ctx).Debug("No payload defined")

	case !v.ExpectEnvelope:
		if v.ErrorPayload == nil {
			v.ErrorPayload = v.Payload
		}

		response.Payload = v.ErrorPayload
		if err = json.Unmarshal(response.Body, response.Payload); err == nil {
			// Extract the payload error if possible
			if errorPayload, ok := response.Payload.(ResponseError); ok {
				if errorPayload.IsError() {
					return errorPayload.Error()
				}
			}
		}
	}

	// Check for a common MSX error formats
	if response.Body != nil {
		oauthErrorPayload := &OAuthErrorDTO{}
		if err = json.Unmarshal(response.Body, oauthErrorPayload); err == nil && oauthErrorPayload.IsError() {
			return oauthErrorPayload.Error()
		}

		errorDtoPayload := &ErrorDTO{}
		if err = json.Unmarshal(response.Body, errorDtoPayload); err == nil && errorDtoPayload.IsError() {
			return errorDtoPayload.Error()
		}

		errorDto2Payload := &ErrorDTO2{}
		if err = json.Unmarshal(response.Body, errorDto2Payload); err == nil && errorDto2Payload.IsError() {
			return errorDto2Payload.Error()
		}

		errorDto3Payload := &ErrorDTO3{}
		if err = json.Unmarshal(response.Body, errorDto3Payload); err == nil && errorDto3Payload.IsError() {
			return errorDto3Payload.Error()
		}
	}

	// Return a generic HTTP status error
	return NewStatusError(response.StatusCode, response.Body)
}

func (v *MsxRequest) UnmarshalSuccess(ctx context.Context, response *MsxResponse) (err error) {
	switch {
	case response.Body == nil:
		logger.WithContext(ctx).Debug("No body to unmarshal")

	case v.ExpectEnvelope:
		// Unmarshal the envelope and payload
		response.Envelope = &MsxEnvelope{Payload: v.Payload}
		if err = json.Unmarshal(response.Body, response.Envelope); err != nil {
			return errors.Wrap(err, "Failed to unmarshal envelope")
		}
		response.Payload = response.Envelope.Payload

	case v.Payload == nil:
		logger.WithContext(ctx).Debug("No payload defined")

	case !v.ExpectEnvelope:
		// Unmarshal the raw payload
		response.Payload = v.Payload
		if err = json.Unmarshal(response.Body, response.Payload); err != nil {
			return errors.Wrap(err, "Failed to unmarshal payload")
		}
	}

	return
}

func (v *MsxRequest) Operation() string {
	return fmt.Sprintf("%s.%s", v.Target.ServiceName, v.Target.EndpointName)
}
