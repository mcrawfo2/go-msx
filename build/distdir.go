package build

import (
	"os"
	"path/filepath"
)

func init() {
	AddTarget("create-dist-dir", "Create the distribution directory for build outputs", CreateDistDir)
	AddTarget("delete-dist-dir", "Remove the distribution directory and all build outputs", DeleteDistDir)
}

func CreateDistDir(args []string) error {
	distDir, err := filepath.Abs(BuildConfig.OutputRoot())
	if err != nil {
		return err
	}
	logger.Infof("Creating distribution directory: %s", distDir)
	return os.MkdirAll(distDir, 0755)
}

func DeleteDistDir(args []string) error {
	distDir, err := filepath.Abs(BuildConfig.OutputRoot())
	if err != nil {
		return err
	}
	logger.Infof("Removing distribution directory: %s", distDir)
	return os.RemoveAll(distDir)
}
