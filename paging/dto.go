// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package paging

import "encoding/json"

// Response

type PaginatedResponse struct {
	Content          interface{}      `json:"content" inject:"Page"`
	HasNext          bool             `json:"hasNext"`
	Size             uint             `json:"size"`
	NumberOfElements uint             `json:"numberOfElements"`
	Number           uint             `json:"number"`
	Pageable         PageableResponse `json:"pageable"`
}

// Paging Response V8
type PaginatedResponseV8 struct {
	Page        int32         `json:"page"`
	PageSize    int32         `json:"pageSize"`
	TotalItems  *int64        `json:"totalItems"`
	HasNext     bool          `json:"hasNext"`
	HasPrevious bool          `json:"hasPrevious"`
	SortBy      string        `json:"sortBy,omitempty"`
	SortOrder   SortDirection `json:"sortOrder,omitempty"`
	Contents    interface{}   `json:"contents" inject:"Page"`
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
