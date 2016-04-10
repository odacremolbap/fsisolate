package image

// PathType path format type
type PathType uint8

// Path types supported
const (
	URLImage PathType = iota
	DirectoryImage
	FileImage
)

// PrepareImage load the image in the root path
func PrepareImage(image string, root string) error {

	return nil
}
