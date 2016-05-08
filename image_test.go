package fsisolate

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

// mockedResource returns a mocked web server and a client that redirects all request to the server
func mockedResource(status int, body []byte) (*httptest.Server, *http.Client) {

	// create server that returns expected response
	s := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
				w.Write(body)
			},
		),
	)

	// create resource that rediects any request to the server
	c := &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse(s.URL)
			},
		},
	}

	return s, c
}

func TestPrepareImage(t *testing.T) {

	var testData = []struct {
		path         string // image path
		root         string // root directory
		status       int    // if URL, status returned
		bodyFile     string // if URL, body returned
		expectedRoot string // expected root returned from prepare
		prepareOK    bool   // whether path should be OK or return an error
	}{
		{"testdata/", "whatever/", 0, "", "testdata/", true},
		{"http://test.url", "testdata/tmp/", 200, "testdata/test.tar", "testdata/tmp/", true},
		{"testdata/test.tar", "testdata/tmp/", 0, "", "testdata/tmp/", true},
		{"*?<notapath", "whatever/", 0, "", "", false},
	}

	for _, td := range testData {

		i := Image{}

		// if an error occurs, we want to stick with unknownPath
		pt, _ := getPathType(td.path)

		if pt == urlPath {
			b, err := ioutil.ReadFile(td.bodyFile)
			if err != nil {
				t.Errorf("Couldn't read file %q to create mocked body for image at %q at directory %q: %s", td.bodyFile, td.path, td.root, err)
				continue
			}
			server, client := mockedResource(td.status, b)
			defer server.Close()

			i.Client = client
		}

		root, err := i.Prepare(td.path, td.root)
		if err != nil {
			if td.prepareOK {
				t.Errorf("Couldn't prepare image at %q at directory %q: %s", td.path, td.root, err)
				continue
			}
			// if failed, but we expected the failure, continue
			continue
		}

		if err == nil && !td.prepareOK {
			t.Errorf("Image prepare for %q to %q should have failed, but did not", td.path, td.root)
			continue
		}

		// remove created dir
		if pt != directoryPath {
			if err := os.RemoveAll(root); err != nil {
				t.Errorf("Couldn't delete root directory %q for test with path %q: %s", root, td.path, err)
			}
		}

		if root != td.expectedRoot {
			t.Errorf("Prepare for image at %q and directory %q returned root %q when expected %q", td.path, td.root, root, td.expectedRoot)
		}
	}

}
