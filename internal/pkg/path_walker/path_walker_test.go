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
			"1.csv",
			"2.csv",
			"3.go",
			"4.xls",
			"5.txt",
		},
		{
			"5.txt",
		},
		{
			"1.csv",
			"2.csv",
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
