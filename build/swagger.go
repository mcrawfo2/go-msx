package build

import (
	"archive/zip"
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/build/maven"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/bmatcuk/doublestar"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func init() {
	AddTarget("install-swagger-ui", "Installs Swagger-UI package", InstallSwaggerUi)
}

type pathMapperFunc func(string) *string

func pathMapper(glob string, strip int) pathMapperFunc {
	return func(fileName string) *string {
		ok, err := doublestar.Match(glob, fileName)
		if err != nil || !ok {
			return nil
		}

		fileNameParts := strings.SplitN(fileName, "/", strip+1)
		return &fileNameParts[strip]
	}
}

func webJarPathFilter(fn pathMapperFunc, fileNames types.StringStack) pathMapperFunc {
	return func(fileName string) *string {
		targetFileName := fn(fileName)
		if targetFileName == nil {
			return targetFileName
		}
		baseFileName := path.Base(*targetFileName)
		if !fileNames.Contains(baseFileName) {
			return nil
		}
		// webjars/swagger-ui/3.22.2/swagger-ui-bundle.js
		fileNameParts := strings.Split(*targetFileName, "/")
		result := path.Join(append(fileNameParts[:2], fileNameParts[3:]...)...)
		return &result
	}
}

func InstallSwaggerUi(args []string) error {
	publicFilesFilter := pathMapper("public/**", 1)
	if err := installStaticJarContents(BuildConfig.Msx.Platform.SwaggerArtifact, publicFilesFilter); err != nil {
		return err
	}

	resourcesWanted := types.StringStack{
		"swagger-ui.css",
		"swagger-ui-bundle.js",
		"swagger-ui-standalone-preset.js",
	}
	resourceFilesFilter := webJarPathFilter(pathMapper("META-INF/resources/**", 2), resourcesWanted)
	return installStaticJarContents(BuildConfig.Msx.Platform.SwaggerWebJar, resourceFilesFilter)
}

func installStaticJarContents(artifactTriple string, fn pathMapperFunc) error {
	version := BuildConfig.Msx.Platform.Version
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
	if err := installFiles(jarFileArtifact, repository, fn); err != nil {
		logger.WithError(err).Error("Failed to extract configs from jar")
	}

	return nil
}

func installFiles(artifact maven.Artifact, repository maven.ArtifactRepository, fn pathMapperFunc) error {
	jarFileData, err := artifact.Retrieve(repository)
	if err != nil {
		return err
	}

	byteReader := bytes.NewReader(jarFileData.Data)
	zipReader, err := zip.NewReader(byteReader, jarFileData.Len())
	if err != nil {
		return err
	}

	outputDir := BuildConfig.OutputStaticPath()
	if err = os.MkdirAll(outputDir, 0755); err != nil {
		logger.WithError(err).Error("Failed to create output directory")
		return err
	}

	for _, file := range zipReader.File {
		transformedFileName := fn(file.Name)
		if transformedFileName == nil {
			continue
		}
		basename := *transformedFileName

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
