// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package vault

// VersionRequest is used to identify versions to Delete, Undelete, and Destroy
type VersionRequest struct {
	Versions []int `json:"versions"`
}

type VersionedWriteRequest struct {
	Options VersionedWriteRequestOptions `json:"options"`
	Data    map[string]interface{}       `json:"data"`
}

type VersionedWriteRequestOptions struct {
	CAS int `json:"cas,omitempty"`
}

type VersionedMetadata struct {
	CurrentVersion int                        `mapstructure:"current_version"`
	OldestVersion  int                        `mapstructure:"oldest_version"`
	MaxVersions    int                        `mapstructure:"max_version"`
	Versions       map[string]MetadataVersion `mapstructure:"versions"`
}

type MetadataVersion struct {
	CreatedTime  string `mapstructure:"created_time"`
	DeletionTime string `mapstructure:"deletion_time"`
	Destroyed    bool   `mapstructure:"destroyed"`
}

type VersionedMetadataRequest struct {
	MaxVersions        int               `json:"max_versions,omitempty"`
	CasRequired        bool              `json:"cas_required,omitempty"`
	DeleteVersionAfter int               `json:"delete_version_after,omitempty"`
	CustomMetadata     map[string]string `json:"custom_metadata,omitempty"`
}
