package config

import "context"

type WatchNotifier chan<- struct{}

func (c WatchNotifier) Notify() {
	c <- struct{}{}
}

type Watcher interface {
	Watch(notification WatchNotifier, ctx context.Context)
}
