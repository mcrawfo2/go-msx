package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"moul.io/banner"
)

func init() {
	OnRootEvent(EventConfigure, PhaseAfter, func(ctx context.Context) error {
		cfg := config.FromContext(ctx)
		if cfg == nil {
			return nil
		}

		bannerEnabled, _ := cfg.BoolOr("banner.enabled", false)
		if !bannerEnabled {
			return nil
		}

		appName, _ := cfg.StringOr("info.app.name", "")
		appVersion, _ := config.FromContext(ctx).StringOr("info.build.version", "")
		if appName != "" {
			bannerText := banner.Inline(appName)
			logger.Infof("\n" + bannerText + "\nVersion: %s\n", appVersion)
		}
		return nil
	})
}
