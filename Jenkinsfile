def REPO_NAME = 'go-msx'
def VMSBLD_CREDENTIALS = '06702990-d4e7-4a80-beaa-6d407896feff'

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
                    assert env.BRANCH_NAME
                    assert env.BUILD_NUMBER
                    assert env.WORKSPACE

                    if (env.sha1) {
                        // When called from GitHub PR Builder
                        env.BRANCH_NAME = env.sha1
                    } else if (env.GIT_BRANCH) {
                        // When called from GitHub Push Notifier
                        env.BRANCH_NAME = env.GIT_BRANCH.replaceAll("origin/", "")
                        currentBuild.description = env.BRANCH_NAME + " (Push)"
                    } else {
                        currentBuild.description = env.BRANCH_NAME + " (Manual)"
                    }
                }
            }
        }

        stage('Checkout') {
            steps {

                checkout([
                    $class                           : 'GitSCM',
                    branches                         : [[name: env.BRANCH_NAME]],
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

                }}
            }
        }

    }
}
