package skel

import (
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"encoding/json"
	"fmt"
	"gopkg.in/pipe.v2"
	"os"
	"path"
)

func init() {
	AddTarget("generate-skel-json", "Create the skel configuration file", GenerateSkelJson)
	AddTarget("generate-build", "Create the build command and configuration", GenerateBuild)
	AddTarget("generate-app", "Create the application command and configuration", GenerateApp)
	AddTarget("generate-migrate", "Create the migrate package", GenerateMigrate)
	AddTarget("generate-local", "Create the local profiles", GenerateLocal)
	AddTarget("generate-dockerfile", "Create a dockerfile for the application", GenerateDockerfile)
	AddTarget("generate-goland", "Create a Goland project for the application", GenerateGoland)
	AddTarget("generate-vscode", "Create a VSCode project for the application", GenerateVsCode)
	AddTarget("generate-kubernetes", "Create production kubernetes manifest templates", GenerateKubernetes)
	AddTarget("generate-manifest", "Create installer manifest templates", GenerateInstallerManifest)
	AddTarget("generate-jenkins", "Create Jenkins CI templates", GenerateJenkinsCi)
	AddTarget("add-go-msx-dependency", "Add msx dependencies", AddGoMsxDependency)
	AddTarget("generate-git", "Create git repository", GenerateGit)
}

func GenerateSkeleton(args []string) error {
	var generators []string

	// Common pre-generators
	generators = append(generators,
		"generate-skel-json",
		"generate-build",
		"generate-app")

	// Archetype-specific generators
	generators = append(generators, archetypes.Generators(skeletonConfig.Archetype)...)

	// Common post-generators
	generators = append(generators,
		"add-go-msx-dependency",
		"generate-local",
		"generate-manifest",
		"generate-dockerfile",
		"generate-goland",
		"generate-vscode",
		"generate-kubernetes",
		"generate-jenkins",
		"generate-git")

	return ExecTargets(generators...)
}

func GenerateSkelJson(args []string) error {
	logger.Info("Generating skel config")

	bytes, err := json.Marshal(skeletonConfig)
	if err != nil {
		return err
	}

	return writeStaticFiles(map[string]StaticFile{
		"Creating skel config": {
			Data:     bytes,
			DestFile: configFileName,
		},
	})
}

func GenerateBuild(args []string) error {
	logger.Info("Generating build command")

	templates := map[string]Template{
		"Creating Makefile": {SourceFile: "Makefile"},
		"Creating build descriptor": {
			SourceFile: "cmd/build/build-${generator}.yml",
			DestFile:   "cmd/build/build.yml",
		},
		"Creating build command source": {
			SourceFile: "cmd/build/build.go.tpl",
			DestFile:   "cmd/build/build.go",
		},
	}

	return renderTemplates(templates)
}

func GenerateInstallerManifest(args []string) error {
	logger.Info("Generating installer manifest")
	return renderTemplates(map[string]Template{
		"Creating pom.xml":      {SourceFile: "manifest/pom.xml"},
		"Creating assembly.xml": {SourceFile: "manifest/assembly.xml"},
		"Creating images manifest": {
			SourceFile: "manifest/resources/manifest-images.yml",
			DestFile:   "manifest/resources/${app.name}-manifest-images.yml",
		},
	})
}

func GenerateJenkinsCi(args []string) error {
	logger.Info("Generating Jenkins CI")
	return renderTemplates(map[string]Template{
		"Creating Jenkinsfile":        {SourceFile: "build/ci/Jenkinsfile"},
		"Creating sonar config":       {SourceFile: "build/ci/sonar-project.properties"},
		"Creating Jenkins job config": {SourceFile: "build/ci/config.xml"},
	})
}

func GenerateLocal(args []string) error {
	logger.Info("Generating local profiles")
	return renderTemplates(map[string]Template{
		"Creating remote profile": {
			SourceFile: "local/profile.remote.yml",
			DestFile:   "local/${app.name}.remote.yml",
		},
	})
}

func GenerateApp(args []string) error {
	logger.Info("Generating application")
	return renderTemplates(map[string]Template{
		"Creating go module definition": {
			SourceFile: "go.mod.tpl",
			DestFile:   "go.mod",
		},
		"Creating README": {
			SourceFile: "README-${generator}.md",
			DestFile:   "README.md",
		},
		"Creating bootstrap configuration": {
			SourceFile: "cmd/app/bootstrap-${generator}.yml",
			DestFile:   "cmd/app/bootstrap.yml",
		},
		"Creating production profile": {
			SourceFile: "cmd/app/profile-${generator}.production.yml",
			DestFile:   "cmd/app/${app.name}.production.yml",
		},
		"Creating beat application entrypoint source": {
			SourceFile: "cmd/app/main-${generator}.go.tpl",
			DestFile:   "cmd/app/main.go",
		},
	})
}

func GenerateMigrate(args []string) error {
	logger.Info("Generating migration scanner")
	err := renderTemplates(map[string]Template{
		"Creating migration root sources": {
			SourceFile: "internal/migrate/migrate.go.tpl",
			DestFile:   "internal/migrate/migrate.go",
		},
		"Creating migration version sources": {
			SourceFile: fmt.Sprintf("internal/migrate/version/migrate_%s.go.tpl", skeletonConfig.Repository),
			DestFile:   "internal/migrate/${app.migrateVersion}/migrate.go",
		},
	})
	if err != nil {
		return err
	}

	err = initializePackageFromFile(
		path.Join(skeletonConfig.TargetDirectory(), "cmd/app/main.go"),
		path.Join(skeletonConfig.AppPackageUrl(), "internal/migrate"))
	if err != nil {
		return err
	}

	err = initializePackageFromFile(
		path.Join(skeletonConfig.TargetDirectory(), "internal/migrate/migrate.go"),
		path.Join(skeletonConfig.AppPackageUrl(), "internal/migrate/"+skeletonConfig.AppMigrateVersion()))
	return err
}

func AddGoMsxDependency(args []string) error {
	logger.Info("Add go-msx dependency")

	targetDirectory := skeletonConfig.TargetDirectory()

	var addDependency = func(name string) pipe.Pipe {
		return exec.WithDir(targetDirectory,
			pipe.Line(
				exec.Info(fmt.Sprintf("- Adding %s to modules", name)),
				pipe.Exec("go", "get", "cto-github.cisco.com/NFV-BU/"+name)))

	}

	pipes := []pipe.Pipe{
		addDependency("go-msx"),
		addDependency("go-msx-build"),
	}

	if skeletonConfig.Archetype == archetypeKeyBeat {
		pipes = append(pipes, addDependency("go-msx-beats"))
	} else if skeletonConfig.Archetype == archetypeKeyServicePack {
		pipes = append(pipes, addDependency("administrationservice"))
		pipes = append(pipes, addDependency("catalogservice"))
	}

	pipes = append(pipes,
		exec.WithDir(targetDirectory,
			pipe.Line(
				exec.Info("- Tidying go modules"),
				pipe.Exec("go", "mod", "tidy")),
		))

	return exec.ExecutePipes(pipes...)
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
		"Creating run configuration: make test": {
			SourceFile: "idea/runConfigurations/make_test.xml",
			DestFile:   ".idea/runConfigurations/make_test.xml",
		},
		"Creating run configuration: make precommit": {
			SourceFile: "idea/runConfigurations/make_precommit.xml",
			DestFile:   ".idea/runConfigurations/make_precommit.xml",
		},
		"Creating run configuration: make dist": {
			SourceFile: "idea/runConfigurations/make_dist.xml",
			DestFile:   ".idea/runConfigurations/make_dist.xml",
		},
		"Creating run configuration: make docker": {
			SourceFile: "idea/runConfigurations/make_docker.xml",
			DestFile:   ".idea/runConfigurations/make_docker.xml",
		},
		"Creating run configuration: make docker-publish": {
			SourceFile: "idea/runConfigurations/make_docker_publish.xml",
			DestFile:   ".idea/runConfigurations/make_docker_publish.xml",
		},
		"Creating run configuration: (local)": {
			SourceFile: "idea/runConfigurations/project__local_.xml",
			DestFile:   ".idea/runConfigurations/${app.name}__local_.xml",
		},
		"Creating run configuration: migrate (local)": {
			SourceFile: "idea/runConfigurations/project_migrate__local_.xml",
			DestFile:   ".idea/runConfigurations/${app.name}_migrate__local_.xml",
		},
		"Creating run configuration: populate (local)": {
			SourceFile: "idea/runConfigurations/project_populate__local_.xml",
			DestFile:   ".idea/runConfigurations/${app.name}_populate__local_.xml",
		},
		"Creating run configuration: (remote)": {
			SourceFile: "idea/runConfigurations/project__remote_.xml",
			DestFile:   ".idea/runConfigurations/${app.name}__remote_.xml",
		},
		"Creating run configuration: migrate (remote)": {
			SourceFile: "idea/runConfigurations/project_migrate__remote_.xml",
			DestFile:   ".idea/runConfigurations/${app.name}_migrate__remote_.xml",
		},
		"Creating run configuration: populate (remote)": {
			SourceFile: "idea/runConfigurations/project_populate__remote_.xml",
			DestFile:   ".idea/runConfigurations/${app.name}_populate__remote_.xml",
		},
	})
}

func GenerateVsCode(args []string) error {
	logger.Info("Generating VSCode project")
	return renderTemplates(map[string]Template{
		"Creating launch configurations": {
			SourceFile: "vscode/launch.json",
			DestFile:   ".vscode/launch.json",
		},
		"Creating task configurations": {
			SourceFile: "vscode/tasks.json",
			DestFile:   ".vscode/tasks.json",
		},
	})
}

func GenerateDockerfile(args []string) error {
	logger.Info("Generating Dockerfile")
	return renderTemplates(map[string]Template{
		"Creating Dockerfile": {SourceFile: "build/package/Dockerfile"},
	})
}

func GenerateKubernetes(args []string) error {
	logger.Info("Generating kubernetes manifest templates")
	return renderTemplates(map[string]Template{
		"Creating deployment template": {
			SourceFile: "deployments/kubernetes-deployment.yml.tpl",
			DestFile:   "deployments/kubernetes/${app.name}-rc.yml.tpl",
		},
		"Creating init template": {
			SourceFile: "deployments/kubernetes-init.yml.tpl",
			DestFile:   "deployments/kubernetes/${app.name}-pod.yml.tpl",
		},
		"Creating pdb template": {
			SourceFile: "deployments/kubernetes-poddisruptionbudget.yml.tpl",
			DestFile:   "deployments/kubernetes/${app.name}-pdb.yml.tpl",
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
	logger.Infof("Target Directory: %s", targetDirectory)

	gitDirectory := path.Join(targetDirectory, ".git")
	if _, err := os.Stat(gitDirectory); err == nil || !os.IsNotExist(err) {
		return err
	}

	return exec.ExecutePipes(
		exec.WithDir(targetDirectory,
			pipe.Line(
				exec.Info("- Initializing git repository"),
				pipe.Exec("git", "init", "."))),
		exec.WithDir(targetDirectory,
			pipe.Line(
				exec.Info("- Staging changes"),
				pipe.Exec("git", "add", "-A"))),
		exec.WithDir(targetDirectory,
			pipe.Line(
				exec.Info("- Committing changes"),
				pipe.Exec("git", "commit", "-m", "Initial Commit")),
		))
}
