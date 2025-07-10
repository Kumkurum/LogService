package grpc_client

import (
	ls "LogService/internal/transport"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"log"
	"os"
	"sync"
	"time"
)

type LoggingClient struct {
	ls.LoggingServiceClient
	nameService string
	socketPath  string //"unix:///tmp/grpc.sock"
	conn        *grpc.ClientConn
	mutex       sync.Mutex
}

func NewLoggingClient(socketPath, nameService string) (*LoggingClient, error) {
	conn, err := grpc.NewClient("unix://"+socketPath, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: 10 * time.Second,
			Backoff:           backoff.DefaultConfig,
		}),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    30 * time.Second,
			Timeout: 10 * time.Second,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect logger : %v", err)
	}
	client := ls.NewLoggingServiceClient(conn)
	// Contact the server and print out its response.
	return &LoggingClient{LoggingServiceClient: client, nameService: nameService, socketPath: socketPath, conn: conn}, nil
}

func (lc *LoggingClient) Close() error {
	return lc.conn.Close()
}

func (lc *LoggingClient) Info(kvPairs ...KeyValue) (err error) {
	if err := lc.isSocketAvailable(); err != nil {
		return err
	}
	message := ConvertToMap(kvPairs...)
	response, err := lc.Logging(context.Background(), &ls.LoggingRequest{Level: ls.LoggingRequest_INFO, Message: message, ServiceName: lc.nameService})
	if err != nil {
		return fmt.Errorf("error in logging : %v", err)
	}
	if response.Result.Code != ls.Error_NONE {
		log.Println("Response:", response.Result)
	}
	return nil
}

func (lc *LoggingClient) Debug(kvPairs ...KeyValue) (err error) {
	if err := lc.isSocketAvailable(); err != nil {
		return err
	}
	message := ConvertToMap(kvPairs...)
	response, err := lc.Logging(context.Background(), &ls.LoggingRequest{Level: ls.LoggingRequest_DEBUG, Message: message, ServiceName: lc.nameService})
	if err != nil {
		return fmt.Errorf("error in logging : %v", err)
	}
	if response.Result.Code != ls.Error_NONE {
		log.Println("Response:", response.Result)
	}
	return nil
}

func (lc *LoggingClient) Warn(kvPairs ...KeyValue) (err error) {
	if err := lc.isSocketAvailable(); err != nil {
		return err
	}
	message := ConvertToMap(kvPairs...)
	response, err := lc.Logging(context.Background(), &ls.LoggingRequest{Level: ls.LoggingRequest_WARN, Message: message, ServiceName: lc.nameService})
	if err != nil {
		return fmt.Errorf("error in logging : %v", err)
	}
	if response.Result.Code != ls.Error_NONE {
		log.Println("Response:", response.Result)
	}
	return nil
}

func (lc *LoggingClient) Critical(kvPairs ...KeyValue) (err error) {
	if err := lc.isSocketAvailable(); err != nil {
		return err
	}
	message := ConvertToMap(kvPairs...)
	response, err := lc.Logging(context.Background(), &ls.LoggingRequest{Level: ls.LoggingRequest_CRITICAL, Message: message, ServiceName: lc.nameService})
	if err != nil {
		return fmt.Errorf("error in logging : %v", err)
	}
	if response.Result.Code != ls.Error_NONE {
		log.Println("Response:", response.Result)
	}
	return nil
}

func (lc *LoggingClient) isSocketAvailable() error {
	fileInfo, err := os.Stat(lc.socketPath)
	if err != nil {
		return err
	}
	// Проверяем, что это socket-файл (не обычный файл или директория)
	if fileInfo.Mode()&os.ModeSocket == 0 {
		return fmt.Errorf("socket %s already exists", lc.socketPath)
	}
	return nil
}

func (lc *LoggingClient) MonitorConnection(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			lc.mutex.Lock()
			if lc.conn != nil {
				lc.mutex.Unlock()
				fmt.Println("Connection established")
				return
			}
			fmt.Println("Connection...")
			lc.conn.Connect()
			lc.mutex.Unlock()
		}
	}
}
