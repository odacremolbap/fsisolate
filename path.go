package fsisolate

import (
	"fmt"
	"net/url"
	"os"
)

// PathType path format type
type pathType uint8

// Path types supported
const (
	unknownPath pathType = iota
	urlPath
	directoryPath
	filePath
)

// GetPathType returns the path type based on a path string
func getPathType(path string) (pathType, error) {

	// check URL format. Only support http or https
	u, err := url.Parse(path)
	if err == nil && (u.Scheme == "http" || u.Scheme == "https") {
		return urlPath, nil
	}

	// check file or directory
	r, err := os.Stat(path)
	if err != nil {
		return unknownPath, fmt.Errorf("Error checking path type: %s", err.Error())
	}

	// check directory
	if r.IsDir() {
		return directoryPath, nil
	}

	// TODO check for other types and don'nt assume not being a URL or Dir, it's a file
	return filePath, nil
}
