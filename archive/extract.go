package archive

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

// ExtractTarball extracts a tarball to a target directory
// Compressed formats are not supported
// TODO if extraction fails halfway, defer to delete remaining files
func ExtractTarball(tarball string, targetDir string) error {

	// check that target directory exists
	_, err := os.Stat(targetDir)
	if err != nil {
		return err
	}

	// check that tarball exists
	tbRead, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer tbRead.Close()

	tr := tar.NewReader(tbRead)
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

		path := filepath.Join(targetDir, header.Name)
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
