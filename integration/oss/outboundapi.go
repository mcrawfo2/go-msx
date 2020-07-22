package oss

import (
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"path"
)

const CredentialTypeNone = "none"
const CredentialTypeBasic = "basic"

const (
	OutboundApiServiceCancellationCharge = "serviceCancellationCharge"
	OutboundApiNotificationUrl           = "notificationUrl"
	OutboundApiPricingOptions            = "pricingoptions"
	OutboundApiServiceAccess             = "serviceAccess"
	OutboundApiAllowedValues             = "allowedValues"
)

var OutboundApiErrDisabled = errors.New("Outbound API endpoint disabled")

type OutboundApi struct {
	ApiName                string             `json:"apiName"`
	BaseContext            string             `json:"baseContext"`
	HttpMethod             string             `json:"httpMethod"`
	Enabled                bool               `json:"enabled"`
	AuthenticationResource string             `json:"authenticationResource"`
	Url                    string             `json:"url"`
	RestCredentialType     *string            `json:"restCredentialType"`
	RestConfigAttributes   map[string]*string `json:"restConfigAttributes"`
}

func (o OutboundApi) Interceptor(next httpclient.DoFunc) httpclient.DoFunc {
	return OutboundApiInterceptor{OutboundApi: o}.intercept(next)
}

type OutboundApiInterceptor struct {
	OutboundApi
}

func (o OutboundApiInterceptor) intercept(next httpclient.DoFunc) httpclient.DoFunc {
	return func(req *http.Request) (*http.Response, error) {
		// Ensure the API is enabled
		if !o.Enabled {
			return nil, OutboundApiErrDisabled
		}

		// Relocate the api onto the api.BaseContext
		relocatedUrl, err := o.relocate(*req.URL)
		if err != nil {
			return nil, err
		}
		req.URL = &relocatedUrl

		// Apply any authentication headers
		err = o.authenticate(req)
		if err != nil {
			return nil, err
		}

		return next(req)
	}
}

func (o OutboundApiInterceptor) relocate(requestUrl url.URL) (url.URL, error) {
	apiBaseUrl, err := url.Parse(o.BaseContext)
	if err != nil {
		return url.URL{}, err
	}

	requestUrl.Host = apiBaseUrl.Host
	requestUrl.Scheme = apiBaseUrl.Scheme
	requestUrl.Path = path.Join(apiBaseUrl.Path, requestUrl.Path)
	return requestUrl, nil
}

func (o OutboundApiInterceptor) authenticate(req *http.Request) error {
	credentialType := types.NewOptionalString(o.RestCredentialType).OrElse(CredentialTypeNone)

	switch credentialType {
	case CredentialTypeBasic:
		req.SetBasicAuth(
			types.NewOptionalString(o.RestConfigAttributes["basic_username"]).OrElse(""),
			types.NewOptionalString(o.RestConfigAttributes["basic_password"]).OrElse(""))
		return nil

	case CredentialTypeNone:
		return nil

	default:
		return errors.Errorf("Authentication method not supported: %s", credentialType)
	}
}
