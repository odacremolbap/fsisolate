package net

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
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

// equalBytes compare two bytes slices
func equalBytes(a, b []byte) bool {

	// nil check
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestDownload(t *testing.T) {

	var testData = []struct {
		url        string // resoruce URL
		status     int    // expected status from URL
		body       []byte // body returned from URL
		directory  string // directory to download to
		file       string // file name of the resource downloaded
		downloadOK bool   // whether the extraction command should succeed
	}{
		{"http://testsite.test/test1", 200, []byte("text"), "testdata/", "test1", true},
		{"http://testsite.test/test2", 200, []byte("text"), "testdata/", "test2", true},
		{"http://testsite.test/test3", 500, []byte("text"), "testdata/", "", false},
		{"http://testsite.test/", 200, []byte("text"), "testdata/", "", true},
	}

	for _, td := range testData {

		server, client := mockedResource(td.status, td.body)
		defer server.Close()

		r := Resource{Client: client}

		download, err := r.Download(td.url, td.directory)

		// if no error ocurred, defer file deletion
		if err == nil {
			defer func(file string) {
				if file != "" {
					if err = os.Remove(file); err != nil {
						t.Errorf("Couldn't delete file %q after test: %s", file, err)
					}
				}
			}(download)
		}

		// check for error downloading resource
		if err != nil {
			if td.downloadOK {
				t.Errorf("Error downloading %q to %q: %s", td.url, td.directory, err)
				continue
			}
			// if an expected error ocurred, continue with next test
			continue
		}

		if err == nil && !td.downloadOK {
			t.Errorf("Downloading %q to %q should have failed, but did not", td.url, td.directory)
			continue
		}

		// check for expected downloaded file name
		if td.file != "" {
			if f := path.Join(td.directory, td.file); download != f {
				t.Errorf("Downloaded file with name %q but expected %q ", download, f)
				continue
			}
		}

		// check for expected file content
		content, err := ioutil.ReadFile(download)
		if err != nil {
			t.Errorf("Cannot read downloaded file %q: %s", download, err)
			continue
		}

		// compare content
		if !equalBytes(content, td.body) {
			t.Errorf("Downloaded file %q contains:\n%v\nbut expected:\n%v\n", download, content, td.body)
		}

	}
}
