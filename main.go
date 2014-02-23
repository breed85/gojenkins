package main

import (
        "flag"
        "fmt"
        "log"
        "os"
        "os/signal"
        "stash.jda.com/scm/~j1014191/gojenkins/slave"
)

var exitCode = 0

func exit() {
        os.Exit(exitCode)
}

func main() {
        defer exit()

        // Set up channel on which to send signal notifications.
        // We must use a buffered channel or risk missing the signal
        // if we're not ready to receive when the signal is sent.
        signalChan := make(chan os.Signal, 1)
        signal.Notify(signalChan, os.Interrupt)

        // quit will be used to allow us to exit on SIGINT or normal completion.
        quit := make(chan bool, 1)
        runError := make(chan error, 1)

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
                        if err := l.Unlock(); nil != err && !os.IsNotExist(err) {
                                log.Fatalf("Failed to remove lock %s: %s", l.Name(), err)
                        }
                }()

                // Ensure the lock will be deleted on a signal from the OS
                go func() {
                        <-signalChan
                        signal.Stop(signalChan)
                        quit <- true
                }()
        }

        // Launch the slave
        go slave.Run(runError)

        // Wait for the slave to complete or a signal to quit.
        select {
        case e := <-runError:
                if nil != e {
                        // Unexpected error
                        log.Print(e)
                        exitCode = 2
                        return
                }
        case <-quit:
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
