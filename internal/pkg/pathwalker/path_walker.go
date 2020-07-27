package pathwalker

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/pkg/errors"
)

// FileList extends list of strings with handy methods to allow check if two lists are equal or not.
type FileList []string

// EqualTo allows to check whether passed File List is equal to ours or not.
// Returns error with basic diff details or nil if lists are identical.
func (fl FileList) EqualTo(other FileList) error {
	if len(fl) != len(other) {
		return errors.Errorf("size of arrays is different: expected %d, got %d", len(fl), len(other))
	}
	for i := range fl {
		if fl[i] != other[i] {
			return errors.Errorf("file name in same position %d is different: expected %q, got %q", i, fl[i], other[i])
		}
	}

	return nil
}

// ListFilesInDir implements a handy shortcut to list files in folder, filtering them by extension, if passed.
// Returns list of base file names.
func ListFilesInDir(rootPath, ext string) ([]string, error) {
	var files []string
	rootBasePath := filepath.Base(rootPath)
	walkerFunc := func(file string, info os.FileInfo, err error) error {
		fileBasePath := filepath.Base(file)
		if (ext == "" || filepath.Ext(file) == ext) && (fileBasePath != rootBasePath) {
			files = append(files, fileBasePath)
		}

		return nil
	}

	err := filepath.Walk(rootPath, walkerFunc)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list files in directory %q", rootPath)
	}
	sort.Strings(files)

	return files, nil
}
