def buildStatus = 0

node {
  properties([disableConcurrentBuilds()])

  // Wipe the workspace so we are building completely clean
  sh "sudo rm -rf *"

  stage ("Checkout") {
    git pool: true, credentialsId: 'd78c94c4-9179-4765-9851-9907b5ef2cc4', url: "git@github.linuxfoundation.org:Engineering/project-management-console.git", branch: "${env.BRANCH_NAME}"
  }

  def project = "pmc"
  def gitCommit = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()
  def shortCommit = gitCommit.take(7)
  def instancePath = "${project}_dev_${env.BUILD_NUMBER}_${shortCommit}"
  def gitAuthor = sh(returnStdout: true, script: 'git log -1 --pretty=format:"%an"').trim()
  def gitCommitMsg = sh(returnStdout: true, script: 'git log -1 --pretty=format:"%s"').trim()

  try {

    stage ("Launching CINCO Instance") {
      sh "lf i create --project=cinco --branch=develop --name='${instancePath}_cinco' -y -d"
    }

    stage ("Launching PMC Instance") {
      sh "lf i create --project=pmc --branch=${env.BRANCH_NAME} --name='${instancePath}' --dep-map='cinco:${instancePath}_cinco' --sequence=jenkins -y -d --no-autorun"
    }

    def workspaceID = sh (script: "cd ${instancePath}; lf i workspace", returnStdout: true).trim()

    stage ("NPM Installation") {
      dir ("${instancePath}") {
        sh "docker exec ${workspaceID} bash -c \"cd /srv/app/src && npm install\""
      }
    }

    stage ("Ionic Installation") {
      dir ("${instancePath}") {
        sh "docker exec ${workspaceID} bash -c \"cd /srv/app/src && npm run build\""
      }
    }

//    stage ("Automated Tests") {
//      dir ("${instancePath}") {
//        try {
//          sh "docker exec ${workspaceID} bash -c \"cd /srv/app && npm run tests\""
//        } finally {
//          step([$class: "JUnitResultArchiver", testResults: "test-results/*.xml"])
//        }
//      }
//    }

    stage("Destroying Instances") {
      sh "lf i rm ${instancePath}/ -y"
      sh "lf i rm ${instancePath}_cinco/ -y"
    }

    if (env.BRANCH_NAME == 'develop') {
      build job: 'PMC - Sandbox', parameters: [string(name: 'SHA', value: "${shortCommit}")], wait: false
    } else if (env.BRANCH_NAME == 'master') {
      build job: 'PMC - Production', parameters: [string(name: 'SHA', value: "${shortCommit}")], wait: false
    }

  } catch(err) {
    buildStatus = 1

    // Making sure we always destroy the instance after each build
    stage("Destroying Instances") {
      sh "lf i rm ${instancePath}/ -y"
      sh "lf i rm ${instancePath}_cinco/ -y"
    }

    throw err

  } finally {
    sh "curl -s https://workflow.engineering.tux.rocks/api/jenkins/github-build/slack -d '{ \
        \"build\": \"${env.BUILD_ID}\", \
        \"build_url\": \"${env.BUILD_URL}\", \
        \"gitAuthor\": \"${gitAuthor}\", \
        \"gitCommit\": \"${shortCommit}\", \
        \"gitCommitMsg\": \"${gitCommitMsg}\", \
        \"job_name\": \"${env.JOB_NAME}\", \
        \"job_url\": \"${env.JOB_URL}\", \
        \"status\": \"${buildStatus}\", \
        \"gitBranch\": \"${env.BRANCH_NAME}\", \
        \"job_base_name\": \"${env.JOB_BASE_NAME}\"}'"
  }

}

