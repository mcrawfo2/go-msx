package skel

import (
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"os"
	"path"
)

func init() {
	AddTarget("generate-build", "Create the build command and configuration", GenerateBuild)
	AddTarget("generate-app", "Create the application command and configuration", GenerateApp)
	AddTarget("generate-dockerfile", "Create a dockerfile for the application", GenerateDockerfile)
	AddTarget("kubernetes-manifest-templates", "Create production kubernetes manifest templates", GenerateKubernetesManifestTemplates)
}

func GenerateSkeleton(args []string) error {
	if err := ConfigureInteractive(nil); err != nil {
		return err
	}
	if err := GenerateBuild(nil); err != nil {
		return err
	}
	if err := GenerateApp(nil); err != nil {
		return err
	}
	if err := GenerateDockerfile(nil); err != nil {
		return err
	}
	if err := GenerateKubernetesManifestTemplates(nil); err != nil {
		return err
	}
	if err := GenerateRepository(nil); err != nil {
		return err
	}
	return nil
}

func GenerateBuild(args []string) error {
	logger.Info("Generating build command")
	return renderTemplates(map[string]Template{
		"Creating Makefile":             {SourceFile: "Makefile"},
		"Creating build descriptor":     {SourceFile: "cmd/build/build.yml"},
		"Creating build command source": {SourceFile: "cmd/build/build.go"},
	})
}

func GenerateApp(args []string) error {
	logger.Info("Generating application")
	return renderTemplates(map[string]Template{
		"Creating go module definition":    {SourceFile: "go.mod"},
		"Creating bootstrap configuration": {SourceFile: "cmd/app/bootstrap.yml"},
		"Creating production profile": {
			SourceFile: "cmd/app/profile.production.yml",
			DestFile:   "cmd/app/${app.name}.production.yml",
		},
		"Creating application entrypoint source": {SourceFile: "cmd/app/main.go"},
	})
}

func GenerateDockerfile(args []string) error {
	logger.Info("Generating Dockerfile")
	return renderTemplates(map[string]Template{
		"Creating Dockerfile": {SourceFile: "docker/Dockerfile"},
	})
}

func GenerateKubernetesManifestTemplates(args []string) error {
	logger.Info("Generating kubernetes manifest templates")
	return renderTemplates(map[string]Template{
		"Creating deployment template": {
			SourceFile: "k8s/kubernetes-deployment.yml.tpl",
			DestFile:   "k8s/${app.name}-rc.yml.tpl",
		},
		"Creating init template": {
			SourceFile: "k8s/kubernetes-init.yml.tpl",
			DestFile:   "k8s/${app.name}-pod.yml.tpl",
		},
		"Creating pdb template": {
			SourceFile: "k8s/kubernetes-poddisruptionbudget.yml.tpl",
			DestFile:   "k8s/${app.name}-pdb.yml.tpl",
		},
	})
}

func GenerateRepository(args []string) error {
	logger.Info("Generating git repository")
	err := renderTemplates(map[string]Template{
		"Creating .gitignore": {
			SourceFile: "gitignore",
			DestFile:   ".gitignore",
		},
	})
	if err != nil {
		return err
	}

	targetDirectory := skeletonConfig.TargetDirectory()
	gitDirectory := path.Join(targetDirectory, ".git")
	if _, err := os.Stat(gitDirectory); err == nil || !os.IsNotExist(err) {
		return err
	}

	logger.Info("- Initializing git repository")
	if err = exec.MustExecuteIn(targetDirectory, "git", "init", "."); err != nil {
		return err
	}
	logger.Info("- Staging changes")
	if err = exec.MustExecuteIn(targetDirectory, "git", "add", "-A"); err != nil {
		return err
	}
	logger.Info("- Committing changes")
	return exec.MustExecuteIn(targetDirectory, "git", "commit", "-m", "Initial commit")
}
