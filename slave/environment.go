package slave

import (
        "github.com/kelseyhightower/envconfig"
        "os"
        "strings"
)

// Spec represents the environment.
type Spec struct {
        // Server should be set to the Jenkins master server URL.
        Server  string

        // Home should be set to the the Jenkins working directory.
        Home    string

        // Name should hold the value of the node name as defined by the Jenkins master.
        Name    string

        // Username should be set to the user to login in to for Jenkins. It only applies if
        // swarm is true.
        Username string

        // Password should be set to the Jenkins user's login password. It only applies if
        // swarm is true.
        Password string

        // Swarm determines whether the swarm client should be used to launch the slave.
        // Valid values are true and false.
        Swarm   bool

        // Swarmversion determines the version of the swarm client to use. It only applies if
        // Swarm is true.
        Swarmversion string

        // Lock determines if a lock file should be created during execution.
        Lock    bool

        // Executors should be set to the number of executors the node should use. It only applies if
        // swarm is true.
        Executors int

        // Mode to set for the node. Valid values are 'normal' (utilize the slave as much as possible)
        // and 'exclusive' (leave this machine for tied jobs only). It only applies if
        // swarm is true.
        Mode    string

        // Labels can be a whitespace separated string of labels for the node. It only applies if
        // swarm is true.
        Labels  string
}

const (
        ENV_SERVER       = "SLAVE_SERVER"       // Environment variable
        ENV_HOME         = "SLAVE_HOME"         // Environment variable
        ENV_NAME         = "SLAVE_NAME"         // Environment variable
        ENV_USERNAME     = "SLAVE_USERNAME"     // Environment variable
        ENV_PASSWORD     = "SLAVE_PASSWORD"     // Environment variable
        ENV_SWARM        = "SLAVE_SWARM"        // Environment variable
        ENV_SWARMVERSION = "SLAVE_SWARMVERSION" // Environment variable
        ENV_LOCK         = "SLAVE_LOCK"         // Environment variable
        ENV_EXECUTORS    = "SLAVE_EXECUTORS"    // Environment variable
        ENV_MODE         = "SLAVE_MODE"         // Environment variable
        ENV_LABELS       = "SLAVE_LABELS"       // Environment variable
)

var spec *Spec = nil

// Environment will read in the environment as defined in Spec. All environment variables will be
// prefixed with "SLAVE_".
func (s *Spec) Environment() (*Spec, error) {
        err := envconfig.Process("slave", s)

        if !s.ValidMode() {
                logger.Printf(
                        "Invalid value for mode [%s]. Valid values are normal and exclusive. Setting to normal.",
                        s.Mode,
                )
                s.Mode = "normal"
        }

        return s, err
}

// File reads the spec from the file listed as the argument.
func (s *Spec) File(name string) (*Spec, error) {
        return nil, nil
}

// ValidMode returns true if the mode is 'normal' or 'exclusive'.
func (s *Spec) ValidMode() (res bool) {
        res = false
        if m := strings.ToLower(s.Mode); m == "normal" || m == "exclusive" {
                res = true
        }
        return
}

// New creates a new instance of the environment spec. The following defaults will be used:
//                host, _ := os.Hostname()
//                dir, _ := os.Getwd()
//                spec = &Spec{
//                        Name:         host,
//                        Home:         dir,
//                        Swarm:        false,
//                        Lock:         false,
//                        Swarmversion: "1.15",
//                        Executors:    2,
//                        Mode:         "normal",
//                }
func NewSpec() *Spec {
        if spec == nil {
                // Set some defaults
                host, _ := os.Hostname()
                dir, _ := os.Getwd()
                spec = &Spec{
                        Name:         host,
                        Home:         dir,
                        Swarm:        false,
                        Lock:         false,
                        Swarmversion: "1.15",
                        Executors:    2,
                        Mode:         "normal",
                }
        }
        return spec
}
