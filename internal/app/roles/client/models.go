package client

import (
	"time"
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
