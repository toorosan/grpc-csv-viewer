package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"grpc-csv-viewer/internal/pkg/logger"
)

func uiHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "Hi there, I am future ui handler!")
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func timeSeriesHandler(w http.ResponseWriter, _ *http.Request) {
	mockTimeSeries := mockTimeSeries()
	err := json.NewEncoder(w).Encode(mockTimeSeries)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func registerHandlers() {
	http.HandleFunc("/api/v1/timeseries", timeSeriesHandler)
	http.HandleFunc("/", uiHandler)
}
