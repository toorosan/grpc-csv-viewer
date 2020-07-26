package main

import (
	"log"

	"grpc-csv-viewer/internal/app/roles/client/rest"

	"github.com/pkg/errors"
)

func main() {
	err := rest.ListenAndServeREST()
	if err != nil {
		log.Fatal(errors.Wrap(err,"failed to listen and serve REST"))
	}
}
