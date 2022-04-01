// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package vault

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCreateTransitKeyRequest(t *testing.T) {
	request := NewCreateTransitKeyRequest()
	assert.Equal(t, KeyTypeAes256Gcm96, request.Type)
	assert.Equal(t, false, *request.AllowPlaintextBackup)
	assert.Equal(t, false, *request.Exportable)
}
