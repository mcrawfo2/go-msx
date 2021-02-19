package fs

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/logtest"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"testing"
)

func TestLoggingFilesystem_Open(t *testing.T) {
	type fields struct {
		Name string
		Fs   http.FileSystem
	}
	type args struct {
		name string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     http.File
		wantErr  bool
		checkLog logtest.Check
	}{
		{
			name: "Simple",
			fields: fields{
				Name: "TestLoggingFilesystem_Open",
				Fs: func() http.FileSystem {
					fs := new(MockFileSystem)
					fs.On("Open", "test.json").Return(nil, nil)
					return fs
				}(),
			},
			args: args{
				"test.json",
			},
			want:    nil,
			wantErr: false,
			checkLog: logtest.Check{
				Filters: []logtest.EntryPredicate{
					logtest.HasLevel(logrus.DebugLevel),
				},
				Validators: []logtest.EntryPredicate{
					logtest.HasMessage("TestLoggingFilesystem_Open.Open(test.json)"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := logtest.RecordLogging()
			logger.SetLevel(logrus.DebugLevel)

			l := LoggingFilesystem{
				Name: tt.fields.Name,
				Fs:   tt.fields.Fs,
			}
			got, err := l.Open(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Open() got = %v, want %v", got, tt.want)
			}
			tt.checkLog.Check(r)
		})
	}
}

func TestRootLoggingFilesystem_Open(t *testing.T) {
	type fields struct {
		Dir http.Dir
		Fs  http.FileSystem
	}
	type args struct {
		name string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      http.File
		wantErr   bool
		checkLogs []logtest.Check
	}{
		{
			name: "Simple",
			fields: fields{
				Dir: http.Dir("/home"),
				Fs: func() http.FileSystem {
					fs := new(MockFileSystem)
					fs.On("Open", "test.json").Return(nil, nil)
					return fs
				}(),
			},
			args: args{
				"test.json",
			},
			want:    nil,
			wantErr: false,
			checkLogs: []logtest.Check{{
				Validators: []logtest.EntryPredicate{
					logtest.HasLevel(logrus.DebugLevel),
					logtest.HasMessage("root.Open(/home : test.json)"),
				},
			}},
		},
		{
			name: "Error",
			fields: fields{
				Dir: http.Dir("/home"),
				Fs: func() http.FileSystem {
					fs := new(MockFileSystem)
					fs.On("Open", "test.json").Return(nil, errors.New("error"))
					return fs
				}(),
			},
			args: args{
				"test.json",
			},
			want:    nil,
			wantErr: true,
			checkLogs: []logtest.Check{
				{
					Filters: []logtest.EntryPredicate{
						logtest.HasMessage("root.Open(/home : test.json)"),
					},
					Validators: []logtest.EntryPredicate{
						logtest.HasLevel(logrus.DebugLevel),
					},
				},
				{
					Filters: []logtest.EntryPredicate{
						logtest.HasMessage("Failed to open /home : test.json"),
					},
					Validators: []logtest.EntryPredicate{
						logtest.HasLevel(logrus.DebugLevel),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.SetLevel(logrus.DebugLevel)
			r := logtest.RecordLogging()

			l := RootLoggingFilesystem{
				Dir: tt.fields.Dir,
				fs:  tt.fields.Fs,
			}
			got, err := l.Open(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Open() got = %v, want %v", got, tt.want)
			}
			for _, checkLog := range tt.checkLogs {
				checkLog.Check(r)
			}
		})
	}
}
