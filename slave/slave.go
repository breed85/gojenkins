// Package slave provides utilities to run a jenkins slave.
package slave

import (
        "fmt"
        "os"
        "os/exec"
        "time"
)

// Run will move execution to the Jenkins working directory as defined in the environment. It will then
// start the Jenkins slave after downloading a new slave jar file from the master.
func Run() error {
        // Change to jenkins directory
        if err := os.Chdir(spec.Jenkinscwd); err != nil {
                return err
        }

        // Attempt to fetch the slave.jar file
        if err := fetch(); err != nil {
                return err
        }

        runslave()

        return nil
}

const (
        INIT_ATTEMPTS = 20              // # of times to run the slave
        INIT_SLEEP    = 2 * time.Second // initial delay between attempts
        MAX_SLEEP     = 5 * time.Minute // max delay between attempts
        MIN_EXEC_TIME = 5 * time.Minute // minimum execution time
        SLAVEFILE     = "slave.jar"     // Jenkins slave file that will be downloaded from the master.
)

// runslave starts the jenkins slave in a loop. It will retry INIT_ATTEMPTS times if the slave
// doesn't run for at least MIN_EXEC_TIME.
func runslave() {
        jnlp := fmt.Sprintf("%s/computer/%s/slave-agent.jnlp", spec.Jenkinsserver, spec.Name)

        attempts := INIT_ATTEMPTS
        sleep := INIT_SLEEP

        for {
                // Setup the command to start Jenkins
                cmd := exec.Command(
                        "java",
                        "-jar",
                        SLAVEFILE,
                        "-jnlpUrl",
                        jnlp,
                )

                cmd.Stderr = logger

                // Start a timer to determine how long it has run for.
                start := time.Now()
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
