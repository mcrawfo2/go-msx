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
