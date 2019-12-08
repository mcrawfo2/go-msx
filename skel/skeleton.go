package skel

import (
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"gopkg.in/pipe.v2"
	"os"
	"path"
)

func init() {
	AddTarget("generate-build", "Create the build command and configuration", GenerateBuild)
	AddTarget("generate-app", "Create the application command and configuration", GenerateApp)
	AddTarget("generate-dockerfile", "Create a dockerfile for the application", GenerateDockerfile)
	AddTarget("generate-goland", "Create a Goland project for the application", GenerateGoland)
	AddTarget("generate-kubernetes", "Create production kubernetes manifest templates", GenerateKubernetes)
	AddTarget("generate-git", "Create git repository", GenerateGit)
}

func GenerateSkeleton(args []string) error {
	if err := ConfigureInteractive(nil); err != nil {
		return err
	}
	return ExecTargets(
		"generate-build",
		"generate-app",
		"generate-dockerfile",
		"generate-goland",
		"generate-kubernetes",
		"generate-git")
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
		"Creating go module definition": {
			SourceFile: "go.mod.tpl",
			DestFile:   "go.mod",
		},
		"Creating bootstrap configuration": {SourceFile: "cmd/app/bootstrap.yml"},
		"Creating production profile": {
			SourceFile: "cmd/app/profile.production.yml",
			DestFile:   "cmd/app/${app.name}.production.yml",
		},
		"Creating remote profile": {
			SourceFile: "local/profile.remote.json5",
			DestFile:   "local/${app.name}.remote.json5",
		},
		"Creating application entrypoint source": {SourceFile: "cmd/app/main.go"},
	})
}

func GenerateGoland(args []string) error {
	logger.Info("Generating Goland project")
	return renderTemplates(map[string]Template{
		"Creating module definition": {
			SourceFile: "idea/project.iml.tpl",
			DestFile:   ".idea/${app.name}.iml",
		},
		"Creating project definition": {
			SourceFile: "idea/modules.xml",
			DestFile:   ".idea/modules.xml",
		},
		"Creating vcs definition": {
			SourceFile: "idea/vcs.xml",
			DestFile:   ".idea/vcs.xml",
		},
		"Creating workspace": {
			SourceFile: "idea/workspace.xml",
			DestFile:   ".idea/workspace.xml",
		},
		"Creating run configuration: make clean": {
			SourceFile: "idea/runConfigurations/make_clean.xml",
			DestFile:   ".idea/runConfigurations/make_clean.xml",
		},
		"Creating run configuration: make dist": {
			SourceFile: "idea/runConfigurations/make_dist.xml",
			DestFile:   ".idea/runConfigurations/make_dist.xml",
		},
		"Creating run configuration: make docker": {
			SourceFile: "idea/runConfigurations/make_docker.xml",
			DestFile:   ".idea/runConfigurations/make_docker.xml",
		},
		"Creating run configuration: local": {
			SourceFile: "idea/runConfigurations/project__local_.xml",
			DestFile:   ".idea/runConfigurations/${app.name}__local_.xml",
		},
		"Creating run configuration: remote": {
			SourceFile: "idea/runConfigurations/project__remote_.xml",
			DestFile:   ".idea/runConfigurations/${app.name}__remote_.xml",
		},
	})
}

func GenerateDockerfile(args []string) error {
	logger.Info("Generating Dockerfile")
	return renderTemplates(map[string]Template{
		"Creating Dockerfile": {SourceFile: "docker/Dockerfile"},
	})
}

func GenerateKubernetes(args []string) error {
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

func GenerateGit(args []string) error {
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

	return exec.ExecutePipesIn(
		targetDirectory,
		pipe.Line(
			exec.Info("- Tidying go modules"),
			pipe.Exec("go", "mod", "tidy")),
		pipe.Line(
			exec.Info("- Initializing git repository"),
			pipe.Exec("git", "init", ".")),
		pipe.Line(
			exec.Info("- Staging changes"),
			pipe.Exec("git", "add", "-A")),
		pipe.Line(
			exec.Info("- Committing changes"),
			pipe.Exec("git", "commit", "-m", "Initial Commit")),
	)
}
