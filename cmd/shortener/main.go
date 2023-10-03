package main

import (
	"log"
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
		log.Panic(err.Error())
	}
}

func run() error {
	flags := config.ConfigureServer()
	newLog, err := logger.NewLogger(flags.Logger)
	if err != nil {
		log.Panic(err.Error())
	}
	defer logger.CloseFileLoger(newLog)
	if err := storage.LoadURL(string(flags.FileStoragePath)); err != nil {
		newLog.Panic(err.Error())
	}
	r := router.Router(flags, newLog)
	address := flags.Host + ":" + strconv.Itoa(flags.Port)
	newLog.Info("Running server", zap.String("address", flags.Host), zap.Int("port", flags.Port))
	return http.ListenAndServe(address, r)
}
