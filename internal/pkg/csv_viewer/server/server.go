package server

import (
	"context"
	"flag"
	"log"

	"grpc-csv-viewer/internal/pkg/csv_viewer"
	"grpc-csv-viewer/internal/pkg/path_walker"

	"github.com/pkg/errors"
)

// NewCSVServer creates new CSV viewer gRPC server.
func NewCSVServer() *csvServer {
	var (
		csvFilesPath = flag.String("csv_files_path", "", "A path to the CSV files with TimeSeries.")
		port         = flag.Int("port", 10000, "The server port")
	)
	flag.Parse()
	s := &csvServer{
		csvFilesPath: *csvFilesPath,
		port:         *port,
	}
	if s.csvPathIsSet() {
		csvFiles, err := path_walker.ListFilesInDir(*csvFilesPath, ".csv")
		if err != nil {
			log.Fatalf(errors.Wrapf(err, "failed to list files in directory %q", *csvFilesPath).Error())
		}
		if len(csvFiles) == 0 {
			log.Fatalf("failed to initialize server: no suitable csv files found in directory %q", csvFilesPath)
		}
		s.defaultFileName = csvFiles[0]
		for _, f := range csvFiles {
			fd, err := gatherFileDetails(f)
			if err != nil {
				log.Fatalf(errors.Wrapf(err, "failed to gather file %q details", mockedCSVFileName).Error())
			}
			s.availableCSVFiles = map[string]*csv_viewer.FileDetails{
				mockedCSVFileName: fd,
			}
		}
	} else {
		s.defaultFileName = mockedCSVFileName
		fd, err := gatherFileDetails("")
		if err != nil {
			log.Fatalf(errors.Wrapf(err, "failed to gather file %q details", mockedCSVFileName).Error())
		}
		s.availableCSVFiles = map[string]*csv_viewer.FileDetails{
			mockedCSVFileName: fd,
		}
	}

	return s
}

type csvServer struct {
	availableCSVFiles map[string]*csv_viewer.FileDetails
	defaultFileName   string
	port              int
	csvFilesPath      string
}

func (s *csvServer) csvPathIsSet() bool {
	return s.csvFilesPath != ""
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
	if query.FileName != "" {
		return s.availableCSVFiles[query.FileName], nil
	} else {
		return s.availableCSVFiles[s.defaultFileName], nil
	}
}

func (s *csvServer) csvValuesFromFile(fileName string) []*csv_viewer.Value {
	return mockValues(fileName)
}

func inRange(value *csv_viewer.Value, filter *csv_viewer.Filter) bool {
	return value.Date < filter.StopDate && value.Date > filter.StartDate
}

func gatherFileDetails(fileName string) (*csv_viewer.FileDetails, error) {
	fd := csv_viewer.FileDetails{
		FileName:  fileName,
		StartDate: 99999999,
		StopDate:  -1,
	}
	if fileName == "" {
		vv := mockValues(mockedCSVFileName)
		for _, v := range vv {
			if v.Date > fd.StopDate {
				fd.StopDate = v.Date
			}
			if v.Date < fd.StartDate {
				fd.StartDate = v.Date
			}
		}
	} else {
		log.Fatalf("failed to process file %q: only mocked values currently supported", fileName)
	}

	return &fd, nil
}
