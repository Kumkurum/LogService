package grpc_service

import (
	"context"
	"log/slog"

	loggingservice "github.com/Kumkurum/LogService/internal/transport"
	"gopkg.in/natefinch/lumberjack.v2"
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
	switch r.Level {
	case loggingservice.LoggingRequest_DEBUG:
		ls.logger.LogAttrs(ctx, slog.LevelDebug, r.Message, slog.String("service", r.ServiceName))
	case loggingservice.LoggingRequest_INFO:
		ls.logger.LogAttrs(ctx, slog.LevelInfo, r.Message, slog.String("service", r.ServiceName))
	case loggingservice.LoggingRequest_WARN:
		ls.logger.LogAttrs(ctx, slog.LevelWarn, r.Message, slog.String("service", r.ServiceName))
	case loggingservice.LoggingRequest_CRITICAL:
		ls.logger.LogAttrs(ctx, slog.LevelError, r.Message, slog.String("service", r.ServiceName))
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
	ls.logger.Info("Start", slog.String("service", "Logger"))
}
func (ls *LoggingService) LogAfterEnd() {
	ls.logger.Info("Stop", slog.String("service", "Logger"))
}
