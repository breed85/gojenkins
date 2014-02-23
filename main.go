package main

import (
        "flag"
        "fmt"
        "log"
        "os"
        "stash.jda.com/scm/~j1014191/gojenkins/slave"
)

var exitCode = 0

func exit() {
        os.Exit(exitCode)
}

func main() {
        defer exit()

        lockFile := flag.Bool("lock", false, "create a lock file during execution")
        flag.Usage = usage
        flag.Parse()

        // Load the environment
        if err := slave.Environment(); err != nil {
                log.Print(err)
                exitCode = 2
                return
        }

        // If a lock file is desired, create it and defer the unlock
        if *lockFile {
                l, err := slave.NewLock()

                if err != nil {
                        log.Printf("Unable to create lock file: %s", err)
                        exitCode = 2
                        return
                }
                // Defer unlocking. This will happen before exit due to LIFO stack of defer statements.
                defer func() {
                        removeerr := l.Unlock()
                        if removeerr != nil {
                                log.Fatal("Failed to remove lock %s", l.Name())
                        }
                }()
        }

        if err := slave.Run(); err != nil {
                // Unexpected error
                log.Print(err)
                exitCode = 2
                return
        }
}

func usage() {
        fmt.Fprintf(os.Stderr,
                "usage: gojenkins\n\n"+
                        "   This command will launch a jenkins slave. There are no command line options, but\n"+
                        "   environment variables can be used to configure each instance.\n\n"+
                        "   SLAVE_JENKINSSERVER    Jenkins master server URL.\n"+
                        "                          Ex. http://master:8080/\n\n"+"   SLAVE_JENKINSCWD       The directory to use to execute the Jenkins slave.\n\n"+
                        "   SLAVE_NAME             Name of the node as defined on the Jenkins master nodes list.\n"+
                        "                          Ex. http://master:8080/computers\n\n",
        )
        flag.PrintDefaults()
        os.Exit(2)
}
