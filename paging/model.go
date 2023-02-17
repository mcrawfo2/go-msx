// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package paging

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"net/url"
	"reflect"
	"strconv"
)

type Response struct {
	Content    interface{}
	Size       uint
	Number     uint
	TotalItems *uint
	Sort       []SortOrder
	State      interface{}
}

func (s Response) Elements() uint {
	if s.Content == nil {
		return 0
	}

	contentValue := reflect.ValueOf(s.Content)
	contentValueKind := contentValue.Kind()
	if contentValueKind == reflect.Ptr {
		contentValue = contentValue.Elem()
		contentValueKind = contentValue.Kind()
	}
	switch contentValueKind {
	case reflect.Slice, reflect.Array:
		return uint(contentValue.Len())
	default:
		return 1
	}
}

func (s Response) HasNext() bool {
	return s.Elements() == s.Size
}

func (s Response) Offset() uint {
	return s.Size * s.Number
}

type Request struct {
	Page  uint
	Size  uint
	Sort  []SortOrder
	State interface{}
}

func (r Request) WithState(state *string) Request {
	r.State = state
	return r
}

func (r Request) QueryParameters() url.Values {
	var result = make(url.Values)
	result.Set("page", strconv.FormatUint(uint64(r.Page), 10))
	result.Set("pageSize", strconv.FormatUint(uint64(r.Size), 10))
	return result
}

func NewRequestFromQuery(page uint, pageSize uint) Request {
	return Request{Page: page, Size: pageSize}
}

type SortDirection string

const (
	SortDirectionAsc  SortDirection = "ASC"
	SortDirectionDesc SortDirection = "DESC"
)

type SortOrder struct {
	Property  string
	Direction SortDirection
}

type SortByOptions struct {
	// DefaultProperty is the API-facing name of the default sort field
	DefaultProperty string
	// Mapping contains a list of API (left) <-> DB (right) field name mappings for sort fields
	Mapping types.StringPairSlice
}
