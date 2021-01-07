package cobraprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var rootCmd = &cobra.Command{
	Use: "root",
	Run: func(cmd *cobra.Command, args []string) {},
}

var secondaryCommand = &cobra.Command{
	Use: "secondary",
	Run: func(cmd *cobra.Command, args []string) {},
}

func TestProvider_Load(t *testing.T) {
	rootCmd.PersistentFlags().Int("aaaa", 1, "Wants an int")
	rootCmd.LocalFlags().Float64("bbbb", 2.2, "Wants a float64")
	secondaryCommand.LocalFlags().String("cccc", "test", "Wants a string")
	rootCmd.AddCommand(secondaryCommand)

	args := []string{
		"--aaaa", "10",
		"--bbbb", "1.1",
		"--cccc", "abc",
		"--profile", "custom",
		"--spring.cloud.nginx.host=172.18.18.70",
	}
	os.Args = append(os.Args[:1], args...)

	cp := NewProvider("commandLine", secondaryCommand, "cli.flag")

	entries, err := cp.Load(context.Background())
	assert.Nil(t, err)
	assert.NotEmpty(t, entries)

	settings := config.MapFromEntries(entries)

	if aaaa, ok := settings["cli.flag.aaaa"]; !ok {
		assert.Fail(t, "aaaa not set")
	} else {
		assert.Equal(t, "10", aaaa)
	}

	if cccc, ok := settings["cli.flag.cccc"]; !ok {
		assert.Fail(t, "cccc not set")
	} else {
		assert.Equal(t, "abc", cccc)
	}

	if bbbb, ok := settings["bbbb"]; !ok {
		assert.Fail(t, "bbbb not set")
	} else {
		assert.Equal(t, "1.1", bbbb)
	}

	if springCloudNginxHost, ok := settings["spring.cloud.nginx.host"]; !ok {
		assert.Fail(t, "spring.cloud.nginx.host not set")
	} else {
		assert.Equal(t, "172.18.18.70", springCloudNginxHost)
	}
}
