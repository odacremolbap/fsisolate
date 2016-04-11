package image

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// downloadImage retrieves a image from internet
// TODO support 302 redirections
func downloadImage(image string, root string) (string, error) {

	imageURL, err := url.Parse(image)
	if err != nil {
		return "", err
	}

	// get name of the file to download
	// if name is empty use a default one
	fileName := imageURL.Path[(strings.LastIndex(imageURL.Path, "/") + 1):]
	if fileName == "" {
		fileName = "isolate-image"
	}

	fileName = filepath.Join(root, fileName)

	// TODO create channel to get progress
	resp, err := http.Get(image)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Remote site returned status %s.", resp.Status)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err = io.Copy(file, resp.Body); err != nil {
		return "", err
	}

	return fileName, nil
}
