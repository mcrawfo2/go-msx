// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package build

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"cto-github.cisco.com/NFV-BU/go-msx/fs"
	"github.com/pkg/errors"
	"github.com/shurcooL/vfsgen"
	"gopkg.in/pipe.v2"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var commandRegexp = regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)'`)

func init() {
	AddTarget("download-generate-deps", "Download generate dependencies", InstallGenerateDependencies)
	AddTarget("generate", "Generate code", GenerateCode)
}

func InstallGenerateDependencies(_ []string) error {
	// Binary output path
	targetPath, err := getGoBin()
	if err != nil {
		return err
	}

	// Book generator
	mdBookUrl, err := getPlatformRustAssetUrl("https://github.com/rust-lang/mdBook/releases/download/v0.4.21/mdbook-v0.4.21")
	if err != nil {
		return err
	}

	// Book mermaid -> image generator
	mdBookMermaidUrl, err := getPlatformRustAssetUrl("https://github.com/badboy/mdbook-mermaid/releases/download/v0.11.2/mdbook-mermaid-v0.11.2")
	if err != nil {
		return err
	}

	script := []pipe.Pipe{
		exec.Info("Downloading generator dependencies"),
		pipe.Script(
			pipe.ChDir(targetPath),
			goInstall("github.com/vektra/mockery/v2@v2.22.1"),
		),
		tarInstall(
			mdBookUrl,
			"mdbook",
			targetPath),
		tarInstall(
			mdBookMermaidUrl,
			"mdbook-mermaid",
			targetPath),
	}

	return exec.ExecutePipes(script...)
}

func getGoBin() (string, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 2048))

	goBinScript := pipe.Line(
		pipe.Exec("go", "env", "GOBIN"),
		pipe.Write(buf))

	if err := pipe.Run(goBinScript); err != nil {
		return "", err
	}

	goBinFolder := strings.TrimSpace(buf.String())
	if goBinFolder != "" {
		logger.Infof("GOBIN=%q", goBinFolder)
		return goBinFolder, nil
	}

	goPathFolder, err := getGoPath()
	if err != nil {
		return "", nil
	}

	return filepath.Join(goPathFolder, "bin"), nil
}

func getGoPath() (string, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 2048))

	goPathScript := pipe.Line(
		pipe.Exec("go", "env", "GOPATH"),
		pipe.Write(buf))

	if err := pipe.Run(goPathScript); err != nil {
		return "", err
	}

	goPathFolder := strings.TrimSpace(buf.String())
	return goPathFolder, nil
}

func tarInstall(url, filename, path string) pipe.Pipe {
	var extractor pipe.Pipe
	var err error
	var tempName string

	switch {
	case strings.HasSuffix(url, ".tar.gz"):
		extractor = pipe.Line(
			exec.ReadUrl(url),
			pipe.Exec("tar", "xzOf", "-", filename),
		)
	case strings.HasSuffix(url, ".tar.bz2"):
		extractor = pipe.Line(
			exec.ReadUrl(url),
			pipe.Exec("tar", "xjOf", "-", filename),
		)
	case strings.HasSuffix(url, ".tar.xz"):
		extractor = pipe.Line(
			exec.ReadUrl(url),
			pipe.Exec("tar", "xJOf", "-", filename),
		)
	case strings.HasSuffix(url, ".zip"):
		if tempName, err = downloadTemp(url, "install.*.zip"); err != nil {
			logger.WithError(err).Error("Failed to download zip file")
			return nil
		}
		extractor = pipe.Script(
			pipe.Exec("unzip", tempName, filename, "-d", path),
			exec.RemoveAll(tempName),
		)
	}

	return pipe.Line(
		exec.Info("Installing %q", filename),
		pipe.MkDirAll(path, 0755),
		extractor,
		pipe.WriteFile(filepath.Join(path, filename), 0755))
}

func GenerateCode(args []string) error {
	for _, p := range BuildConfig.Generate {
		var err error
		logger.Infof("Generating path '%s'", p.Path)
		if p.VfsGen != nil {
			err = generateCodePathVfs(p)
		} else if len(p.BuiltIn) != 0 {
			err = generateBuiltin(p)
		} else {
			err = generateCodePathCommand(p)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func generateBuiltin(p Generate) error {
	for _, name := range p.BuiltIn {
		target := Target(name)
		if err := target.RunE(target, nil); err != nil {
			return err
		}
	}
	return nil
}

func generateCodePathCommand(p Generate) error {
	if p.Command == "" {
		p.Command = "go generate"
	}

	args := commandRegexp.FindAllString(p.Command, -1)
	return exec.ExecutePipes(
		exec.WithDir(p.Path,
			exec.Exec(args[0], args[1:])))
}

func generateCodePathVfs(p Generate) (err error) {
	root, err := getGenerateRootDir(p)
	if err != nil {
		return err
	}

	targetFs, err := fs.NewGlobFileSystem(http.Dir(root), p.VfsGen.Includes, p.VfsGen.Excludes)
	if err != nil {
		return err
	}

	projectRoot, err := os.Getwd()
	if err != nil {
		return err
	}

	fileName := path.Join(projectRoot, p.Path, p.VfsGen.Filename)

	packageName := path.Base(p.Path)
	return vfsgen.Generate(targetFs, vfsgen.Options{
		Filename:     fileName,
		PackageName:  packageName,
		VariableName: p.VfsGen.VariableName,
	})
}

func getGenerateRootDir(p Generate) (string, error) {
	abs, err := filepath.Abs(p.Path)
	if err != nil {
		return "", err
	}

	if p.VfsGen.Root != "" {
		abs, err = filepath.Abs(filepath.Join(abs, p.VfsGen.Root))
		if err != nil {
			return "", err
		}
	}

	return abs, nil
}

func getPlatformRustAssetUrl(baseUrl string) (string, error) {
	var triple = ""
	var suffix = ""
	switch runtime.GOOS {
	case "darwin":
		triple = "x86_64-apple-darwin"
		suffix = ".tar.gz"
	case "linux":
		triple = "x86_64-unknown-linux-gnu"
		suffix = ".tar.gz"
	case "windows":
		triple = "x86_64-pc-windows-msvc"
		suffix = ".zip"
	default:
		return baseUrl, errors.Errorf("mdbook not supported on this platform %q", runtime.GOOS)
	}

	return baseUrl + "-" + triple + suffix, nil
}

func downloadTemp(url string, tempName string) (string, error) {
	var writer *os.File
	var response *http.Response
	var err error

	if writer, err = os.CreateTemp("", tempName); err != nil {
		return "", err
	} else if response, err = http.DefaultClient.Get(url); err != nil {
		return "", err
	} else if _, err = io.Copy(writer, response.Body); err != nil {
		return "", err
	} else if err = writer.Close(); err != nil {
		return "", err
	}

	return writer.Name(), nil
}
