package client

import (
	"context"
	"fmt"
	"io"
	"time"

	"grpc-csv-viewer/internal/pkg/csvviewer"
	"grpc-csv-viewer/internal/pkg/logger"
)

// listFiles lists available CSV files.
func listFiles(client csvviewer.CSVViewerClient) (resultChan chan *csvviewer.FileDetails) {
	logger.Infof("Listing files")
	resultChan = make(chan *csvviewer.FileDetails)
	// Running goroutine to enable stream processing through channel.
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		stream, err := client.ListFiles(ctx, &csvviewer.FilesQuery{})
		if err != nil {
			logger.Fatalf("%v.ListFiles(_) = _, %v", client, err)
		}
		for {
			value, err := stream.Recv()
			if err == io.EOF {
				close(resultChan)

				return
			}
			if err != nil {
				// ToDo: maybe make resultChan an interface and send error to it?
				logger.Fatalf("%v.ListFiles(_) = _, %v", client, err)
			}
			resultChan <- value
		}
	}()

	return resultChan
}

// listValues lists all the values per given filter.
func listValues(client csvviewer.CSVViewerClient, filter *csvviewer.Filter) (resultChan chan *csvviewer.Value) {
	timeFrame := ""
	if filter.GetStartDate() != 0 && filter.GetStopDate() != 0 {
		timeFrame = fmt.Sprintf(" within time frame: [%q - %q]", filter.StartDate, filter.StopDate)
	}
	logger.Infof("Looking for %q file values%s", filter.GetFileName(), timeFrame)
	resultChan = make(chan *csvviewer.Value)
	// Running goroutine to enable stream processing through channel.
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		stream, err := client.ListValues(ctx, filter)
		if err != nil {
			logger.Fatalf("%v.ListValues(_) = _, %v", client, err)
		}
		for {
			value, err := stream.Recv()
			if err == io.EOF {
				close(resultChan)

				return
			}
			if err != nil {
				// ToDo: maybe make resultChan an interface and send error to it?
				logger.Fatalf("%v.ListValues(_) = _, %v", client, err)
			}
			resultChan <- value
		}
	}()

	return resultChan
}
