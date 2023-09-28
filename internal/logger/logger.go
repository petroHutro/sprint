package logger

import (
	"net/http"
	"sprint/internal/config"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	responseData struct {
		status int
		size   int
	}
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func CloseFileLoger(logger *zap.Logger) {
	logger.Sync()
}

func newFileLogger(logFile string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{logFile}
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func newConsoleLogger() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()
	return logger, nil
}

func newMultiLogger(filePath string) (*zap.Logger, error) {
	consoleLogger, err := newConsoleLogger()
	if err != nil {
		return nil, err
	}
	fileLogger, err := newFileLogger(filePath)
	if err != nil {
		return nil, err
	}
	multiCore := zapcore.NewTee(fileLogger.Core(), consoleLogger.Core())
	logger := zap.New(multiCore)
	return logger, nil
}

func NewLogger(conf config.Logger) (*zap.Logger, error) {
	if conf.MultiFlag {
		logger, err := newMultiLogger(conf.FilePath)
		if err != nil {
			return nil, err
		}
		return logger, nil
	} else if conf.FileFlag {
		logger, err := newFileLogger(conf.FilePath)
		if err != nil {
			return nil, err
		}
		return logger, nil
	}
	logger, err := newConsoleLogger()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func LoggingMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			responseData := &responseData{
				status: 0,
				size:   0,
			}
			lw := loggingResponseWriter{
				ResponseWriter: w,
				responseData:   responseData,
			}
			next.ServeHTTP(&lw, r)
			duration := time.Since(start)
			fields := []zap.Field{
				zap.String("uri", r.RequestURI),
				zap.String("method", r.Method),
				zap.Int("status", responseData.status),
				zap.Duration("duration", duration),
				zap.Int("size", responseData.size),
			}
			logger.Info("Received request", fields...)
		})
	}
}
