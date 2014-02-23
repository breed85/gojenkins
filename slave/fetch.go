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
        message := e.message
        if e.err != nil {
                message = fmt.Sprintf("%s: %s", message, e.err)
        }
        return message
}

type getter func(string) ([]byte, error)

// fetch is the main entry point for fetching the slave.jar.
func fetch() error {
        return fetchfn(get)
}

// fetchfn retrieves the slave.jar from the jenkins master.
// To support testing, the function accepts a getter that will be called to retrieve
// the content of the slave file. For testing purposes, a mock getter can be passed
// without the need to setup a mock http server.
func fetchfn(fn getter) error {
        jarUrl := fmt.Sprintf("%s/jnlpJars/slave.jar", spec.Jenkinsserver)

        // Read slave.jar content
        content, err := fn(jarUrl)
        if err != nil {
                return fetchError{"Failed to read response", err}
        }

        // Create destination file, replacing it if it exists
        f, err := os.Create(SLAVEFILE)
        if err != nil {
                return fetchError{"Error creating slave.jar", err}
        }
        defer f.Close()

        // Write slave.jar
        if _, err := f.Write(content); err != nil {
                return fetchError{"Failed to write slave.jar", err}
        }

        return nil
}

// get retrieves the content at the specified URL. The content will be returned as a []byte.
// If an error occurs it will be passed back to the caller.
func get(url string) ([]byte, error) {
        resp, err := http.Get(url)
        if err != nil {
                return nil, err
        } else if resp.StatusCode != http.StatusOK {
                return nil, fetchError{fmt.Sprintf("%d: %s", resp.StatusCode, resp.Status), nil}
        }
        defer resp.Body.Close()

        return ioutil.ReadAll(resp.Body)
}
