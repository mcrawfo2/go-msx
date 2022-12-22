// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package v8

import (
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"strings"
)

var ErrUnknownSortOrder = errors.New("Unknown sort order")

type PagingSortingInputs struct {
	PagingInputs
	SortingInputs
}

type PagingInputs struct {
	Page     int `req:"query" default:"0" example:"0" minimum:"0" format:"int32" required:"true" reference:"page"`
	PageSize int `req:"query" default:"100" example:"100" minimum:"1" format:"int32" required:"true" reference:"pageSize"`
}

type SortingInputs struct {
	SortBy    string `req:"query" default:"" optional:"true" reference:"sortBy"`
	SortOrder string `req:"query" default:"asc" enum:"asc,desc" optional:"true" reference:"sortOrder"`
}

type PagingOutputs struct {
	Paging PagingResponse `resp:"paging"`
}

type PagingResponse struct {
	Page        int         `json:"page" format:"int32"`
	PageSize    int         `json:"pageSize" format:"int32"`
	TotalItems  *int        `json:"totalItems,omitempty"`
	HasNext     bool        `json:"hasNext"`
	HasPrevious bool        `json:"hasPrevious"`
	SortBy      string      `json:"sortBy,omitempty"`
	SortOrder   string      `json:"sortOrder,omitempty" enum:"asc,desc"`
	Contents    interface{} `json:"contents" inject:"Page"`
}

func (r PagingResponse) Example() interface{} {
	return PagingResponse{
		Page:        0,
		PageSize:    10,
		TotalItems:  types.NewIntPtr(100),
		HasNext:     true,
		HasPrevious: false,
		SortBy:      "tenantId",
		SortOrder:   SortDirectionAsc,
		Contents:    []interface{}{},
	}
}

type SortDirection string

func (_ SortDirection) Enum() []any {
	return []any{
		SortDirectionAsc,
		SortDirectionDesc,
	}
}

const (
	SortDirectionAsc  = "ASC"
	SortDirectionDesc = "DESC"
)

type PagingConverter struct{}

func (c PagingConverter) FromPagingInputs(pageReq PagingInputs) (result paging.Request) {
	result.Page = uint(pageReq.Page)
	result.Size = uint(pageReq.PageSize)
	return
}

func (c PagingConverter) FromPagingSortingInputs(pageReq PagingSortingInputs) (result paging.Request) {
	result = c.FromSortingInputs(pageReq.SortingInputs)
	result.Page = uint(pageReq.Page)
	result.Size = uint(pageReq.PageSize)
	return
}

func (c PagingConverter) FromSortingInputs(sortReq SortingInputs) (result paging.Request) {
	if sortReq.SortBy != "" {
		sortResult := paging.SortOrder{
			Property:  sortReq.SortBy,
			Direction: paging.SortDirection(strings.ToUpper(sortReq.SortOrder)),
		}
		result.Sort = append(result.Sort, sortResult)
	}
	return
}

func (c PagingConverter) FromPageSortQuery(page, pageSize int, sortBy, sortOrder string) (request paging.Request, err error) {
	request = c.FromPageQuery(page, pageSize)

	sortIn, err := c.FromSortQuery(sortBy, sortOrder)
	if err != nil {
		return
	}
	request.Sort = []paging.SortOrder{sortIn}

	return
}

func (c PagingConverter) FromPageQuery(page, pageSize int) paging.Request {
	return paging.Request{
		Page: uint(page),
		Size: uint(pageSize),
	}
}

func (c PagingConverter) FromSortQuery(sortBy, sortOrder string) (paging.SortOrder, error) {
	if sortBy == "" {
		return paging.SortOrder{}, nil
	}

	if sortOrder == "" {
		sortOrder = string(paging.SortDirectionAsc)
	}

	var sortDirection paging.SortDirection
	switch paging.SortDirection(strings.ToUpper(sortOrder)) {
	case paging.SortDirectionDesc:
		sortDirection = paging.SortDirectionDesc
	case paging.SortDirectionAsc:
		sortDirection = paging.SortDirectionAsc
	default:
		return paging.SortOrder{}, errors.Wrap(ErrUnknownSortOrder, sortOrder)
	}

	return paging.SortOrder{
		Property:  sortBy,
		Direction: sortDirection,
	}, nil
}

func (c PagingConverter) ToPagingResponse(pout paging.Response) PagingResponse {
	presp := PagingResponse{
		Page:        int(pout.Number),
		PageSize:    int(pout.Size),
		HasNext:     pout.HasNext(),
		HasPrevious: pout.Offset() > 0,
	}

	totalItems := pout.TotalItems
	if pout.TotalItems != nil {
		localTotalItems := int(*totalItems)
		presp.TotalItems = &localTotalItems
	}

	if pout.Sort != nil && len(pout.Sort) == 1 {
		presp.SortBy = pout.Sort[0].Property
		presp.SortOrder = string(pout.Sort[0].Direction)
	}

	return presp
}
