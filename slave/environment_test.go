package slave

import (
        "os"
        "testing"
)

var (
        server  = "http://jenkins/server"
        cwd     = "c:\\jenkins\\workspace"
        name    = "slavename"
)

func TestEnvironment(t *testing.T) {
        // Setup fake environment

        os.Setenv(ENV_JENKINSSERVER, server)
        os.Setenv(ENV_JENKINSCWD, cwd)
        os.Setenv(ENV_NAME, name)

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

func TestGetJenkinsServer(t *testing.T) {
        testGetS(t, GetJenkinsServer, ENV_JENKINSSERVER, server)
}

func TestGetJenkinsCwd(t *testing.T) {
        testGetS(t, GetJenkinsCwd, ENV_JENKINSCWD, cwd)
}

func TestGetJenkinsName(t *testing.T) {
        testGetS(t, GetName, ENV_NAME, name)
}

type gets func() string

func testGetS(t *testing.T, fn gets, env string, exp string) {
        // Uninitialed value
        spec = Spec{}
        v := fn()

        if len(v) != 0 {
                t.Errorf("Expected nil value, but got %s", v)
        }

        // Initiated value
        os.Setenv(env, exp)

        if err := Environment(); err != nil {
                t.Error(err)
        }

        v = fn()

        if v != exp {
                t.Errorf("Expected: %s, but got %s", exp, v)
        }
}
