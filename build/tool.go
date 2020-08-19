package build

import (
	"archive/tar"
	"compress/gzip"
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func init() {
	AddTarget("build-tool", "Build the binary tool", BuildTool)
	AddTarget("publish-tool", "Publish the binary tool", PublishTool)
}

func BuildTool(args []string) error {
	if BuildConfig.Tool.Name == "" {
		return errors.New("Tool name not specified.  Please provide tool.name in build.yml")
	}
	if BuildConfig.Tool.Cmd == "" {
		return errors.New("Tool command name not specified.  Please provide tool.cmd in build.yml")
	}

	if err := buildToolForOs("linux"); err != nil {
		return errors.Wrapf(err, "Failed to build linux binary")
	}

	if err := buildToolForOs("darwin"); err != nil {
		return errors.Wrapf(err, "Failed to build darwin binary")
	}

	return nil
}

func buildToolForOs(goos string) error {
	logger.Infof("Build tool %s for %s", BuildConfig.Tool.Name, goos)

	env := BuildConfig.Go.Environment()
	env["goos"] = goos
	env["goflags"] = ""

	buildArgs := []string{
		"build",
		"-o", path.Join(BuildConfig.OutputToolPath(), BuildConfig.Tool.Name, goos, BuildConfig.Tool.Cmd),
	}

	builderFlags := strings.Fields(os.Getenv("BUILDER_FLAGS"))

	sourceFile := strings.Fields(path.Join("cmd", BuildConfig.Tool.Cmd, "main.go"))

	if _, err := os.Stat(sourceFile[0]); os.IsNotExist(err) {
		sourceFile = strings.Fields(path.Join("cmd", BuildConfig.Tool.Cmd, BuildConfig.Tool.Cmd+".go"))
		if _, err := os.Stat(sourceFile[0]); os.IsNotExist(err) {
			return errors.Errorf("Could not locate main.go or %s.go in cmd/%s",
				BuildConfig.Tool.Cmd,
				BuildConfig.Tool.Cmd)
		}
	}

	err := exec.ExecutePipes(
		exec.WithEnv(env,
			exec.Exec(
				"go",
				buildArgs,
				builderFlags,
				sourceFile)))
	if err != nil {
		return err
	}

	logger.Infof("Successfully built tool %s for %s", BuildConfig.Tool.Name, goos)

	return nil
}

func PublishTool(args []string) error {
	// Ensure we have a tool name
	if BuildConfig.Tool.Name == "" {
		return errors.New("Tool name not specified.  Please provide tool.name in build.yml")
	}

	// Ensure we can log into artifactory
	if BuildConfig.Binaries.Username == "" || BuildConfig.Binaries.Password == "" {
		return errors.New("Artifactory username or password unset: Please supply " +
			"ARTIFACTORY_USERNAME/ARTIFACTORY_PASSWORD environment variables")
	}

	// Ensure we know where artifactory is
	if BuildConfig.Binaries.Repository == "" {
		return errors.New("Artifactory repository unset: Please supply " +
			"artifactory.repository configuration setting in build.yml")
	}

	logger.Info("Publishing tool packages")

	if err := publishToolForOs("linux"); err != nil {
		return errors.Wrapf(err, "Failed to publish tool for linux")
	}

	if err := publishToolForOs("darwin"); err != nil {
		return errors.Wrapf(err, "Failed to publish tool for darwin")
	}

	logger.Info("Successfully published tool packages.")

	return nil
}

func publishToolForOs(goos string) error {
	if err := packageToolForOs(goos); err != nil {
		return errors.Wrapf(err, "Failed to package tool %q for %q", BuildConfig.Tool.Name, goos)
	}

	sourceFile := BuildConfig.Tool.PublishArtifactPath(goos)
	uploadUrl := BuildConfig.Tool.PublishUrl(goos)

	err := uploadArtifactory(sourceFile, uploadUrl)
	if err != nil {
		return err
	}

	uploadUrl = BuildConfig.Tool.PublishLatestUrl(goos)
	return uploadArtifactory(sourceFile, uploadUrl)
}

func packageToolForOs(goos string) error {
	logger.Infof("Creating tool package %q for %s", BuildConfig.Tool.Name, goos)

	// Create the output directory
	toolPath := filepath.Join(BuildConfig.OutputToolPath(), BuildConfig.Tool.Name)
	err := os.MkdirAll(toolPath, 0755)
	if err != nil {
		return err
	}

	toolSourcePath := filepath.Join(toolPath, goos)

	// Generate the tarball into the output directory
	outputFileName := fmt.Sprintf("%s-%s-%s.tar.gz",
		BuildConfig.Tool.Name,
		goos,
		BuildConfig.Build.Number)
	outputFile := filepath.Join(toolPath, outputFileName)
	fw, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer fw.Close()

	gw := gzip.NewWriter(fw)
	tw := tar.NewWriter(gw)

	err = filepath.Walk(toolSourcePath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			logger.WithError(err).Errorf("Failed to process entry at path %q", p)
			return nil
		}

		if info.IsDir() {
			return nil
		} else if !info.Mode().IsRegular() {
			logger.Warnf("Ignoring non-regular file at path %q", p)
			return nil
		}

		subpath := strings.TrimPrefix(p, toolSourcePath+"/")

		logger.Infof("Adding %q", subpath)

		hdr := &tar.Header{
			Name:    subpath,
			Size:    info.Size(),
			Mode:    int64(info.Mode().Perm()),
			ModTime: info.ModTime(),
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return errors.Wrapf(err, "Failed to write tar file header for %q", p)
		}

		data, err := ioutil.ReadFile(p)
		if err != nil {
			return errors.Wrapf(err, "Failed to read file data from %q", p)
		}

		if _, err := tw.Write(data); err != nil {
			return errors.Wrapf(err, "Failed to write file body for %q", p)
		}

		return nil
	})

	if err != nil {
		return err
	}

	err = tw.Close()
	if err != nil {
		return err
	}

	err = gw.Close()
	if err != nil {
		return err
	}

	logger.Infof("Successfully created tool package %q for %q", BuildConfig.Tool.Name, goos)
	return nil

}
