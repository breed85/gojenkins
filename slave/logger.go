package slave

import (
        "log"
        "os"
)

// Logger supports capturing output from the Jenkins slave process.
type Logger struct {
        *log.Logger
}

var logger *Logger

func init() {
        logger = &Logger{log.New(os.Stderr, "Slave: ", log.Ldate|log.Ltime|log.Lshortfile)}
}

func (l *Logger) Write(p []byte) (n int, err error) {
        l.Print(string(p))
        return len(p), nil
}
