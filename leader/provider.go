package leader

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/pkg/errors"
)

type LeadershipProvider interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	MasterKey(ctx context.Context) string
	IsLeader(ctx context.Context, key string) bool
	ReleaseLeadership(ctx context.Context, key string)
}

var (
	logger                          = log.NewLogger("msx.leader")
	leadershipProvider              LeadershipProvider
	ErrLeadershipProviderNotDefined = errors.New("Leadership provider not registered")
)

func RegisterLeadershipProvider(provider LeadershipProvider) {
	if provider != nil {
		leadershipProvider = provider
	}
}

func IsLeadershipProviderRegistered() bool {
	return leadershipProvider != nil
}

func IsLeader(ctx context.Context, key string) (bool, error) {
	if !IsLeadershipProviderRegistered() {
		return false, ErrLeadershipProviderNotDefined
	}

	return leadershipProvider.IsLeader(ctx, key), nil
}

func IsMasterLeader(ctx context.Context) (bool, error) {
	if !IsLeadershipProviderRegistered() {
		return false, ErrLeadershipProviderNotDefined
	}

	masterKey := leadershipProvider.MasterKey(ctx)
	return leadershipProvider.IsLeader(ctx, masterKey), nil
}

func ReleaseLeadership(ctx context.Context, key string) error {
	if !IsLeadershipProviderRegistered() {
		return ErrLeadershipProviderNotDefined
	}

	leadershipProvider.ReleaseLeadership(ctx, key)
	return nil
}

func ReleaseMasterLeadership(ctx context.Context, key string) error {
	if !IsLeadershipProviderRegistered() {
		return ErrLeadershipProviderNotDefined
	}

	masterKey := leadershipProvider.MasterKey(ctx)
	leadershipProvider.ReleaseLeadership(ctx, masterKey)
	return nil
}

func Start(ctx context.Context) error {
	if !IsLeadershipProviderRegistered() {
		return ErrLeadershipProviderNotDefined
	}

	logger.WithContext(ctx).Info("Starting leadership election")
	return leadershipProvider.Start(ctx)
}

func Stop(ctx context.Context) error {
	if !IsLeadershipProviderRegistered() {
		return ErrLeadershipProviderNotDefined
	}

	logger.WithContext(ctx).Info("Stopping leadership election")
	return leadershipProvider.Stop(ctx)
}
