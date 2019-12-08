package build

import (
	"archive/zip"
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/build/maven"
	"github.com/bmatcuk/doublestar"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func init() {
	AddTarget("install-swagger-ui", "Installs Swagger-UI package", InstallSwaggerUi)
}

func InstallSwaggerUi(args []string) error {
	version := BuildConfig.Msx.Platform.Version
	artifactTriple := BuildConfig.Msx.Platform.SwaggerArtifact
	artifact := maven.NewArtifactDescriptor(artifactTriple)
	if artifact.Version == "" {
		artifact = artifact.WithVersion(version)
	}

	descriptorPtr, err := maven.ResolveArtifactVersion(artifact)
	if err != nil {
		return err
	}
	descriptor := *descriptorPtr

	// Resolve the Jar file
	repository := maven.NewDefaultHttpRepository()
	artifactFactory, err := maven.NewArtifactFactory(descriptor, repository)
	if err != nil {
		return err
	}

	logger.Infof("Installing configs from %s", descriptor.Triple())
	jarFileArtifact := artifactFactory.CreateArtifact(descriptor.JarFileName())

	// extract public directory to output directory
	if err := installPublicFolder(jarFileArtifact, repository); err != nil {
		logger.WithError(err).Error("Failed to extract configs from jar")
	}

	return nil
}

func installPublicFolder(artifact maven.Artifact, repository maven.ArtifactRepository) error {
	jarFileData, err := artifact.Retrieve(repository)
	if err != nil {
		return err
	}

	byteReader := bytes.NewReader(jarFileData.Data)
	zipReader, err := zip.NewReader(byteReader, jarFileData.Len())
	if err != nil {
		return err
	}

	outputDir := BuildConfig.App.OutputStaticPath()
	if err = os.MkdirAll(outputDir, 0755); err != nil {
		logger.WithError(err).Error("Failed to create output directory")
		return err
	}

	fileGlob := "public/**"
	for _, file := range zipReader.File {
		var ok bool
		ok, err = doublestar.Match(fileGlob, file.Name)
		if err != nil || !ok {
			continue
		}

		filenameParts := strings.SplitN(file.Name, "/", 2)
		basename := filenameParts[1]

		logger.Infof("Installing %s to %s", basename, outputDir)

		readCloser, err := file.Open()
		if err != nil {
			logger.WithError(err).Error("Failed to open file")
			continue
		}

		fileBytes, err := ioutil.ReadAll(readCloser)
		if err = readCloser.Close(); err != nil {
			logger.WithError(err).Error("Failed to read file")
			continue
		}

		outputFilePath := path.Join(outputDir, basename)
		outputDirectory := path.Join(outputDir, path.Dir(basename))
		if err = os.MkdirAll(outputDirectory, 0755); err != nil {
			logger.WithError(err).Error("Failed to create directory")
			continue
		}

		if strings.HasSuffix(file.Name, "/") {
			continue
		}

		if err = ioutil.WriteFile(outputFilePath, fileBytes, 0644); err != nil {
			logger.WithError(err).Error("Failed to write file")
			continue
		}
	}

	return nil
}