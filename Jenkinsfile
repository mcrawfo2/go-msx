def REPO_NAME = 'go-msx'
def ARTIFACTORY_CREDENTIALS = 'f79d92f8-694b-4c29-b477-faeabcef86cb'
def GIT_CREDENTIALS = 'msx-jenkins-gen-ssh-key'
def GITHUB_CREDENTIALS = 'msx-jenkins-gen-token-secret-text'
def GITHUB_APP_ID_SONARQUBE = '7'
def GITHUB_APP_NAME_SONARQUBE = 'engit-sonar-int-gen-GPK'
def SONARQUBE_CREDENTIALS = 'SONARQUBE_GPK_ACCESS_TOKEN'
def SONARQUBE_INSTALLATION = 'GPK SonarQube'
def SLACK_BOT_CREDENTIALS = 'slack-bot-secret'
def SLACK_CHANNEL = 'go-msx-build'
def SOURCEGRAPH_ENDPOINT = 'https://sourcegraph.infra.ciscomsx.com'
def SOURCEGRAPH_CREDENTIALS = 'sourcegraph-lsif-token'
def TRIGGER = ''
def TRIGGER_URL = ''
def TRIGGER_SUFFIX = ''
def TRUNK = 'master'
def TOOLS = ''
def TEST_SUMMARY = null

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
                        TRIGGER = currentBuild.description
                        TRIGGER_URL = "https://cto-github.cisco.com/NFV-BU/${REPO_NAME}/tree/" + env.BRANCH_NAME
                    } else if (env.sha1) {
                        // When called from GitHub PR Builder
                        env.BRANCH_NAME = env.sha1
                        TRIGGER = "PR #" + ghprbPullId
                        TRIGGER_SUFFIX = ": " + env.ghprbPullTitle
                        TRIGGER_URL = env.ghprbPullLink
                    } else if (env.GIT_BRANCH) {
                        // When called from GitHub Push Notifier
                        env.BRANCH_NAME = env.GIT_BRANCH.replaceAll("origin/", "")
                        currentBuild.description = env.BRANCH_NAME + " (Push)"
                        TRIGGER = currentBuild.description
                        TRIGGER_URL = "https://cto-github.cisco.com/NFV-BU/${REPO_NAME}/tree/" + env.BRANCH_NAME
                    }

                    TOOLS = env.WORKSPACE + "/bin"
                    sh "mkdir -p ${TOOLS}"
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
                        credentialsId: GIT_CREDENTIALS,
                        url          : "git@cto-github.cisco.com:NFV-BU/${REPO_NAME}.git",
                        refspec:       '+refs/pull/*:refs/remotes/origin/pr/* +refs/heads/*:refs/remotes/origin/*'
                    ]]
                ])

            }
        }

        stage('Perform Build') {
            steps {
                sshagent([GIT_CREDENTIALS]) {
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

                    script {
                        TEST_SUMMARY = junit("${REPO_NAME}/test/junit-report.xml")
                    }

                }
            }
        }

        stage('Index') {
            steps {
                sshagent([GIT_CREDENTIALS]) {
                    withEnv([
                        "WORKSPACE=$WORKSPACE/$REPO_NAME",
                        "SRC_ENDPOINT=${SOURCEGRAPH_ENDPOINT}"
                    ]) { dir("$WORKSPACE") {
                    withCredentials([string(credentialsId: GITHUB_CREDENTIALS, variable: 'GITHUB_ACCESS_TOKEN')]) {
                    withCredentials([string(credentialsId: SOURCEGRAPH_CREDENTIALS, variable: 'SRC_ACCESS_TOKEN')]) {
                    script {
                        if (env.BRANCH_NAME == TRUNK || env.BRANCH_NAME == 'master') {
                            // Download sourcegraph go indexer
                            sh "curl -L https://github.com/sourcegraph/lsif-go/releases/download/v1.2.0/src_linux_amd64 -o ${TOOLS}/lsif-go"
                            sh "chmod +x ${TOOLS}/lsif-go"

                            // Download sourcegraph cli
                            sh "curl -L ${SOURCEGRAPH_ENDPOINT}/.api/src-cli/src_linux_amd64 -o ${TOOLS}/src"
                            sh "chmod +x ${TOOLS}/src"

                                // Index go
                            sh "${TOOLS}/lsif-go"
                            sh "${TOOLS}/src lsif upload -github-token=${GITHUB_ACCESS_TOKEN} -ignore-upload-failure"
                        }
                    }}}}}
                }
            }
        }

        stage('Static Analysis') {
            steps {
                withEnv(["WORKSPACE=$WORKSPACE/$REPO_NAME"]) { dir("$WORKSPACE") {
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
                sshagent([GIT_CREDENTIALS]) {
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

    post {
        always {
            script {
                def RESULT_EMOJIS = [
                    "SUCCESS": ":white_check_mark:",
                    "UNSTABLE": ":jenkins_unstable:",
                    "FAILURE": ":x:",
                    "ABORTED": ":abort:"
                ]

                if (SLACK_CHANNEL != "") {
                    def RESULT_EMOJI = RESULT_EMOJIS[currentBuild.currentResult]
                    def startTimeString = new Date(currentBuild.startTimeInMillis).format("MMM dd HH:mm")
                    def durationString = currentBuild.durationString.replace(" and counting", "")

                    def section = { label, text -> ["type": "section", "fields": [ [ "type": "mrkdwn", "text": "*" + label + ":*\n" + text ] ] ] }
                    def section2 = { label1, text1, label2, text2 -> ["type": "section", "fields": [ [ "type": "mrkdwn", "text": "*" + label1 + ":*\n" + text1 ], [ "type": "mrkdwn", "text": "*" + label2 + ":*\n" + text2 ] ] ] }
                    def header = { text -> ["type": "header", "text": [ "type": "plain_text", "emoji": true, "text": text ] ] }

                    def blocks = [
                        header("${RESULT_EMOJI} ${REPO_NAME} #${env.BUILD_NUMBER}"),
                        section2(
                            "Job", "<${env.RUN_DISPLAY_URL}|${REPO_NAME}/${env.BUILD_NUMBER}>",
                            "Trigger", "<${TRIGGER_URL}|${TRIGGER}>${TRIGGER_SUFFIX}"),
                        section("Built", "${startTimeString} - _${durationString}_")
                    ]

                    def testSummary = "Unknown Result"
                    if (TEST_SUMMARY != null) {
                        testSummary = "_Passed_: ${TEST_SUMMARY.passCount} " +
                            "/ ${TEST_SUMMARY.totalCount}\n" +
                            "_Failed_: ${TEST_SUMMARY.failCount}, " +
                            "_Skipped_: ${TEST_SUMMARY.skipCount}"
                    }
                    blocks.add(section("Tests", testSummary))

                    slackSend(
                        channel: SLACK_CHANNEL,
                        blocks: blocks,
                        botUser: true,
                        tokenCredentialId: SLACK_BOT_CREDENTIALS)
                }
            }
        }
    }

}
