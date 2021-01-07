package config

import (
	"context"
)

type Watcher struct {
	cache   *Cache
	entries ProviderEntries
	notify  chan WatcherNotification
}

func (w *Watcher) Watch(ctx context.Context) {
	defer w.Close()

	// We haven't invalidated the cache yet, so this will always succeed
	w.entries, _ = w.cache.Load(ctx)

	go w.cache.Run(ctx)

	for w.cache.Notify() != nil {
		select {
		case <-ctx.Done():
			return

		case _, ok := <-w.cache.Notify():
			if !ok {
				return
			}

			w.cache.Invalidate()
			w.reload(ctx)
		}
	}
}

func (w *Watcher) reload(ctx context.Context) {
	entries, err := w.cache.Load(ctx)
	if err != nil {
		notification := WatcherNotification{
			Provider: w.cache.Provider(),
			Entries:  entries,
			Error:    err,
		}

		w.notify <- notification
		return
	}

	notification := w.compare(entries)

	switch {
	case len(notification.Delta) > 0:
		w.entries = entries.Clone()
		w.notify <- notification

	case notification.Error != nil:
		w.notify <- notification

	}
}

func (w *Watcher) compare(entries ProviderEntries) WatcherNotification {
	err := entries.Validate()

	if err != nil {
		return WatcherNotification{
			Provider: w.cache.Provider(),
			Entries:  entries,
			Error:    err,
		}
	}

	delta := w.entries.Compare(entries)

	return WatcherNotification{
		Provider: w.cache.Provider(),
		Entries:  entries,
		Delta:    delta,
	}
}

func (w *Watcher) Close() {
	// Close the output
	if w.notify != nil {
		close(w.notify)
		w.notify = nil
	}
}

func (w *Watcher) Notify() <-chan WatcherNotification {
	return w.notify
}

func NewWatcher(cache *Cache) *Watcher {
	return &Watcher{
		cache:   cache,
		entries: nil,
		notify:  make(chan WatcherNotification, 1),
	}
}

type WatcherNotification struct {
	Provider Provider
	Error    error
	Entries  ProviderEntries
	Delta    ProviderDelta
}
