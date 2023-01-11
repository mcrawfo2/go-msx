package npm

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"github.com/bmatcuk/doublestar"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func InstallNodePackageContents(packageName string, version string, packageDir, outputDir string) error {
	repo := NewHttpRepository()
	p, err := repo.PackageInfo(packageName, version)
	if err != nil {
		return err
	}

	// extract build directory to output directory
	stripDir := path.Clean("/" + packageDir)
	mapper := pathMapper(path.Join(packageDir, "**/*"), pathStrip(stripDir))
	return installTarFiles(p.Dist.Tarball, mapper, outputDir)
}

func pathStrip(p string) int {
	splits := strings.Split(p, "/")
	return len(splits) - 1
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

func installTarFiles(packageFileUrl string, mapper pathMapperFunc, outputDir string) error {
	logger.Infof("Downloading %q", packageFileUrl)
	repo := NewHttpRepository()
	packageFileBytes, err := repo.Get(packageFileUrl)
	if err != nil {
		return err
	}

	bytesReader := bytes.NewReader(packageFileBytes)
	gzipReader, err := gzip.NewReader(bytesReader)
	if err != nil {
		return err
	}
	tarReader := tar.NewReader(gzipReader)

	if err = os.MkdirAll(outputDir, 0755); err != nil {
		logger.WithError(err).Error("Failed to create output directory")
		return err
	}

	var header *tar.Header
	var count int
	for {
		header, err = tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		outputFileName := mapper(header.Name)
		if outputFileName == nil {
			continue
		}

		*outputFileName = strings.TrimPrefix(*outputFileName, "/")
		outputPath := filepath.Join(outputDir, *outputFileName)
		logger.Infof("Installing %s to %s", header.Name, outputPath)

		switch header.Typeflag {
		case tar.TypeReg:
			// Create the directory
			outputParent := filepath.Dir(outputPath)
			err = os.MkdirAll(outputParent, 0755)
			if err != nil {
				return err
			}

			var data []byte
			if data, err = ioutil.ReadAll(tarReader); err != nil {
				return err
			}

			if err = ioutil.WriteFile(outputPath, data, 0644); err != nil {
				return err
			}

			count++
		}
	}

	logger.Infof("Installed %d files", count)

	return nil
}
