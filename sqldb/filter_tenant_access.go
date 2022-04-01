// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sqldb

import (
	"cto-github.cisco.com/NFV-BU/go-msx/rbac"
	db "cto-github.cisco.com/NFV-BU/go-msx/sqldb/prepared"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

type TenantAccessFilter struct {
	Access             rbac.TenantAccess
	TenantIdColumnName string
}

func (t TenantAccessFilter) ValidateTenantAccess(tenantId uuid.UUID) error {
	return t.Access.ValidateTenantAccess(db.ToApiUuid(tenantId))
}

func (t TenantAccessFilter) Filter() goqu.Expression {
	switch {
	case t.Access.Unfiltered:
		// User has all tenants
		return nil

	default:
		// User has access to the records for their tenants
		return goqu.I(t.TenantIdColumnName).In(t.tenants()...)
	}
}

func (t TenantAccessFilter) tenants() (result []interface{}) {
	for _, tenantId := range t.Access.TenantIds {
		result = append(result, db.ToModelUuid(tenantId))
	}
	return
}

func NewTenantAccessFilter(access rbac.TenantAccess, tenantIdColumnName string) TenantAccessFilter {
	return TenantAccessFilter{
		Access:             access,
		TenantIdColumnName: tenantIdColumnName,
	}
}
