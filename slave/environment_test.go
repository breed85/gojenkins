package slave

import (
        "os"
        "testing"
)

func TestEnvironment(t *testing.T) {
        // Setup fake environment
        var (
                server  = "http://jenkins/server"
                cwd     = "c:\\jenkins\\workspace"
                name    = "slavename"
        )

        os.Setenv("SLAVE_JENKINSSERVER", server)
        os.Setenv("SLAVE_JENKINSCWD", cwd)
        os.Setenv("SLAVE_NAME", name)

        if err := Environment(); err != nil {
                t.Error(err)
        }

        if server != spec.Jenkinsserver {
                t.Error("Expected ", server, " got ", spec.Jenkinsserver)
        }
        if cwd != spec.Jenkinscwd {
                t.Error("Expected ", cwd, " got ", spec.Jenkinscwd)
        }
        if name != spec.Name {
                t.Error("Expected ", name, " got ", spec.Name)
        }
}
