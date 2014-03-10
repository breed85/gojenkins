// Package slave provides utilities to run a jenkins slave.
package slave

import (
        "os"
        "os/exec"
        "strings"
        "time"
)

// Connector interface describes a slave that can be run via the Run command.
type Connector interface {
        // Url returns the URL to retrieve the file needed to run the slave such as a JAR file.
        Url() string

        // File returns the name of the file that will be created from the content at Url().
        File() string

        // Overwrite returns a bool. True indicates that the File() should be overwritten with each execution.
        // False indicates that the File() should be reused on each execution.
        Overwrite() bool

        // Command builds an exec.Cmd object that will be run to start the slave.
        Command() *exec.Cmd

        // Restart provides a channel to listen for a restart
        Restart() chan bool
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
        restart := c.Restart()

        for {
                // Setup the command to start Jenkins
                cmd := c.Command()
                if logger != nil {
                        cmd.Stderr = logger
                }

                // Start a timer to determine how long it has run for.
                start := time.Now()
                printf("Executing %s", strings.Join(cmd.Args, " "))
                if err := cmd.Start(); err != nil {
                        print(err)
                }

                finished := make(chan bool)
                go func() {
                        cmd.Wait()
                        finished <- true
                }()

                select {
                case <-finished:
                        // If we ran for a while and shut down unexpectedly, do not try to restart.
                        if time.Since(start) > MIN_EXEC_TIME {
                                print("Jenkins shut down unexpectedly, but ran for a decent time. Quitting.")
                                return
                        }

                        // While we have additional attempts left, sleep and attempt to start Jenkins again.
                        if attempts > 0 {
                                printf("Jenkins aborted rather quickly. Will try again in %d seconds.", sleep/time.Second)
                                printf("%d attempts remaining.", attempts)
                                attempts--

                                time.Sleep(sleep)

                                sleep *= 2
                                if sleep > MAX_SLEEP {
                                        sleep = MAX_SLEEP
                                }

                        } else {
                                printf("Failed to start Jenkins in %d attempts. Quitting.", INIT_ATTEMPTS)
                                return
                        }
                case res, ok := <-restart:
                        if !ok {
                                // The monitor was closed, quit
                                cmd.Process.Kill()
                                break
                        }
                        if res {
                                // Got a message to restart
                                print("Restarting...")
                                cmd.Process.Kill()
                                attempts = INIT_ATTEMPTS
                        }
                }

        }
}

func print(v ...interface{}) {
        if logger != nil {
                logger.Print(v...)
        }
}

func printf(format string, v ...interface{}) {
        if logger != nil {
                logger.Printf(format, v...)
        }
}
