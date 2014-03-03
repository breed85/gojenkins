package slave

import (
        "testing"
)

func TestSwarmUrl(t *testing.T) {
        NewSpec().Environment()
        s := &Swarm{}

        exp := "http://maven.jenkins-ci.org/content/repositories/releases/org/jenkins-ci/plugins/swarm-client/1.15/swarm-client-1.15-jar-with-dependencies.jar"
        if url := s.Url(); url != exp {
                t.Log("Expected: %s, Got: %s", exp, url)
        }

        spec = nil
}

func TestSwarmOverwrite(t *testing.T) {
        s := &Swarm{}

        if s.Overwrite() {
                t.Log("Expected false, Got true.")
                t.Fail()
        }
}

func TestSwarmCommand(t *testing.T) {
        s := &Swarm{}

        exp := []string{
                "java",
                "-jar",
                "swarm-client-1.15-jar-with-dependencies.jar",
                "-username",
                "testuser",
                "-password",
                "testpass",
                "-master",
                "http://test",
                "-labels",
                "label1 label2",
                "-mode",
                "normal",
                "-name",
                "testname",
                "-executors",
                "2",
        }
        env, _ := NewSpec().Environment()
        env.Username = "testuser"
        env.Password = "testpass"
        env.Server = "http://test"
        env.Labels = "label1 label2"
        env.Name = "testname"

        cmd := s.Command()

        for i, v := range cmd.Args {
                if v != exp[i] {
                        t.Logf("Expected: %s, Got: %s", exp[i], v)
                }
        }
}
