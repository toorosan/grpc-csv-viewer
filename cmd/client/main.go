package main

import (
	"flag"
	"time"

	"grpc-csv-viewer/internal/app/roles/client/rest"
	"grpc-csv-viewer/internal/pkg/csvviewer/client"
	"grpc-csv-viewer/internal/pkg/logger"

	"github.com/pkg/errors"
)

func main() {
	var (
		err            error
		idleTimeout    = flag.Duration("idle_timeout", time.Second*30, "The max amount of time without communication to keep connection up.")
		logLevel       = flag.String("log_level", "info", "Severity level filter for logging messages.")
		restListenAddr = flag.String("rest_listen_addr", ":8081", "The address to bind to listen for requests to gRPC Client REST API in the format of host:port")
		serverAddr     = flag.String("server_addr", "localhost:8082", "The address of gRPC Server to connect to in the format of host:port")
	)
	flag.Parse()
	logger.SetLevel(logger.LoggingLevel(*logLevel))

	logger.Info("configuring gRPC connection")
	terminationChan := make(chan bool, 1)
	err = client.ConfigureConnectionParameters(client.ConnectionConfig{
		IdleTimeout:     *idleTimeout,
		ServerAddr:      *serverAddr,
		TerminationChan: &terminationChan,
	})
	if err != nil {
		logger.Fatal(errors.Wrap(err, "failed to configure gRPC connection").Error())
	}

	err = rest.ListenAndServeREST(*restListenAddr)
	if err != nil {
		logger.Fatal(errors.Wrap(err, "failed to listen and serve REST").Error())
	}
	close(terminationChan)
}
