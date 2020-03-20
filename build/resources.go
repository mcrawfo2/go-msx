package build

import (
	"github.com/bmatcuk/doublestar"
	copypkg "github.com/otiai10/copy"
	"path/filepath"
	"strings"
)

func init() {
	AddTarget("install-resources", "Installs Resources", InstallResources)
}

func InstallResources(args []string) error {
	files, err := collectIncludedResources()
	if err != nil {
		return err
	}

	for _, inputFilePath := range files {
		outputFilePath := filepath.Join(BuildConfig.OutputResourcesPath(), inputFilePath)
		logger.Infof("Copying %s to %s", inputFilePath, outputFilePath)
		if err := copypkg.Copy(inputFilePath, outputFilePath); err != nil {
			return err
		}
	}

	return nil
}

func collectIncludedResources() ([]string, error) {
	var results []string
	for _, inc := range BuildConfig.Resources.Includes {
		if strings.HasPrefix(inc, "/") {
			inc = inc[1:]
		}

		incFiles, err := doublestar.Glob(inc)
		if err != nil {
			return nil, err
		}

		for _, incFile := range incFiles {
			excluded, err := excludeFilteredResource(incFile)
			if err != nil {
				return nil, err
			}

			if !excluded {
				results = append(results, incFile)
			}
		}
	}

	return results, nil
}

func excludeFilteredResource(included string) (bool, error) {
	excludes := BuildConfig.Resources.Excludes
	excludes = append(excludes, "/dist/**", "/test/**", "/local/**", "/vendor/**")
	for _, exc := range excludes {
		if strings.HasPrefix(exc, "/") {
			exc = exc[1:]
		}

		matches, err := doublestar.Match(exc, included)
		if err != nil {
			return false, err
		} else if matches {
			return true, err
		}
	}
	return false, nil
}
