package main

import (
	"net/http"
	"sprint/internal/config"
	"sprint/internal/logger"
	"sprint/internal/router"
	"sprint/internal/storage"
	"strconv"

	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	flags := config.ConfigureServer()
	log, err := logger.NewLogger(flags.Logger)
	if err != nil {
		panic(err)
	}
	defer logger.CloseFileLoger(log)
	storage.LoadURL(string(flags.FileStoragePath))
	log.Info("Running server", zap.String("address", "local"))
	r := router.Router(flags, log)
	address := flags.Host + ":" + strconv.Itoa(flags.Port)
	return http.ListenAndServe(address, r)
}
