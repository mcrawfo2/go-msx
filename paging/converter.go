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
	response.Content = types.NewOptional(dataResponse).OrElse(response.Content)
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
