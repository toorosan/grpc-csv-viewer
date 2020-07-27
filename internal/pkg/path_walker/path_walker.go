package path_walker

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/pkg/errors"
)

type FileList []string

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
func ListFilesInDir(rootPath, ext string) ([]string, error) {
	var files []string
	rootPathBase := filepath.Base(rootPath)
	walkerFunc := func(file string, info os.FileInfo, err error) error {
		filepath.Abs(file)
		if (ext == "" || filepath.Ext(file) == ext) && (filepath.Base(file) != rootPathBase) {
			files = append(files, file)
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
