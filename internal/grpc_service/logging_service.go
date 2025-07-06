package grpc_service

import (
	loggingservice "LoggingService/internal/transport"
	"context"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
)

type LoggingService struct {
	sensorsStorage loggingservice.LoggingServiceServer
	logFile        *lumberjack.Logger
	logger         *slog.Logger
}

func NewLoggingService(fileName string, maxSize int, isCompress bool) *LoggingService {
	logFile := &lumberjack.Logger{
		Filename: fileName, //"logs/server.log"
		MaxSize:  maxSize,  // MB
		Compress: isCompress,
	}
	logger := slog.New(slog.NewJSONHandler(logFile, nil))
	return &LoggingService{
		logFile: logFile,
		logger:  logger,
	}
}

func (l *LoggingService) GetSensors(ctx context.Context, r *loggingservice.LoggingRequest) (*loggingservice.LoggingResponse, error) {

	attrs := make([]slog.Attr, 0, len(r.Message)+1)
	attrs = append(attrs, slog.String("service", r.ServiceName))
	for k, v := range r.Message {
		attrs = append(attrs, slog.String(k, v))
	}
	switch r.Level {
	case loggingservice.LoggingRequest_DEBUG:
		l.logger.LogAttrs(ctx, slog.LevelDebug, r.ServiceName, attrs...)
	case loggingservice.LoggingRequest_INFO:
		l.logger.LogAttrs(ctx, slog.LevelInfo, r.ServiceName, attrs...)
	case loggingservice.LoggingRequest_WARN:
		l.logger.LogAttrs(ctx, slog.LevelWarn, r.ServiceName, attrs...)
	case loggingservice.LoggingRequest_CRITICAL:
		l.logger.LogAttrs(ctx, slog.LevelError, r.ServiceName, attrs...)
	default:
		return &loggingservice.LoggingResponse{
			Result: &loggingservice.Error{Code: loggingservice.Error_ERROR},
		}, nil
	}
	return &loggingservice.LoggingResponse{
		Result: &loggingservice.Error{Code: loggingservice.Error_NONE},
	}, nil
}
