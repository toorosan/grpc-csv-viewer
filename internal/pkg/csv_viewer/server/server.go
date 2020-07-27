package server

import (
	"context"
	"log"
	"path/filepath"

	"grpc-csv-viewer/internal/pkg/csv_reader"
	"grpc-csv-viewer/internal/pkg/csv_viewer"
	"grpc-csv-viewer/internal/pkg/path_walker"

	"github.com/pkg/errors"
)

// NewCSVServer creates new CSV viewer gRPC server.
func NewCSVServer(csvFilesPath string) *csvServer {
	s := &csvServer{
		csvFilesPath: csvFilesPath,
	}
	if s.csvPathIsSet() {
		csvFiles, err := path_walker.ListFilesInDir(s.csvFilesPath, ".csv")
		if err != nil {
			log.Fatalf(errors.Wrapf(err, "failed to list files in directory %q", s.csvFilesPath).Error())
		}
		if len(csvFiles) == 0 {
			log.Fatalf("failed to initialize server: no suitable csv files found in directory %q", s.csvFilesPath)
		}
		s.defaultFileName = csvFiles[0]
		for _, f := range csvFiles {
			fd, err := s.gatherFileDetails(f)
			if err != nil {
				log.Fatalf(errors.Wrapf(err, "failed to gather file %q details", f).Error())
			}
			s.availableCSVFiles = map[string]*fileDetailsWithTimeSeries{
				fd.FileName: fd,
			}
		}
	} else {
		s.defaultFileName = mockedCSVFileName
		fd, err := s.gatherFileDetails("")
		if err != nil {
			log.Fatalf(errors.Wrapf(err, "failed to gather file %q details", mockedCSVFileName).Error())
		}
		s.availableCSVFiles = map[string]*fileDetailsWithTimeSeries{
			mockedCSVFileName: fd,
		}
	}

	return s
}

type fileDetailsWithTimeSeries struct {
	*csv_viewer.FileDetails

	// Values are loaded for now at the service initialization, aka eager initialization.
	// provides better performance on small amounts of data, but requires a lot of memory.
	// Consider using lazy initialization and even access to the dataset by cursor
	// if bigger files are required to be processed.
	values []*csv_viewer.Value
}

type csvServer struct {
	availableCSVFiles map[string]*fileDetailsWithTimeSeries
	csvFilesPath      string
	defaultFileName   string
}

func (s *csvServer) ListValues(filter *csv_viewer.Filter, stream csv_viewer.CSVViewer_ListValuesServer) error {
	if filter == nil {
		filter = &csv_viewer.Filter{}
	}
	if filter.FileName == "" {
		filter.FileName = s.defaultFileName
	}
	if filter.StartDate == 0 {
		filter.StartDate = s.availableCSVFiles[s.defaultFileName].StartDate
	}
	if filter.StopDate == 0 {
		filter.StopDate = s.availableCSVFiles[s.defaultFileName].StopDate
	}
	for _, value := range s.csvValuesFromFile(filter.FileName) {
		if inRange(value, filter) {
			err := stream.Send(value)
			if err != nil {
				return errors.Wrapf(err, "failed to send value %#v to the stream", value)
			}
		}
	}

	return nil
}

func (s *csvServer) GetFileDetails(ctx context.Context, query *csv_viewer.FileQuery) (*csv_viewer.FileDetails, error) {
	if query.FileName == "" {
		return s.availableCSVFiles[s.defaultFileName].FileDetails, nil
	}

	if s.availableCSVFiles[query.FileName] != nil {
		return s.availableCSVFiles[query.FileName].FileDetails, nil
	}

	return nil, nil
}

func (s *csvServer) csvPathIsSet() bool {
	return s.csvFilesPath != ""
}

func (s *csvServer) csvValuesFromFile(fileName string) []*csv_viewer.Value {
	if s.availableCSVFiles[fileName] != nil {
		return s.availableCSVFiles[fileName].values
	}

	return nil
}

func inRange(value *csv_viewer.Value, filter *csv_viewer.Filter) bool {
	return value.Date < filter.StopDate && value.Date > filter.StartDate
}

func (s *csvServer) gatherFileDetails(baseFileName string) (*fileDetailsWithTimeSeries, error) {
	fd := fileDetailsWithTimeSeries{
		FileDetails: &csv_viewer.FileDetails{
			FileName:  baseFileName,
			StartDate: 9999999999,
			StopDate:  -1,
		},
	}
	if baseFileName == "" {
		fd.FileName = mockedCSVFileName
		fd.values = mockValues(mockedCSVFileName)
	} else {
		vv, err := csv_reader.ReadTimeSeriesFromCSV(filepath.Join(s.csvFilesPath, baseFileName))
		if err != nil {
			log.Fatalf(errors.Wrapf(err, "failed to process file %q: ", baseFileName).Error())
		}
		for i := range vv.TimeSeries {
			fd.values = append(fd.values, &csv_viewer.Value{Date: vv.TimeSeries[i].Date.Unix(), Value: vv.TimeSeries[i].Value})
		}
	}
	for _, v := range fd.values {
		if v.Date > fd.StopDate {
			fd.StopDate = v.Date
		}
		if v.Date < fd.StartDate {
			fd.StartDate = v.Date
		}
	}

	return &fd, nil
}
