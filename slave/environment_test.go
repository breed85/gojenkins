package slave

import (
        "os"
        "strconv"
        "testing"
)

func TestEnvironment(t *testing.T) {
        var (
                server       = "http://jenkins/server"
                home         = "c:\\jenkins\\workspace"
                name         = "slavename"
                username     = "jenkins-slave"
                password     = "pass"
                swarm        = true
                swarmversion = "1.16"
                lock         = true
                executors    = 4
                mode         = "exclusive"
                labels       = "label1 label2"
        )

        // Test defaults
        os.Clearenv()
        s, err := Environment()
        if nil != err {
                t.Error(err)
        }

        host, _ := os.Hostname()
        dir, _ := os.Getwd()

        expectString(t, host, s.Name)
        expectString(t, dir, s.Home)
        expectBool(t, false, s.Swarm)
        expectBool(t, false, s.Lock)
        expectString(t, "1.15", s.Swarmversion)
        expectInt(t, 2, s.Executors)
        expectString(t, "normal", s.Mode)

        // Setup fake environment
        os.Clearenv()
        spec = nil
        os.Setenv(ENV_SERVER, server)
        os.Setenv(ENV_HOME, home)
        os.Setenv(ENV_NAME, name)
        os.Setenv(ENV_USERNAME, username)
        os.Setenv(ENV_PASSWORD, password)
        os.Setenv(ENV_SWARM, strconv.FormatBool(swarm))
        os.Setenv(ENV_SWARMVERSION, swarmversion)
        os.Setenv(ENV_LOCK, strconv.FormatBool(lock))
        os.Setenv(ENV_EXECUTORS, strconv.Itoa(executors))
        os.Setenv(ENV_MODE, mode)
        os.Setenv(ENV_LABELS, labels)

        s, err = Environment()
        if nil != err {
                t.Error(err)
        }

        expectString(t, server, s.Server)
        expectString(t, home, s.Home)
        expectString(t, name, s.Name)
        expectString(t, username, s.Username)
        expectString(t, password, s.Password)
        expectBool(t, swarm, s.Swarm)
        expectString(t, swarmversion, s.Swarmversion)
        expectBool(t, lock, s.Lock)
        expectInt(t, executors, s.Executors)
        expectString(t, mode, s.Mode)
        expectString(t, labels, s.Labels)

        // Cleanup
        spec = nil
        os.Clearenv()
}

func expectString(t *testing.T, exp, act string) {
        if exp != act {
                t.Logf("Expected: %s, Got: %s", exp, act)
                t.Fail()
        }
}

func expectBool(t *testing.T, exp, act bool) {
        if exp != act {
                t.Logf("Expected: %s, Got: %s", exp, act)
                t.Fail()
        }
}

func expectInt(t *testing.T, exp, act int) {
        if exp != act {
                t.Logf("Expected: %d, Got: %d", exp, act)
                t.Fail()
        }
}
