// Declarative Jenkins pipeline for building and testing the promoter
@Library('jenkins-pipelines') _

def errors = []
def registry = "${env.DOCKER_REGISTRY}"
def imageName
def imageId
def newVersion
def released = []

pipeline {
    agent any
    options{
      // Whole pipeline timeout
      timeout(time: 15, unit: 'MINUTES')

      // Build properties
      disableConcurrentBuilds()

      // delete old artefacts to avoid filling the disk
      buildDiscarder(logRotator(numToKeepStr: '25', artifactNumToKeepStr: '10'))

      gitLabConnection('gitlab')
    }

    stages {
        stage('Notify GitLab') {
            steps {
                script {
                    notifications.tellBuildRunning()
                    goUtils.createNetrc()
                    git.checkBranchName()
                }
            }
        }

        stage('Build') {
            steps {
                script {
                    imageName = registry + "spa-server:${BRANCH_NAME}"
                    image = goUtils.buildGoImage(imageName, "Dockerfile", "", true)
                    imageId = image.id
                }
            }
            post {
                failure {
                    script {
                        errors << 'Failed to build the image'
                    }
                }
            }
        }
        stage("Test") {
            parallel {
                stage('Static Analysis') {
                    steps {
                        sh './run_static_analysis.sh'
                    }

                    post {
                        failure {
                            script {
                                errors << 'Static Analysis Failed'
                            }
                        }
                    }
                }
            }
        }
        stage('Tag & Deploy') {
            when {
                branch 'master'
                expression {
                    git.version().hasChanges && currentBuild.resultIsBetterOrEqualTo('SUCCESS')
                }
            }
            steps {
                script {
                    newVersion = git.updateVersion()

                    image.push('latest')
                    image.push(newVersion)
                    released << newVersion
                    log.good("Released image ${imageName} with version ${newVersion} and latest")
                }
            }
        }
    }
    post {
        always {
            script {
                    notifications.tellBuildResult(errors, released)
            }
            sh 'scripts/cleanup_modules.sh'
            sh "rm .netrc || true"
        }
    }
}
