package fsisolate

import (
	"github.com/odacremolbap/fsisolate/image"
	"github.com/odacremolbap/fsisolate/runtime"
)

// Prepare prepares the filesystem structure to start a chrooted execution
func Prepare(imagePath string, root string) (*runtime.ChrootedProcess, error) {

	// get image into root, and return the new root directory
	root, err := image.PrepareImage(imagePath, root)
	if err != nil {
		return nil, err
	}

	// create the chroot process structure
	return runtime.NewChrootProcess(root), nil
}
