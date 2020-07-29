package server

import (
	"context"
	"testing"

	"grpc-csv-viewer/internal/pkg/csvviewer"
)

func TestNewCSVServerNoConfig(t *testing.T) {
	// given
	// no option flags were passed
	// when
	s := NewCSVServer("").(*csvServer)

	// then
	if len(s.availableCSVFiles) != 1 {
		t.Fatalf("failed to verify only mocked cvs file is configured")
	}
	if s.availableCSVFiles[mockedCSVFileName] == nil {
		t.Fatalf("failed to verify mocked cvs file is configured properly")
	}
	if s.availableCSVFiles[mockedCSVFileName].FileName != mockedCSVFileName {
		t.Fatalf("failed to verify only mocked cvs file is configured")
	}
	minDate := int64(99999999999)
	maxDate := int64(-1)
	for _, v := range mockedValuesCache {
		if v.Date < minDate {
			minDate = v.Date
		}
		if v.Date > maxDate {
			maxDate = v.Date
		}
	}
	if s.firstFileName != mockedCSVFileName {
		t.Fatalf("failed to verify mocked cvs file is set as default one if no other files are configured")
	}
	if s.availableCSVFiles[mockedCSVFileName].StartDate != minDate {
		t.Fatalf("failed to verify mocked cvs file is properly configured: start date expected %d, got %d", minDate, s.availableCSVFiles[mockedCSVFileName].StartDate)
	}
	if s.availableCSVFiles[mockedCSVFileName].StopDate != maxDate {
		t.Fatalf("failed to verify mocked cvs file is properly configured: stop date expected %d, got %d", maxDate, s.availableCSVFiles[mockedCSVFileName].StopDate)
	}
	if len(s.availableCSVFiles[mockedCSVFileName].values) != len(mockedValuesCache) {
		t.Fatalf("failed to verify mocked cvs file is properly configured: size of time series expected %d, got %d", len(mockedValuesCache), len(s.availableCSVFiles[mockedCSVFileName].values))
	}

	// when
	// request existing file details
	fileDetails, err := s.GetFileDetails(context.Background(), &csvviewer.FileQuery{})
	if err != nil {
		t.Fatalf("failed to query file details for default file, got: %v", err)
	}
	if fileDetails.FileName != mockedCSVFileName {
		t.Fatal("failed to validate default returned file is the mocked one in case  if no other files are configured")
	}

	// when
	// request non-existing file details
	fileDetails, err = s.GetFileDetails(context.Background(), &csvviewer.FileQuery{FileName: "some-non-existing-file"})
	if err != nil {
		t.Fatalf("failed to query non-existing file details, got: %v", err)
	}
	if fileDetails != nil {
		t.Fatal("failed to validate empty details returned for non-existing file name")
	}
}

func TestNewCSVServerConfiguredRealCSVFiles(t *testing.T) {
	// given
	expectedFileName := "3.csv"
	cvsFilesPath := "./test_data"
	// when
	s := NewCSVServer(cvsFilesPath).(*csvServer)

	// then
	if len(s.availableCSVFiles) != 1 {
		t.Fatalf("failed to verify real cvs file is configured")
	}
	if s.availableCSVFiles[mockedCSVFileName] != nil {
		t.Fatalf("failed to verify mocked cvs file is configured properly")
	}
	if s.availableCSVFiles[expectedFileName].FileName != expectedFileName {
		t.Fatalf("failed to verify real cvs file is configured properly")
	}
	if s.firstFileName != expectedFileName {
		t.Fatalf("failed to verify real cvs file is set as default one if no other files are configured")
	}
	// Check file `test_data/3.csv` for more details about the following 3 validations below:
	minDate := int64(1546301700)
	if s.availableCSVFiles[expectedFileName].StartDate != minDate {
		t.Fatalf("failed to verify real cvs file is properly configured: start date expected %d, got %d", minDate, s.availableCSVFiles[expectedFileName].StartDate)
	}
	maxDate := int64(1546317000)
	if s.availableCSVFiles[expectedFileName].StopDate != maxDate {
		t.Fatalf("failed to verify real cvs file is properly configured: stop date expected %d, got %d", maxDate, s.availableCSVFiles[expectedFileName].StopDate)
	}
	if len(s.availableCSVFiles[expectedFileName].values) != 18 {
		t.Fatalf("failed to verify real cvs file is properly configured: size of time series expected %d, got %d", len(mockedValuesCache), len(s.availableCSVFiles[expectedFileName].values))
	}

	// when
	// request existing file details
	fileDetails, err := s.GetFileDetails(context.Background(), &csvviewer.FileQuery{})
	if err != nil {
		t.Fatalf("failed to query file details for default file, got: %v", err)
	}
	if fileDetails.GetFileName() != expectedFileName {
		t.Fatal("failed to validate default returned file is the default real one in case if no other files are configured")
	}

	// when
	// request non-existing file details
	fileDetails, err = s.GetFileDetails(context.Background(), &csvviewer.FileQuery{FileName: "some-non-existing-file"})
	if err != nil {
		t.Fatalf("failed to query non-existing file details, got: %v", err)
	}
	if fileDetails != nil {
		t.Fatal("failed to validate empty details returned for non-existing file name")
	}
}
