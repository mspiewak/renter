pipeline { 
    agent any 
    stages {
        stage('Build') { 
            steps { 
                docker build -t mspiewak/renter
            }
        }
        stage('Deploy') {
            steps {
                sh 'echo hello'
            }
        }
    }
}