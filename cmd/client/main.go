package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"grpc-csv-viewer/internal/app/roles/client"
)

func uiHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hi there, I am future ui handler!")
	if err != nil {
		log.Fatal(err)
	}
}

func timeSeriesHandler(w http.ResponseWriter, r *http.Request) {
	mockTimeSeries := client.TimeSeries{
		FileName:  "mocked-values.csv",
		StartDate: time.Now().Add(-time.Hour * 3),
		StopDate:  time.Now(),
		Values: []client.SeriesItem{
			{
				Date:  time.Now().Add(-time.Hour * 3),
				Value: 1,
			},
			{
				Date:  time.Now().Add(-time.Hour * 2),
				Value: 2,
			},
			{
				Date:  time.Now().Add(-time.Hour * 1),
				Value: 4,
			},
			{
				Date:  time.Now(),
				Value: 1,
			},
		},
	}
	err := json.NewEncoder(w).Encode(mockTimeSeries)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/api/timeseries", timeSeriesHandler)
	http.HandleFunc("/", uiHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
