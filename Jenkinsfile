def REPO_NAME = 'go-msx'
def ARTIFACTORY_CREDENTIALS = 'f79d92f8-694b-4c29-b477-faeabcef86cb'
def VMSBLD_CREDENTIALS = 'msx-jenkins-gen-ssh-key'
def GITHUB_CREDENTIALS = 'msx-jenkins-gen-token-secret-text'
def GITHUB_APP_ID_SONARQUBE = '7'
def GITHUB_APP_NAME_SONARQUBE = 'engit-sonar-int-gen-GPK'
def SONARQUBE_CREDENTIALS = 'SONARQUBE_GPK_ACCESS_TOKEN'
def SONARQUBE_INSTALLATION = 'GPK SonarQube'
def TRUNK = 'master'

pipeline {

    agent {
        label env.SLAVE_LABEL
    }

    /* Only keep the 10 most recent builds. */
    options {
        buildDiscarder(logRotator(numToKeepStr: '10', artifactNumToKeepStr: '10'))
    }

    stages {

        stage('Preparation') {
            steps {
                deleteDir()

                script {
                    assert env.BUILD_NUMBER
                    assert env.WORKSPACE

                    if (env.BRANCH_NAME) {
                        currentBuild.description = env.BRANCH_NAME + " (Manual)"
                    } else if (env.sha1) {
                        // When called from GitHub PR Builder
                        env.BRANCH_NAME = env.sha1
                    } else if (env.GIT_BRANCH) {
                        // When called from GitHub Push Notifier
                        env.BRANCH_NAME = env.GIT_BRANCH.replaceAll("origin/", "")
                        currentBuild.description = env.BRANCH_NAME + " (Push)"
                    }
                }
            }
        }

        stage('Checkout') {
            steps {

                checkout([
                    $class                           : 'GitSCM',
                    branches                         : [[name: env.BRANCH_NAME ?: TRUNK]],
                    doGenerateSubmoduleConfigurations: false,
                    extensions                       : [[
                        $class: 'RelativeTargetDirectory',
                        relativeTargetDir: REPO_NAME,
                    ]],
                    userRemoteConfigs                : [[
                        credentialsId: VMSBLD_CREDENTIALS,
                        url          : "git@cto-github.cisco.com:NFV-BU/${REPO_NAME}.git",
                        refspec:       '+refs/pull/*:refs/remotes/origin/pr/* +refs/heads/*:refs/remotes/origin/*'
                    ]]
                ])

            }
        }

        stage('Perform Build') {
            steps {
                sshagent([VMSBLD_CREDENTIALS]) {
                    withEnv([
                        "GOPATH=${env.WORKSPACE}/go",
                        "GOPRIVATE=cto-github.cisco.com/NFV-BU",
                        "GOPROXY=https://engci-maven.cisco.com/artifactory/go/,https://proxy.golang.org,direct",
                        "PATH+GOBIN=${env.WORKSPACE}/go/bin",
                        "WORKSPACE=$WORKSPACE/$REPO_NAME"
                    ]) { dir ("$WORKSPACE") {
                        sh 'git config --global url."git@cto-github.cisco.com:".insteadOf "https://cto-github.cisco.com/"'
                        sh 'make test'
                    }}
                }
            }
        }

        stage('Static Analysis') {
            steps {
                withEnv(["WORKSPACE=$WORKSPACE/$REPO_NAME"]) { dir("$WORKSPACE") {

                    junit 'test/junit-report.xml'

                    publishCoverage adapters: [coberturaAdapter('test/cobertura-coverage.xml')], sourceFileResolver: sourceFiles('NEVER_STORE')


                    withCredentials([string(credentialsId: GITHUB_CREDENTIALS, variable: 'GITHUB_ACCESS_TOKEN')]) {
                        withSonarQubeEnv(credentialsId: SONARQUBE_CREDENTIALS, installationName: SONARQUBE_INSTALLATION) {
                            script {
                                def sonarProperties = [
                                    'userHome': env.WORKSPACE,
                                    'links.ci': env.JOB_URL,
                                    'github.oauth': env.GITHUB_ACCESS_TOKEN,
                                    'sonar.alm.github.app.id': GITHUB_APP_ID_SONARQUBE,
                                    'sonar.alm.github.app.name': GITHUB_APP_NAME_SONARQUBE,
                                ]

                                if (env.BRANCH_NAME != TRUNK && env.BRANCH_NAME != 'master') {
                                    if (env.ghprbTargetBranch) {
                                        sonarProperties['pullrequest.github.repository'] = "NFV-BU/${REPO_NAME}"
                                        sonarProperties['pullrequest.provider'] = "github"
                                        sonarProperties['pullrequest.key'] = env.ghprbPullId
                                        sonarProperties['pullrequest.branch'] = env.ghprbSourceBranch
                                        sonarProperties['pullrequest.base'] = env.ghprbTargetBranch
                                    } else {
                                        sonarProperties['branch.name'] = env.BRANCH_NAME
                                        sonarProperties['branch.target'] = TRUNK
                                    }
                                }

                                def sonarHome = tool name: 'sonarscaner', type: 'hudson.plugins.sonar.SonarRunnerInstallation'
                                def sonarCommand = "$sonarHome/bin/sonar-scanner"
                                sonarProperties.each { key, value -> sonarCommand = sonarCommand + " -Dsonar.$key=$value" }
                                sonarCommand = sonarCommand + " -Dproject.settings=sonar-project.properties"
                                sh sonarCommand
                            }
                        }
                    }

                }}
            }
        }

        stage('Skel') {
            steps {
                sshagent([VMSBLD_CREDENTIALS]) {
                    withEnv([
                        "GOPATH=${env.WORKSPACE}/go",
                        "GOPRIVATE=cto-github.cisco.com/NFV-BU",
                        "GOPROXY=https://engci-maven.cisco.com/artifactory/go/,https://proxy.golang.org,direct",
                        "PATH+GOBIN=${env.WORKSPACE}/go/bin",
                        "WORKSPACE=$WORKSPACE/$REPO_NAME"]) { dir("$WORKSPACE") { withCredentials([usernamePassword(
                    credentialsId: ARTIFACTORY_CREDENTIALS,
                    passwordVariable: 'ARTIFACTORY_PASSWORD',
                    usernameVariable: 'ARTIFACTORY_USERNAME')]) {
                    script {
                        if (env.BRANCH_NAME == TRUNK) {
                            sh "make skel"
                            sh "make publish-skel"
                        }
                    }
                }}}}
            }
        }

    }
}
