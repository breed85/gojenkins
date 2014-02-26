package slave

import (
        "fmt"
        "os/exec"
        "strconv"
)

// Swarm is an object that represents a Jenkins swarm client connector for connecting to a Jenkins master.
type Swarm struct{}

func (s *Swarm) Url() string {
        return fmt.Sprintf("http://maven.jenkins-ci.org/content/repositories/releases/org/jenkins-ci/plugins/swarm-client/%s/%s", spec.Swarmversion, s.File())
}

func (s *Swarm) File() string {
        return fmt.Sprintf("swarm-client-%s-jar-with-dependencies.jar", spec.Swarmversion)
}

func (s *Swarm) Overwrite() bool {
        return false
}

func (s *Swarm) Command() *exec.Cmd {
        env, _ := Environment()

        args := []string{"-jar", s.File()}

        if len(env.Username) > 0 {
                args = append(args, "-username", env.Username)
        }

        if len(env.Password) > 0 {
                args = append(args, "-password", env.Password)
        }

        if len(env.Server) > 0 {
                args = append(args, "-master", env.Server)
        }

        if len(env.Labels) > 0 {
                args = append(args, "-labels", env.Labels)
        }

        args = append(args,
                "-mode",
                env.Mode,
                "-name",
                env.Name,
                "-executors",
                strconv.Itoa(env.Executors),
        )

        return exec.Command("java", args...)
}
