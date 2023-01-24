// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSlice_RemoveAll(t *testing.T) {
	type testCase[I any] struct {
		name    string
		s       Slice[I]
		indices []int
		want    Slice[I]
	}
	tests := []testCase[int]{
		{
			name:    "First",
			s:       []int{1, 2, 3},
			indices: []int{0},
			want:    []int{2, 3},
		},
		{
			name:    "Second",
			s:       []int{1, 2, 3},
			indices: []int{1},
			want:    []int{1, 3},
		},
		{
			name:    "Third",
			s:       []int{1, 2, 3},
			indices: []int{2},
			want:    []int{1, 2},
		},
		{
			name:    "Consecutive",
			s:       []int{1, 2, 3},
			indices: []int{1, 2},
			want:    []int{1},
		},
		{
			name:    "Disjoint",
			s:       []int{1, 2, 3},
			indices: []int{0, 2},
			want:    []int{2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.s.RemoveAll(tt.indices...))
		})
	}
}
