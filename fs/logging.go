package fs

import (
	"net/http"
)

type LoggingFilesystem struct {
	Name string
	Fs   http.FileSystem
}

func (l LoggingFilesystem) Open(name string) (http.File, error) {
	logger.Debugf("%s.Open(%s)", l.Name, name)
	return l.Fs.Open(name)
}

type RootLoggingFilesystem struct {
	Dir http.Dir
	fs  http.FileSystem
}

func (l RootLoggingFilesystem) Open(name string) (http.File, error) {
	logger.Debugf("root.Open(%s : %s)", l.Dir, name)
	fs := l.fs
	if fs == nil {
		fs = l.Dir
	}
	f, err := fs.Open(name)
	if err != nil {
		logger.WithError(err).Debugf("Failed to open %s : %s", l.Dir, name)
	}
	return f, err
}
