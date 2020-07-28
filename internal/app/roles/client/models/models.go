package models

import (
	"math"
	"time"

	"grpc-csv-viewer/internal/pkg/csvviewer"
)

// SeriesItem defines data type for a single Time Series value.
type SeriesItem struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}

// FileDetails defines data structure to be replied to UI per request to list available files.
type FileDetails struct {
	FileName  string    `json:"fileName"`
	StartDate time.Time `json:"startDate"`
	StopDate  time.Time `json:"stopDate"`
}

// FileDetailsFromGRPC prepares UI-compatible FileDetails object from raw gRPC response.
func FileDetailsFromGRPC(rawFiles []*csvviewer.FileDetails) []FileDetails {
	fd := make([]FileDetails, len(rawFiles))
	for i := range rawFiles {
		fd[i] = FileDetails{
			FileName:  rawFiles[i].FileName,
			StartDate: time.Unix(rawFiles[i].StartDate, 0),
			StopDate:  time.Unix(rawFiles[i].StopDate, 0),
		}
	}

	return fd
}

// TimeSeries defines data structure to be replied to UI per TimeSeries request.
type TimeSeries struct {
	FileDetails
	Values []SeriesItem `json:"values"`
}

// TimeSeriesFromRawValues prepares UI-compatible TimeSeries object from raw gRPC values.
func TimeSeriesFromRawValues(values []*csvviewer.Value) TimeSeries {
	ts := TimeSeries{
		Values: make([]SeriesItem, len(values)),
	}
	for i := range values {
		// NaN and Inf values processing.
		// ToDo: check if we can skip them at all, as those values look like broken ones.
		v := values[i].Value
		switch {
		case math.IsNaN(v):
			v = 0
		case math.IsInf(v, 0):
			v = 0
		}
		ts.Values[i] = SeriesItem{
			Date:  time.Unix(values[i].Date, 0),
			Value: v,
		}
	}

	return ts
}
