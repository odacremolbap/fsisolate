package image

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

// PathType path format type
type PathType uint8

// Path types supported
const (
	Unknown PathType = iota
	URLImage
	DirectoryImage
	FileImage
)

// PrepareImage prepares the directory to isolate with chroot
// If image is a URL it gets downloaded to root/image and extracted to root/root
// If image is a tarball file it gets extracted to root/root
// If image is a directory (and root is empty) that directory will be the new root
// TODO for URL, check that root and image directories doesn't exists
// TODO for tarballs, check that root directory doesn't exists
func PrepareImage(image string, root string) (string, error) {

	ipt, err := getImagePathType(image)
	if err != nil {
		return "", err
	}

	// if it's an URL, download the file to a local folder
	if ipt == URLImage {
		imageDir := filepath.Join(root, "image")
		os.Mkdir(imageDir, 0777)

		// reuse image varaible, so we can use 'image' to access local file instead of URL
		// this is OK as long as we don't have to use again the URL
		image, err = downloadImage(image, imageDir)
		if err != nil {
			return "", err
		}
	}

	// if it's a file (downloaded or not), decompress to new root
	if ipt == URLImage || ipt == FileImage {

		// append 'root' to the root dir argument
		root = filepath.Join(root, "root")
		os.Mkdir(root, 0777)

		err := extractImage(image, root)
		if err != nil {
			return "", err
		}
		return root, nil
	}

	// if it's a directory discard root value and use the image's directory
	return image, nil
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
		return Unknown, fmt.Errorf("Error checking image type: %s", err.Error())
	}

	// check directory
	if r.IsDir() {
		return DirectoryImage, nil
	}

	// TODO check for other types and don'nt assume not being a URL or Dir, it's a file
	return FileImage, nil
}
