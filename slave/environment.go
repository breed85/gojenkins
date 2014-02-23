package slave

import (
        "github.com/kelseyhightower/envconfig"
)

// Spec represents the environment.
type Spec struct {
        // Jenkinsserver holds the value of env variable SLAVE_JENKINSSERVER.
        // It should be set to the Jenkins master server URL.
        Jenkinsserver string

        // Jenkinscwd holds the value of env variable SLAVE_JENKINSCWD.
        // It should be set to the path where the Jenkins working directory.
        Jenkinscwd string

        // Name holds the value of env variable SLAVE_NAME.
        // It should hold the value of the node name as defined by the Jenkins master.
        Name    string
}

const (
        ENV_JENKINSSERVER = "SLAVE_JENKINSSERVER" // Environment variable
        ENV_JENKINSCWD    = "SLAVE_JENKINSCWD"    // Environment variable
        ENV_NAME          = "SLAVE_NAME"          // Environment variable
)

var spec Spec

// Environment will read in the environment as defined in Spec. All environment variables will be
// prefixed with "SLAVE_".
func Environment() error {
        return envconfig.Process("slave", &spec)
}

// GetJenkinsServer returns the value of Jenkinsserver
func GetJenkinsServer() string {
        return spec.Jenkinsserver
}

// GetJenkinsCwd returns the value of Jenkinscwd
func GetJenkinsCwd() string {
        return spec.Jenkinscwd
}

// GetName returns the value of Name
func GetName() string {
        return spec.Name
}
