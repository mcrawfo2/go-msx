// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package types

type Slice[I any] []I

func (s Slice[I]) AnySlice() (results []any) {
	results = make([]any, 0, len(s))
	for _, v := range s {
		results = append(results, v)
	}
	return results
}

func (s Slice[I]) RemoveAll(indices ...int) (results Slice[I]) {
	results = make([]I, len(s)-len(indices), len(s))

	var out = 0
	var in = 0

	for _, removal := range indices {
		if removal > in {
			copy(results[out:], s[in:removal])
			out = out + removal - in
		}
		in = removal + 1
	}
	if in < len(s) {
		copy(results[out:], s[in:])
	}

	return
}

type ComparableSlice[I comparable] []I

func (s ComparableSlice[I]) Contains(value I) bool {
	for _, item := range s {
		if item == value {
			return true
		}
	}
	return false
}

func (s ComparableSlice[I]) ContainsAny(values ...I) bool {
	for _, item := range s {
		for _, value := range values {
			if item == value {
				return true
			}
		}
	}
	return false
}
