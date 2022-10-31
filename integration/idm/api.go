// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package idm

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/paging"
)

//go:generate mockery --inpackage --name=Api --structname=MockIdm
type Api interface {
	IsTokenActive() (*integration.MsxResponse, error)

	GetMyProvider() (*integration.MsxResponse, error)
	GetProviderByName(providerName string) (*integration.MsxResponse, error)

	GetProviderExtensionByName(name string) (*integration.MsxResponse, error)

	// Deprecated: Use v8/TenantsApiService (or newer)
	GetTenantById(tenantId string) (*integration.MsxResponse, error)
	// Deprecated: Tenants should generally be retrieved by ID, not tenantName.  The underlying REST API is
	// current retired and due for decommissioning.
	GetTenantByName(tenantName string) (*integration.MsxResponse, error)

	GetTenantByIdV8(tenantId string) (*integration.MsxResponse, error)

	// Deprecated: Use v8/UsersApiService (or newer).
	GetUserById(userId string) (*integration.MsxResponse, error)
	GetUserByIdV8(userId string) (*integration.MsxResponse, error)

	GetRoles(resolvePermissionNames bool, p paging.Request) (*integration.MsxResponse, error)
	CreateRole(dbinstaller bool, body RoleCreateRequest) (*integration.MsxResponse, error)
	UpdateRole(dbinstaller bool, body RoleUpdateRequest) (*integration.MsxResponse, error)
	DeleteRole(roleName string) (*integration.MsxResponse, error)

	GetCapabilities(p paging.Request) (*integration.MsxResponse, error)
	BatchCreateCapabilities(populator bool, owner string, capabilities []CapabilityCreateRequest) (*integration.MsxResponse, error)
	BatchUpdateCapabilities(populator bool, owner string, capabilities []CapabilityUpdateRequest) (*integration.MsxResponse, error)
	DeleteCapability(populator bool, owner string, name string) (*integration.MsxResponse, error)
}

type IdmApi = Api
