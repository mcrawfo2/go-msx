package integration

import (
	"bytes"
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/discovery"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"text/template"
)

var (
	logger = log.NewLogger("msx.integration")
)

type MsxServiceEndpoint struct {
	Method string
	Path   string
}

type MsxService struct {
	serviceName string
	contextPath string
	endpoints   map[string]MsxServiceEndpoint
	ctx         context.Context
}

func (v *MsxService) NewHttpRequest(r *MsxRequest) (*http.Request, error) {
	var req *http.Request
	var err error
	var buf io.Reader

	endpoint, ok := v.endpoints[r.EndpointName]
	if !ok {
		return nil, errors.New("No endpoint defined: " + r.EndpointName)
	}

	fullUrl, err := v.endpointUrl(r.EndpointName, r.EndpointParameters, r.QueryParameters)
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

	return req, nil
}

func (v *MsxService) discoverService() (*discovery.ServiceInstance, error) {
	logger.Debugf("Discovering %s", v.serviceName)
	if instances, err := discovery.Discover(v.ctx, v.serviceName, true); err != nil {
		return nil, err
	} else if len(instances) == 0 {
		return nil, errors.New(fmt.Sprintf("No healthy instances of %s found", v.serviceName))
	} else {
		return instances.SelectRandom(), nil
	}
}

func (v *MsxService) endpointUrl(endpointName string, variables map[string]string, queryParameters url.Values) (string, error) {
	endpoint, ok := v.endpoints[endpointName]
	if !ok {
		return "", errors.New("No endpoint defined: " + endpointName)
	}

	subPathTemplate, err := template.New(endpointName).Parse(endpoint.Path)
	if err != nil {
		return "", errors.Wrap(err, "Failed to parse endpoint url template: "+endpointName)
	}

	var pathBuffer bytes.Buffer
	pathBuffer.WriteString("http://")

	if serviceInstance, err := v.discoverService(); err != nil {
		return "", errors.Wrap(err, "Failed to discover service instance")
	} else {
		pathBuffer.WriteString(serviceInstance.Host)
		pathBuffer.WriteString(":")
		pathBuffer.WriteString(strconv.Itoa(serviceInstance.Port))
	}

	pathBuffer.WriteString(v.contextPath)

	err = subPathTemplate.Execute(&pathBuffer, variables)
	if err != nil {
		return "", errors.Wrap(err, "Failed to fill url template: "+endpointName)
	}

	if len(queryParameters) > 0 {
		pathBuffer.WriteString("?")
		pathBuffer.WriteString(queryParameters.Encode())
	}

	result := pathBuffer.String()
	return result, nil
}

func (v *MsxService) newHttpClient() (*http.Client, error) {
	// TODO: context http client factory
	return http.DefaultClient, nil
}

func (v *MsxService) Execute(request *MsxRequest) (response *MsxResponse, err error) {
	httpRequest, err := v.NewHttpRequest(request)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create request")
	}

	httpClient, err := v.newHttpClient()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create http client")
	}

	response = &MsxResponse{}

	var resp *http.Response
	if resp, err = httpClient.Do(httpRequest); err != nil {
		return nil, errors.Wrap(err, "Failed to execute request")
	}

	defer resp.Body.Close()
	if response.Body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, errors.Wrap(err, "Failed to read response body")
	}

	response.BodyString = string(response.Body)
	response.StatusCode = resp.StatusCode
	response.Status = resp.Status
	response.Headers = resp.Header

	if response.StatusCode > 399 {
		// Fully log the response
		logger.Errorf("%s : %s", response.Status, httpRequest.URL.String())
		var responseBytes []byte
		responseBytes, _ = json.Marshal(response)
		logger.Error(string(responseBytes))

		return response, v.UnmarshalError(request, response)
	} else {
		logger.Infof("%s : %s", response.Status, httpRequest.URL.String())

		err = v.UnmarshalSuccess(request, response)
		return
	}
}

func (v *MsxService) UnmarshalError(request *MsxRequest, response *MsxResponse) (err error) {
	switch {
	case response.Body == nil:
		logger.Debug("No body to unmarshal")

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
		logger.Debug("No payload defined")

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
			return errorDtoPayload.Error()
		}
	}

	// Return a generic HTTP status error
	return NewStatusError(response.StatusCode, response.Body)
}

func (v *MsxService) UnmarshalSuccess(request *MsxRequest, response *MsxResponse) (err error) {
	switch {
	case response.Body == nil:
		logger.Debug("No body to unmarshal")

	case request.ExpectEnvelope:
		// Unmarshal the envelope and payload
		response.Payload = request.Payload
		response.Envelope = &MsxEnvelope{Payload: response.Payload}
		if err = json.Unmarshal(response.Body, response.Envelope); err != nil {
			return errors.Wrap(err, "Failed to unmarshal envelope")
		}

	case request.Payload == nil:
		logger.Debug("No payload defined")

	case !request.ExpectEnvelope:
		// Unmarshal the raw payload
		response.Payload = request.Payload
		if err = json.Unmarshal(response.Body, response.Payload); err != nil {
			return errors.Wrap(err, "Failed to unmarshal payload")
		}
	}

	return
}

func NewMsxService(ctx context.Context, serviceName, contextPath string, endpoints map[string]MsxServiceEndpoint) *MsxService {
	return &MsxService{
		serviceName: serviceName,
		contextPath: contextPath,
		endpoints:   endpoints,
		ctx:         ctx,
	}
}
