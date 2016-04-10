package fsisolate

import (
	"github.com/odacremolbap/fsisolate/image"
	"github.com/odacremolbap/fsisolate/runtime"
	log "github.com/sirupsen/logrus"
)

// PrepareChrootedProcess prepares the filesystem structure to start a chrooted execution
func PrepareChrootedProcess(imagePath string, root string, exec string, args []string) (*runtime.ChrootedProcess, error) {
	log.Debugf(`preparing new chroot isolated environment
		image:     %s
		fs root:   %s
		command:   %s
		arguments: %s`, imagePath, root, exec, args)

	// get image into root, and return the new root directory
	root, err := image.PrepareImage(imagePath, root)
	if err != nil {
		return nil, err
	}

	// create the chroot process structure
	return runtime.NewChrootProcess(root), nil
}
