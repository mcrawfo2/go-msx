package integration

import (
	"bytes"
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient/discoveryinterceptor"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient/loginterceptor"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient/rpinterceptor"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient/traceinterceptor"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"text/template"
)

var (
	logger = log.NewLogger("msx.integration")
)

type MsxServiceEndpoint struct {
	Method string
	Path   string
}

func (e MsxServiceEndpoint) Url(endpointName, serviceName string, variables map[string]string, queryParameters url.Values) (string, error) {
	subPathTemplate, err := template.New(endpointName).Parse(e.Path)
	if err != nil {
		return "", errors.Wrap(err, "Failed to parse e url template")
	}

	var pathBuffer bytes.Buffer
	pathBuffer.WriteString("http://")
	pathBuffer.WriteString(serviceName)

	err = subPathTemplate.Execute(&pathBuffer, variables)
	if err != nil {
		return "", errors.Wrap(err, "Failed to fill url template")
	}

	if len(queryParameters) > 0 {
		pathBuffer.WriteString("?")
		pathBuffer.WriteString(queryParameters.Encode())
	}

	result := pathBuffer.String()
	return result, nil
}

type MsxService struct {
	serviceName      string
	endpoints        map[string]MsxServiceEndpoint
	resourceProvider bool
	ctx              context.Context
}

func (v *MsxService) newHttpRequest(r *MsxRequest) (*http.Request, error) {
	var req *http.Request
	var err error
	var buf io.Reader

	endpoint, ok := v.endpoints[r.EndpointName]
	if !ok {
		return nil, errors.New("No endpoint defined: " + r.EndpointName)
	}

	fullUrl, err := endpoint.Url(
		r.EndpointName,
		v.serviceName,
		r.EndpointParameters,
		r.QueryParameters)
	if err != nil {
		return nil, err
	}

	if r.Body != nil {
		buf = bytes.NewBuffer(r.Body)
	} else {
		buf = http.NoBody
	}

	req, err = http.NewRequest(endpoint.Method, fullUrl, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	if r.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if !r.NoToken {
		userContext := security.UserContextFromContext(v.ctx)
		if userContext.Token == "" {
			logger.Warn("Token required but not set")
		} else {
			req.Header.Set("Authorization", "Bearer "+userContext.Token)
		}
	}

	for k, vs := range r.Headers {
		for _, v := range vs {
			req.Header.Set(k, v)
		}
	}

	return req.WithContext(v.ctx), nil
}

func (v *MsxService) newHttpClientDo() (httpclient.DoFunc, error) {
	factory := httpclient.FactoryFromContext(v.ctx)
	if factory == nil {
		return nil, errors.New("Failed to retrieve http client factory from context")
	}

	httpClient := factory.NewHttpClient()
	httpClientDo := loginterceptor.NewInterceptor(httpClient.Do)
	if !v.resourceProvider {
		httpClientDo = discoveryinterceptor.NewInterceptor(httpClientDo)
	} else {
		httpClientDo = rpinterceptor.NewInterceptor(httpClientDo)
	}
	httpClientDo = traceinterceptor.NewInterceptor(httpClientDo)

	return httpClientDo, nil
}

func (v *MsxService) Execute(request *MsxRequest) (response *MsxResponse, err error) {
	httpRequest, err := v.newHttpRequest(request)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create request")
	}

	httpRequest = httpRequest.WithContext(httpclient.ContextWithOperationName(
		httpRequest.Context(),
		v.serviceName+"."+request.EndpointName))

	httpClientDo, err := v.newHttpClientDo()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create http client")
	}

	var resp *http.Response
	if resp, err = httpClientDo(httpRequest); err != nil {
		return nil, errors.Wrap(err, "Failed to execute request")
	}
	defer resp.Body.Close()

	ctx := httpRequest.Context()

	response = &MsxResponse{}
	if response.Body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, errors.Wrap(err, "Failed to read response body")
	}

	response.BodyString = string(response.Body)
	response.StatusCode = resp.StatusCode
	response.Status = resp.Status
	response.Headers = resp.Header

	if response.StatusCode > 399 {
		err = v.UnmarshalError(ctx, request, response)
	} else {
		err = v.UnmarshalSuccess(ctx, request, response)
	}
	return
}

func (v *MsxService) UnmarshalError(ctx context.Context, request *MsxRequest, response *MsxResponse) (err error) {
	switch {
	case response.Body == nil:
		logger.WithContext(ctx).Debug("No body to unmarshal")

	case request.ExpectEnvelope:
		response.Envelope = &MsxEnvelope{}
		if err = json.Unmarshal(response.Body, response.Envelope); err != nil {
			return errors.Wrap(err, "Failed to unmarshal envelope")
		}

		// Extract the envelope error if possible
		if response.Envelope.IsError() {
			return response.Envelope.Error()
		}

	case request.ErrorPayload == nil && request.Payload == nil:
		logger.WithContext(ctx).Debug("No payload defined")

	case !request.ExpectEnvelope:
		if request.ErrorPayload == nil {
			request.ErrorPayload = request.Payload
		}
		response.Payload = request.ErrorPayload
		if err = json.Unmarshal(response.Body, response.Payload); err != nil {
			return errors.Wrap(err, "Failed to unmarshal payload")
		}

		// Extract the payload error if possible
		if errorPayload, ok := response.Payload.(ResponseError); ok {
			if errorPayload.IsError() {
				return errorPayload.Error()
			}
		}
	}

	// Check for a spring authentication error
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
	}

	// Return a generic HTTP status error
	return NewStatusError(response.StatusCode, response.Body)
}

func (v *MsxService) UnmarshalSuccess(ctx context.Context, request *MsxRequest, response *MsxResponse) (err error) {
	switch {
	case response.Body == nil:
		logger.WithContext(ctx).Debug("No body to unmarshal")

	case request.ExpectEnvelope:
		// Unmarshal the envelope and payload
		response.Payload = request.Payload
		response.Envelope = &MsxEnvelope{Payload: response.Payload}
		if err = json.Unmarshal(response.Body, response.Envelope); err != nil {
			return errors.Wrap(err, "Failed to unmarshal envelope")
		}

	case request.Payload == nil:
		logger.WithContext(ctx).Debug("No payload defined")

	case !request.ExpectEnvelope:
		// Unmarshal the raw payload
		response.Payload = request.Payload
		if err = json.Unmarshal(response.Body, response.Payload); err != nil {
			return errors.Wrap(err, "Failed to unmarshal payload")
		}
	}

	return
}

func (v *MsxService) Operation(request *MsxRequest) string {
	return fmt.Sprintf("%s.%s", v.serviceName, request.EndpointName)
}

func NewMsxService(ctx context.Context, serviceName string, endpoints map[string]MsxServiceEndpoint) *MsxService {
	return &MsxService{
		serviceName: serviceName,
		endpoints:   endpoints,
		ctx:         ctx,
	}
}

func NewMsxServiceResourceProvider(ctx context.Context, serviceName string, endpoints map[string]MsxServiceEndpoint) *MsxService {
	return &MsxService{
		serviceName:      serviceName,
		endpoints:        endpoints,
		resourceProvider: true,
		ctx:              ctx,
	}
}
