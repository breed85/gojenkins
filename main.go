package main

import (
        "flag"
        "fmt"
        "log"
        "os"
        "os/signal"
        "stash.jda.com/scm/~j1014191/gojenkins/slave"
        "time"
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
        // We catch the Interrupt to ensure that app stops and not just the
        // process that is launched.
        signalChan := make(chan os.Signal, 1)
        signal.Notify(signalChan, os.Interrupt)
        // quit will be used to allow us to exit on SIGINT.
        quit := make(chan bool, 1)
        go func() {
                // Block until a signal is caught
                <-signalChan
                quit <- true
        }()

        // runError will be used to notify of any errors that occur while running the slave.
        runError := make(chan error, 1)

        // Load the environment
        s := slave.NewSpec()
        env, err := s.Environment()
        if err != nil {
                log.Print(err)
                exitCode = 2
                return
        }

        // Create the flags
        f := flags{env}
        f.Load()

        flag.Usage = usage
        flag.Parse()

        if !f.ValidMode() {
                log.Printf(
                        "Invalid value for mode [%s]. Valid values are 'normal' and 'exclusive'.",
                        f.Mode,
                )
                exitCode = 2
                return
        }

        // If a lock file is desired, create it and defer the unlock
        if f.Lock {
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

        }

        var mon chan bool
        if len(f.File) > 0 {
                // Read the JSON config first
                c, err := os.Open(f.File)
                if err != nil {
                        log.Printf("Error: %s", err)
                        exitCode = 2
                        return
                }
                _, err = f.Json(c)
                if err != nil {
                        log.Printf("Error: %s", err)
                        exitCode = 2
                        c.Close()
                        return
                }
                c.Close()

                // Begin monitoring the config every minute
                mon = f.Monitor(f.File, time.Minute)
        }

        // Launch the slave
        if f.Swarm {
                go slave.Run(&slave.Swarm{mon}, runError)
        } else {
                go slave.Run(&slave.Jnlp{mon}, runError)
        }

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
                close(mon)
        }
}

func usage() {
        fmt.Fprintf(os.Stderr,
                "usage: gojenkins\n\n"+
                        "   This command will launch a jenkins slave. Command line options override the environment.\n"+
                        "   All command line options can be set via environment variables prefixed with 'SLAVE_'.\n\n",
        )
        flag.PrintDefaults()
        os.Exit(2)
}

// flags is a wrapper for the environment spec allowing the CLI to interact with the environment.
type flags struct {
        *slave.Spec
}

// Load sets up the command line options. CLI options override any environment variables set.
func (f *flags) Load() {
        flag.StringVar(&f.Server, "server", f.Server, "\n\tJenkins server to use. Ex. http://localhost:8080\n")
        flag.StringVar(&f.Home, "home", f.Home, "\n\tJenkins working directory.\n")
        flag.BoolVar(&f.Swarm, "swarm", f.Swarm, "\n\tUse swarm client to connect to jenkins.\n")
        flag.StringVar(&f.Swarmversion, "swarmversion", f.Swarmversion, "\n\tVersion of swarm client to use. Requires -swarm\n")
        flag.BoolVar(&f.Lock, "lock", f.Lock, "\n\tCreate a lock file during execution in -home directory with name [name].lock\n")
        flag.StringVar(&f.Username, "username", f.Username, "\n\tUsername to log into the jenkins system. Requires -swarm\n")
        flag.StringVar(&f.Password, "password", f.Password, "\n\tPassword to log into the jenkins system. Requires -swarm\n")
        flag.StringVar(&f.Name, "name", f.Name, "\n\tName of the host on Jenkins. When used with -swarm, the name will be used to create a node.\n\tOtherwise, the node [name] must exist on the master already.\n")
        flag.StringVar(&f.Mode, "mode", f.Mode, "\n\tMode to set for the slave node. Valid values are 'normal' (utilize the slave as much as possible)\n\tor 'exclusive' (leave this machine for tied jobs only). Requires -swarm\n")
        flag.StringVar(&f.Labels, "labels", f.Labels, "\n\tLabels to apply to the node. Requires -swarm. Can be a space separated list.\n")
        flag.IntVar(&f.Executors, "executors", f.Executors, "\n\tNumber of executors to use for the node. Requires -swarm\n")
        flag.StringVar(&f.File, "file", f.File, "\n\tConfig file to use. The content should be json. The file will be automatically monitored for changes.\n\tAny settings in the file will take precedence over command line flags.\n")
}
