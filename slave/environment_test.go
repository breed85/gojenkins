package slave

import (
        "os"
        "strconv"
        "strings"
        "testing"
        "time"
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
        s, err := NewSpec().Environment()
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

        // Test changed values
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

        s, err = NewSpec().Environment()
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

        // Test mode wrong
        os.Clearenv()
        spec = nil
        os.Setenv(ENV_MODE, "wrong")

        s, err = NewSpec().Environment()
        if nil != err {
                t.Error(err)
        }

        expectString(t, "normal", s.Mode)

        // Cleanup
        spec = nil
        os.Clearenv()
}

func TestJson(t *testing.T) {
        spec = nil
        s := NewSpec()

        // Test bad Json
        r := strings.NewReader("{ malformed: json")
        res, err := s.Json(r)
        if err == nil {
                t.Log("Expected error on malformed json.")
                t.Fail()
        }

        expectBool(t, false, res)

        // Test good Json
        r = strings.NewReader(`
            {
                "Server": "http://test",
                "Home": "c:\\jenkins",
                "Name": "testslave",
                "Username": "slaveuser",
                "Password": "slavepass",
                "Swarm": true,
                "Swarmversion": "1.16",
                "Lock": true,
                "Executors": 8,
                "Mode": "exclusive",
                "Labels": "Label1 Label2"
            }`)
        res, err = s.Json(r)
        if err != nil {
                t.Log("Error: ", err)
                t.Fail()
        }

        //Expect the result to be true for changed data
        expectBool(t, true, res)

        // Check that the Spec is correctly updated
        expectString(t, "http://test", s.Server)
        expectString(t, "c:\\jenkins", s.Home)
        expectString(t, "testslave", s.Name)
        expectString(t, "slaveuser", s.Username)
        expectString(t, "slavepass", s.Password)
        expectBool(t, true, s.Swarm)
        expectString(t, "1.16", s.Swarmversion)
        expectBool(t, true, s.Lock)
        expectInt(t, 8, s.Executors)
        expectString(t, "exclusive", s.Mode)
        expectString(t, "Label1 Label2", s.Labels)

        // Test for no updates on no changes
        // Rewind the reader to the beginning of the buffer.
        r.Seek(0, os.SEEK_SET)
        res, err = s.Json(r)
        if err != nil {
                t.Log("Error: ", err)
                t.Fail()
        }
        expectBool(t, false, res)
}

func TestMonitor(t *testing.T) {
        spec = nil
        s := NewSpec()
        name := os.TempDir() + "/TestMonitor.config"
        f, _ := os.Create(name)
        defer os.Remove(name)
        f.Write([]byte(`
            {
                "Server": "http://test",
                "Home": "c:\\jenkins",
                "Name": "testslave",
                "Username": "slaveuser",
                "Password": "slavepass",
                "Swarm": true,
                "Swarmversion": "1.16",
                "Lock": true,
                "Executors": 8,
                "Mode": "exclusive",
                "Labels": "Label1 Label2"
            }`))

        ch := s.Monitor(name, time.Millisecond*30)
        expectBool(t, true, <-ch)

        // Check that the Spec is correctly updated
        expectString(t, "http://test", s.Server)
        expectString(t, "c:\\jenkins", s.Home)
        expectString(t, "testslave", s.Name)
        expectString(t, "slaveuser", s.Username)
        expectString(t, "slavepass", s.Password)
        expectBool(t, true, s.Swarm)
        expectString(t, "1.16", s.Swarmversion)
        expectBool(t, true, s.Lock)
        expectInt(t, 8, s.Executors)
        expectString(t, "exclusive", s.Mode)
        expectString(t, "Label1 Label2", s.Labels)

        close(ch)
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
