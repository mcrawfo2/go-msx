// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package testhelpers

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/pmezard/go-difflib/difflib"
	"reflect"
)

func typeAndKind(v interface{}) (reflect.Type, reflect.Kind) {
	t := reflect.TypeOf(v)
	k := t.Kind()

	if k == reflect.Ptr {
		t = t.Elem()
		k = t.Kind()
	}
	return t, k
}

// Diff returns a diff of both values as long as both are of the same type and
// are a struct, map, slice or array. Otherwise it returns an empty string.
func Diff(expected interface{}, actual interface{}) string {
	return DiffWithConfig(expected, actual, spewConfig)
}

func DiffWithConfig(expected, actual any, spewConfig spew.ConfigState) string {
	if expected == nil || actual == nil {
		return ""
	}

	et, ek := typeAndKind(expected)
	at, _ := typeAndKind(actual)

	if et != at {
		return fmt.Sprintf("Types do not match: expected %q, actual %q", et.String(), at.String())
	}

	switch ek {
	case reflect.String, reflect.Bool, reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
		break

	default:
		return ""

	}

	e := spewConfig.Sdump(expected)
	a := spewConfig.Sdump(actual)

	diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(e),
		B:        difflib.SplitLines(a),
		FromFile: "Expected",
		FromDate: "",
		ToFile:   "Actual",
		ToDate:   "",
		Context:  1,
	})

	return diff
}

func Dump(actual interface{}) string {
	return DumpWithConfig(actual, spewConfig)
}

func DumpWithConfig(actual any, spewConfig spew.ConfigState) string {
	if actual == nil {
		return "<nil>"
	}

	_, ak := typeAndKind(actual)
	if ak != reflect.Struct && ak != reflect.Map && ak != reflect.Slice && ak != reflect.Array {
		return fmt.Sprintf("%+v", actual)
	}

	return spewConfig.Sdump(actual)
}

var spewConfig = spew.ConfigState{
	Indent:                  " ",
	DisablePointerAddresses: true,
	DisableCapacities:       true,
	SortKeys:                true,
}
