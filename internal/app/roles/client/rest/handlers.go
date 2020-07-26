package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func uiHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "Hi there, I am future ui handler!")
	if err != nil {
		log.Fatal(err)
	}
}

func timeSeriesHandler(w http.ResponseWriter, _ *http.Request) {
	mockTimeSeries := mockTimeSeries()
	err := json.NewEncoder(w).Encode(mockTimeSeries)
	if err != nil {
		log.Fatal(err)
	}
}

func registerHandlers() {
	http.HandleFunc("/api/v1/timeseries", timeSeriesHandler)
	http.HandleFunc("/", uiHandler)
}
