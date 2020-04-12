package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/leader"
	"cto-github.cisco.com/NFV-BU/go-msx/leader/consulprovider"
)

func init() {
	OnEvent(EventConfigure, PhaseAfter, registerLeadershipProvider)
	OnEvent(EventReady, PhaseBefore, startLeadershipElection)
	OnEvent(EventStop, PhaseBefore, stopLeadershipElection)
}

func registerLeadershipProvider(ctx context.Context) error {
	logger.Info("Registering consul leadership provider")
	leadershipProvider, err := consulprovider.NewLeadershipProvider(ctx)
	if err == consulprovider.ErrDisabled {
		logger.Info(err)
	} else if err != nil {
		return err
	} else if leadershipProvider != nil {
		leader.RegisterLeadershipProvider(leadershipProvider)
	}

	return nil
}

func startLeadershipElection(ctx context.Context) error {
	if err := leader.Start(ctx); err != nil && err != leader.ErrLeadershipProviderNotDefined {
		return err
	}
	return nil
}

func stopLeadershipElection(ctx context.Context) error {
	if err := leader.Stop(ctx); err != nil && err != leader.ErrLeadershipProviderNotDefined {
		logger.Error(err)
	}
	return nil
}
