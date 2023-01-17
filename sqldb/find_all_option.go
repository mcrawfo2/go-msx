// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
	"github.com/doug-martin/goqu/v9"
	"strings"
)

type FindAllOption = func(ds *goqu.SelectDataset, pgReq paging.Request) (*goqu.SelectDataset, paging.Request)

func Where(where ...WhereOption) FindAllOption {
	return func(ds *goqu.SelectDataset, pgReq paging.Request) (*goqu.SelectDataset, paging.Request) {
		for _, w := range where {
			ds = ds.Where(w.Expression())
		}
		return ds, pgReq
	}
}

func Keys(keys KeysOption) FindAllOption {
	return func(ds *goqu.SelectDataset, pgReq paging.Request) (*goqu.SelectDataset, paging.Request) {
		ds = ds.Where(keys)
		return ds, pgReq
	}
}

func Distinct(distinct ...string) FindAllOption {
	return func(ds *goqu.SelectDataset, pgReq paging.Request) (*goqu.SelectDataset, paging.Request) {
		ds = ds.Distinct(strings.Join(distinct, ","))
		return ds, pgReq
	}
}

func Sort(sortOrders []paging.SortOrder) FindAllOption {
	return func(ds *goqu.SelectDataset, pgReq paging.Request) (*goqu.SelectDataset, paging.Request) {
		for _, sortOrder := range sortOrders {
			ident := goqu.I(sortOrder.Property)
			switch sortOrder.Direction {
			case paging.SortDirectionDesc:
				ds = ds.OrderAppend(ident.Desc())
			default:
				ds = ds.OrderAppend(ident.Asc())
			}
		}

		pgReq.Sort = sortOrders

		return ds, pgReq
	}
}

func Paging(pagingRequest paging.Request) FindAllOption {
	return func(ds *goqu.SelectDataset, pgReq paging.Request) (*goqu.SelectDataset, paging.Request) {
		if pagingRequest.Size > 0 {
			ds = ds.
				Limit(pagingRequest.Size).
				Offset(pagingRequest.Page * pagingRequest.Size)
		}

		for _, sortOrder := range pagingRequest.Sort {
			ident := goqu.I(sortOrder.Property)
			switch sortOrder.Direction {
			case paging.SortDirectionDesc:
				ds = ds.OrderAppend(ident.Desc())
			default:
				ds = ds.OrderAppend(ident.Asc())
			}
		}

		pgReq.Size = pagingRequest.Size
		pgReq.Page = pagingRequest.Page
		pgReq.Sort = pagingRequest.Sort

		return ds, pgReq
	}
}
