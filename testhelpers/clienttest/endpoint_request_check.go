package clienttest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/tidwall/gjson"
	"reflect"
	"strings"
	"testing"
)

type EndpointRequestVerifier func(t *testing.T, req *integration.MsxEndpointRequest)

type EndpointRequestCheck struct {
	Validators []EndpointRequestPredicate
}

func (r EndpointRequestCheck) Check(req *integration.MsxEndpointRequest) []error {
	var results []error

	for _, predicate := range r.Validators {
		if !predicate.Matches(req) {
			results = append(results, EndpointRequestCheckError{
				Validator: predicate,
			})
		}
	}

	return results

}

type EndpointRequestCheckError struct {
	Validator EndpointRequestPredicate
}

func (c EndpointRequestCheckError) Error() string {
	return fmt.Sprintf("Failed Request validator: %s", c.Validator.Description)
}

type EndpointRequestPredicate struct {
	Description string
	Matches     func(*integration.MsxEndpointRequest) bool
}

func EndpointRequestHasQueryParam(name, value string) EndpointRequestPredicate {
	return EndpointRequestPredicate{
		Description: fmt.Sprintf("request.QueryParameters.Has(%q, %q)", name, value),
		Matches: func(request *integration.MsxEndpointRequest) bool {
			values := request.QueryParameters[name]
			return types.StringStack(values).Contains(value)
		},
	}
}

func EndpointRequestHasName(name string) EndpointRequestPredicate {
	return EndpointRequestPredicate{
		Description: fmt.Sprintf("request.EndpointName == %q", name),
		Matches: func(request *integration.MsxEndpointRequest) bool {
			return request.EndpointName == name
		},
	}
}

func EndpointRequestHasToken(token bool) EndpointRequestPredicate {
	return EndpointRequestPredicate{
		Description: fmt.Sprintf("request.NoToken != %v", token),
		Matches: func(request *integration.MsxEndpointRequest) bool {
			return request.NoToken == !token
		},
	}
}

func EndpointRequestHasExpectEnvelope(envelope bool) EndpointRequestPredicate {
	return EndpointRequestPredicate{
		Description: fmt.Sprintf("request.ExpectEnvelope == %v", envelope),
		Matches: func(request *integration.MsxEndpointRequest) bool {
			return request.ExpectEnvelope == envelope
		},
	}
}

func EndpointRequestHasEndpointParameter(name, value string) EndpointRequestPredicate {
	return EndpointRequestPredicate{
		Description: fmt.Sprintf("request.EndpointParameters[%q] == %q", name, value),
		Matches: func(request *integration.MsxEndpointRequest) bool {
			return request.EndpointParameters[name] == value
		},
	}
}

func EndpointRequestHasHeader(name, value string) EndpointRequestPredicate {
	return EndpointRequestPredicate{
		Description: fmt.Sprintf("request.Headers[%q] contains %q", name, value),
		Matches: func(request *integration.MsxEndpointRequest) bool {
			for _, header := range request.Headers[name] {
				if header == value {
					return true
				}
			}
			return false
		},
	}
}

func EndpointRequestHasBodySubstring(substring string) EndpointRequestPredicate {
	return EndpointRequestPredicate{
		Description: fmt.Sprintf("endpointRequest.Body contains %q", substring),
		Matches: func(request *integration.MsxEndpointRequest) bool {
			body := string(request.Body)
			return strings.Contains(body, substring)
		},
	}
}

func EndpointRequestHasBodyJson(path string, eq func(gjson.Result) bool) EndpointRequestPredicate {
	return EndpointRequestPredicate{
		Description: fmt.Sprintf("endpointRequest.Body[%q] passes fn", path),
		Matches: func(request *integration.MsxEndpointRequest) bool {
			return eq(gjson.GetBytes(request.Body, path))
		},
	}
}

func EndpointRequestHasBodyJsonValue(path string, value interface{}) EndpointRequestPredicate {
	return EndpointRequestPredicate{
		Description: fmt.Sprintf("endpointRequest.Body[%q] == %v", path, value),
		Matches: func(request *integration.MsxEndpointRequest) bool {
			actualValue := gjson.GetBytes(request.Body, path).Value()
			return reflect.DeepEqual(actualValue, value)
		},
	}
}

func EndpointRequestHasBodyJsonValueType(path string, kind gjson.Type) EndpointRequestPredicate {
	return EndpointRequestPredicate{
		Description: fmt.Sprintf("endpointRequest.Body[%q] == %s", path, kind.String()),
		Matches: func(request *integration.MsxEndpointRequest) bool {
			return gjson.GetBytes(request.Body, path).Type == kind
		},
	}
}
