package main

import (
	"context"
	cliConfig "cto-github.cisco.com/NFV-BU/go-msx/cli/pflag"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/support/log"
	"github.com/spf13/pflag"
	"os"
	"time"
)

var logger = log.StandardLogger()

func init() {
	config.SetStaticConfig(map[string]string{
		"application.name": "config",
	})
}

func options() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	flagSet.Bool("options", false, "Print options")
	flagSet.Bool( "version", false, "Print version")
	flagSet.String( "log_level", "debug", "Logging level {panic, fatal, error, warning, info, debug}")
	flagSet.Int( "port", 8080, "Server `Port`")
	flagSet.String( "cert", "/etc/ssl/cert/ssl.crt", "Cert File `Cert`")
	flagSet.String( "pkey", "/etc/ssl/key/ssl.key", "Private-Key File `Pkey`")
	flagSet.Int64( "timeout", 1000, "Timeout in milliseconds, (10ms to 24 hours)")
	flagSet.Int64( "read_timeout", 2000, "Read timeout in milliseconds (10ms to 24 hours)")
	flagSet.Int64("write_timeout", 3000, "Write timeout in milliseconds (10ms to 24 hours)")
	flagSet.Int("max_header_size", 512, "Max header size in bytes")
	flagSet.Int64( "max_request_size", 10480, "Max request size in bytes")
	return flagSet
}

func main() {
	flagSet := options()
	cliConfig.RegisterConfigProvider(flagSet, "cli.flag.")

	ctx, cancel := context.WithTimeout(context.Background(), 90 * time.Second)
	defer cancel()


	if err := config.Bootstrap().Load(ctx); err != nil {
		logger.Fatal(err)
	}

	config.RegisterRemoteConfigProviders()
	if err := config.Application().Load(ctx); err != nil {
		logger.Fatal(err)
	}

	logger.Info("Dumping application configuration")
	config.Application().Each(func(name, value string) {
		logger.Infof("%s: %s", name, value)
	})
}
