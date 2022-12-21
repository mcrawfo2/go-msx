// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package openapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"github.com/swaggest/openapi-go/openapi3"
	"reflect"
	"testing"
)

func TestPopulateFieldsFromTags(t *testing.T) {
	var tests = []struct {
		name    string
		tag     string
		want    *openapi3.Schema
		wantErr bool
	}{
		{
			name: "Title",
			tag:  `title:"some-title"`,
			want: new(openapi3.Schema).WithTitle("some-title"),
		},
		{
			name: "MultipleOf",
			tag:  `multipleOf:"10"`,
			want: new(openapi3.Schema).WithMultipleOf(10),
		},
		{
			name: "ExclusiveMaximum",
			tag:  `exclusiveMaximum:"true"`,
			want: new(openapi3.Schema).WithExclusiveMaximum(true),
		},
		{
			name: "Pattern",
			tag:  `pattern:"\\d"`,
			want: new(openapi3.Schema).WithPattern(`\d`),
		},
		{
			name: "Default",
			tag:  `default:"abc"`,
			want: new(openapi3.Schema).WithDefault(types.OptionalOf("abc").ValueInterface()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := new(openapi3.Schema).WithType(openapi3.SchemaTypeString)
			tt.want = tt.want.WithType(openapi3.SchemaTypeString)
			err := PopulateFieldsFromTags(s, reflect.StructTag(tt.tag))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.True(t,
				reflect.DeepEqual(tt.want, s),
				testhelpers.Diff(tt.want, s))
		})
	}

}
