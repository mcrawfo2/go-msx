package config

import (
	"path"
	"strings"
)

func NewFileProvider(name, fileName string) Provider {
	fileExt := strings.ToLower(path.Ext(fileName))
	switch fileExt {
	case ".yml", ".yaml":
		return NewCachedLoader(NewYAMLFile(name, fileName))
	case ".ini":
		return NewCachedLoader(NewINIFile(name, fileName))
	case ".json", ".json5":
		return NewCachedLoader(NewJSONFile(name, fileName))
	case ".toml":
		return NewCachedLoader(NewTOMLFile(name, fileName))
	case ".properties":
		return NewCachedLoader(NewPropertiesFile(name, fileName))
	default:
		logger.Error("Unknown config file extension: ", fileExt)
		return nil
	}
}
