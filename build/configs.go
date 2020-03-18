package build

import (
	"archive/zip"
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/build/maven"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/bmatcuk/doublestar"
	copypkg "github.com/otiai10/copy"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func init() {
	AddTarget("install-executable-configs", "Copy configured files to distribution config directory", InstallExecutableConfigs)
	AddTarget("install-extra-configs", "Copy custom files to distribution config directory", InstallExtraConfigs)
	AddTarget("install-entrypoint", "Copy custom entrypoint to distribution root directory", InstallEntryPoint)
	AddTarget("install-dependency-configs", "Download dependency config files to distribution config directory", InstallDependencyConfigs)
}

func InstallExecutableConfigs(args []string) error {
	if len(args) > 0 {
		return errors.New("Custom configs must be installed with `install-extra-configs`")
	}

	return installConfigs(
		BuildConfig.InputCommandRoot(),
		BuildConfig.OutputConfigPath(),
		BuildConfig.Executable.ConfigFiles)
}

func InstallExtraConfigs(args []string) error {
	inputDir, _ := filepath.Abs(".")
	outputDir := BuildConfig.OutputConfigPath()
	return installConfigs(inputDir, outputDir, args)
}

func InstallEntryPoint(args []string) error {
	if len(args) > 1 {
		return errors.New("Only one entrypoint script may be specified")
	}

	inputFilePath, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}
	return installEntryPoint(inputFilePath)
}

func InstallDependencyConfigs([]string) error {
	version := BuildConfig.Msx.Platform.Version
	downloaded := &types.StringStack{}
	artifacts := BuildConfig.Msx.Platform.ParentArtifacts

	for _, artifactTriple := range artifacts {
		artifact := maven.NewArtifactDescriptor(artifactTriple)
		if artifact.Version == "" {
			artifact = artifact.WithVersion(version)
		}

		descriptorPtr, err := maven.ResolveArtifactVersion(artifact)
		if err != nil {
			return err
		}
		descriptor := *descriptorPtr

		if err = installDependencyConfigs(descriptor, downloaded); err != nil {
			return err
		}
	}

	return nil
}

func installDependencyConfigs(descriptor maven.ArtifactDescriptor, downloaded *types.StringStack) (err error) {
	if downloaded.Contains(descriptor.ArtifactId) {
		return nil
	} else {
		*downloaded = append(*downloaded, descriptor.ArtifactId)
	}

	if descriptor.Scope == "test" {
		return nil
	}

	// Match the include pattern
	pattern := strings.ReplaceAll(BuildConfig.Msx.Platform.IncludeGroups, ".", "/")
	if ok, err := doublestar.Match(pattern, descriptor.GroupPath()); err != nil {
		return err
	} else if !ok {
		return nil
	}

	// Resolve the POM file
	repository := maven.NewDefaultHttpRepository()
	artifactFactory, err := maven.NewArtifactFactory(descriptor, repository)
	if err != nil {
		return err
	}

	pomFileArtifact := artifactFactory.CreateArtifact(descriptor.PomFileName())
	pomFileData, err := pomFileArtifact.Retrieve(repository)
	if err != nil {
		return err
	}

	pomFile, err := maven.NewPomFile(pomFileData)
	if err != nil {
		return err
	}

	if pomFile.Packaging != "pom" {
		logger.Infof("Installing configs from %s", descriptor.Triple())
		jarFileArtifact := artifactFactory.CreateArtifact(descriptor.JarFileName())

		// extract defaults-*.properties config files to output directory
		if err := installJarConfigs(jarFileArtifact, repository); err != nil {
			logger.WithError(err).Error("Failed to extract configs from jar")
		}

	}

	// Install direct dependencies
	for _, dependency := range pomFile.Dependencies {
		dependency = dependency.WithVersion(descriptor.Version)
		if err = installDependencyConfigs(dependency, downloaded); err != nil {
			return err
		}
	}

	// Install parent
	parent := pomFile.Parent.WithVersion(descriptor.Version)
	if err = installDependencyConfigs(parent, downloaded); err != nil {
		return err
	}

	return nil
}

func installJarConfigs(artifact maven.Artifact, repository maven.ArtifactRepository) error {
	jarFileData, err := artifact.Retrieve(repository)
	if err != nil {
		return err
	}

	byteReader := bytes.NewReader(jarFileData.Data)
	zipReader, err := zip.NewReader(byteReader, jarFileData.Len())
	if err != nil {
		return err
	}

	outputDir := BuildConfig.OutputConfigPath()
	if err = os.MkdirAll(outputDir, 0755); err != nil {
		logger.WithError(err).Error("Failed to create output directory")
		return err
	}

	fileGlob := "**/defaults-*.properties"
	for _, file := range zipReader.File {
		var ok bool
		ok, err = doublestar.Match(fileGlob, file.Name)
		if err != nil || !ok {
			continue
		}
		basename := path.Base(file.Name)

		logger.Infof("Installing %s to %s", basename, outputDir)

		readCloser, err := file.Open()
		if err != nil {
			logger.WithError(err).Error("Failed to open file")
			continue
		}

		bytes, err := ioutil.ReadAll(readCloser)
		if err = readCloser.Close(); err != nil {
			logger.WithError(err).Error("Failed to read file")
			continue
		}

		outputFilePath := path.Join(outputDir, basename)

		if err = ioutil.WriteFile(outputFilePath, bytes, 0644); err != nil {
			logger.WithError(err).Error("Failed to write config file")
			continue
		}
	}

	return nil
}

func installConfigs(inputDir, outputDir string, files []string) error {
	logger.Infof("Source directory: %s", inputDir)
	logger.Infof("Destination directory: %s", outputDir)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	for _, v := range files {
		inputFilePath := v
		if !filepath.IsAbs(v) {
			inputFilePath = filepath.Join(inputDir, v)
		}
		if strings.HasSuffix(v, "/") {
			return errors.New("Configs must be files: " + v)
		}
		outputFileName := path.Base(v)
		if outputFileName == "." || strings.HasSuffix(outputFileName, "/") {
			return errors.New("Configs must be files: " + v)
		}
		outputFilePath := filepath.Join(outputDir, path.Base(v))

		logger.Infof("Copying %s to %s", inputFilePath, outputFilePath)
		if err := copypkg.Copy(inputFilePath, outputFilePath); err != nil {
			return err
		}
	}
	return nil
}

func installEntryPoint(inputFilePath string) error {
	outputFilePath := path.Join(BuildConfig.OutputRoot(), "entrypoint.sh")
	logger.Infof("Source file: %s", inputFilePath)
	logger.Infof("Destination file: %s", outputFilePath)

	if err := copypkg.Copy(inputFilePath, outputFilePath); err != nil {
		return err
	}

	logger.Infof("Copying %s to %s", inputFilePath, outputFilePath)
	return os.Chmod(outputFilePath, 0755)
}
