package client

import (
	"log"
	"sync"
	"time"

	"grpc-csv-viewer/internal/pkg/csvviewer"
	"grpc-csv-viewer/internal/pkg/logger"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type ConnectionConfig struct {
	// Global client connection configuration parameters.
	// Defaults to 30 seconds.
	IdleTimeout time.Duration

	// Keeps gRPC server address to connect to.
	// Defaults to 'localhost:8082'
	ServerAddr string

	// Chanel used to send termination signal to the gRPC client.
	// Mandatory for client initialization to ensure connection could be dropped at anytime.
	TerminationChan *chan bool
}

var (
	// keeps global gRPC client instance for idleTimeout time
	client csvviewer.CSVViewerClient

	// Global client connection configuration parameters.
	connectionConfig ConnectionConfig = ConnectionConfig{
		IdleTimeout: 30 * time.Second,
		ServerAddr:  "localhost:8082",
	}
	// Global status of connection keeping last communication time.
	// Used to terminate connection being kept idle for too long.
	lastCommunicationTime time.Time = time.Now()

	// Mutex to ensure only one connection initialization is possible at a time.
	mut sync.Mutex

	// Allows/disables connection establishing.
	enabled bool
)

// ConfigureConnectionParameters is used to configure gRPC connection parameters for singletone gRPC client.
func ConfigureConnectionParameters(cfg ConnectionConfig) error {
	if cfg.TerminationChan == nil {
		return errors.New("termination channel is missing in passed connection configuration")
	}
	connectionConfig.TerminationChan = cfg.TerminationChan
	if cfg.IdleTimeout != 0 {
		connectionConfig.IdleTimeout = cfg.IdleTimeout
	}
	if cfg.ServerAddr != "" {
		connectionConfig.ServerAddr = cfg.ServerAddr
	}
	enabled = true
	logger.Info("gRPC connection is configured, establishing first connection")

	return ensureConnection()
}

// ListValues lists values from gRPC server.
// ToDo: add proper values filter as input parameter.
func ListValues() []*csvviewer.Value {
	err := ensureConnection()
	if err != nil {
		logger.Error(errors.Wrap(err, "failed to ensure grpc connection").Error())
	}

	var result []*csvviewer.Value
	for v := range listValues(client, &csvviewer.Filter{}) {
		result = append(result, v)
	}

	return result
}

// PrintValues prints values from gRPC server.
func PrintValues() {
	err := ensureConnection()
	if err != nil {
		logger.Error(errors.Wrap(err, "failed to ensure grpc connection").Error())
	}

	for v := range listValues(client, &csvviewer.Filter{}) {
		log.Println(v)
	}
}

// ensureConnection establishes temporary connection to the gRPC server.
// Uses parameters from connectionConfig.
func ensureConnection() error {
	mut.Lock()
	defer mut.Unlock()
	if !enabled {
		return errors.New("connection is not allowed as termination signal was received.")
	}
	lastCommunicationTime = time.Now()
	if client != nil {
		return nil
	}
	threadWaitGroup := &sync.WaitGroup{}
	threadWaitGroup.Add(1)
	go func() {
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithInsecure()) // todo: drop once TLS will be enabled
		opts = append(opts, grpc.WithBlock())
		conn, err := grpc.Dial(connectionConfig.ServerAddr, opts...)
		if err != nil {
			logger.Fatalf(errors.Wrapf(err, "failed to dial grpc server %q", connectionConfig.ServerAddr).Error())
		}
		defer func() {
			e := conn.Close()
			if e != nil {
				logger.Errorf(errors.Wrap(e, "failed to close connection to grpc server").Error())
			}
		}()

		client = csvviewer.NewCSVViewerClient(conn)
		defer func() {
			// When this function returns, this means connection is dropped,
			// and so we need to nullify client to signalize about it.
			client = nil
		}()
		lastCommunicationTime = time.Now()
		// Send signal about connection establishment outside of goroutine.
		threadWaitGroup.Done()
		// Define ticker to check connection idle status every 1 second.
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		logger.Infof("initiating communication channel healthcheck loop", connectionConfig.ServerAddr)
		for {
			select {
			case <-*connectionConfig.TerminationChan:
				logger.Info("grpc client received termination signal, stopping communication channel")
				enabled = false

				return
			case <-ticker.C:
				logger.Debugf("grpc client healthheck, communication channel is idle for %s", time.Second*time.Duration(time.Now().Unix()-lastCommunicationTime.Unix()))
				if lastCommunicationTime.Add(connectionConfig.IdleTimeout).Unix() < time.Now().Unix() {
					logger.Infof("grpc client was idle too long (> %s), closing communication channel until further notice", connectionConfig.IdleTimeout)

					return
				}
			}
		}
	}()

	threadWaitGroup.Wait()
	logger.Infof("connection to grpc server %q established successfully, communications are allowed", connectionConfig.ServerAddr)

	return nil
}
