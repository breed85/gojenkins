package main

import (
        "flag"
        "log"
        "os"
        "stash.jda.com/scm/~j1014191/gojenkins/slave"
)

func main() {
        flag.Usage = slave.Usage
        flag.Parse()

        if err := slave.Environment(); err != nil {
                log.Fatal(err)
                os.Exit(2)
        }

        slave.Run()
}
