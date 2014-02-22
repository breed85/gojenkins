package slave

import (
        "fmt"
        "io/ioutil"
        "net/http"
        "os"
)

type fetchError struct {
        message string
        err     error
}

func (e fetchError) Error() string {
        return fmt.Sprintf("%s: %s", e.message, e.err)
}

type getter func(string) ([]byte, error)

func fetch() error {
        return fetchfn(get)
}

func fetchfn(fn getter) error {
        jarUrl := fmt.Sprintf("%s/jnlpJars/slave.jar", spec.Jenkinsserver)

        // Create destination file, replacing it if it exists
        f, err := os.Create(SLAVEFILE)
        if err != nil {
                return fetchError{"Error creating slave jar", err}
        }
        defer f.Close()

        // Read slave.jar content
        content, err := fn(jarUrl)
        if err != nil {
                return fetchError{"Failed to read response", err}
        }

        // Write slave.jar
        if _, err := f.Write(content); err != nil {
                return fetchError{"Failed to write slave.jar", err}
        }

        return nil
}

func get(url string) ([]byte, error) {
        resp, err := http.Get(url)
        if err != nil {
                return nil, err
        }
        defer resp.Body.Close()

        return ioutil.ReadAll(resp.Body)
}
