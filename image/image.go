package image

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// PathType path format type
type PathType uint8

// Path types supported
const (
	URLImage PathType = iota
	DirectoryImage
	FileImage
)

// PrepareImage prepares the directory to isolate
// If image is a URL it gets downloaded to root/image and extracted to root/root
// If image is a tarball file it gets extracted to root/root
// If image is a directory (and root is empty) that directory will be the new root
// TODO for URL, check that root and image directories doesn't exists
// TODO for tarballs, check that root directory doesn't exists
// TODO if the image parameter is a directory, there's no need to call PrepareImage. Anyway, we are supporting DirectoryImage
func PrepareImage(image string, root string) error {
	log.Debugf("preparing image %s \ninto %s", image, root)
	ipt, err := getImagePathType(image)
	if err != nil {
		return err
	}

	// if it's an URL, download the file to a local folder
	if ipt == URLImage {
		imageDir := filepath.Join(root, "image")
		os.Mkdir(imageDir, 0777)

		// reuse image varaible, so we can use 'image' to access local file instead of URL
		// this is OK as long as we don't have to use again the URL
		image, err = downloadImage(image, imageDir)
		if err != nil {
			return err
		}
	}

	// if it's a file (downloaded or not), decompress to new root
	if ipt == URLImage || ipt == FileImage {

		rootDir := filepath.Join(root, "root")
		os.Mkdir(rootDir, 0777)

		err := extractImage(image, rootDir)
		if err != nil {
			return err
		}
	}

	if ipt == DirectoryImage {
		// if it's a directory, we expect that image == root or root is ""
		if root != "" && root != image {
			// TODO we could copy image contents to root. For now let's return an error
			return fmt.Errorf("If image is a directory, new root must be that same directory of empty.")
		}
	}

	return nil
}

// getImagePathType check image type based on the image path string
func getImagePathType(image string) (PathType, error) {

	// check URL format. Only support http or https
	imageURL, err := url.Parse(image)
	if err == nil && (imageURL.Scheme == "http" || imageURL.Scheme == "https") {
		return URLImage, nil
	}

	// check file or directory
	r, err := os.Stat(image)
	if err != nil {
		return 0, fmt.Errorf("Error checking image type: %s", err.Error())
	}

	// check directory
	if r.IsDir() {
		return DirectoryImage, nil
	}

	// TODO check for other types and don¡nt assume not being a URL or Dir, it's a file
	return FileImage, nil
}
