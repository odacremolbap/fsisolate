package fsisolate

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/odacremolbap/fsisolate/archive"
	"github.com/odacremolbap/fsisolate/net"
)

// PathType path format type
type pathType uint8

// Path types supported
const (
	unknownImage pathType = iota
	netImage
	directoryImage
	fileImage
)

// Image manages isolation images
type Image struct {
	Path   string       // path to the image, can be directory, local tarball file or URL to tarball
	Root   string       // directory where the new root will be placed. Non used if Image.path is a directory
	Client *http.Client // http configured client to download image in case path type is URLImage
}

// Prepare prepares the directory to isolate with chroot
// If Path is a URL the image will get downloaded to a temporary directory and extracted to root
// If Path is a tarball file it will be extracted to root
// If image is a directory that directory will be the new root. Image.Root value won't be used
func (i *Image) Prepare() (string, error) {

	ptype, err := i.fillImagePathType()

	// if it's an URL, download the file to a local folder
	if ptype == unknownImage || err != nil {
		return "", fmt.Errorf("Cannot prepare image: image path format unknown")
	}

	// if it's a directory discard root value and use the image's directory
	if ptype == directoryImage {
		return i.Path, nil
	}

	// real local image file.
	var localFilePath string

	// if it's an URL, download the file to a local temp folder
	// TODO is currently out of the scope: image cache management.
	// For now download the tarball into a temp dir, extract to root, and delete.
	if ptype == netImage {

		// generate temp dir
		var downloadDir string
		downloadDir, err = ioutil.TempDir("", "fsisolate")
		if err != nil {
			return "", err
		}

		// download image to temp dir
		r := net.Resource{}
		localFilePath, err = r.Download(i.Path, downloadDir)
		if err != nil {
			return "", err
		}
	} else {
		// if it's a file assign to local file var
		localFilePath = i.Path
	}

	os.Mkdir(i.Root, 0777)

	err = archive.ExtractTarball(localFilePath, i.Root)
	if err != nil {
		return "", err
	}
	return i.Root, nil

}

// fillImagePathType check image type based on the image path
func (i *Image) fillImagePathType() (pathType, error) {

	// check URL format. Only support http or https
	imageURL, err := url.Parse(i.Path)
	if err == nil && (imageURL.Scheme == "http" || imageURL.Scheme == "https") {
		return netImage, nil
	}

	// check file or directory
	r, err := os.Stat(i.Path)
	if err != nil {
		return unknownImage, fmt.Errorf("Error checking image type: %s", err.Error())
	}

	// check directory
	if r.IsDir() {
		return directoryImage, nil
	}

	// TODO check for other types and don'nt assume not being a URL or Dir, it's a file
	return fileImage, nil
}
