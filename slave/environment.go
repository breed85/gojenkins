package slave

import (
        "github.com/kelseyhightower/envconfig"
)

type Spec struct {
        Jenkinsserver string
        Jenkinscwd    string
        Name          string
        Jarurl        string
}

var spec Spec

func Environment() error {
        return envconfig.Process("slave", &spec)
}
