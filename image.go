package fsisolate

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/odacremolbap/fsisolate/archive"
	"github.com/odacremolbap/fsisolate/net"
)

// Image manages isolation images
type Image struct {
	Client *http.Client // http configured client to download image in case path type is URLImage
}

// Prepare prepares the directory to isolate with chroot
// path: the path to the image. can be directory, local tarball file or URL to tarball
// root: directory where the new root will be placed. Non used if Image.path is a directory
// If path is a URL the image will get downloaded to a temporary directory and extracted to root
// If path is a tarball file it will be extracted to root
// If path is a directory that directory will be the new root. Image.Root value won't be used
func (i *Image) Prepare(path, root string) (string, error) {

	ptype, err := getPathType(path)

	// if it's an URL, download the file to a local folder
	if ptype == unknownPath || err != nil {
		return "", fmt.Errorf("Cannot prepare image: image path format unknown")
	}

	// if it's a directory discard root value and use the image's directory
	if ptype == directoryPath {
		return path, nil
	}

	// real local image file.
	var localFilePath string

	// if it's an URL, download the file to a local temp folder
	// TODO is currently out of the scope: image cache management.
	// For now download the tarball into a temp dir, extract to root, and delete.
	if ptype == urlPath {

		// generate temp dir
		var downloadDir string
		downloadDir, err = ioutil.TempDir("", "fsisolate")
		if err != nil {
			return "", err
		}

		// download image to temp dir
		r := net.Resource{}
		localFilePath, err = r.Download(path, downloadDir)
		if err != nil {
			return "", err
		}
	} else {
		// if it's a file assign to local file var
		localFilePath = path
	}

	os.Mkdir(root, 0777)

	err = archive.ExtractTarball(localFilePath, root)
	if err != nil {
		return "", err
	}
	return root, nil

}
