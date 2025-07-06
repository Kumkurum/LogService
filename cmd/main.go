package main

import (
	"LoggingService/internal/grpc_service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	s := grpc.NewServer()
	grpcService := &grpc_service.LoggingService{}
	homeSyncGrpc.RegisterHomeSyncGrpcServiceServer(s, grpcService)
	// Открыть порт 50051 для приема сообщений
	lis, err := net.Listen("tcp", ":"+addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Начать цикл приема и обработку запросов
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
