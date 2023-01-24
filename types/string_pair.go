// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package types

type StringPair struct {
	Left  string
	Right string
}

type StringPairSlice []StringPair

func (s StringPairSlice) MapToRight(left string) (string, bool) {
	for _, pair := range s {
		if pair.Left == left {
			return pair.Right, true
		}
	}
	return left, false
}

func (s StringPairSlice) MapToLeft(right string) (string, bool) {
	for _, pair := range s {
		if pair.Right == right {
			return pair.Left, true
		}
	}
	return right, false
}
