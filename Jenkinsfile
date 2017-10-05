#!groovy

node() {
  try {

    withCredentials([file(credentialsId: 'npm_nexus_auth', variable: 'NEXUS_AUTH')]) {
      withEnv(['NEXUS_EMAIL=engineering@linuxfoundation.org']) {
        docker.withRegistry("https://433610389961.dkr.ecr.us-west-2.amazonaws.com", "ecr:us-west-2:jenkins-aws-user") {
          docker.image('433610389961.dkr.ecr.us-west-2.amazonaws.com/npm-publish:latest').inside {

            stage('Checkout') {
              checkout scm
            }

            stage('Build') {
              sh "npm install"
              sh "npm run build"
            }

            stage('Publish NPM snapshot') {
              def currentVersion = sh(returnStdout: true, script: "npm version | grep \"{\" | tr -s ':'  | cut -d \"'\" -f 4").trim()
              def newVersion = "${currentVersion}-${buildNumber}"
              sh "npm version ${newVersion} --no-git-tag-version && npm publish --tag next"
            }

          }
        }
      }
    }

    currentBuild.result = 'SUCCESS'

  } catch (error) {
    currentBuild.result = 'FAILED' // Important for the Stash notifier!
    throw error
  }
}