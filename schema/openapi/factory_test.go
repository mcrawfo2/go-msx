// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package openapi

import (
	"github.com/stretchr/testify/assert"
	"github.com/swaggest/jsonschema-go"
	"testing"
)

func TestMapSchema(t *testing.T) {
	s := MapSchema(StringSchema())
	assert.NotNil(t, s)
}

func TestNewType(t *testing.T) {
	s := NewType(jsonschema.Null)
	assert.NotNil(t, s)
}
