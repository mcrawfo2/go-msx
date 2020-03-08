package consulprovider

import (
	"context"
	"github.com/pkg/errors"
	"sync"
)

var ErrDisabled = errors.New("Leadership provider disabled")
var ErrAlreadyStarted = errors.New("Leadership provider already started")

type LeadershipProvider struct {
	cfg         *ConsulLeaderElectionConfig
	childCtx    context.Context
	childCancel context.CancelFunc
	workers     map[string]*LeadershipInitiator
	started     bool
	mtx         sync.Mutex
}

func (l *LeadershipProvider) isStarted() bool {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	return l.started
}

func (l *LeadershipProvider) setStarted(started bool) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.started = started
}

func (l *LeadershipProvider) MasterKey(ctx context.Context) string {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	return l.cfg.DefaultMasterKey
}

func (l *LeadershipProvider) Start(ctx context.Context) error {
	err := func() error {
		l.mtx.Lock()
		defer l.mtx.Unlock()
		if l.started {
			return ErrAlreadyStarted
		}
		l.started = true
		return nil
	}()
	if err != nil {
		return err
	}

	l.workers = make(map[string]*LeadershipInitiator)
	l.childCtx, l.childCancel = context.WithCancel(ctx)

	for _, properties := range l.cfg.LeaderProperties {
		worker := NewLeadershipInitiator(l.childCtx, properties)
		go worker.Start()
		l.workers[properties.Key] = worker
	}

	return nil
}

func (l *LeadershipProvider) Stop(ctx context.Context) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if !l.started {
		return nil
	}

	if l.childCancel != nil {
		l.childCancel()
		l.childCancel = nil
	}
	l.childCtx = nil
	l.started = false
	l.workers = nil
	return nil
}

func (l *LeadershipProvider) worker(key string) *LeadershipInitiator {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if l.workers == nil {
		return nil
	}

	worker, _ := l.workers[key]
	return worker
}

func (l *LeadershipProvider) IsLeader(ctx context.Context, key string) bool {
	logger.WithContext(ctx).Debugf("Checking leadership session for key %q", key)

	worker := l.worker(key)
	if worker != nil {
		return worker.IsLeader(ctx)
	}

	return false
}

func (l *LeadershipProvider) ReleaseLeadership(ctx context.Context, key string) {
	logger.WithContext(ctx).Infof("Releasing leadership session for key %q", key)

	worker := l.worker(key)
	if worker != nil {
		worker.Release(ctx)
	}
}

func NewLeadershipProvider(ctx context.Context) (*LeadershipProvider, error) {
	leaderElectionConfig, err := NewConsulLeaderElectionConfig(ctx)
	if err != nil {
		return nil, err
	}

	if !leaderElectionConfig.Enabled {
		return nil, ErrDisabled
	}

	return &LeadershipProvider{
		cfg: leaderElectionConfig,
	}, nil
}
