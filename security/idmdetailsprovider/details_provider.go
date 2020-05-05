package idmdetailsprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cache/lru"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
)

const (
	configRootIdmTokenDetailsProvider = "security.token.details"
)

var (
	logger = log.NewLogger("msx.security.idmdetailsprovider")
)

type IdmTokenDetailsProviderConfig struct {
	Fast         bool `config:"default=false"`
	ActiveCache  lru.CacheConfig
	DetailsCache lru.CacheConfig
}

type TokenDetailsProvider struct {
	cfg          IdmTokenDetailsProviderConfig
	detailsCache lru.Cache
	activeCache  lru.Cache
	fetcher      detailsFetcher
}

func (t *TokenDetailsProvider) IsTokenActive(ctx context.Context) (active bool, err error) {
	token := security.UserContextFromContext(ctx).Token
	if token == "" {
		return false, security.ErrTokenNotFound
	}

	activeInterface, exists := t.activeCache.Get(token)
	if !exists {
		if t.cfg.Fast {
			active, err = t.fetchActiveCombined(ctx, token)
		} else {
			active, err = t.fetchActiveSeparate(ctx, token)
		}
	} else {
		active = activeInterface.(bool)
	}

	if err != nil {
		logger.WithContext(ctx).WithError(err).Error("Failed to load token details")
	}

	return active, nil
}

func (t *TokenDetailsProvider) TokenDetails(ctx context.Context) (*security.UserContextDetails, error) {
	token := security.UserContextFromContext(ctx).Token
	if token == "" {
		return nil, security.ErrTokenNotFound
	}

	detailsInterface, exists := t.detailsCache.Get(token)
	if exists {
		return detailsInterface.(*security.UserContextDetails), nil
	}

	return t.fetchDetails(ctx, token)
}

func (t *TokenDetailsProvider) fetchActiveCombined(ctx context.Context, token string) (active bool, err error) {
	_, exists := t.detailsCache.Get(token)
	if exists {
		return t.fetchActiveSeparate(ctx, token)
	}

	details, err := t.fetchDetails(ctx, token)
	if err != nil {
		return false, err
	}
	return details.Active, nil
}

func (t *TokenDetailsProvider) fetchActiveSeparate(ctx context.Context, token string) (active bool, err error) {
	logger.WithContext(ctx).Info("Verifying token active")

	active, err = t.fetcher.FetchActive(ctx)
	if err != nil {
		return
	}
	t.activeCache.Set(token, active)
	return active, nil
}

func (t *TokenDetailsProvider) fetchDetails(ctx context.Context, token string) (details *security.UserContextDetails, err error) {
	logger.WithContext(ctx).Info("Loading token details")

	details, err = t.fetcher.FetchDetails(ctx)
	if err != nil {
		return nil, err
	}

	t.detailsCache.Set(token, details)
	t.activeCache.Set(token, details.Active)
	return details, nil
}

func RegisterTokenDetailsProvider(ctx context.Context) error {
	logger.Info("Registering IDM token details provider")
	var cfg IdmTokenDetailsProviderConfig
	if err := config.FromContext(ctx).Populate(&cfg, configRootIdmTokenDetailsProvider); err != nil {
		return err
	}

	var fetcher detailsFetcher
	if cfg.Fast {
		fetcher = new(fastDetailsFetcher)
	} else {
		fetcher = new(slowDetailsFetcher)
	}

	security.SetTokenDetailsProvider(&TokenDetailsProvider{
		cfg:          cfg,
		fetcher:      fetcher,
		activeCache:  lru.NewCacheFromConfig(&cfg.ActiveCache),
		detailsCache: lru.NewCacheFromConfig(&cfg.DetailsCache),
	})

	return nil
}
