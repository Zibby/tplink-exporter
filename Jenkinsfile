pipeline {
  environment {
    registry = "zibby/tplink-exporter"
    registryCredential = 'f8a79f84-5ad0-43e4-b32c-87e2c6001a62'
    dockerImage = ''
  }
  agent any
  stages {
    stage('Clone Git') {
      steps {
        git 'https://github.com/Zibby/tplink-exporter'
      }
    }
    stage('Build Image') {
      steps {
        script {
          dockerImage = docker.build registry + "_" + "$BUILD_NUMBER"
        }
      }
    }
    stage('Push Image') {
      steps {
        script {
          docker.withRegistry( '', registryCredential) {
            dockerImage.push()
          }
        }
      }
    }
    stage('clean_up') {
      steps {
        sh "docker rmi $registry:$BUILD_NUMBER"
      }
    }
  }
}