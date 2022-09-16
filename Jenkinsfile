#!/usr/bin/env groovy

 pipeline {

     agent any
     environment {
        registry = "registry.hyneo.ru/hyneo-core"
        registryCredential = "nexusadmin"
        dockerImage = ''
     }

     stages {
         stage('Build') {
             steps {
              script {
                dockerImage = docker.build registry + ":$BUILD_NUMBER"
                }
             }
         }
         stage('Push registry nexus'){
             steps{
                 script{
                     docker.withRegistry('https://registry.hyneo.ru', registryCredential ) {
                         dockerImage.push()
                         dockerImage.push('latest')
                     }
                 }
             }
        }
//         stage ('Deploy') {
//             steps{
//                 sshagent(credentials : ['launch']) {
//                 sh '''
//                     [ -d ~/.ssh ] || mkdir ~/.ssh && chmod 0700 ~/.ssh
//                     ssh-keyscan -t rsa,dsa -p 11 mc.hyneo.ru >> ~/.ssh/known_hosts
//                     ssh -p 11 suro@mc.hyneo.ru 'cd ./core && docker compose pull core && docker compose stop core && docker compose up -d core'
//                     '''
//                 }
//             }
//         }
        stage('Dangling Images') {
            steps {
                sh 'docker images -q -f dangling=true | xargs --no-run-if-empty docker rmi'
            }
        }
     }
}