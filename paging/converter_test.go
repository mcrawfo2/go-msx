// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package paging

import (
	"reflect"
	"testing"
	"time"
)

type Payload struct {
	Id         string
	CreatedBy  string
	CreatedOn  time.Time
	ModifiedBy string
	ModifiedOn time.Time
	ClosedBy   string
	ClosedOn   time.Time
}

func TestConverter_ResponseToPaginatedResponseV8(t *testing.T) {
	uintOneZero := uint(0)
	uintOneHundred := uint(100)
	int64Zero := int64(0)
	int64OneHundred := int64(100)

	payload := Payload{
		Id:         "",
		CreatedBy:  "",
		CreatedOn:  time.Time{},
		ModifiedBy: "",
		ModifiedOn: time.Time{},
		ClosedBy:   "",
		ClosedOn:   time.Time{},
	}

	var contents []interface{}
	contents = append(contents, payload)

	type args struct {
		response Response
		objects  []interface{}
	}

	tests := []struct {
		name string
		args args
		want PaginatedResponseV8
	}{
		{
			name: "TotalItems_nil",
			args: args{
				response: Response{
					Size:   1,
					Number: 10,
				},
				objects: contents,
			},
			want: PaginatedResponseV8{
				Page:        10,
				PageSize:    1,
				TotalItems:  nil,
				HasNext:     false,
				HasPrevious: true,
				SortBy:      "",
				SortOrder:   "",
				Contents:    contents,
			},
		},
		{
			name: "TotalItems_0",
			args: args{
				response: Response{
					Size:       1,
					Number:     10,
					TotalItems: &uintOneZero,
				},
				objects: contents,
			},
			want: PaginatedResponseV8{
				Page:        10,
				PageSize:    1,
				TotalItems:  &int64Zero,
				HasNext:     false,
				HasPrevious: true,
				SortBy:      "",
				SortOrder:   "",
				Contents:    contents,
			},
		},
		{
			name: "TotalItems_100",
			args: args{
				response: Response{
					Size:       1,
					Number:     10,
					TotalItems: &uintOneHundred,
				},
				objects: contents,
			},
			want: PaginatedResponseV8{
				Page:        10,
				PageSize:    1,
				TotalItems:  &int64OneHundred,
				HasNext:     false,
				HasPrevious: true,
				SortBy:      "",
				SortOrder:   "",
				Contents:    contents,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Converter{}
			if got := c.ResponseToPaginatedResponseV8(tt.args.response, tt.args.objects); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResponseToPaginatedResponseV8() = %v, want %v", got, tt.want)
			}
		})
	}
}
