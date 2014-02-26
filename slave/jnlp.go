package slave

import (
        "fmt"
        "os/exec"
)

// Jnlp struct represents a Jnlp connector to a Jenkins master.
type Jnlp struct{}

func (j *Jnlp) Url() string {
        return fmt.Sprintf("%s/jnlpjars/%s", spec.Server, j.File())
}

func (j *Jnlp) File() string {
        return "slave.jar"
}

func (j *Jnlp) Overwrite() bool {
        return true
}

func (j *Jnlp) Command() *exec.Cmd {
        return exec.Command(
                "java",
                "-jar",
                j.File(),
                "-jnlpurl",
                fmt.Sprintf("%s/computer/%s/slave-agent.jnlp", spec.Server, spec.Name),
        )
}
