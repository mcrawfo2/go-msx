package config

import (
	"context"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/pkg/errors"
)

type INIFile struct {
	name   string
	path   string
	reader ContentReader
}

func (f *INIFile) Description() string {
	return fmt.Sprintf("%s: [%s]", f.name, f.path)
}

func (f *INIFile) Load(ctx context.Context) (map[string]string, error) {
	logger.Infof("Loading %s INI config: %s", f.name, f.path)

	settings := map[string]string{}

	bytes, err := f.reader()
	if err != nil {
		return nil, err
	}

	file, err := ini.Load(bytes)
	if err != nil {
		return nil, err
	}

	for _, section := range file.Sections() {
		for _, key := range section.Keys() {
			token := fmt.Sprintf("%s.%s", section.Name(), key.Name())
			settings[token] = key.String()
		}
	}

	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "Failed to load INI config")
	}

	return settings, nil
}

func NewINIFile(name, path string, reader ContentReader) *INIFile {
	return &INIFile{
		name:   name,
		path:   path,
		reader: reader,
	}
}
