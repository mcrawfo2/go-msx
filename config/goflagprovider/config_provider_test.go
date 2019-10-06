package goflagprovider

import (
	"context"
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfigProvider_Load(t *testing.T) {
	flagset := flag.NewFlagSet("TestConfigProvider", flag.ContinueOnError)

	flagset.Int("aaaa", 0, "wants int")
	flagset.Float64("bbbb", 1.0, "wants float64")
	flagset.String("eeee", "test", "wants string string")

	args := []string{
		"--aaaa", "10",
		"--bbbb", "1.1",
		"--profile", "custom",
		"--Spring.Cloud.nGinx.HOST=172.18.18.70",
	}
	os.Args = append(os.Args[:1], args...)

	cp := &ConfigProvider{
		flagSet: flagset,
		prefix:  "cli.flag.",
	}
	settings, err := cp.Load(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, settings)

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
