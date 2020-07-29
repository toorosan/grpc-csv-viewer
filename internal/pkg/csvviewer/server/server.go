package server

import (
	"context"
	"path/filepath"

	"grpc-csv-viewer/internal/pkg/csvreader"
	"grpc-csv-viewer/internal/pkg/csvviewer"
	"grpc-csv-viewer/internal/pkg/logger"
	"grpc-csv-viewer/internal/pkg/pathwalker"

	"github.com/pkg/errors"
)

// NewCSVServer creates new CSV viewer gRPC server.
func NewCSVServer(csvFilesPath string) csvviewer.CSVViewerServer {
	s := &csvServer{
		csvFilesPath: csvFilesPath,
	}
	s.availableCSVFiles = map[string]*fileDetailsWithTimeSeries{}
	if s.csvPathIsSet() {
		csvFiles, err := pathwalker.ListFilesInDir(s.csvFilesPath, ".csv")
		if err != nil {
			logger.Fatalf(errors.Wrapf(err, "failed to list files in directory %q", s.csvFilesPath).Error())
		}
		if len(csvFiles) == 0 {
			logger.Fatalf("failed to initialize server: no suitable csv files found in directory %q", s.csvFilesPath)
		}
		s.firstFileName = csvFiles[0]
		for _, f := range csvFiles {
			fd, err := s.gatherFileDetails(f)
			if err != nil {
				logger.Fatalf(errors.Wrapf(err, "failed to gather file %q details", f).Error())
			}
			s.availableCSVFiles[fd.FileName] = fd
			s.availableCSVFilesList = append(s.availableCSVFilesList, fd)
		}
	} else {
		s.firstFileName = mockedCSVFileName
		fd, err := s.gatherFileDetails("")
		if err != nil {
			logger.Fatalf(errors.Wrapf(err, "failed to gather file %q details", mockedCSVFileName).Error())
		}
		s.availableCSVFiles = map[string]*fileDetailsWithTimeSeries{
			mockedCSVFileName: fd,
		}
	}

	return s
}

type fileDetailsWithTimeSeries struct {
	*csvviewer.FileDetails

	// Values are loaded for now at the service initialization, aka eager initialization.
	// provides better performance on small amounts of data, but requires a lot of memory.
	// Consider using lazy initialization and even access to the dataset by cursor
	// if bigger files are required to be processed.
	values []*csvviewer.Value
}

type csvServer struct {
	availableCSVFiles     map[string]*fileDetailsWithTimeSeries
	availableCSVFilesList []*fileDetailsWithTimeSeries
	csvFilesPath          string
	firstFileName         string
}

func (s *csvServer) ListFiles(query *csvviewer.FilesQuery, stream csvviewer.CSVViewer_ListFilesServer) error {
	logger.Debugf("requested ListFiles")
	for _, f := range s.availableCSVFilesList {
		logger.Debugf("sending FileDetails %#v", f.FileDetails)
		err := stream.Send(f.FileDetails)
		if err != nil {
			return errors.Wrapf(err, "failed to send value %#v to the stream", f.FileDetails)
		}
	}

	return nil
}

func (s *csvServer) ListValues(filter *csvviewer.Filter, stream csvviewer.CSVViewer_ListValuesServer) error {
	if filter == nil {
		filter = &csvviewer.Filter{}
	}
	logger.Debugf("requested ListValues with the filter %#v", filter)
	if filter.FileName == "" {
		filter.FileName = s.firstFileName
	}
	if filter.StartDate == 0 {
		filter.StartDate = s.availableCSVFiles[s.firstFileName].StartDate
	}
	if filter.StopDate == 0 {
		filter.StopDate = s.availableCSVFiles[s.firstFileName].StopDate
	}
	for _, value := range s.csvValuesFromFile(filter.FileName) {
		if inRange(value, filter) {
			logger.Debugf("sending Value %#v", value)
			err := stream.Send(value)
			if err != nil {
				return errors.Wrapf(err, "failed to send value %#v to the stream", value)
			}
		}
	}

	return nil
}

func (s *csvServer) GetFileDetails(ctx context.Context, query *csvviewer.FileQuery) (*csvviewer.FileDetails, error) {
	logger.Debugf("requested GetFileDetails with the query %#v", query)
	if query.GetFileName() == "" {
		return s.availableCSVFiles[s.firstFileName].FileDetails, nil
	}

	if s.availableCSVFiles[query.GetFileName()] != nil {
		return s.availableCSVFiles[query.GetFileName()].FileDetails, nil
	}

	return nil, nil
}

func (s *csvServer) csvPathIsSet() bool {
	return s.csvFilesPath != ""
}

func (s *csvServer) csvValuesFromFile(fileName string) []*csvviewer.Value {
	if s.availableCSVFiles[fileName] != nil {
		return s.availableCSVFiles[fileName].values
	}

	return nil
}

func (s *csvServer) gatherFileDetails(baseFileName string) (*fileDetailsWithTimeSeries, error) {
	fd := fileDetailsWithTimeSeries{
		FileDetails: &csvviewer.FileDetails{
			FileName:  baseFileName,
			StartDate: 99999999999,
			StopDate:  -1,
		},
	}
	if baseFileName == "" {
		fd.FileName = mockedCSVFileName
		fd.values = mockValues(mockedCSVFileName)
	} else {
		vv, err := csvreader.ReadTimeSeriesFromCSV(filepath.Join(s.csvFilesPath, baseFileName))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to process file %q: ", baseFileName)
		}
		for i := range vv.TimeSeries {
			fd.values = append(fd.values, &csvviewer.Value{Date: vv.TimeSeries[i].Date.Unix(), Value: vv.TimeSeries[i].Value})
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

func inRange(value *csvviewer.Value, filter *csvviewer.Filter) bool {
	return value.Date < filter.StopDate && value.Date > filter.StartDate
}
