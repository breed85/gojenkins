package slave

import (
        "testing"
)

func TestJnlpUrl(t *testing.T) {
        env, _ := NewSpec().Environment()
        env.Server = "http://test"
        j := &Jnlp{}
        s := j.Url()

        exp := env.Server + "/jnlpjars/slave.jar"
        if s != exp {
                t.Logf("Expected: %s Got: %s", exp, s)
                t.Fail()
        }

        spec = nil
}

func TestJnlpOverwrite(t *testing.T) {
        j := &Jnlp{}

        b := j.Overwrite()

        if b != true {
                t.Log("Expected true. Got false")
                t.Fail()
        }
}

func TestJnlpCommand(t *testing.T) {
        env, _ := NewSpec().Environment()
        env.Server = "http://test"
        env.Name = "testname"

        j := &Jnlp{}
        cmd := j.Command()

        exp := []string{
                "java",
                "-jar",
                "slave.jar",
                "-jnlpurl",
                "http://test/computer/testname/slave-agent.jnlp",
        }

        for i, v := range cmd.Args {
                if v != exp[i] {
                        t.Log("Expected: %s, Got: %s", exp[i], v)
                        t.Fail()
                }
        }
}
