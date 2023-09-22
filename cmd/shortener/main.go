package main

import (
	"net/http"
	"sprint/internal/config"
	"sprint/internal/logger"
	"sprint/internal/router"
	"strconv"

	"go.uber.org/zap"
)

func main() {
	fileLogger, err := logger.NewFileLogger("app.log")
	if err != nil {
		panic(err)
	}
	defer logger.CloseFileLoger(fileLogger)
	multiLogger := logger.CreateMultiLogger(fileLogger)
	multiLogger.Info("Running server", zap.String("address", "local"))
	if err := run(multiLogger); err != nil {
		panic(err)
	}
}

func run(loggers *zap.Logger) error {
	flags := config.ConfigureServer()
	r := router.Router(flags, loggers)
	address := flags.Host + ":" + strconv.Itoa(flags.Port)
	return http.ListenAndServe(address, r)
}
