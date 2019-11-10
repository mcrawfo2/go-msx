package config

import (
	"context"
	"github.com/radovskyb/watcher"
	"time"
)

type FileWatchInvalidator struct {
	provider CachedLoader
	fileName string
}

func (i *FileWatchInvalidator) Description() string {
	return i.provider.Description()
}

func (i *FileWatchInvalidator) Load(ctx context.Context) (map[string]string, error) {
	return i.provider.Load(ctx)
}

func (i *FileWatchInvalidator) Watch(n WatchNotifier, ctx context.Context) {
	w := watcher.New()

	go func() {
		for {
			select {
			case event := <-w.Event:
				logger.Debugf("File event received for %s: %v", i.fileName, event)
				i.provider.Invalidate()
				n.Notify()
			case err := <-w.Error:
				logger.Error(err.Error())
			case <-w.Closed:
				return
			case <-ctx.Done():
				w.Close()
			}
		}
	}()

	if err := w.Start(time.Millisecond * 1000); err != nil {
		logger.Error(err.Error())
	}
}

func NewFileWatchInvalidator(provider CachedLoader, fileName string) *FileWatchInvalidator {
	result := &FileWatchInvalidator{
		provider: provider,
		fileName: fileName,
	}

	return result
}
