// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package skel

import (
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/text"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/pipe.v2"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var ErrNoTemplates = errors.Errorf("no templates")

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
	AddTarget("generate-harness", "Create production harness files", GenerateHarness)
	AddTarget("generate-deployment-variables", "Create deployment variables manifest", GenerateDeploymentVariables)
	AddTarget("generate-manifest", "Create installer manifest templates", GenerateInstallerManifest)
	AddTarget("generate-jenkins", "Create Jenkins CI templates", GenerateJenkinsCi)
	AddTarget("add-go-msx-dependency", "Add msx dependencies", AddGoMsxDependency)
	AddTarget("generate-git", "Create git repository", GenerateGit)
	AddTarget("generate-github", "Create github configuration files", GenerateGithub)
	AddTarget("generate-spui", "Create service pack UI", GenerateSPUI)
}

var defaultPreGenerators = []string{
	"generate-skel-json",
	"generate-build",
	"generate-app",
	"generate-test",
}

var defaultPostGenerators = []string{
	"generate-deployment-variables",
	"add-go-msx-dependency",
	"generate-local",
	"generate-manifest",
	"generate-dockerfile",
	"generate-goland",
	"generate-vscode",
	"generate-jenkins",
	"generate-github",
	"generate-git",
}

type prePost struct {
	pre  []string
	post []string
}

var prePostByArchetype = map[string]prePost{
	archetypeKeySPUI: {pre: []string{}, post: []string{}}, // override with nothing for SPUI
}

func prePostGenerators(target string) (preAndPost prePost) { // everything else uses the defaults
	prp, found := prePostByArchetype[target]
	if !found {
		prp.pre = defaultPreGenerators
		prp.post = defaultPostGenerators
	}
	return prp
}

// GenerateSkeleton is the root command
func GenerateSkeleton(_ []string) error {
	var generators []string

	preAndPost := prePostGenerators(skeletonConfig.Archetype)
	generators = append(generators, preAndPost.pre...)

	// Archetype-specific generators
	generators = append(generators, archetypes.Generators(skeletonConfig.Archetype)...)

	generators = append(generators, preAndPost.post...)

	logger.Infof("Using archetype: %s", skeletonConfig.Archetype)
	logger.Infof("Generators will be: %s", generators)

	return ExecTargets(generators...)
}

func GenerateSkelJson(_ []string) error {
	logger.Info("Generating skel config")

	noTargetDirectory := *skeletonConfig
	noTargetDirectory.TargetParent = ""
	noTargetDirectory.TargetDir = ""

	bytes, err := json.MarshalIndent(noTargetDirectory, "", "    ")
	if err != nil {
		return err
	}

	template := Template{
		Name:       "Creating skel config",
		DestFile:   projectConfigFileName,
		SourceData: bytes,
		Format:     text.FileFormatJson,
	}

	return template.Render(NewRenderOptions())
}

func GenerateBuild(_ []string) error {
	logger.Info("Generating build command")

	templates := TemplateSet{
		{
			Name:       "Creating Makefile",
			SourceFile: "Makefile",
			Format:     text.FileFormatMakefile,
		},
		{
			Name:       "Creating build descriptor",
			SourceFile: "cmd/build/build-${generator}.yml",
			DestFile:   "cmd/build/build.yml",
			Format:     text.FileFormatYaml,
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
			Format:     text.FileFormatXml,
		},
		{
			Name:       "Creating assembly.xml",
			SourceFile: "manifest/assembly.xml",
			Format:     text.FileFormatXml,
		},
		{
			Name:       "Creating gitignore",
			SourceFile: "manifest/gitignore",
			DestFile:   "manifest/.gitignore",
			Format:     text.FileFormatDocker,
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
			Format:     text.FileFormatGroovy,
		},
		{
			Name:       "Creating sonar config",
			SourceFile: "build/ci/sonar-project.properties",
			Format:     text.FileFormatProperties,
		},
		{
			Name:       "Creating Jenkins job config",
			SourceFile: "build/ci/config.xml",
			Format:     text.FileFormatXml,
		},
		{
			Name:       "Creating checks file",
			SourceFile: "build/ci/checks.yml",
			Format:     text.FileFormatYaml,
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
		Format:     text.FileFormatYaml,
	}

	return template.Render(NewRenderOptions())
}

func GenerateApp(_ []string) error {
	if IsExternal {
		if err := FindPublicGoMsxVersion(); err != nil {
			return err
		}
	}

	logger.Info("Generating application")
	templates := TemplateSet{
		{
			Name:       "Creating go module definition",
			SourceFile: "go.mod.tpl",
			DestFile:   "go.mod",
			Format:     text.FileFormatGoMod,
		},
		{
			Name:       "Creating README",
			SourceFile: "README-${generator}.md",
			DestFile:   "README.md",
			Format:     text.FileFormatMarkdown,
		},
		{
			Name:       "Creating bootstrap configuration",
			SourceFile: "cmd/app/bootstrap-${generator}.yml",
			DestFile:   "cmd/app/bootstrap.yml",
			Format:     text.FileFormatYaml,
		},
		{
			Name:       "Creating production profile",
			SourceFile: "cmd/app/profile-${generator}.production.yml",
			DestFile:   "cmd/app/${app.name}.production.yml",
			Format:     text.FileFormatYaml,
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

	err = InitializePackageFromFile(
		path.Join(skeletonConfig.TargetDirectory(), "cmd/app/main.go"),
		path.Join(skeletonConfig.AppPackageUrl(), "internal/migrate"))
	if err != nil {
		return err
	}

	err = InitializePackageFromFile(
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

var GoMsxVersion string

func FindPublicGoMsxVersion() error {
	resp, err := http.Get("https://api.github.com/repos/mcrawfo2/go-msx/tags")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.Errorf("Failed to retrieve GitHub tags, returned %d", resp.StatusCode)
	}

	type GithubTag struct {
		Name string `json:"name"`
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve tags")
	}

	var tags []GithubTag
	if err = json.Unmarshal(bodyBytes, &tags); err != nil {
		return errors.Wrap(err, "Failed to decode tags")
	}

	var highestTagVersion types.Version
	var highestTagName string
	for _, tag := range tags {
		if !strings.HasPrefix(tag.Name, "v") {

		}

		tagName := strings.TrimPrefix(tag.Name, "v")
		tagVersion, err := types.NewVersion(tagName)
		if err != nil {
			continue
		}

		if highestTagName == "" || highestTagVersion.Lt(tagVersion) {
			highestTagVersion = tagVersion
			highestTagName = tag.Name
		}
		break
	}

	if highestTagName == "" {
		return errors.New("Failed to locate any recent version tags")
	}

	GoMsxVersion = highestTagName
	return nil
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
	var pipes = []pipe.Pipe{
		addDependency("go-msx"),
	}

	if !IsExternal {
		pipes = append(pipes, addDependency("go-msx-build"))
		pipes = append(pipes, addDependency("go-msx-populator"))

		if skeletonConfig.Archetype == archetypeKeyBeat {
			pipes = append(pipes, addDependency("go-msx-beats"))
		} else if skeletonConfig.Archetype == archetypeKeyServicePack {
			pipes = append(pipes, addDependency("administrationservice"))
			pipes = append(pipes, addDependency("catalogservice"))
		}
	}

	pipes = append(pipes,
		exec.WithDir(targetDirectory,
			pipe.Line(
				exec.Info("- Tidying go modules"),
				pipe.Exec("go", "mod", "tidy")),
		))

	return exec.ExecutePipes(pipes...)
}

func AddDependencies(deps []string) error {
	logger.Info("Adding dependencies")

	targetDirectory := skeletonConfig.TargetDirectory()

	var addDependency = func(name string) pipe.Pipe {
		return exec.WithDir(targetDirectory,
			pipe.Line(
				exec.Info(fmt.Sprintf("- Adding %s to modules", name)),
				pipe.Exec("go", "get", name)))

	}

	var pipes []pipe.Pipe
	for _, dep := range deps {
		pipes = append(pipes, addDependency(dep))
	}

	return exec.ExecutePipes(pipes...)
}

func GenerateGoland(_ []string) error {
	logger.Info("Generating Goland project")
	templates := TemplateSet{
		{
			Name:       "Creating module definition",
			SourceFile: "idea/project.iml.tpl",
			DestFile:   ".idea/${app.name}.iml",
			Format:     text.FileFormatXml,
		},
		{
			Name:       "Creating project definition",
			SourceFile: "idea/modules.xml",
			DestFile:   ".idea/modules.xml",
			Format:     text.FileFormatXml,
		},
		{
			Name:       "Creating vcs definition",
			SourceFile: "idea/vcs.xml",
			DestFile:   ".idea/vcs.xml",
			Format:     text.FileFormatXml,
		},
		{
			Name:       "Creating workspace",
			SourceFile: "idea/workspace.xml",
			DestFile:   ".idea/workspace.xml",
			Format:     text.FileFormatXml,
		},
		{
			Name:       "Creating run configuration: make clean",
			DestFile:   ".idea/runConfigurations/make_clean.xml",
			SourceFile: "idea/runConfigurations/make_clean.xml",
			Format:     text.FileFormatXml,
		},
		{
			Name:       "Creating run configuration: make test",
			SourceFile: "idea/runConfigurations/make_test.xml",
			DestFile:   ".idea/runConfigurations/make_test.xml",
			Format:     text.FileFormatXml,
		},
		{
			Name:       "Creating run configuration: make precommit",
			SourceFile: "idea/runConfigurations/make_precommit.xml",
			DestFile:   ".idea/runConfigurations/make_precommit.xml",
			Format:     text.FileFormatXml,
		},
		{
			Name:       "Creating run configuration: make dist",
			SourceFile: "idea/runConfigurations/make_dist.xml",
			DestFile:   ".idea/runConfigurations/make_dist.xml",
			Format:     text.FileFormatXml,
		},
		{
			Name:       "Creating run configuration: make docker",
			SourceFile: "idea/runConfigurations/make_docker.xml",
			DestFile:   ".idea/runConfigurations/make_docker.xml",
			Format:     text.FileFormatXml,
		},
		{
			Name:       "Creating run configuration: make docker-publish",
			SourceFile: "idea/runConfigurations/make_docker_publish.xml",
			DestFile:   ".idea/runConfigurations/make_docker_publish.xml",
			Format:     text.FileFormatXml,
		},
		{
			Name:       "Creating run configuration: (local)",
			SourceFile: "idea/runConfigurations/project__local_.xml",
			DestFile:   ".idea/runConfigurations/${app.name}__local_.xml",
			Format:     text.FileFormatXml,
		},
		{
			Name:       "Creating run configuration: migrate (local)",
			SourceFile: "idea/runConfigurations/project_migrate__local_.xml",
			DestFile:   ".idea/runConfigurations/${app.name}_migrate__local_.xml",
			Format:     text.FileFormatXml,
		},
		{
			Name:       "Creating run configuration: populate (local)",
			SourceFile: "idea/runConfigurations/project_populate__local_.xml",
			DestFile:   ".idea/runConfigurations/${app.name}_populate__local_.xml",
			Format:     text.FileFormatXml,
		},
		{
			Name:       "Creating run configuration: (remote)",
			SourceFile: "idea/runConfigurations/project__remote_.xml",
			DestFile:   ".idea/runConfigurations/${app.name}__remote_.xml",
			Format:     text.FileFormatXml,
		},
		{
			Name:       "Creating run configuration: migrate (remote)",
			SourceFile: "idea/runConfigurations/project_migrate__remote_.xml",
			DestFile:   ".idea/runConfigurations/${app.name}_migrate__remote_.xml",
			Format:     text.FileFormatXml,
		},
		{
			Name:       "Creating run configuration: populate (remote)",
			SourceFile: "idea/runConfigurations/project_populate__remote_.xml",
			DestFile:   ".idea/runConfigurations/${app.name}_populate__remote_.xml",
			Format:     text.FileFormatXml,
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
			Format:     text.FileFormatJson,
		},
		{
			Name:       "Creating task configurations",
			SourceFile: "vscode/tasks.json",
			DestFile:   ".vscode/tasks.json",
			Format:     text.FileFormatJson,
		},
	}

	return templates.Render(NewRenderOptions())
}

func GenerateDockerfile(_ []string) error {
	logger.Info("Generating Dockerfile")
	templates := TemplateSet{
		{
			Name:       "Creating Release Dockerfile",
			SourceFile: "build/package/Dockerfile",
			Format:     text.FileFormatDocker,
		},
		{
			Name:       "Creating Debug Dockerfile",
			SourceFile: "build/package/Dockerfile.debug",
			Format:     text.FileFormatDocker,
		},
	}
	if skeletonConfig.Archetype != archetypeKeyBeat {
		templates = append(templates, Template{
			Name:       "Creating docker entrypoint",
			SourceFile: "build/package/docker-entrypoint.sh",
			Format:     text.FileFormatBash,
		})
	}

	return templates.Render(NewRenderOptions())
}

func GenerateKubernetes(_ []string) error {
	logger.Infof("Generating kubernetes %s manifest templates", skeletonConfig.Archetype)
	var templates TemplateSet
	if skeletonConfig.Archetype == archetypeKeyApp || skeletonConfig.Archetype == archetypeKeyServicePack {
		templates = TemplateSet{
			{
				Name:       "Creating deployment template",
				SourceFile: "deployments/kubernetes-rc.yml.tpl",
				DestFile:   "deployments/kubernetes/${app.name}-rc.yml.tpl",
				Format:     text.FileFormatYaml,
			},
			{
				Name:       "Creating migrate template",
				SourceFile: "deployments/kubernetes-pod.yml.tpl",
				DestFile:   "deployments/kubernetes/${app.name}-pod.yml.tpl",
				Format:     text.FileFormatYaml,
			},
			{
				Name:       "Creating populate template",
				SourceFile: "deployments/kubernetes-meta.yml.tpl",
				DestFile:   "deployments/kubernetes/${app.name}-meta.yml.tpl",
				Format:     text.FileFormatYaml,
			},
			{
				Name:       "Creating pdb template",
				SourceFile: "deployments/kubernetes-pdb.yml.tpl",
				DestFile:   "deployments/kubernetes/${app.name}-pdb.yml.tpl",
				Format:     text.FileFormatYaml,
			},
			{
				Name:       "Creating Skaffold file",
				SourceFile: "deployments/kubernetes/skaffold.yaml.tpl",
				DestFile:   "skaffold.yaml",
				Format:     text.FileFormatYaml,
			},
			{
				Name:       "Creating K8S minivms deployment manifest for skaffold",
				SourceFile: "deployments/kubernetes/minivms/minivms-deployment.yaml.tpl",
				DestFile:   "deployments/kubernetes/minivms/${app.name}-deployment.yaml",
				Format:     text.FileFormatYaml,
			},
			{
				Name:       "Creating K8S msxlite deployment manifest for skaffold",
				SourceFile: "deployments/kubernetes/msxlite/msxlite-deployment.yaml.tpl",
				DestFile:   "deployments/kubernetes/msxlite/${app.name}-deployment.yaml",
				Format:     text.FileFormatYaml,
			},
		}
	} else if skeletonConfig.Archetype == archetypeKeyBeat {
		templates = TemplateSet{
			{
				Name:       "Creating deployment template",
				SourceFile: "deployments/beats/kubernetes-ps.yml.tpl",
				DestFile:   "deployments/kubernetes/${app.name}-ps.yml.tpl",
				Format:     text.FileFormatYaml,
			},
			{
				Name:       "Creating config map template",
				SourceFile: "deployments/beats/kubernetes-cm.yml.tpl",
				DestFile:   "deployments/kubernetes/${app.name}-cm.yml.tpl",
				Format:     text.FileFormatYaml,
			},
			{
				Name:       "Creating pdb template",
				SourceFile: "deployments/beats/kubernetes-pdb.yml.tpl",
				DestFile:   "deployments/kubernetes/${app.name}-pdb.yml.tpl",
				Format:     text.FileFormatYaml,
			},
		}

	} else {
		err := fmt.Errorf("kubernetes %s: %w", skeletonConfig.Archetype, ErrNoTemplates)
		logger.Errorf("%s", err)
		return err
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
			Format:     text.FileFormatYaml,
		},
	}

	return templates.Render(NewRenderOptions())
}

func GenerateHarness(_ []string) error {
	logger.Info("Generating harness")
	templates := TemplateSet{
		{
			Name:       "Creating harness manifest",
			SourceFile: "harness/service.yaml.tpl",
			DestFile:   "deployments/harness/service.yaml",
			Format:     text.FileFormatYaml,
		},
		{
			Name:       "Creating harness helm values",
			SourceFile: "harness/helm-chart-values.yaml.tpl",
			DestFile:   "deployments/harness/values_${app.name}.yaml",
			Format:     text.FileFormatYaml,
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
		Format:     text.FileFormatDocker,
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
				pipe.Exec("git", "init", "--initial-branch="+skeletonConfig.Trunk))),
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

func GenerateGithub(_ []string) error {
	logger.Info("Generating github repository configuration files")
	templates := TemplateSet{
		{
			Name:       "Creating Pull Request template",
			SourceFile: "github/pr-template.md",
			DestFile:   ".github/PULL_REQUEST_TEMPLATE.md",
			Format:     text.FileFormatMarkdown,
		},
	}

	return templates.Render(NewRenderOptions())
}

func GoGenerate(targetDirectory string) error {
	logger.Info("Executing go generate in " + targetDirectory)

	pipes := []pipe.Pipe{
		exec.WithDir(targetDirectory,
			pipe.Line(
				exec.Info("- Generating"),
				pipe.Exec("go", "generate")),
		),
	}

	return exec.ExecutePipes(pipes...)
}

func TempDir() string {
	ws := os.Getenv("WORKSPACE")
	if ws != "" {
		return ws
	}

	return os.TempDir()
}

func GenerateSPUI(_ []string) error {
	logger.Info("Generating service pack UI")

	projName := skeletonConfig.AppName + "-ui"

	targetDirectory := filepath.Join(skeletonConfig.TargetParent, projName)

	skeletonConfig.TargetDir = targetDirectory
	logger.Infof("Target Directory: %s", targetDirectory)
	err := exec.ExecutePipes(pipe.MkDirAll(targetDirectory, 0755))
	if err != nil {
		logger.Warnf("failed to create target dir: %s", targetDirectory)
		return err
	}

	skeletonConfig.AppUUID = uuid.NewString()

	generatorDir := filepath.Join(TempDir(), uuid.NewString()) // the npm generator app loaded here
	logger.Infof("Generator Directory: %s", generatorDir)

	// the create-project script rimrafs its target dir, :# ,
	// so we need to generate elsewhere and then copy in to place
	tmpTargetDir := filepath.Join(TempDir(), uuid.NewString())
	logger.Infof("Temp target Directory: %s", tmpTargetDir)

	err = exec.ExecutePipes(
		exec.WithDir(".",
			pipe.Line(
				exec.Info("- Cloning Angular Template to "+generatorDir),
				pipe.Exec("git", "clone",
					"https://github.com/CiscoDevNet/angular9-msx-service-pack-ui-generator",
					generatorDir))),
	)
	if err != nil {
		logger.Warn("failed to clone template Angular 9 Tenant Centric Service Pack Sample")
		return err
	}

	err = exec.ExecutePipes(
		exec.WithDir(generatorDir,
			pipe.Line(
				exec.Info("- Creating Angular Project"),
				pipe.Exec("npm", "run", "create-project", "--",
					"-project-name="+projName,
					"-project-description=\""+skeletonConfig.AppDescription+"\"",
					"-project-uuid="+skeletonConfig.AppUUID,
					"-output-dir="+tmpTargetDir))), // skeletonConfig.TargetParent
	)
	if err != nil {
		logger.Warn("npm failed to create project")
		return err
	}

	err = exec.ExecutePipes(
		exec.WithDir(tmpTargetDir,
			pipe.Line(
				exec.Info("- Copying generated project from %s to %s", tmpTargetDir, skeletonConfig.TargetParent),
				pipe.Exec("cp", "-Rv", projName, skeletonConfig.TargetParent),
				pipe.Write(os.Stdout))),
	)
	if err != nil {
		logger.Warn("failed to copy generated project")
		return err
	}

	templates := SPUITemplates(path.Join("spui", "patch"), "")
	if err := templates.RenderTo(targetDirectory, NewRenderOptions()); err != nil {
		return err
	}

	return nil
}

func SPUITemplates(srcroot, dstroot string) TemplateSet {
	return TemplateSet{
		{Name: "Overlaying license",
			SourceFile: path.Join(srcroot, "LICENSE.md"),
			DestFile:   path.Join(dstroot, "LICENSE.md"),
			Format:     text.FileFormatMarkdown},
		{Name: "Overlaying jenkins file",
			SourceFile: path.Join(srcroot, "becomesbin/ci/Jenkinsfile"),
			DestFile:   path.Join(dstroot, "bin/ci/Jenkinsfile"),
			Format:     text.FileFormatJenkins},
		{Name: "Overlaying sonar properties",
			SourceFile: path.Join(srcroot, "becomesbin/ci/sonar-project.properties"),
			DestFile:   path.Join(dstroot, "bin/ci/sonar-project.properties"),
			Format:     text.FileFormatProperties},
		{Name: "Overlaying conformance script",
			SourceFile: path.Join(srcroot, "becomesbin/conformance.sh"),
			DestFile:   path.Join(dstroot, "bin/conformance.sh"),
			Format:     text.FileFormatBash},
		{Name: "Overlaying docker build script",
			SourceFile: path.Join(srcroot, "becomesbin/docker-build.sh"),
			DestFile:   path.Join(dstroot, "bin/docker-build.sh"),
			Format:     text.FileFormatBash},
		{Name: "Overlaying docker clean script",
			SourceFile: path.Join(srcroot, "becomesbin/docker-clean.sh"),
			DestFile:   path.Join(dstroot, "bin/docker-clean.sh"),
			Format:     text.FileFormatBash},
		{Name: "Overlaying docker push script",
			SourceFile: path.Join(srcroot, "becomesbin/docker-push.sh"),
			DestFile:   path.Join(dstroot, "bin/docker-push.sh"),
			Format:     text.FileFormatBash},
		{Name: "Overlaying package script",
			SourceFile: path.Join(srcroot, "becomesbin/package.sh"),
			DestFile:   path.Join(dstroot, "bin/package.sh"),
			Format:     text.FileFormatBash},
		{Name: "Overlaying dockerfile",
			SourceFile: path.Join(srcroot, "becomesbin/package/Dockerfile"),
			DestFile:   path.Join(dstroot, "bin/package/Dockerfile"),
			Format:     text.FileFormatDocker},
		{Name: "Overlaying publish script",
			SourceFile: path.Join(srcroot, "becomesbin/publish.sh"),
			DestFile:   path.Join(dstroot, "bin/publish.sh"),
			Format:     text.FileFormatBash},
		{Name: "Overlaying vars script",
			SourceFile: path.Join(srcroot, "becomesbin/vars.sh"),
			DestFile:   path.Join(dstroot, "bin/vars.sh"),
			Format:     text.FileFormatBash},
		{Name: "Overlaying jest configuration",
			SourceFile: path.Join(srcroot, "jest.config.js"),
			DestFile:   path.Join(dstroot, "jest.config.js"),
			Format:     text.FileFormatJavaScript},
		{Name: "Overlaying jest init",
			SourceFile: path.Join(srcroot, "jest.init.js"),
			DestFile:   path.Join(dstroot, "jest.init.js"),
			Format:     text.FileFormatJavaScript},
		{Name: "Overlaying package lock",
			SourceFile: path.Join(srcroot, "package-lock.json"),
			DestFile:   path.Join(dstroot, "package-lock.json"),
			Format:     text.FileFormatJson},
		{Name: "Replacing package file",
			SourceFile: path.Join(srcroot, "package.json"),
			DestFile:   path.Join(dstroot, "package.json"),
			Format:     text.FileFormatJson},
		{Name: "Overlaying empty module",
			SourceFile: path.Join(srcroot, "src/spec-helpers/empty-module.js"),
			DestFile:   path.Join(dstroot, "src/spec-helpers/empty-module.js"),
			Format:     text.FileFormatJson},
		{Name: "Overlaying api client",
			SourceFile: path.Join(srcroot, "src/spec-helpers/mocks/api-client.ts"),
			DestFile:   path.Join(dstroot, "src/spec-helpers/mocks/api-client.ts"),
			Format:     text.FileFormatTypeScript},
		{Name: "Overlaying mock index",
			SourceFile: path.Join(srcroot, "src/spec-helpers/mocks/index.ts"),
			DestFile:   path.Join(dstroot, "src/spec-helpers/mocks/index.ts"),
			Format:     text.FileFormatTypeScript},
		{Name: "Overlaying mock monitor",
			SourceFile: path.Join(srcroot, "src/spec-helpers/mocks/mock-monitor.service.ts"),
			DestFile:   path.Join(dstroot, "src/spec-helpers/mocks/mock-monitor.service.ts"),
			Format:     text.FileFormatTypeScript},
		{Name: "Overlaying html transformers",
			SourceFile: path.Join(srcroot, "src/spec-helpers/transformers/html.js"),
			DestFile:   path.Join(dstroot, "src/spec-helpers/transformers/html.js"),
			Format:     text.FileFormatJavaScript},
		{Name: "Overlaying eslint ignore",
			SourceFile: path.Join(srcroot, ".eslintignore"),
			DestFile:   path.Join(dstroot, ".eslintignore"),
			Format:     text.FileFormatOther},
		{Name: "Overlaying eslint config",
			SourceFile: path.Join(srcroot, ".eslintrc.json"),
			DestFile:   path.Join(dstroot, ".eslintrc.json"),
			Format:     text.FileFormatJson},
		{Name: "Overlaying gitignore",
			SourceFile: path.Join(srcroot, ".gitignore"),
			DestFile:   path.Join(dstroot, ".gitignore"),
			Format:     text.FileFormatOther},
		{Name: "Overlaying npm run commands",
			SourceFile: path.Join(srcroot, ".npmrc"),
			DestFile:   path.Join(dstroot, ".npmrc"),
			Format:     text.FileFormatOther},
		{Name: "Overlaying checks file",
			SourceFile: path.Join(srcroot, ".checks.yml"),
			DestFile:   path.Join(dstroot, ".checks.yml"),
			Format:     text.FileFormatOther},
	}
}
