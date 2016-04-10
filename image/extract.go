package image

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// extractImage retrieves a image from internet
// TODO not all tar formats supported (zipped, ...)
// TODO if extraction fails halfway, defer to delete remaining files
func extractImage(image string, root string) error {
	log.Debugf("decompressing image %s \ninto %s", image, root)

	imgRead, err := os.Open(image)
	if err != nil {
		return err
	}
	defer imgRead.Close()

	tr := tar.NewReader(imgRead)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			// exit loop at end of stream
			break
		}

		// extract error
		if err != nil {
			return err
		}

		path := filepath.Join(root, header.Name)
		fi := header.FileInfo()

		// restore dir
		if fi.IsDir() {
			if err = os.MkdirAll(path, fi.Mode()); err != nil {
				return err
			}
			continue
		}

		// restore file
		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fi.Mode())
		if err != nil {
			return err
		}

		// write file contents
		defer file.Close()
		_, err = io.Copy(file, tr)
		if err != nil {
			return err
		}

	}
	return nil
}
