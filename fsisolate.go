package fsisolate

// Prepare prepares the filesystem structure to start a chrooted execution
func Prepare(imagePath string, root string) (*ChrootedProcess, error) {

	// use default values
	img := Image{}

	// prepare, download if URL
	// root returns the new root where the image is going to be executed
	realRoot, err := img.Prepare(imagePath, root)
	if err != nil {
		return nil, err
	}

	// get image into root, and return the new root directory
	// root, err = image.PrepareImage(imagePath, root)
	// if err != nil {
	// 	return nil, err
	// }

	// create the chroot process structure
	return NewChrootProcess(realRoot), nil
}
