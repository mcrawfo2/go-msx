package consulprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/consul"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"sync"
	"time"
)

var ErrInitiatorAlreadyStarted = errors.New("LeadershipInitiator already started")

type LeadershipInitiator struct {
	properties    LeaderProperties
	parentCtx     context.Context
	acquireCtx    context.Context
	acquireCancel context.CancelFunc
	renewCtx      context.Context
	renewCancel   context.CancelFunc
	started       bool
	acquired      bool
	sessionId     string
	mtx           sync.Mutex
}

func (l *LeadershipInitiator) setAcquired(acquired bool) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.acquired = acquired
}

func (l *LeadershipInitiator) isAcquired() bool {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	return l.acquired
}

func (l *LeadershipInitiator) setStarted(started bool) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.started = started
}

func (l *LeadershipInitiator) isStarted() bool {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	return l.started
}

func (l *LeadershipInitiator) createAcquireContext() {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.acquireCtx, l.acquireCancel = context.WithCancel(l.parentCtx)
}

func (l *LeadershipInitiator) clearAcquireContext() {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.acquireCtx, l.acquireCancel = nil, nil
}

func (l *LeadershipInitiator) cancelAcquireContext() {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if l.acquireCancel != nil {
		l.acquireCancel()
	}
}

func (l *LeadershipInitiator) createRenewContext() {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.renewCtx, l.renewCancel = context.WithCancel(l.acquireCtx)
}

func (l *LeadershipInitiator) clearRenewContext() {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.renewCtx, l.renewCancel = nil, nil
}

func (l *LeadershipInitiator) cancelRenewContext() {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if l.renewCancel != nil {
		l.renewCancel()
	}
}

func (l *LeadershipInitiator) loop() {
	ticker := time.NewTicker(time.Duration(l.properties.BusyWaitMillis) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-l.acquireCtx.Done():
			logger.
				WithContext(l.acquireCtx).
				WithError(l.acquireCtx.Err()).
				Warnf("Leader election loop stopped for key %q", l.properties.Key)
			return
		case <-ticker.C:
			l.acquire(l.acquireCtx)
		}
	}
}

func (l *LeadershipInitiator) acquire(ctx context.Context) {
	var acquired bool
	err := consul.PoolFromContext(ctx).WithConnection(func(connection *consul.Connection) (err error) {
		l.sessionId, _, err = connection.Client().Session().Create(&api.SessionEntry{
			Name:     l.properties.Key,
			Behavior: "release",
			TTL:      (5 * time.Duration(l.properties.HeartBeatMillis) * time.Millisecond).String(),
		}, nil)

		if err != nil {
			return err
		}

		acquired, _, err = connection.Client().KV().Acquire(&api.KVPair{
			Key:     l.properties.Key, // distributed lock
			Value:   []byte(l.sessionId),
			Session: l.sessionId,
		}, nil)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.WithContext(ctx).WithError(err).Errorf("Failed to acquire leadership")
		return
	}

	if !acquired {
		logger.WithContext(ctx).Debugf("Leadership lock unavailable")
		return
	}

	logger.WithContext(ctx).Infof("Leadership lock acquired")
	l.setAcquired(acquired)
	defer l.setAcquired(false)

	l.createRenewContext()
	defer l.clearRenewContext()

	l.heartbeat()
}

func (l *LeadershipInitiator) heartbeat() {
	ticker := time.NewTicker(time.Duration(l.properties.HeartBeatMillis) * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-l.renewCtx.Done():
			logger.
				WithContext(l.renewCtx).
				WithError(l.renewCtx.Err()).
				Warnf("Leader election heartbeat loop stopped for key %q", l.properties.Key)
			return
		case <-ticker.C:
			err := l.renew(l.renewCtx)
			if err != nil {
				logger.WithContext(l.renewCtx).WithError(err).Errorf("Lost leadership")
				return
			}
		}
	}
}

func (l *LeadershipInitiator) renew(ctx context.Context) error {
	return consul.PoolFromContext(ctx).WithConnection(func(connection *consul.Connection) (err error) {
		sessionEntry, _, err := connection.Client().Session().Renew(l.sessionId, nil)
		if err == nil && sessionEntry == nil {
			return errors.New("Consul session invalidated")
		}
		return err
	})
}

func (l *LeadershipInitiator) IsLeader(context.Context) bool {
	return l.isAcquired()
}

func (l *LeadershipInitiator) Release(context.Context) {
	l.cancelRenewContext()
}

func (l *LeadershipInitiator) Stop() {
	l.cancelRenewContext()
	l.cancelAcquireContext()
}

func (l *LeadershipInitiator) Start() {
	err := func() error {
		l.mtx.Lock()
		defer l.mtx.Unlock()

		if l.acquireCtx != nil || l.acquireCancel != nil || l.started {
			return ErrInitiatorAlreadyStarted
		}

		l.started = true
		return nil
	}()
	defer l.setStarted(false)

	if err != nil {
		return
	}

	l.createAcquireContext()
	defer l.clearAcquireContext()

	l.loop()
}

func NewLeadershipInitiator(ctx context.Context, properties LeaderProperties) *LeadershipInitiator {
	return &LeadershipInitiator{
		parentCtx:  ctx,
		properties: properties,
	}
}
