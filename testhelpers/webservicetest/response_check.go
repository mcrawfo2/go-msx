// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package webservicetest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/tidwall/gjson"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type ResponseVerifier func(t *testing.T, resp *httptest.ResponseRecorder)

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
	return fmt.Sprintf("Failed response validator: %s\n%s", c.Validator.Description)
}

type ResponsePredicate struct {
	Description string
	Matches     func(*httptest.ResponseRecorder) bool
	Diff        func(*httptest.ResponseRecorder) string
}

func ResponseHasStatus(status int) ResponsePredicate {
	return ResponsePredicate{
		Description: fmt.Sprintf("response.Code == %d", status),
		Matches: func(resp *httptest.ResponseRecorder) bool {
			return resp.Code == status
		},
		Diff: func(resp *httptest.ResponseRecorder) string {
			return testhelpers.Diff(status, resp.Code)
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
		Diff: func(resp *httptest.ResponseRecorder) string {
			return testhelpers.Diff([]string{value}, resp.Header().Values(header))
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
		Diff: func(resp *httptest.ResponseRecorder) string {
			return testhelpers.Diff(substring, resp.Body.String())
		},
	}
}

func ResponseHasBodyJson(path string, eq func(gjson.Result) bool) ResponsePredicate {
	return ResponsePredicate{
		Description: fmt.Sprintf("response.Body[%q] passes fn", path),
		Matches: func(resp *httptest.ResponseRecorder) bool {
			return eq(gjson.GetBytes(resp.Body.Bytes(), path))
		},
		Diff: func(resp *httptest.ResponseRecorder) string {
			return ""
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
		Diff: func(resp *httptest.ResponseRecorder) string {
			actualValue := gjson.GetBytes(resp.Body.Bytes(), path).Value()
			return testhelpers.Diff(value, actualValue)
		},
	}
}

func ResponseHasBodyJsonValueType(path string, kind gjson.Type) ResponsePredicate {
	return ResponsePredicate{
		Description: fmt.Sprintf("response.Body[%q] == %s", path, kind.String()),
		Matches: func(resp *httptest.ResponseRecorder) bool {
			return gjson.GetBytes(resp.Body.Bytes(), path).Type == kind
		},
		Diff: func(resp *httptest.ResponseRecorder) string {
			actualType := gjson.GetBytes(resp.Body.Bytes(), path).Type
			return testhelpers.Diff(kind, actualType)
		},
	}
}
