package slave

import (
        "io"
        "log"
)

// Logger supports capturing output from the Jenkins slave process.
type Logger struct {
        *log.Logger
}

var logger *Logger

func InitLog(w io.Writer) {
        logger = &Logger{log.New(w, "Slave: ", log.Ldate|log.Ltime|log.Lshortfile)}
}

func (l *Logger) Write(p []byte) (n int, err error) {
        l.Print(string(p))
        return len(p), nil
}
