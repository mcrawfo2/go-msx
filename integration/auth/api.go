// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package auth

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

//go:generate mockery --inpackage --name=Api --structname=MockAuth
type Api interface {
	Login(user, password string) (*integration.MsxResponse, error)
	SwitchContext(accessToken string, userId types.UUID) (*integration.MsxResponse, error)
	Logout() (*integration.MsxResponse, error)

	GetTokenDetails(noDetails bool) (*integration.MsxResponse, error)
	GetTokenKeys() (keys JsonWebKeys, response *integration.MsxResponse, err error)

	GetTenantHierarchyRoot() (*integration.MsxResponse, error)
	GetTenantHierarchyParent(tenantId types.UUID) (*integration.MsxResponse, error)
	GetTenantHierarchyAncestors(tenantId types.UUID) (*integration.MsxResponse, []types.UUID, error)
	GetTenantHierarchyChildren(tenantId types.UUID) (*integration.MsxResponse, []types.UUID, error)
	GetTenantHierarchyDescendants(tenantId types.UUID) (*integration.MsxResponse, []types.UUID, error)
}

type AuthApi = Api
