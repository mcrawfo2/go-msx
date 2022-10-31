// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package idm

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

type ProviderResponse struct {
	Name             string     `json:"name"`
	DisplayName      *string    `json:"displayName"`
	Description      *string    `json:"description"`
	ProvidersID      types.UUID `json:"providersId"`
	Email            *string    `json:"email"`
	NotificationType *string    `json:"notificationType"`
	Locale           *string    `json:"locale"`
	AnyConnectURL    *string    `json:"anyConnectURL"`
}

type ProviderExtensionResponse struct {
	Name          string   `json:"name"`
	AllowedValues []string `json:"allowedValues"`
	DisplayType   string   `json:"displayType"`
	Type          string   `json:"type"`
	Label         string   `json:"label"`
	Value         string   `json:"value"`
}

type TenantResponse struct {
	TenantId          types.UUID  `json:"tenantId"`
	ParentId          *types.UUID `json:"parentId"`
	ProviderId        types.UUID  `json:"providerId"`
	ProviderName      string      `json:"providerName"`
	TenantName        string      `json:"tenantName"`
	DisplayName       string      `json:"displayName"`
	Image             string      `json:"image"`
	Email             *string     `json:"email"`
	CreatedOn         int64       `json:"createdOn"`
	ModifiedOn        int64       `json:"lastUpdated"`
	Suspended         bool        `json:"suspended"`
	TenantDescription string      `json:"tenantDescription"`
	URL               string      `json:"url"`
	TenantGroupName   interface{} `json:"tenantGroupName"`
	TenantExtension   struct {
		Parameters interface{} `json:"parameters"`
	} `json:"tenantExtension"`
}

type TenantResponseV8 struct {
	TenantId         types.UUID  `json:"id"`
	ParentId         *types.UUID `json:"parentId"`
	Name             string      `json:"name"`
	Image            *string     `json:"image"`
	Email            *string     `json:"email"`
	Description      string      `json:"description"`
	URL              string      `json:"url"`
	Suspended        bool        `json:"suspended"`
	NumberOfChildren int64       `json:"numberOfChildren"`
	CreatedOn        types.Time  `json:"createdOn"`
	ModifiedOn       types.Time  `json:"lastUpdated"`
}

type RoleCreateRequest struct {
	CapabilityAddList    []string `json:"capabilityAddList"`
	CapabilityRemoveList []string `json:"capabilityRemoveList"`
	CapabilityList       []string `json:"capabilitylist"`
	Description          string   `json:"description"`
	DisplayName          string   `json:"displayName"`
	IsSeeded             string   `json:"isSeeded"`
	Owner                string   `json:"owner"`
	ResourceDescriptor   string   `json:"resourceDescriptor"`
	RoleName             string   `json:"roleName"`
	Roleid               string   `json:"roleid"`
}

type RoleUpdateRequest RoleCreateRequest

type RoleResponse struct {
	CapabilityList     []string `json:"capabilitylist"`
	Description        string   `json:"description"`
	DisplayName        string   `json:"displayName"`
	Href               string   `json:"href"`
	IsSeeded           string   `json:"isSeeded"`
	Owner              string   `json:"owner"`
	ResourceDescriptor string   `json:"resourceDescriptor"`
	RoleName           string   `json:"roleName"`
	RoleId             string   `json:"roleid"`
	Status             string   `json:"status"`
}

type RoleListResponse struct {
	Roles []RoleResponse `json:"roles"`
}

type CapabilityCreateRequest struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	DisplayName string `json:"displayName"`
	IsDefault   string `json:"isDefault"`
	Name        string `json:"name"`
	ObjectName  string `json:"objectName"`
	Operation   string `json:"operation"`
}

type CapabilityBatchCreateRequest struct {
	Capabilities []CapabilityCreateRequest `json:"capabilities"`
}

type CapabilityUpdateRequest CapabilityCreateRequest

type CapabilityBatchUpdateRequest struct {
	Capabilities []CapabilityUpdateRequest `json:"capabilities"`
}

type CapabilityResponse struct {
	CapabilityCreateRequest
	Resources []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Method      string `json:"method"`
		Endpoint    string `json:"endpoint"`
	} `json:"resources"`
}

type CapabilityListResponse struct {
	Capabilities []CapabilityResponse `json:"capabilities"`
}

type UserResponse struct {
	Id            string               `json:"id"`
	Email         string               `json:"email"`
	UserId        string               `json:"userId"`
	Name          string               `json:"name"`
	FirstName     string               `json:"firstName"`
	LastName      string               `json:"lastName"`
	Status        string               `json:"status,omitempty"`
	PwdPolicyname string               `json:"pwdPolicyname,omitempty"`
	Locale        string               `json:"locale,omitempty"`
	Seeded        bool                 `json:"isSeeded,omitempty"`
	Deleted       string               `json:"isDeleted,omitempty"`
	Tenants       []UserTenantResponse `json:"tenants,omitempty"`
	Roles         []UserRoleResponse   `json:"roles,omitempty"`
}

type UserResponseV8 struct {
	Id            string       `json:"id"`
	Email         string       `json:"email"`
	Username      string       `json:"username"`
	FirstName     string       `json:"firstName"`
	LastName      string       `json:"lastName"`
	Status        string       `json:"status"`
	PwdPolicyname string       `json:"passwordPolicyName"`
	Locale        string       `json:"locale"`
	Deleted       string       `json:"deleted"`
	TenantIds     []types.UUID `json:"tenantIds"`
	RolesIds      []types.UUID `json:"roleIds"`
}

type UserRoleResponse struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type UserTenantResponse struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}
