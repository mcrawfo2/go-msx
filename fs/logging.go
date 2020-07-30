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
	Fs http.Dir
}

func (l RootLoggingFilesystem) Open(name string) (http.File, error) {
	logger.Debugf("root.Open(%s : %s)", l.Fs, name)
	f, err := l.Fs.Open(name)
	if err != nil {
		logger.WithError(err).Debugf("Failed to open %s : %s", l.Fs, name)
	}
	return f, err
}
