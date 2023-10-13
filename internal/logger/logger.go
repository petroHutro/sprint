package logger

import (
	"fmt"
	"net/http"
	"sprint/internal/config"

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
type Logger struct {
	logger *zap.Logger
}

var Log Logger

func (l *Logger) Error(msg string, fields ...interface{}) {
	l.logger.Error(fmt.Sprintf(msg, fields...))
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, fields...))
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.logger.Panic(msg, fields...)
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func (l *Logger) Shutdown() {
	l.logger.Sync()
}

func newFileLogger(logFile string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{logFile}
	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("cannot create file logger Build: %w", err)
	}
	return logger, nil
}

func newConsoleLogger() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = ""
	config.EncoderConfig.EncodeCaller = nil
	config.DisableStacktrace = true
	logger, _ := config.Build()
	return logger, nil
}

func newMultiLogger(filePath string) (*zap.Logger, error) {
	consoleLogger, err := newConsoleLogger()
	if err != nil {
		return nil, fmt.Errorf("cannot create consol logger in multi logger: %w", err)
	}
	fileLogger, err := newFileLogger(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot create file logger in multi logger: %w", err)
	}
	multiCore := zapcore.NewTee(fileLogger.Core(), consoleLogger.Core())
	logger := zap.New(multiCore)
	return logger, nil
}

func InitLogger(conf config.Logger) error {
	if conf.MultiFlag {
		logger, err := newMultiLogger(conf.FilePath)
		if err != nil {
			return fmt.Errorf("cannot create multi logger: %w", err)
		}
		Log.logger = logger
		return nil
	} else if conf.FileFlag {
		logger, err := newFileLogger(conf.FilePath)
		if err != nil {
			return fmt.Errorf("cannot create file logger: %w", err)

		}
		Log.logger = logger
		return nil
	}
	logger, err := newConsoleLogger()
	if err != nil {
		return fmt.Errorf("cannot create consol logger: %w", err)
	}
	Log.logger = logger
	return nil
}
