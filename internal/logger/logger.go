package logger

import (
	"net/http"
	"os"
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

type Logger struct {
	logger *zap.Logger
}

func CloseFileLoger(logger *Logger) {
	logger.logger.Sync()
}

func NewFileLogger(logFile string) (*Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{logFile}
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{logger: logger}, nil
}

func NewConsoleLogger() (*Logger, error) {
	config := zap.NewDevelopmentConfig()
	logger, _ := config.Build()
	logger.Info("Running server", zap.String("address", "local"))
	return &Logger{logger: logger}, nil
}

// func (l *FileLogger) Error(msg string, fields ...zap.Field) {
// 	l.logger.Error(msg, fields...)
// }

// func (l *FileLogger) Info(msg string, fields ...zap.Field) {
// 	l.logger.Info(msg, fields...)
// }

func CreateMultiLogger(fileLogger *Logger) *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg.EncoderConfig),
		zapcore.Lock(zapcore.AddSync(os.Stdout)),
		cfg.Level,
	)

	multiLogger := zap.New(zapcore.NewTee(consoleCore, fileLogger.logger.Core()))
	return multiLogger
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

// var Log *zap.Logger = zap.NewNop()

// func Initialize(level string) error {
// 	lvl, err := zap.ParseAtomicLevel(level)
// 	if err != nil {
// 		return err
// 	}
// 	// config := zap.Config добавить !!!
// 	// создаём новую конфигурацию логера
// 	cfg := zap.NewProductionConfig()
// 	cfg.Level = lvl
// 	zl, err := cfg.Build()
// 	if err != nil {
// 		return err
// 	}
// 	// устанавливаем синглтон
// 	Log = zl
// 	return nil
// }

// RequestLogger — middleware-логер для входящих HTTP-запросов.
// func RequestLogger(h http.HandlerFunc) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		Log.Debug("got incoming HTTP request",
// 			zap.String("method", r.Method),
// 			zap.String("path", r.URL.Path),
// 		)
// 		h(w, r)
// 	})
// }

// func WithLogging(h http.Handler) http.Handler {
// 	logFn := func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		responseData := &responseData{
// 			status: 0,
// 			size:   0,
// 		}
// 		lw := loggingResponseWriter{
// 			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
// 			responseData:   responseData,
// 		}
// 		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter
// 		duration := time.Since(start)
// 		sugar := Log.Sugar()
// 		sugar.Infoln(
// 			"uri", r.RequestURI,
// 			"method", r.Method,
// 			"status", responseData.status, // получаем перехваченный код статуса ответа
// 			"duration", duration,
// 			"size", responseData.size, // получаем перехваченный размер ответа
// 		)
// 	}
// 	return http.HandlerFunc(logFn)
// }

// type Logger interface {
// 	Info(msg string, fields ...zap.Field)
// 	Error(msg string, fields ...zap.Field)
// }

// type DevelopmentLogger struct {
// 	logger *zap.SugaredLogger
// }

// func NewDevelopmentLogger() Logger {
// 	config := zap.NewDevelopmentConfig()
// 	logger, _ := config.Build()
// 	// logger.Sugar().Infoln()
// 	logger.Info("Running server", zap.String("address", "local"))
// 	return &DevelopmentLogger{logger: logger.Sugar()}
// }

// func (l *DevelopmentLogger) Info(msg string, fields ...zap.Field) {
// 	l.logger.Info(msg, fields)
// }

// func (l *DevelopmentLogger) Error(msg string, fields ...zap.Field) {
// 	l.logger.Error(msg, fields)
// }

// type FileLogger struct {
// 	logger *zap.SugaredLogger
// }

// func NewFileLogger() Logger {
// 	config := zap.NewProductionConfig()
// 	config.OutputPaths = []string{"app.log"} // Здесь указываем путь к файлу логов
// 	logger, _ := config.Build()
// 	logger.Info("Running server", zap.String("address", "local"))
// 	return &FileLogger{logger: logger.Sugar()}
// }

// func (l *FileLogger) Info(msg string, fields ...zap.Field) {
// 	l.logger.Info(msg, fields)
// }

// func (l *FileLogger) Error(msg string, fields ...zap.Field) {
// 	l.logger.Error(msg, fields)
// }

// // type MyService struct {
// // 	Logger Logger
// // }

// // func NewMyService(logger Logger) *MyService {
// // 	return &MyService{Logger: logger}
// // }

// // func (s *MyService) DoSomething() {
// // 	s.Logger.Info("Doing something...")
// // 	// Ваша логика здесь
// // }

// // func WithLogging(logger Logger) http.Handler {
// // 	logFn := func(w http.ResponseWriter, r *http.Request) {
// // 		start := time.Now()

// // 		responseData := &responseData{
// // 			status: 0,
// // 			size:   0,
// // 		}
// // 		lw := loggingResponseWriter{
// // 			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
// // 			responseData:   responseData,
// // 		}
// // 		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

// // 		duration := time.Since(start)
// // 		sugar := Log.Sugar()
// // 		sugar.Infoln(
// // 			"uri", r.RequestURI,
// // 			"method", r.Method,
// // 			"status", responseData.status, // получаем перехваченный код статуса ответа
// // 			"duration", duration,
// // 			"size", responseData.size, // получаем перехваченный размер ответа
// // 		)
// // 	}
// // 	return http.HandlerFunc(logFn)
// // }

// func LogMiddleware(log *MultiLogger) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			// Логируем запрос
// 			start := time.Now()

// 			responseData := &responseData{
// 				status: 0,
// 				size:   0,
// 			}
// 			lw := loggingResponseWriter{
// 				ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
// 				responseData:   responseData,
// 			}
// 			next.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

// 			duration := time.Since(start)

// 			// fields := []zap.Field{
// 			// 	zap.String("uri", r.RequestURI),
// 			// 	zap.String("method", r.Method),
// 			// 	zap.Int("status", responseData.status),
// 			// 	zap.Duration("duration", duration),
// 			// 	zap.Int("size", responseData.size),
// 			// }
// 			// log.Info("Received request", fields...)

// 			logFields := &LogFields{
// 				URI:      r.RequestURI,
// 				Method:   r.Method,
// 				Status:   responseData.status,
// 				Duration: duration,
// 				Size:     responseData.size,
// 			}
// 			fmt.Println(r.Method)
// 			log.Info("Your message", zap.Any("fields", logFields.ToZapFields()))

// 			// for _, logger := range loggers {
// 			// 	logger.Info("Received request", fields...)
// 			// }

// 			// logger.Info(fmt.Sprintf("Received request: %s %s", r.Method, r.URL.Path))
// 			// Передаем запрос следующему обработчику
// 			// next.ServeHTTP(w, r)
// 		})
// 	}
// }

// type MultiLogger struct {
// 	Loggers []Logger
// }

// func NewMultiLogger(loggers ...Logger) *MultiLogger {
// 	return &MultiLogger{Loggers: loggers}
// }

// func (m *MultiLogger) Info(msg string, fields ...zap.Field) {
// 	for _, logger := range m.Loggers {
// 		logger.Info(msg, fields...)
// 	}
// }

// func (m *MultiLogger) Error(msg string, fields ...zap.Field) {
// 	for _, logger := range m.Loggers {
// 		logger.Error(msg, fields...)
// 	}
// }

// type LogFields struct {
// 	URI      string
// 	Method   string
// 	Status   int
// 	Duration time.Duration
// 	Size     int
// }

// func (lf *LogFields) ToZapFields() []zap.Field {
// 	return []zap.Field{
// 		zap.String("uri", lf.URI),
// 		zap.String("method", lf.Method),
// 		zap.Int("status", lf.Status),
// 		zap.Duration("duration", lf.Duration),
// 		zap.Int("size", lf.Size),
// 	}
// }

// type Logger interface {
// 	Error(msg string, fields ...zap.Field)
// 	Info(msg string, fields ...zap.Field)
// }
