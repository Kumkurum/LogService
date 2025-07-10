package main

import (
	"LoggingService/internal/grpc_service"
	logging_service "LoggingService/internal/transport"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var socketPath string //"tcp" : ":50055"
	var version, help bool
	var fileName string
	var isCompress bool
	var maxSize, maxBackups int
	flag.StringVar(&socketPath, "socketPath", "/tmp/grpc.sock", "address of grpc service, default: /tmp/grpc.sock")
	flag.BoolVar(&version, "version", false, "Version service")
	flag.BoolVar(&help, "help", false, "Help how to use service")

	flag.StringVar(&fileName, "fileName", "log/service.log", "path to log file, default log/service.log")
	flag.BoolVar(&isCompress, "isCompress", true, "is need compress log file, default true")
	flag.IntVar(&maxSize, "maxSize", 100, "max size of log file, default 100")
	flag.IntVar(&maxBackups, "maxBackups", 3, "max size of count backups log files, default 3")
	flag.Parse()
	if version {
		fmt.Println("Version 0.0.1")
		return
	}
	if help {
		fmt.Println("This is a service for logging messages from other services")
		fmt.Println("flag socketPath to set port for grpc service, default = 50055")
		fmt.Println("network - is unix socket")
		fmt.Println("just easy to use!")
		return
	}
	// Удаляем старый socket, если он существует
	if err := os.RemoveAll(socketPath); err != nil {
		log.Fatalf("Failed to remove old socket: %v", err)
	}

	// Создаём слушатель для Unix Socket
	lis, err := net.Listen("unix", socketPath)
	defer func(lis net.Listener) {
		err := lis.Close()
		log.Println("Closing socket")
		if err != nil {
			fmt.Printf("Failed to close socket: %v", err)
		}
	}(lis)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stop
		lis.Close()           // Закрываем сокет
		os.Remove(socketPath) // Явное удаление
		os.Exit(0)
	}()

	// Настраиваем права (чтобы другие сервисы могли подключаться)
	if err := os.Chmod(socketPath, 0777); err != nil {
		log.Fatalf("Failed to set socket permissions: %v", err)
	}

	s := grpc.NewServer()
	loggingService := grpc_service.NewLoggingService(fileName, maxSize, maxBackups, isCompress)
	loggingService.LogAfterStart()
	defer loggingService.LogAfterEnd()

	logging_service.RegisterLoggingServiceServer(s, loggingService)
	_ = s.Serve(lis)
}
