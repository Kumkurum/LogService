package test

import (
	"context"
	ls "github.com/Kumkurum/LogService/internal/transport"
	logger "github.com/Kumkurum/LogService/pkg/log_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	// Указываем путь к Unix Socket
	socketPath := "unix:///tmp/grpc.sock"

	// Настраиваем подключение
	conn, err := grpc.NewClient(socketPath, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Создаём клиент
	client := ls.NewLoggingServiceClient(conn)

	// Вызываем RPC-метод
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	message := make(map[string]string)
	message["message"] = "hello world"
	resp, err := client.Logging(ctx, &ls.LoggingRequest{Level: ls.LoggingRequest_INFO, Message: message, ServiceName: "test"})
	if err != nil {
		log.Fatalf("RPC failed: %v", err)
	}

	log.Println("Response:", resp.Result)

}

func TestLoggingClient(t *testing.T) {
	socketPath := "/tmp/grpc.sock"
	client, err := logger.NewLoggingClient(socketPath, "TEST")
	defer client.Close()
	if err != nil {
		log.Printf("Failed to connect %v", err)
	}

	err = client.Info(logger.KeyValue{Value: "first", Key: "INFO"}, logger.KeyValue{Value: "first", Key: "test"}, logger.KeyValue{Value: "first", Key: "test"}, logger.KeyValue{Value: "first", Key: "test"})
	if err != nil {
		log.Printf("Failed to connect %v", err)
	}
	err = client.Debug(logger.KeyValue{Value: "first", Key: "DEBUG"})
	if err != nil {
		log.Printf("Failed to connect %v", err)
	}
	time.Sleep(5 * time.Second)
	err = client.Debug(logger.KeyValue{Value: "first", Key: "DEBUG"})
	if err != nil {
		log.Printf("Failed to connect %v", err)
	} else {
		log.Printf("Debug logging client connected")
	}
	time.Sleep(5 * time.Second)
	err = client.Debug(logger.KeyValue{Value: "first", Key: "DEBUG"})
	if err != nil {
		log.Printf("Failed to connect %v", err)
	} else {
		log.Printf("Debug logging client connected")
	}
	time.Sleep(5 * time.Second)
	err = client.Warn(logger.KeyValue{Value: "first", Key: "WARN"})
	if err != nil {
		log.Printf("Failed to connect %v", err)
	} else {
		log.Printf("Debug logging client connected")
	}
	time.Sleep(5 * time.Second)
	err = client.Critical(logger.KeyValue{Value: "first", Key: "CRITICAL"})
	if err != nil {
		log.Printf("Failed to connect %v", err)
	} else {
		log.Printf("Debug logging client connected")
	}
}
