package config

import (
	"context"
	"github.com/pkg/errors"
	"github.com/radovskyb/watcher"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func NewMockFileWatcher(filename string) *MockFileWatcher {
	closedChan := make(chan struct{})

	mockFileWatcher := new(MockFileWatcher)
	mockFileWatcher.On("Closed").Return((<-chan struct{})(closedChan))
	mockFileWatcher.On("Close").Run(func(args mock.Arguments) {
		closedChan <- struct{}{}
	})

	return mockFileWatcher
}

func TestFileWatcher_Implementations(t *testing.T) {
	var _ FileWatcher = new(MockFileWatcher)
	var _ FileWatcher = new(WatcherFileWatcher)
}

func TestFileNotifier_Notify(t *testing.T) {
	var ch = make(chan struct{})

	type fields struct {
		filename       string
		notify         chan struct{}
		watcherFactory FileWatcherFactory
	}
	tests := []struct {
		name   string
		fields fields
		want   <-chan struct{}
	}{
		{
			name: "Simple",
			fields: fields{
				filename:       "testdata/config.json",
				notify:         ch,
				watcherFactory: func(filename string) (FileWatcher, error) {
					return new(MockFileWatcher), nil
				},
			},
			want: ch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &FileNotifier{
				filename:       tt.fields.filename,
				notify:         tt.fields.notify,
				watcherFactory: tt.fields.watcherFactory,
			}
			if got := i.Notify(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Notify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileNotifier_Run(t *testing.T) {
	var ch = make(chan struct{})
	var ctx context.Context = context.Background()

	type fields struct {
		filename       string
		notify         chan struct{}
		watcherFactory FileWatcherFactory
	}
	tests := []struct {
		name   string
		fields fields
		want   <-chan struct{}
	}{
		{
			name: "Simple",
			fields: fields{
				filename:       "testdata/config.json",
				notify:         ch,
				watcherFactory: func(filename string) (FileWatcher, error) {
					events := make(chan watcher.Event)
					errs := make(chan error)
					mockFileWatcher := NewMockFileWatcher(filename)
					mockFileWatcher.On("Event").Return((<-chan watcher.Event)(events))
					mockFileWatcher.On("Error").Return((<-chan error)(errs))
					mockFileWatcher.On("Start", mock.AnythingOfType("time.Duration")).Run(func(args mock.Arguments) {
						// Send event
						go func() {
							errs <- errors.New("some error")
							events <- watcher.Event{
								Op:       watcher.Write,
								Path:     filename,
								OldPath:  filename,
								FileInfo: nil,
							}
						}()
						// Wait for event
						go func() {
							<- ch
							mockFileWatcher.Close()
						}()
					}).Return(errors.New("some error"))

					return mockFileWatcher, nil
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileNotifier := &FileNotifier{
				filename:       tt.fields.filename,
				notify:         tt.fields.notify,
				watcherFactory: tt.fields.watcherFactory,
			}

			fileNotifier.Run(ctx)
		})
	}
}

func TestNewWatcherFileWatcher(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Simple",
			args:    args{
				filename: "testdata/config.json",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWatcherFileWatcher(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWatcherFileWatcher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.want {
				t.Errorf("NewWatcherFileWatcher() got = %v, want %v", got, tt.want)
			}
		})
	}
}