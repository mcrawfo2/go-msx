package build

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"cto-github.cisco.com/NFV-BU/go-msx/fs"
	"fmt"
	"github.com/shurcooL/vfsgen"
	"gopkg.in/pipe.v2"
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

func InstallGenerateDependencies(args []string) error {
	goPath, err := getGoPath()
	if err != nil {
		return err
	}

	goos := strings.Title(runtime.GOOS)

	script := pipe.Script(
		exec.Info("Downloading generator dependencies"),
		pipe.ChDir(filepath.Join(goPath, "bin")),
		tarInstall(fmt.Sprintf(
			"https://github.com/vektra/mockery/releases/download/v2.3.0/mockery_2.3.0_%s_x86_64.tar.gz", goos),
			"mockery"),
		goInstall("bou.ke/staticfiles"),
		pipe.Write(os.Stdout),
	)

	return pipe.Run(script)
}

func getGoPath() (string, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 512))

	goBinScript := pipe.Line(
		pipe.Exec("go", "env", "GOBIN"),
		pipe.Write(buf))

	if err := pipe.Run(goBinScript); err != nil {
		return "", err
	}

	goBinFolder := strings.TrimSpace(buf.String())
	if goBinFolder != "" {
		logger.Infof("%q", goBinFolder)
		return goBinFolder, nil
	}

	buf.Truncate(0)

	goPathScript := pipe.Line(
		pipe.Exec("go", "env", "GOPATH"),
		pipe.Write(buf))

	if err := pipe.Run(goPathScript); err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}

func goInstall(packageName string) pipe.Pipe {
	return pipe.Exec("go", "get", packageName)
}

func tarInstall(url, filename string) pipe.Pipe {
	return pipe.Line(
		pipe.Exec("curl", "-L", url),
		pipe.Exec("tar", "-xz"),
		pipe.WriteFile(filename, 0755))
}

func GenerateCode(args []string) error {
	for _, p := range BuildConfig.Generate {
		var err error
		logger.Infof("Generating path '%s'", p.Path)
		if p.VfsGen != nil {
			err = generateCodePathVfs(p)
		} else {
			err = generateCodePathCommand(p)
		}
		if err != nil {
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
