package slave

import (
        "encoding/json"
        "github.com/kelseyhightower/envconfig"
        "io"
        "os"
        "reflect"
        "strings"
        "time"
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

        // File is the name of the config file to load the settings from. This can be used as an
        // alternative to environment variables or command line flags.
        File    string

        // Log is the name of the file to log output to. The file will be truncated with each execution.
        Log     string
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
        ENV_FILE         = "SLAVE_FILE"         // Environment variable
        ENV_LOG          = "SLAVE_LOG"          // Environment variable
)

var spec *Spec = nil

// Environment will read in the environment as defined in Spec. All environment variables will be
// prefixed with "SLAVE_".
func (s *Spec) Environment() (*Spec, error) {
        if err := envconfig.Process("slave", s); err != nil {
                return nil, err
        }

        if !s.ValidMode() {
                // If the mode is not expected, default it to normal.
                s.Mode = "normal"
        }

        return s, nil
}

// Json reads the spec from the provided io.Reader and returns true
// if the underlying spec was updated.
func (s *Spec) Json(r io.Reader) (changed bool, err error) {
        changed = false
        // Decode the JSON.
        dec := json.NewDecoder(r)
        tmpSpec := s.copy()
        err = dec.Decode(tmpSpec)
        if err != nil {
                return
        }

        // Use reflection to compare the original spec to the new one. If data
        // has changed for a field, update the original spec.
        tmpValue := reflect.ValueOf(tmpSpec).Elem()
        specValue := reflect.ValueOf(s).Elem()
        specType := specValue.Type()

        for i := 0; i < specValue.NumField(); i++ {
                field := specValue.Field(i)
                fieldType := specType.Field(i)
                value := tmpValue.FieldByName(fieldType.Name)
                if value.Type().AssignableTo(fieldType.Type) && field.CanSet() &&
                        value.Interface() != field.Interface() {
                        field.Set(value)
                        changed = true
                }
        }

        return
}

// Monitor watches the file 'name' for changes in the environment. If a change
// is detected, true is sent on the channel returned and the Spec is
// updated. Close the channel to stop monitoring.
// At least two consecutive errors must occur when attempting to handle the JSON
// before the routine will panic. This should handle a settings file that is being
// changed.
func (s *Spec) Monitor(name string, dur time.Duration) (ch chan bool) {
        const MAX_ERR = 2
        ch = make(chan bool, 1)

        go func() {
                // Check for updates every tick
                t := time.NewTicker(dur)
                errors := 0
                for {
                        select {
                        case <-t.C:
                                res, err := s.jsonFile(name)
                                if err != nil {
                                        errors++
                                        if errors >= MAX_ERR {
                                                panic(err)
                                        }
                                }
                                if res {
                                        ch <- res
                                }
                                errors = 0
                        case r, ok := <-ch:
                                if !ok {
                                        t.Stop()
                                        break
                                }
                                ch <- r
                        }
                }
        }()

        return
}

// jsonFile opens a file and uses the JSON environment to interpret it's contents
// and apply them to the spec. Returns the result of running s.Json on the file.
func (s *Spec) jsonFile(name string) (bool, error) {
        f, err := os.Open(name)
        if err != nil {
                return false, err
        }
        defer f.Close()

        res, err := s.Json(f)
        if err != nil {
                return false, err
        }
        return res, nil
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

func (s *Spec) copy() *Spec {
        copy := &Spec{}
        *copy = *s
        return copy
}
