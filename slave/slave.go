package slave

import (
        "flag"
        "fmt"
        "log"
        "os"
)

func Usage() {
        fmt.Fprintf(os.Stderr, "usage: TODO\n")
        flag.PrintDefaults()
        os.Exit(2)
}

func Run() {
        // Entring Main loop

        if err := os.Chdir(spec.Jenkinscwd); err != nil {
                log.Fatal(err)
                os.Exit(2)
        }

        fetch()
}
