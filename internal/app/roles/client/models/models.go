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

// TimeSeries defines data structure to be replied to UI per TimeSeries request.
type TimeSeries struct {
	FileName  string       `json:"fileName"`
	StartDate time.Time    `json:"startDate"`
	StopDate  time.Time    `json:"stopDate"`
	Values    []SeriesItem `json:"values"`
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
