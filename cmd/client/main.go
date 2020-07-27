package main

import (
	"grpc-csv-viewer/internal/app/roles/client/rest"
	"grpc-csv-viewer/internal/pkg/logger"

	"github.com/pkg/errors"
)

func main() {
	err := rest.ListenAndServeREST()
	if err != nil {
		logger.Fatal(errors.Wrap(err, "failed to listen and serve REST").Error())
	}
}
