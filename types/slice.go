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
