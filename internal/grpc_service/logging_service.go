package grpc_service

import (
	"context"
	loggingservice "github.com/Kumkurum/LogService/internal/transport"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
)

type LoggingService struct {
	loggingservice.LoggingServiceServer
	logFile *lumberjack.Logger
	logger  *slog.Logger
}

func NewLoggingService(fileName string, maxSize, maxBackups int, isCompress bool) *LoggingService {
	logFile := &lumberjack.Logger{
		Filename:   fileName, // Path
		MaxSize:    maxSize,  // MB
		MaxBackups: maxBackups,
		Compress:   isCompress,
	}
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(logFile, opts))

	return &LoggingService{
		logFile: logFile,
		logger:  logger,
	}
}

func (ls *LoggingService) Logging(ctx context.Context, r *loggingservice.LoggingRequest) (*loggingservice.LoggingResponse, error) {

	attrs := make([]slog.Attr, 0, len(r.Message)+1)
	attrs = append(attrs, slog.String("service", r.ServiceName))
	for k, v := range r.Message {
		attrs = append(attrs, slog.String(k, v))
	}
	switch r.Level {
	case loggingservice.LoggingRequest_DEBUG:
		ls.logger.LogAttrs(ctx, slog.LevelDebug, r.ServiceName, attrs...)
	case loggingservice.LoggingRequest_INFO:
		ls.logger.LogAttrs(ctx, slog.LevelInfo, r.ServiceName, attrs...)
	case loggingservice.LoggingRequest_WARN:
		ls.logger.LogAttrs(ctx, slog.LevelWarn, r.ServiceName, attrs...)
	case loggingservice.LoggingRequest_CRITICAL:
		ls.logger.LogAttrs(ctx, slog.LevelError, r.ServiceName, attrs...)
	default:
		return &loggingservice.LoggingResponse{
			Result: &loggingservice.Error{Code: loggingservice.Error_ERROR},
		}, nil
	}
	return &loggingservice.LoggingResponse{
		Result: &loggingservice.Error{Code: loggingservice.Error_NONE},
	}, nil
}
func (ls *LoggingService) LogAfterStart() {
	ls.logger.Info("Logger start!")
}
func (ls *LoggingService) LogAfterEnd() {
	ls.logger.Info("Logger end!")
}
