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
	AddTarget("generate-test", "Create a sample test", GenerateTest)
	AddTarget("generate-local", "Create the local profiles", GenerateLocal)
	AddTarget("generate-dockerfile", "Create a dockerfile for the application", GenerateDockerfile)
	AddTarget("generate-goland", "Create a Goland project for the application", GenerateGoland)
	AddTarget("generate-vscode", "Create a VSCode project for the application", GenerateVsCode)
	AddTarget("generate-kubernetes", "Create production kubernetes manifest templates", GenerateKubernetes)
	AddTarget("generate-deployment-variables", "Create deployment variables manifest", GenerateDeploymentVariables)
	AddTarget("generate-manifest", "Create installer manifest templates", GenerateInstallerManifest)
	AddTarget("generate-jenkins", "Create Jenkins CI templates", GenerateJenkinsCi)
	AddTarget("add-go-msx-dependency", "Add msx dependencies", AddGoMsxDependency)
	AddTarget("generate-git", "Create git repository", GenerateGit)
	AddTarget("generate-webservices", "Create web services from swagger manifest", GenerateDomainOpenApi)
}

// Root command
func GenerateSkeleton(_ []string) error {
	var generators []string
	// Common pre-generators
	generators = append(generators,
		"generate-skel-json",
		"generate-build",
		"generate-app",
		"generate-test")

	// Archetype-specific generators
	generators = append(generators, archetypes.Generators(skeletonConfig.Archetype)...)

	// Common post-generators
	generators = append(generators,
		"generate-deployment-variables",
		"add-go-msx-dependency",
		"generate-local",
		"generate-manifest",
		"generate-dockerfile",
		"generate-goland",
		"generate-vscode",
		"generate-jenkins",
		"generate-git")

	return ExecTargets(generators...)
}

func GenerateSkelJson(_ []string) error {
	logger.Info("Generating skel config")

	bytes, err := json.Marshal(skeletonConfig)
	if err != nil {
		return err
	}

	template := Template{
		Name:       "Creating skel config",
		DestFile:   projectConfigFileName,
		SourceData: bytes,
		Format:     FileFormatJson,
	}

	return template.Render(NewRenderOptions())
}

func GenerateBuild(_ []string) error {
	logger.Info("Generating build command")

	templates := TemplateSet{
		{
			Name:       "Creating Makefile",
			SourceFile: "Makefile",
			Format:     FileFormatMakefile,
		},
		{
			Name:       "Creating build descriptor",
			SourceFile: "cmd/build/build-${generator}.yml",
			DestFile:   "cmd/build/build.yml",
			Format:     FileFormatYaml,
		},
		{
			Name:       "Creating build command source",
			SourceFile: "cmd/build/build.go.tpl",
			DestFile:   "cmd/build/build.go",
		},
	}

	return templates.Render(NewRenderOptions())
}

func GenerateInstallerManifest(_ []string) error {
	logger.Info("Generating installer manifest")
	templates := TemplateSet{
		{
			Name:       "Creating pom.xml",
			SourceFile: "manifest/pom.xml",
			Format:     FileFormatXml,
		},
		{
			Name:       "Creating assembly.xml",
			SourceFile: "manifest/assembly.xml",
			Format:     FileFormatXml,
		},
		{
			Name:       "Creating gitignore",
			SourceFile: "manifest/gitignore",
			DestFile:   "manifest/.gitignore",
			Format:     FileFormatDocker,
		},
	}

	return templates.Render(NewRenderOptions())
}

func GenerateJenkinsCi(_ []string) error {
	logger.Info("Generating Jenkins CI")

	templates := TemplateSet{
		{
			Name:       "Creating Jenkinsfile",
			SourceFile: iff(hasUI(), "build/ci/Jenkinsfile-ui", "build/ci/Jenkinsfile"),
			DestFile:   "build/ci/Jenkinsfile",
			Format:     FileFormatGroovy,
		},
		{
			Name:       "Creating sonar config",
			SourceFile: "build/ci/sonar-project.properties",
			Format:     FileFormatProperties,
		},
		{
			Name:       "Creating Jenkins job config",
			SourceFile: "build/ci/config.xml",
			Format:     FileFormatXml,
		},
	}

	return templates.Render(NewRenderOptions())
}

func GenerateLocal(_ []string) error {
	logger.Info("Generating local profiles")
	template := Template{
		Name:       "Creating remote profile",
		SourceFile: "local/profile.remote.yml",
		DestFile:   "local/${app.name}.remote.yml",
		Format:     FileFormatYaml,
	}

	return template.Render(NewRenderOptions())
}

func GenerateApp(_ []string) error {
	logger.Info("Generating application")
	templates := TemplateSet{
		{
			Name:       "Creating go module definition",
			SourceFile: "go.mod.tpl",
			DestFile:   "go.mod",
			Format:     FileFormatGoMod,
		},
		{
			Name:       "Creating README",
			SourceFile: "README-${generator}.md",
			DestFile:   "README.md",
			Format:     FileFormatMarkdown,
		},
		{
			Name:       "Creating bootstrap configuration",
			SourceFile: "cmd/app/bootstrap-${generator}.yml",
			DestFile:   "cmd/app/bootstrap.yml",
			Format:     FileFormatYaml,
		},
		{
			Name:       "Creating production profile",
			SourceFile: "cmd/app/profile-${generator}.production.yml",
			DestFile:   "cmd/app/${app.name}.production.yml",
			Format:     FileFormatYaml,
		},
		{
			Name:       "Creating application entrypoint source",
			SourceFile: "cmd/app/main-${generator}.go.tpl",
			DestFile:   "cmd/app/main.go",
		},
	}

	return templates.Render(NewRenderOptions())
}

func GenerateMigrate(_ []string) error {
	logger.Info("Generating migration scanner")

	templates := TemplateSet{
		{
			Name:       "Creating migration root sources",
			SourceFile: "internal/migrate/migrate.go.tpl",
			DestFile:   "internal/migrate/migrate.go",
		},
		{
			Name:       "Creating migration version sources",
			SourceFile: fmt.Sprintf("internal/migrate/version/migrate_%s.go.tpl", skeletonConfig.Repository),
			DestFile:   "internal/migrate/${app.migrateVersion}/migrate.go",
		},
	}

	err := templates.Render(NewRenderOptions())
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

func GenerateTest(_ []string) error {
	logger.Info("Generating test")

	templates := TemplateSet{
		{
			Name:       "Creating test sources",
			SourceFile: "internal/empty_test.go",
			DestFile:   "internal/empty_test.go",
		},
	}

	return templates.Render(NewRenderOptions())
}

func AddGoMsxDependency(_ []string) error {
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

func GenerateGoland(_ []string) error {
	logger.Info("Generating Goland project")
	templates := TemplateSet{
		{
			Name:       "Creating module definition",
			SourceFile: "idea/project.iml.tpl",
			DestFile:   ".idea/${app.name}.iml",
			Format:     FileFormatXml,
		},
		{
			Name:       "Creating project definition",
			SourceFile: "idea/modules.xml",
			DestFile:   ".idea/modules.xml",
			Format:     FileFormatXml,
		},
		{
			Name:       "Creating vcs definition",
			SourceFile: "idea/vcs.xml",
			DestFile:   ".idea/vcs.xml",
			Format:     FileFormatXml,
		},
		{
			Name:       "Creating workspace",
			SourceFile: "idea/workspace.xml",
			DestFile:   ".idea/workspace.xml",
			Format:     FileFormatXml,
		},
		{
			Name:       "Creating run configuration: make clean",
			DestFile:   ".idea/runConfigurations/make_clean.xml",
			SourceFile: "idea/runConfigurations/make_clean.xml",
			Format:     FileFormatXml,
		},
		{
			Name:       "Creating run configuration: make test",
			SourceFile: "idea/runConfigurations/make_test.xml",
			DestFile:   ".idea/runConfigurations/make_test.xml",
			Format:     FileFormatXml,
		},
		{
			Name:       "Creating run configuration: make precommit",
			SourceFile: "idea/runConfigurations/make_precommit.xml",
			DestFile:   ".idea/runConfigurations/make_precommit.xml",
			Format:     FileFormatXml,
		},
		{
			Name:       "Creating run configuration: make dist",
			SourceFile: "idea/runConfigurations/make_dist.xml",
			DestFile:   ".idea/runConfigurations/make_dist.xml",
			Format:     FileFormatXml,
		},
		{
			Name:       "Creating run configuration: make docker",
			SourceFile: "idea/runConfigurations/make_docker.xml",
			DestFile:   ".idea/runConfigurations/make_docker.xml",
			Format:     FileFormatXml,
		},
		{
			Name:       "Creating run configuration: make docker-publish",
			SourceFile: "idea/runConfigurations/make_docker_publish.xml",
			DestFile:   ".idea/runConfigurations/make_docker_publish.xml",
			Format:     FileFormatXml,
		},
		{
			Name:       "Creating run configuration: (local)",
			SourceFile: "idea/runConfigurations/project__local_.xml",
			DestFile:   ".idea/runConfigurations/${app.name}__local_.xml",
			Format:     FileFormatXml,
		},
		{
			Name:       "Creating run configuration: migrate (local)",
			SourceFile: "idea/runConfigurations/project_migrate__local_.xml",
			DestFile:   ".idea/runConfigurations/${app.name}_migrate__local_.xml",
			Format:     FileFormatXml,
		},
		{
			Name:       "Creating run configuration: populate (local)",
			SourceFile: "idea/runConfigurations/project_populate__local_.xml",
			DestFile:   ".idea/runConfigurations/${app.name}_populate__local_.xml",
			Format:     FileFormatXml,
		},
		{
			Name:       "Creating run configuration: (remote)",
			SourceFile: "idea/runConfigurations/project__remote_.xml",
			DestFile:   ".idea/runConfigurations/${app.name}__remote_.xml",
			Format:     FileFormatXml,
		},
		{
			Name:       "Creating run configuration: migrate (remote)",
			SourceFile: "idea/runConfigurations/project_migrate__remote_.xml",
			DestFile:   ".idea/runConfigurations/${app.name}_migrate__remote_.xml",
			Format:     FileFormatXml,
		},
		{
			Name:       "Creating run configuration: populate (remote)",
			SourceFile: "idea/runConfigurations/project_populate__remote_.xml",
			DestFile:   ".idea/runConfigurations/${app.name}_populate__remote_.xml",
			Format:     FileFormatXml,
		},
	}

	return templates.Render(NewRenderOptions())
}

func GenerateVsCode(_ []string) error {
	logger.Info("Generating VSCode project")
	templates := TemplateSet{
		{
			Name:       "Creating launch configurations",
			SourceFile: "vscode/launch.json",
			DestFile:   ".vscode/launch.json",
			Format:     FileFormatJson,
		},
		{
			Name:       "Creating task configurations",
			SourceFile: "vscode/tasks.json",
			DestFile:   ".vscode/tasks.json",
			Format:     FileFormatJson,
		},
	}

	return templates.Render(NewRenderOptions())
}

func GenerateDockerfile(_ []string) error {
	logger.Info("Generating Dockerfile")
	template := Template{
		Name:       "Creating Dockerfile",
		SourceFile: "build/package/Dockerfile",
		Format:     FileFormatDocker,
	}

	return template.Render(NewRenderOptions())
}

func GenerateKubernetes(_ []string) error {
	logger.Info("Generating kubernetes manifest templates")
	templates := TemplateSet{
		{
			Name:       "Creating deployment template",
			SourceFile: "deployments/kubernetes-rc.yml.tpl",
			DestFile:   "deployments/kubernetes/${app.name}-rc.yml.tpl",
			Format:     FileFormatYaml,
		},
		{
			Name:       "Creating migrate template",
			SourceFile: "deployments/kubernetes-pod.yml.tpl",
			DestFile:   "deployments/kubernetes/${app.name}-pod.yml.tpl",
			Format:     FileFormatYaml,
		},
		{
			Name:       "Creating populate template",
			SourceFile: "deployments/kubernetes-meta.yml.tpl",
			DestFile:   "deployments/kubernetes/${app.name}-meta.yml.tpl",
			Format:     FileFormatYaml,
		},
		{
			Name:       "Creating pdb template",
			SourceFile: "deployments/kubernetes-pdb.yml.tpl",
			DestFile:   "deployments/kubernetes/${app.name}-pdb.yml.tpl",
			Format:     FileFormatYaml,
		},
	}

	return templates.Render(NewRenderOptions())
}

func GenerateKubernetesForBeats(_ []string) error {
	logger.Info("Generating kubernetes manifest templates for beats")
	templates := TemplateSet{
		{
			Name:       "Creating deployment template",
			SourceFile: "deployments/beats/kubernetes-ps.yml.tpl",
			DestFile:   "deployments/kubernetes/${app.name}-ps.yml.tpl",
			Format:     FileFormatYaml,
		},
		{
			Name:       "Creating config map template",
			SourceFile: "deployments/beats/kubernetes-cm.yml.tpl",
			DestFile:   "deployments/kubernetes/${app.name}-cm.yml.tpl",
			Format:     FileFormatYaml,
		},
		{
			Name:       "Creating pdb template",
			SourceFile: "deployments/beats/kubernetes-pdb.yml.tpl",
			DestFile:   "deployments/kubernetes/${app.name}-pdb.yml.tpl",
			Format:     FileFormatYaml,
		},
	}

	return templates.Render(NewRenderOptions())
}

func GenerateDeploymentVariables(_ []string) error {
	logger.Info("Generating deployment variables")
	templates := TemplateSet{
		{
			Name:       "Creating deployment variables",
			SourceFile: "deployments/deployment_variables.yml",
			DestFile:   "deployments/kubernetes/${deployment.group}_deployment_variables.yml",
			Format:     FileFormatYaml,
		},
	}

	return templates.Render(NewRenderOptions())
}

func GenerateGit(_ []string) error {
	logger.Info("Generating git repository")
	template := Template{
		Name:       "Creating git ignore list",
		SourceFile: "gitignore",
		DestFile:   ".gitignore",
		Format:     FileFormatDocker,
	}
	if err := template.Render(NewRenderOptions()); err != nil {
		return err
	}

	targetDirectory := skeletonConfig.TargetDirectory()
	logger.Infof("Target Directory: %s", targetDirectory)

	gitDirectory := path.Join(targetDirectory, ".git")
	if _, err := os.Stat(gitDirectory); err == nil || !os.IsNotExist(err) {
		logger.Warn(".git directory exists.  Not recreating.")
		return err
	}

	gitRepositoryUrl := fmt.Sprintf(
		"git@cto-github.cisco.com:NFV-BU/%s.git",
		skeletonConfig.AppName)

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
		),
		exec.WithDir(targetDirectory,
			pipe.Line(
				exec.Info("- Setting origin"),
				pipe.Exec("git", "remote", "add", "origin", gitRepositoryUrl)),
		))
}
