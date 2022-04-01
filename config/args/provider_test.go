// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package args

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestProvider_Load(t *testing.T) {
	flags := types.StringStack{"aaaa", "bbbb"}

	args := []string{
		"--aaaa", "10",
		"--bbbb", "1.1",
		"--profile", "custom",
		"--spring.cloud.nginx.host=172.18.18.70",
	}
	os.Args = append([]string{"TestProvider_Load"}, args...)

	cp := NewCommandLineExtrasProvider(nil, func(name string) bool {
		return flags.Contains(name)
	})

	entries, err := cp.Load(context.Background())
	assert.Nil(t, err)
	assert.NotEmpty(t, entries)

	settings := config.MapFromEntries(entries)

	if _, ok := settings["aaaa"]; ok {
		assert.Fail(t, "aaaa set")
	}

	if _, ok := settings["bbbb"]; ok {
		assert.Fail(t, "bbbb set")
	}

	if profile, ok := settings["profile"]; !ok {
		assert.Fail(t, "profile not set")
	} else {
		assert.Equal(t, "custom", profile)
	}

	if springCloudNginxHost, ok := settings["spring.cloud.nginx.host"]; !ok {
		assert.Fail(t, "spring.cloud.nginx.host not set")
	} else {
		assert.Equal(t, "172.18.18.70", springCloudNginxHost)
	}
}
