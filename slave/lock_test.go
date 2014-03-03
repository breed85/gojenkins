package slave

import (
        "os"
        "testing"
)

func TestLock(t *testing.T) {
        // Setup
        os.Clearenv()
        os.Setenv(ENV_NAME, "host")
        os.Setenv(ENV_HOME, ".")
        // Clear package variable spec to ensure our Env is picked up.
        spec = nil
        NewSpec().Environment()

        // Test lock creation
        f, err := NewLock()
        if err != nil {
                t.Error(err)
        }
        defer os.Remove(f.Name())

        // Fail if the file was not created
        if _, err = os.Stat(f.Name()); os.IsNotExist(err) {
                t.Log("Failed to create file")
                t.Fail()
        }

        // Cleanup
        os.Clearenv()
        spec = nil
}

func TestUnlock(t *testing.T) {
        f, err := os.Create("unlocktest.lock")
        if err != nil {
                t.Error(err)
        }
        l := &LockFile{f}

        l.Unlock()

        if _, err := os.Stat(l.Name()); !os.IsNotExist((err)) {
                t.Log("Failed to unlock file")
                t.Fail()
        }
}
