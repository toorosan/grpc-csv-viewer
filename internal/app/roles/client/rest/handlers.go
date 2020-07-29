package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"grpc-csv-viewer/internal/app/roles/client/models"
	"grpc-csv-viewer/internal/pkg/csvviewer"
	"grpc-csv-viewer/internal/pkg/csvviewer/client"
	"grpc-csv-viewer/internal/pkg/logger"
)

func getQueryParameter(r *http.Request, name string) string {
	pp, ok := r.URL.Query()[name]

	if !ok || len(pp[0]) == 0 {
		return ""
	}

	return pp[0]
}

func getInt64QueryParameter(r *http.Request, name string) (result int64) {
	strValue := getQueryParameter(r, name)
	if strValue != "" {
		var err error
		result, err = strconv.ParseInt(strValue, 10, 64)
		if err != nil {
			logger.Fatalf("failed to convert %q parameter value %q to int64", name, strValue)
		}
	}

	return result
}

func listFilesHandler(w http.ResponseWriter, _ *http.Request) {
	rawFiles := client.ListFiles()

	err := json.NewEncoder(w).Encode(models.FileDetailsFromGRPC(rawFiles))
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func timeSeriesHandler(w http.ResponseWriter, r *http.Request) {
	filter := &csvviewer.Filter{
		FileName:  getQueryParameter(r, "fileName"),
		StartDate: getInt64QueryParameter(r, "startDate"),
		StopDate:  getInt64QueryParameter(r, "stopDate"),
	}
	rawTimeSeries := client.ListValues(filter)
	response := models.TimeSeriesFromRawValues(rawTimeSeries)
	if filter.GetFileName() != "" {
		response.FileName = filter.GetFileName()
	}
	if filter.GetStartDate() != 0 {
		response.StartDate = time.Unix(filter.GetStartDate(), 0)
	}
	if filter.GetStopDate() != 0 {
		response.StopDate = time.Unix(filter.GetStopDate(), 0)
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func registerHandlers() {
	v1BasePath := "/api/v1"
	http.HandleFunc(v1BasePath+"/timeseries", timeSeriesHandler)
	http.HandleFunc(v1BasePath+"/files", listFilesHandler)
}
