node {
  properties([pipelineTriggers([pollSCM('* * * * *')]),disableConcurrentBuilds()])

  // Wipe the workspace so we are building completely clean
  sh "sudo rm -rf *"
  sh "docker network prune -f"
  sh "docker container prune -f"
  sh "docker image prune -f"

  stage ("Checkout") {
    git pool: true, url: "git@github.linuxfoundation.org:Engineering/project-management-console.git", branch: "${env.BRANCH_NAME}"
  }

  def project = "cinco"
  def gitCommit = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()
  def shortCommit = gitCommit.take(6)
  def instancePath = "${project}_dev_${env.BUILD_NUMBER}_${shortCommit}"

  try {

    stage ("Launching CINCO Instance") {
      sh "lf i create --project=cinco --branch=develop --name='${instancePath}_cinco' -y -d"
    }

    stage ("Launching Project Management Console") {
      sh "lf i create --project=pmc --branch=develop --name='${instancePath}' --dep-map=cinco:${instancePath}_cinco --sequence=jenkins -y -d --no-autorun"
    }

    def workspaceID = sh (script: "cd ${instancePath}; lf i workspace", returnStdout: true).trim()

    stage ("NPM Install") {
      dir ("${instancePath}") {
        sh "docker exec ${workspaceID} bash -c \"cd /srv/app/src && npm install\""
      }
    }

    stage ("Ionic Install") {
      dir ("${instancePath}") {
        sh "docker exec ${workspaceID} bash -c \"cd /srv/app/src && npm run build\""
      }
    }

    stage("Destroying Instances") {
      sh "lf i rm ${instancePath}_cinco/ -y"
      sh "lf i rm ${instancePath}/ -y"
    }

  } catch(err) {
    // Making sure we always destroy the instance after each build
    stage("Destroying Instance") {
      sh "lf i rm ${instancePath}_cinco/ -y"
      sh "lf i rm ${instancePath}/ -y"
    }

    throw err

  }

}

