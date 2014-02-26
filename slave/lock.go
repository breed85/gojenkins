package slave

import (
        "fmt"
        "os"
)

// LockFile wraps os.File
type LockFile struct {
        *os.File
}

// NewLock will create a lock file based on the slave name with ".lock" appended in the Jenkins working directory.
func NewLock() (*LockFile, error) {
        name := fmt.Sprintf("%s.lock", spec.Name)

        // Move the Jenkins working directory
        if err := os.Chdir(spec.Home); err != nil {
                return nil, err
        }

        // Create the lock file
        f, err := os.Create(name)
        if err != nil {
                return nil, err
        }

        return &LockFile{f}, nil
}

// Unlock removes the lock file
func (l *LockFile) Unlock() error {
        l.Close()
        err := os.Remove(l.Name())
        return err
}
