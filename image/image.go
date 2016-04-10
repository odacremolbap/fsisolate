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
	URLImage PathType = iota
	DirectoryImage
	FileImage
)

// PrepareImage load the image in the root path
// TODO check that root and image directories doesn't exists
// TODO if the image parameter is a directory, there's no need to call PrepareImage, anyway, we are supporting DirectoryImage
func PrepareImage(image string, root string) error {

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
		return 0, fmt.Errorf("Unknown image type. %s", err.Error())
	}

	// check directory
	if r.IsDir() {
		return DirectoryImage, nil
	}

	// TODO check for other types and don¡nt assume not being a URL or Dir, it's a file
	return FileImage, nil
}
