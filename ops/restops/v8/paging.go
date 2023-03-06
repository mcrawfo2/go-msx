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
var ErrUnknownSortBy = errors.New("Unknown sort by")

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

type PagingConverter struct {
	SortByOptions paging.SortByOptions
}

// FromPagingSortOrder maps input sort names to their respective db column names
func (c PagingConverter) FromPagingSortOrder(sort []paging.SortOrder) ([]paging.SortOrder, error) {
	if len(sort) == 0 && c.SortByOptions.DefaultProperty != "" {
		sort = []paging.SortOrder{{
			Property:  c.SortByOptions.DefaultProperty,
			Direction: SortDirectionAsc,
		}}
	}

	if len(c.SortByOptions.Mapping) == 0 {
		// No sort options defined, free-for-all
		return sort, nil
	}

	for i, sortOrder := range sort {
		mappedProperty, ok := c.SortByOptions.Mapping.MapToRight(sortOrder.Property)
		if !ok {
			return nil, errors.Wrap(ErrUnknownSortBy, sortOrder.Property)
		}
		sort[i].Property = mappedProperty
	}
	return sort, nil
}

// ToPagingSortOrder maps db column names to their respective input sort names
func (c PagingConverter) ToPagingSortOrder(sort []paging.SortOrder) ([]paging.SortOrder, error) {
	if len(c.SortByOptions.Mapping) == 0 {
		// No sort options defined, free-for-all
		return sort, nil
	}

	for i, sortOrder := range sort {
		mappedProperty, ok := c.SortByOptions.Mapping.MapToLeft(sortOrder.Property)
		if !ok {
			return nil, errors.Wrap(ErrUnknownSortBy, sortOrder.Property)
		}
		sort[i].Property = mappedProperty
	}
	return sort, nil
}

func (c PagingConverter) FromPagingInputs(pageReq PagingInputs) (result paging.Request) {
	result.Page = uint(pageReq.Page)
	result.Size = uint(pageReq.PageSize)
	return
}

func (c PagingConverter) FromPagingSortingInputs(pageReq PagingSortingInputs) (result paging.Request, err error) {
	result, err = c.FromSortingInputs(pageReq.SortingInputs)
	if err != nil {
		return
	}
	result.Page = uint(pageReq.Page)
	result.Size = uint(pageReq.PageSize)
	return
}

func (c PagingConverter) FromSortingInputs(sortReq SortingInputs) (result paging.Request, err error) {
	if sortReq.SortBy != "" {
		sortResult := paging.SortOrder{
			Property:  sortReq.SortBy,
			Direction: paging.SortDirection(strings.ToUpper(sortReq.SortOrder)),
		}
		result.Sort = append(result.Sort, sortResult)
		result.Sort, err = c.FromPagingSortOrder(result.Sort)
	}
	return
}

func (c PagingConverter) FromPageSortQuery(page, pageSize int, sortBy, sortOrder string) (request paging.Request, err error) {
	request = c.FromPageQuery(page, pageSize)
	request.Sort, err = c.FromSortQuery(sortBy, sortOrder)
	return
}

func (c PagingConverter) FromPageQuery(page, pageSize int) paging.Request {
	return paging.Request{
		Page: uint(page),
		Size: uint(pageSize),
	}
}

func (c PagingConverter) FromSortQuery(sortBy, sortOrder string) ([]paging.SortOrder, error) {
	if sortBy == "" {
		return nil, nil
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
		return nil, errors.Wrap(ErrUnknownSortOrder, sortOrder)
	}

	result := []paging.SortOrder{{
		Property:  sortBy,
		Direction: sortDirection,
	}}

	return c.ToPagingSortOrder(result)
}

func (c PagingConverter) ToPagingResponse(pout paging.Response) (response PagingResponse, err error) {
	response = PagingResponse{
		Page:        int(pout.Number),
		PageSize:    int(pout.Size),
		HasNext:     pout.HasNext(),
		HasPrevious: pout.Offset() > 0,
	}

	totalItems := pout.TotalItems
	if pout.TotalItems != nil {
		localTotalItems := int(*totalItems)
		response.TotalItems = &localTotalItems
	}

	pout.Sort, err = c.ToPagingSortOrder(pout.Sort)
	if err != nil {
		return
	}

	if pout.Sort != nil && len(pout.Sort) == 1 {
		response.SortBy = pout.Sort[0].Property
		response.SortOrder = strings.ToLower(string(pout.Sort[0].Direction))
	}

	return
}
