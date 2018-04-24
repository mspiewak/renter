pipeline { 
    agent any 
    stages {
        stage('Build') { 
            steps { 
                sh 'docker build -t mspiewak/renter' 
            }
        }
        stage('Deploy') {
            steps {
                sh 'echo hello'
            }
        }
    }
}