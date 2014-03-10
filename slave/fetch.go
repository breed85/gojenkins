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
func fetch(c Connector) error {
        return fetchfn(get, c)
}

// fetchfn TODO
func fetchfn(fn getter, c Connector) error {
        // Get out if we are not going to overwrite and the file already exists.
        if _, err := os.Stat(c.File()); !c.Overwrite() && !os.IsNotExist(err) {
                return nil
        }

        // Read content
        content, err := fn(c.Url())
        if err != nil {
                return fetchError{"Failed to read response", err}
        }

        // Create destination file, replacing it if it exists
        file, err := os.Create(c.File())
        if err != nil {
                return fetchError{fmt.Sprintf("Error creating %s", c.File()), err}
        }
        defer file.Close()

        // Write slave.jar
        if _, err := file.Write(content); err != nil {
                return fetchError{fmt.Sprintf("Failed to write %s", c.File()), err}
        }

        return nil
}

// get retrieves the content at the specified URL. The content will be returned as a []byte.
// If an error occurs it will be passed back to the caller.
func get(url string) ([]byte, error) {
        if logger != nil {
                logger.Printf("GET %s", url)
        }

        resp, err := http.Get(url)
        if err != nil {
                return nil, err
        } else if resp.StatusCode != http.StatusOK {
                return nil, fetchError{fmt.Sprintf("%s Resp: %s", url, resp.Status), err}
        }
        defer resp.Body.Close()

        return ioutil.ReadAll(resp.Body)
}
