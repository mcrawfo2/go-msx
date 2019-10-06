package pflagprovider

import (
	"context"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfigProvider_Load(t *testing.T) {
	flagset := pflag.NewFlagSet("TestConfigProvider", pflag.ContinueOnError)
	flagset.ParseErrorsWhitelist.UnknownFlags = true

	flagset.Int("aaaa", 0, "wants int")
	flagset.Float64("bbbb", 1.0, "wants float64")
	flagset.StringSlice("cccc", []string{"test"}, "wants string slice")
	flagset.StringArray("dddd", []string{"test"}, "wants string array")
	flagset.String("eeee", "test", "wants string string")

	args := []string{
		"--aaaa", "10",
		"--bbbb", "1.1",
		"--cccc", "a,b,c",
		"--dddd", "a",
		"--dddd", "b",
		"--dddd", "c",
		"--profile", "custom",
		"--Spring.Cloud.nGinx.HOST=172.18.18.70",
	}
	os.Args = append(os.Args[:1], args...)

	err := flagset.Parse(args)
	assert.Nil(t, err)

	cp := &ConfigProvider{
		flagset: flagset,
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

	if cccc, ok := settings["cli.flag.cccc"]; !ok {
		assert.Fail(t, "cccc not set")
	} else {
		assert.Equal(t, "[a,b,c]", cccc)
	}

	if dddd, ok := settings["cli.flag.dddd"]; !ok {
		assert.Fail(t, "dddd not set")
	} else {
		assert.Equal(t, "[a,b,c]", dddd)
	}

	if eeee, ok := settings["cli.flag.eeee"]; !ok {
		assert.Fail(t, "eeee not set")
	} else {
		assert.Equal(t, "test", eeee)
	}

	if springCloudNginxHost, ok := settings["spring.cloud.nginx.host"]; !ok {
		assert.Fail(t, "spring.cloud.nginx.host not set")
	} else {
		assert.Equal(t, "172.18.18.70", springCloudNginxHost)
	}
}
