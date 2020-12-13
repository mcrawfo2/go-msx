package webservicetest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/tidwall/gjson"
	"net/http/httptest"
	"reflect"
	"strings"
)

type ResponseCheck struct {
	Validators []ResponsePredicate
}

func (r ResponseCheck) Check(resp *httptest.ResponseRecorder) []error {
	var results []error

	for _, predicate := range r.Validators {
		if !predicate.Matches(resp) {
			results = append(results, ResponseCheckError{
				Validator: predicate,
			})
		}
	}

	return results

}

type ResponseCheckError struct {
	Validator ResponsePredicate
}

func (c ResponseCheckError) Error() string {
	return fmt.Sprintf("Failed response validator: %s", c.Validator.Description)
}

type ResponsePredicate struct {
	Description string
	Matches     func(*httptest.ResponseRecorder) bool
}

func ResponseHasStatus(status int) ResponsePredicate {
	return ResponsePredicate{
		Description: fmt.Sprintf("response.Code == %d", status),
		Matches: func(resp *httptest.ResponseRecorder) bool {
			return resp.Code == status
		},
	}
}

func ResponseHasHeader(header, value string) ResponsePredicate {
	return ResponsePredicate{
		Description: fmt.Sprintf("response.Header[%s] contains %q", header, value),
		Matches: func(resp *httptest.ResponseRecorder) bool {
			values := types.StringStack(resp.Header().Values(header))
			return values.Contains(value)
		},
	}
}

func ResponseHasBodySubstring(substring string) ResponsePredicate {
	return ResponsePredicate{
		Description: fmt.Sprintf("response.Body contains %q", substring),
		Matches: func(resp *httptest.ResponseRecorder) bool {
			body := resp.Body.String()
			return strings.Contains(body, substring)
		},
	}
}

func ResponseHasBodyJson(path string, eq func(gjson.Result) bool) ResponsePredicate {
	return ResponsePredicate{
		Description: fmt.Sprintf("response.Body[%q] passes fn", path),
		Matches: func(resp *httptest.ResponseRecorder) bool {
			return eq(gjson.GetBytes(resp.Body.Bytes(), path))
		},
	}
}

func ResponseHasBodyJsonValue(path string, value interface{}) ResponsePredicate {
	return ResponsePredicate{
		Description: fmt.Sprintf("response.Body[%q] == %v", path, value),
		Matches: func(resp *httptest.ResponseRecorder) bool {
			actualValue := gjson.GetBytes(resp.Body.Bytes(), path).Value()
			return reflect.DeepEqual(actualValue, value)
		},
	}
}
