package main

import (
	"flag"
	"fmt"
	"net"

	pb "grpc-csv-viewer/internal/pkg/csv_viewer"
	"grpc-csv-viewer/internal/pkg/csv_viewer/server"
	"grpc-csv-viewer/internal/pkg/logger"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func main() {
	var (
		csvFilesPath = flag.String("csv_files_path", "", "A path to the CSV files with TimeSeries.")
		port         = flag.Int("port", 10000, "The server port")
		logLevel     = flag.String("log_level", "info", "The server port")
	)
	flag.Parse()

	logger.SetLevel(logger.LoggingLevel(*logLevel))

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		logger.Fatalf(errors.Wrapf(err, "failed to start listeninig gRPC requests on port %d", port).Error())
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterCSVViewerServer(grpcServer, server.NewCSVServer(*csvFilesPath))

	err = grpcServer.Serve(lis)
	if err != nil {
		logger.Fatal(errors.Wrap(err, "failed to start serving gRPC requests").Error())
	}
}
