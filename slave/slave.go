// Package slave provides utilities to run a jenkins slave.
package slave

import (
        "os"
        "os/exec"
        "strings"
        "time"
)

type Connector interface {
        Url() string
        File() string
        Overwrite() bool
        Command() *exec.Cmd
}

// Run will move execution to the Jenkins working directory as defined in the environment.
// Fetch the file for the connector and run the slave.
func Run(c Connector, res chan<- error) {
        // Change to jenkins directory
        if err := os.Chdir(spec.Home); err != nil {
                res <- err
                return
        }

        // Attempt to fetch the connector's file
        if err := fetch(c); err != nil {
                res <- err
                return
        }

        // Run the connector
        runslave(c)

        res <- nil
}

const (
        INIT_ATTEMPTS = 10              // # of times to run the slave
        INIT_SLEEP    = 2 * time.Second // initial delay between attempts
        MAX_SLEEP     = 5 * time.Minute // max delay between attempts
        MIN_EXEC_TIME = 5 * time.Minute // minimum execution time
)

// runslave starts the jenkins slave in a loop. It will retry INIT_ATTEMPTS times if the slave
// doesn't run for at least MIN_EXEC_TIME.
func runslave(c Connector) {
        attempts := INIT_ATTEMPTS
        sleep := INIT_SLEEP

        for {
                // Setup the command to start Jenkins
                cmd := c.Command()
                cmd.Stderr = logger

                // Start a timer to determine how long it has run for.
                start := time.Now()
                logger.Printf("Executing %s", strings.Join(cmd.Args, " "))
                if err := cmd.Run(); err != nil {
                        logger.Println(err)
                }

                // If we ran for a while and shut down unexpectedly, do not try to restart.
                if time.Since(start) > MIN_EXEC_TIME {
                        logger.Print("Jenkins shut down unexpectedly, but ran for a decent time. Quitting.")
                        return
                }

                // While we have additional attempts left, sleep and attempt to start Jenkins again.
                if attempts > 0 {
                        logger.Printf("Jenkins aborted rather quickly. Will try again in %d seconds.", sleep/time.Second)
                        logger.Printf("%d attempts remaining.", attempts)
                        attempts--

                        time.Sleep(sleep)

                        sleep *= 2
                        if sleep > MAX_SLEEP {
                                sleep = MAX_SLEEP
                        }

                } else {
                        logger.Printf("Failed to start Jenkins in %d attempts. Quitting.", INIT_ATTEMPTS)
                        return
                }
        }
}
