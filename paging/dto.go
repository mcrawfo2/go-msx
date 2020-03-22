package paging

import "encoding/json"

// Response

type PaginatedResponse struct {
	Content          interface{}      `json:"content"`
	HasNext          bool             `json:"hasNext"`
	Size             uint             `json:"size"`
	NumberOfElements uint             `json:"numberOfElements"`
	Number           uint             `json:"number"`
	Pageable         PageableResponse `json:"pageable"`
}

type PageableResponse struct {
	Page  uint         `json:"page"`
	Size  uint         `json:"size"`
	Sort  SortResponse `json:"sort"`
	State interface{}  `json:"pagingState"`
}

type SortResponse struct {
	Orders []SortOrderResponse
}

func (s SortResponse) MarshalJSON() ([]byte, error) {
	if len(s.Orders) == 0 {
		return json.Marshal(map[string]interface{}{
			"sorted": false,
		})
	}

	return json.Marshal(map[string]interface{}{
		"orders": s.Orders,
	})
}

type SortOrderResponse struct {
	Property  string        `json:"property"`
	Direction SortDirection `json:"direction"`
}

// Request

type PageableRequest struct {
	Page  uint        `json:"page"`
	Size  uint        `json:"size"`
	Sort  SortRequest `json:"sort"`
	State interface{} `json:"pagingState"`
}

type SortRequest struct {
	Properties []string           `json:"properties"`
	Direction  SortDirection      `json:"direction"`
	Orders     []SortOrderRequest `json:"orders"`
}

type SortOrderRequest SortOrderResponse
