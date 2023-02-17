// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package v2

import (
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
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
	Content          any              `json:"content" inject:"Page"`
	HasNext          bool             `json:"hasNext"`          // more pages?
	Size             uint             `json:"size"`             // requested page size
	NumberOfElements uint             `json:"numberOfElements"` // current page length
	Number           uint             `json:"number"`           // current page number
	Pageable         PageableResponse `json:"pageable"`
}

func (r PagingResponse) Example() any {
	return PagingResponse{
		Content:          []any{},
		HasNext:          true,
		Size:             10,
		NumberOfElements: 10,
		Number:           3,
		Pageable: PageableResponse{
			Page: 3,
			Size: 10,
			Sort: SortResponse{
				Orders: []SortOrderResponse{
					{
						Property:  "someProperty",
						Direction: SortDirectionAsc,
					},
				},
			},
		},
	}
}

type PageableResponse struct {
	Page uint         `json:"page"` // pageNumber
	Size uint         `json:"size"` // pageSize
	Sort SortResponse `json:"sort"`
}

type SortResponse struct {
	Orders []SortOrderResponse
}

type SortOrderResponse struct {
	Property  string        `json:"property"`
	Direction SortDirection `json:"direction"`
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
		result.Sort, err = c.ToPagingSortOrder(result.Sort)
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
		Content:          pout.Content,
		HasNext:          pout.HasNext(),
		Size:             pout.Size,
		NumberOfElements: pout.Elements(),
		Number:           pout.Number,
		Pageable: PageableResponse{
			Page: pout.Number,
			Size: pout.Size,
		},
	}

	pout.Sort, err = c.ToPagingSortOrder(pout.Sort)
	if err != nil {
		return
	}

	if len(pout.Sort) > 0 {
		response.Pageable.Sort = SortResponse{
			Orders: c.ToSortOrderResponses(pout.Sort),
		}
	}

	return response, nil
}

func (c PagingConverter) ToSortOrderResponses(orders []paging.SortOrder) []SortOrderResponse {
	var results []SortOrderResponse
	for _, order := range orders {
		results = append(results, c.ToSortOrderResponse(order))
	}
	return results
}

func (c PagingConverter) ToSortOrderResponse(order paging.SortOrder) SortOrderResponse {
	return SortOrderResponse{
		Property:  order.Property,
		Direction: SortDirection(order.Direction),
	}
}
