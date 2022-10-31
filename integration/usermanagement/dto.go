// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package usermanagement

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/auth"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/idm"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/secrets"
)

/* idm dtos */

type ProviderResponse = idm.ProviderResponse

type ProviderExtensionResponse = idm.ProviderExtensionResponse

type TenantResponse = idm.TenantResponse

type TenantResponseV8 = idm.TenantResponseV8

type RoleCreateRequest = idm.RoleCreateRequest

type RoleUpdateRequest = idm.RoleUpdateRequest // alias of RoleCreateRequest

type RoleResponse = idm.RoleResponse

type RoleListResponse = idm.RoleListResponse

type CapabilityCreateRequest = idm.CapabilityCreateRequest

type CapabilityUpdateRequest = idm.CapabilityUpdateRequest // alias of CapabilityCreateRequest

type CapabilityBatchCreateRequest = idm.CapabilityBatchCreateRequest

type CapabilityBatchUpdateRequest = idm.CapabilityBatchUpdateRequest

type CapabilityResponse = idm.CapabilityResponse

type CapabilityListResponse = idm.CapabilityListResponse

type UserResponse = idm.UserResponse

type UserResponseV8 = idm.UserResponseV8

type UserRoleResponse = idm.UserRoleResponse

type UserTenantResponse = idm.UserTenantResponse

/* auth dtos */

type TokenDetails = auth.TokenDetails

type LoginResponse = auth.LoginResponse

type JsonWebKey = auth.JsonWebKey

type JsonWebKeys = auth.JsonWebKeys

/* secrets dtos */

type Pojo integration.Pojo
type PojoArray integration.PojoArray

type EncryptSecretsDTO = secrets.EncryptSecretsDTO

type GetSecretRequestDTO = secrets.GetSecretRequestDTO

type GenerateSecretRequestDTO = secrets.GenerateSecretRequestDTO

type SecretPolicySetRequest = secrets.SecretPolicySetRequest

type SecretPolicyResponse = secrets.SecretPolicyResponse

type AgingRule = secrets.AgingRule

type CharacterRule = secrets.CharacterRule

type DictionaryRule = secrets.DictionaryRule

type HistoryRule = secrets.HistoryRule

type KeyRule = secrets.KeyRule

type LengthRule = secrets.LengthRule

type SecretsResponse = secrets.SecretsResponse
