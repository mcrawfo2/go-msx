// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package goflagprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestProvider_Load(t *testing.T) {
	flagset := flag.NewFlagSet("TestProvider", flag.ContinueOnError)

	flagset.Int("aaaa", 0, "wants int")
	flagset.Float64("bbbb", 1.0, "wants float64")
	flagset.String("eeee", "test", "wants string string")

	args := []string{
		"--aaaa", "10",
		"--bbbb", "1.1",
		"--profile", "custom",
		"--spring.cloud.nginx.host=172.18.18.70",
	}
	os.Args = append([]string{"TestProvider_Load"}, args...)

	cp := NewProvider("test", flagset, "cli.flag")
	entries, err := cp.Load(context.Background())
	assert.Nil(t, err)
	assert.NotEmpty(t, entries)

	settings := config.MapFromEntries(entries)

	if aaaa, ok := settings["cli.flag.aaaa"]; !ok {
		assert.Fail(t, "aaaa not set")
	} else {
		assert.Equal(t, "10", aaaa)
	}

	if eeee, ok := settings["cli.flag.eeee"]; !ok {
		assert.Fail(t, "eeee not set")
	} else {
		assert.Equal(t, "test", eeee)
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
