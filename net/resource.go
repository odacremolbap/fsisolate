package net

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Resource exposes helper methods to access internet resoruces
type Resource struct {
	Client *http.Client
}

// Download downloads a resource from internet
// TODO support 302 redirections
func (r *Resource) Download(resourceURL string, directory string) (string, error) {

	// use a default
	if r.Client == nil {
		r.Client = &http.Client{}
	}

	parsedURL, err := url.Parse(resourceURL)
	if err != nil {
		return "", err
	}

	// get name of the resource to download
	// if name is empty use a default one
	// TODO remove querystring if present
	fileName := parsedURL.Path[(strings.LastIndex(parsedURL.Path, "/") + 1):]

	resp, err := r.Client.Get(resourceURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Remote site returned status %s.", resp.Status)
	}

	var file *os.File
	if fileName != "" {
		file, err = os.Create(filepath.Join(directory, fileName))
	} else {
		// if no remote file name could be extracted, generate a file name in the target directory
		file, err = ioutil.TempFile(directory, "isolate")
	}

	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err = io.Copy(file, resp.Body); err != nil {
		return "", err
	}

	return file.Name(), nil
}
