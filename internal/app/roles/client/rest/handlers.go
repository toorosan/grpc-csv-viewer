package rest

import (
	"encoding/json"
	"net/http"

	"grpc-csv-viewer/internal/app/roles/client/models"
	"grpc-csv-viewer/internal/pkg/csvviewer/client"
	"grpc-csv-viewer/internal/pkg/logger"
)

func listFilesHandler(w http.ResponseWriter, _ *http.Request) {
	rawFiles := client.ListFiles()

	err := json.NewEncoder(w).Encode(models.FileDetailsFromGRPC(rawFiles))
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func timeSeriesHandler(w http.ResponseWriter, _ *http.Request) {
	rawTimeSeries := client.ListValues()

	err := json.NewEncoder(w).Encode(models.TimeSeriesFromRawValues(rawTimeSeries))
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func registerHandlers() {
	v1BasePath := "/api/v1"
	http.HandleFunc(v1BasePath+"/timeseries", timeSeriesHandler)
	http.HandleFunc(v1BasePath+"/files", listFilesHandler)
}
