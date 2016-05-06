package archive

import (
	"os"
	"path"
	"testing"
)

func TestExtractTarBall(t *testing.T) {

	var tarballData = []struct {
		Tarball     string // tarball file to test
		Dir         string // dir to extract the tarball file. Should be a different dir for each case
		ItemCreated string // relative path from Dir to an item created
		ExtractOK   bool   // whether the extraction command should succeed
	}{
		{"testdata/text.tar", "testdata/test1", "text", true},
		{"testdata/text.tar.gz", "testdata/test2", "", false},
		{"testdata/simple.tar", "testdata/test3", "loop-linux", true},
	}

	for _, tb := range tarballData {

		// if directory exists remove it and create it again
		if tb.Dir != "" {
			if err := os.RemoveAll(tb.Dir); err != nil {
				t.Errorf("Couldn't delete directory %q before actual testing: %s", tb.Dir, err)
				continue
			}

			if err := os.Mkdir(tb.Dir, 0777); err != nil {
				t.Errorf("Error creating test directory %q before testing: %s", tb.Dir, err)
				continue
			}

			defer func(dir string) {
				// delete test directory
				if dir != "" {
					if err := os.RemoveAll(dir); err != nil {
						t.Errorf("Couldn't delete directory %q after test: %s", dir, err)
					}
				}
			}(tb.Dir)
		}

		// extract tarball to test directory
		err := ExtractTarball(tb.Tarball, tb.Dir)
		if err != nil && tb.ExtractOK {
			t.Errorf("Error extracting tarball %q to %q: %s", tb.Tarball, tb.Dir, err)
			continue
		} else if err == nil && !tb.ExtractOK {
			t.Errorf("Extraction of tarball %q to %q should have failed, but did not.", tb.Tarball, tb.Dir)
			continue
		}

		// check extraction
		if tb.ItemCreated != "" {
			item := path.Join(tb.Dir, tb.ItemCreated)
			_, err = os.Stat(item)
			if err != nil {
				t.Errorf("Couldn't read extracted file %q: %s", item, err)
			}
		}

	}
}
