package slave

import (
        "fmt"
        "log"
        "os"
        "os/exec"
        "time"
)

var logger *log.Logger

func init() {
        logger = log.New(os.Stderr, "Slave: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Usage() {
        fmt.Fprintf(os.Stderr,
                "usage: gojenkins\n\n"+
                        "   This command will launch a jenkins slave. There are no command line options, but\n"+
                        "   environment variables can be used to configure each instance.\n\n"+
                        "   SLAVE_JENKINSSERVER    Jenkins master server URL.\n"+
                        "                          Ex. http://master:8080/\n\n"+
                        "   SLAVE_JENKINSCWD       The directory to use to execute the Jenkins slave.\n\n"+
                        "   SLAVE_NAME             Name of the node as defined on the Jenkins master nodes list.\n"+
                        "                          Ex. http://master:8080/computers\n\n",
        )
        os.Exit(2)
}

func Run() {
        // Change to jenkins directory
        if err := os.Chdir(spec.Jenkinscwd); err != nil {
                logger.Fatal(err)
                os.Exit(2)
        }

        fetch()

        runslave()

}

const SLAVEFILE = "slave.jar"

func runslave() {
        jnlp := fmt.Sprintf("%s/computer/%s/slave-agent.jnlp")

        const (
                INIT_ATTEMPTS = 20              // # of times to run the slave
                INIT_SLEEP    = 2 * time.Second // initial delay between attempts
                MAX_SLEEP     = 5 * time.Minute // max delay between attempts
                MIN_EXEC_TIME = 5 * time.Minute // minimum execution time
        )

        attempts := INIT_ATTEMPTS
        sleep := INIT_SLEEP

        for {
                cmd := exec.Command(
                        "java",
                        "-jar",
                        SLAVEFILE,
                        "-jnlpUrl",
                        jnlp,
                )

                start := time.Now()

                if err := cmd.Run(); err != nil {
                        logger.Println(err)
                }

                if time.Since(start) > MIN_EXEC_TIME {
                        logger.Print("Jenkins shut down unexpectedly, but ran for a decent time. Quitting.")
                        return
                }

                if attempts > 0 {
                        attempts--
                        logger.Printf("Jenkins aborted rather quickly. Will try again in %d seconds.", sleep/time.Second)
                        logger.Printf("%d attempts remaining.", attempts)

                        time.Sleep(sleep)

                        sleep *= 2
                        if sleep > MAX_SLEEP {
                                sleep = MAX_SLEEP
                        }

                } else {
                        return
                }
        }
}
