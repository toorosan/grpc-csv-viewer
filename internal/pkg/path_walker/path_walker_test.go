package path_walker

import (
	"testing"
)

func TestListFilesInDir(t *testing.T) {
	// given
	filesPath := "./test_data"
	extArrays := []string{
		"",
		".txt",
		".csv",
		".xyz",
		"go",
	}

	expectedFiles := []FileList{
		{
			"test_data/1.csv",
			"test_data/2.csv",
			"test_data/3.go",
			"test_data/4.xls",
			"test_data/5.txt",
		},
		{
			"test_data/5.txt",
		},
		{
			"test_data/1.csv",
			"test_data/2.csv",
		},
		{},
		{},
	}

	var ff []string
	var err error
	for i := range extArrays {
		// when
		ff, err = ListFilesInDir(filesPath, extArrays[i])
		if err != nil {
			t.Fatalf("failed to gather list of files from dir %q and for extension %q: %v", filesPath, extArrays[i], err)
		}

		// then
		err = expectedFiles[i].EqualTo(ff)
		if err != nil {
			t.Fatalf("failed to validate list of files in dir %q and for extension %q: expected %v, got %v. Error: %v", filesPath, extArrays[i], expectedFiles[i], ff, err)
		}
	}
}
