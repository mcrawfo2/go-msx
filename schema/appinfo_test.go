// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package schema

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestAppInfoFromConfig(t *testing.T) {
	values := map[string]string{
		configKeyAppName:        "application-name",
		configKeyAppDescription: "application-description",
		configKeyBuildVersion:   "build-version",
		configKeyDisplayName:    "display-name",
	}

	cfg := configtest.NewInMemoryConfig(values)

	want := &AppInfo{
		Name:        values[configKeyAppName],
		DisplayName: values[configKeyDisplayName],
		Description: values[configKeyAppDescription],
		Version:     values[configKeyBuildVersion],
	}

	got, err := AppInfoFromConfig(cfg)
	assert.NoError(t, err)
	assert.True(t,
		reflect.DeepEqual(got, want),
		testhelpers.Diff(want, got))
}
