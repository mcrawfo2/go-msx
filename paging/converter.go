// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package paging

import "cto-github.cisco.com/NFV-BU/go-msx/types"

type Converter struct{}

func (c Converter) RequestFromQuery(page, pageSize uint) Request {
	return Request{
		Page: page,
		Size: pageSize,
	}
}

func (c Converter) RequestFromPageableRequest(request PageableRequest) Request {
	return Request{
		Page:  request.Page,
		Size:  request.Size,
		Sort:  c.SortFromSortRequest(request.Sort),
		State: request.State,
	}
}

func (c Converter) SortFromSortRequest(request SortRequest) (results []SortOrder) {
	for _, order := range request.Orders {
		results = append(results, c.SortOrderFromSortOrderRequest(order))
	}

	for _, property := range request.Properties {
		results = append(results, c.SortOrderFromSortOrderRequest(SortOrderRequest{
			Property:  property,
			Direction: request.Direction,
		}))
	}

	return
}

func (c Converter) SortOrderFromSortOrderRequest(request SortOrderRequest) SortOrder {
	return SortOrder{
		Property:  request.Property,
		Direction: request.Direction,
	}
}

func (c Converter) ResponseToPaginatedResponse(response Response, dataResponse interface{}) PaginatedResponse {
	response.Content = types.OptionalOf(dataResponse).OrElse(response.Content)
	orders := c.SortOrderListToSortOrderResponseList(response.Sort)
	return PaginatedResponse{
		Content:          response.Content,
		HasNext:          response.HasNext(),
		Size:             response.Size,
		NumberOfElements: response.Elements(),
		Number:           response.Number,
		Pageable: PageableResponse{
			Page: response.Number,
			Size: response.Size,
			Sort: SortResponse{
				Orders: orders,
			},
			State: response.State,
		},
	}
}

func (c Converter) ResponseToPaginatedResponseV8(response Response, objects []interface{}) PaginatedResponseV8 {

	if objects == nil {
		objects = make([]interface{}, 0)
	}

	presp := PaginatedResponseV8{
		Page:        int32(response.Number),
		PageSize:    int32(response.Size),
		HasNext:     response.HasNext(),
		HasPrevious: response.Offset() > 0,
		Contents:    objects,
	}

	totalItems := response.TotalItems
	if totalItems != nil {
		localTotalItems := int64(*totalItems)
		presp.TotalItems = &localTotalItems
	}

	//TODO: Support for multiple fields of {sortBy, sortOrder}
	if response.Sort != nil && len(response.Sort) == 1 {
		presp.SortBy = response.Sort[0].Property
		presp.SortOrder = response.Sort[0].Direction
	}

	return presp
}

func (c Converter) SortOrderToSortOrderResponse(order SortOrder) SortOrderResponse {
	return SortOrderResponse{
		Property:  order.Property,
		Direction: order.Direction,
	}
}

func (c Converter) SortOrderListToSortOrderResponseList(orders []SortOrder) []SortOrderResponse {
	var results []SortOrderResponse
	for _, order := range orders {
		results = append(results, c.SortOrderToSortOrderResponse(order))
	}
	return results
}

func ToApiSingleSortBy(orders []SortOrder) string {
	if len(orders) != 1 {
		return ""
	}
	return orders[0].Property
}

func ToApiSingleSortOrder(orders []SortOrder) string {
	if len(orders) != 1 {
		return ""
	}
	return string(orders[0].Direction)
}
