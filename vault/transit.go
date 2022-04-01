// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package vault

const KeyTypeAes256Gcm96 = "aes256-gcm96"

type CreateTransitKeyRequest struct {
	Type                 string `json:"type"`
	Exportable           *bool  `json:"exportable,omitempty"`
	AllowPlaintextBackup *bool  `json:"allow_plaintext_backup,omitempty"`
}

func NewCreateTransitKeyRequest() CreateTransitKeyRequest {
	no := false
	return CreateTransitKeyRequest{
		Type:                 KeyTypeAes256Gcm96,
		Exportable:           &no,
		AllowPlaintextBackup: &no,
	}
}
