package slave

import (
        "fmt"
        "io/ioutil"
        "net/http"
        "net/http/httptest"
        "os"
        "os/exec"
        "testing"
)

type mockError struct {
        int
}

func (e mockError) Error() string {
        return fmt.Sprintf("%d", e)
}

func TestFetchError(t *testing.T) {
        err := fetchError{"Test Error", &mockError{2}}

        if err.Error() != "Test Error: {2}" {
                t.Error("Expected: Test Error: {2}", "Got: ", err.Error())
        }
}

var fetch_data = "Test content"

func mockGet(url string) ([]byte, error) {
        return []byte(fetch_data), nil
}

type testFetcher struct{}

func (f *testFetcher) Url() string {
        return "http://test/"
}

func (f *testFetcher) File() string {
        return "test_fetch0123.txt"
}

func (f *testFetcher) Overwrite() bool {
        return true
}

func (f *testFetcher) Command() *exec.Cmd {
        return nil
}

func TestFetchfn(t *testing.T) {
        // Setup environment
        os.Clearenv()
        Environment()
        fetcher := &testFetcher{}

        // Test
        fetchfn(mockGet, fetcher)

        f, err := os.Open(fetcher.File())
        if os.IsNotExist(err) {
                t.Error(err)
        }
        defer f.Close()
        defer os.Remove(fetcher.File())

        content, err := ioutil.ReadAll(f)
        if string(content) != fetch_data {
                t.Errorf("Expected: %s  Got: %s", fetch_data, string(content))
        }

        // Cleanup
        spec = nil
}

func TestGet(t *testing.T) {
        ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                fmt.Fprint(w, fetch_data)
        }))
        defer ts.Close()

        res, err := get(ts.URL)
        if err != nil {
                t.Fatal(err)
        }

        if string(res) != fetch_data {
                t.Errorf("Expected: %s, Got: %s", fetch_data, string(res))
        }
}
