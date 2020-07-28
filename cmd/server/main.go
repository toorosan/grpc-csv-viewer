package main

import (
	"flag"
	"net"

	pb "grpc-csv-viewer/internal/pkg/csvviewer"
	"grpc-csv-viewer/internal/pkg/csvviewer/server"
	"grpc-csv-viewer/internal/pkg/logger"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func main() {
	var (
		csvFilesPath = flag.String("csv_files_path", "", "A path to the CSV files with TimeSeries.")
		listenAddr   = flag.String("listen_addr", ":8082", "The address to bind to listen for requests to gRPC Server in the format of host:port")
		logLevel     = flag.String("log_level", "info", "Severity level filter for logging messages.")
	)
	flag.Parse()

	logger.SetLevel(logger.LoggingLevel(*logLevel))

	lis, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		logger.Fatalf(errors.Wrapf(err, "failed to start listening gRPC requests on address %q", *listenAddr).Error())
	}
	logger.Infof("started listening for gRPC requests on address %q", *listenAddr)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterCSVViewerServer(grpcServer, server.NewCSVServer(*csvFilesPath))

	err = grpcServer.Serve(lis)
	if err != nil {
		logger.Fatal(errors.Wrap(err, "failed to start serving gRPC requests").Error())
	}
}
