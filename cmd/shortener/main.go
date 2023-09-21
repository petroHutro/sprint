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
	if err := logger.Initialize(); err != nil {
		panic(err)
	}
	logger.Log.Info("Running server", zap.String("address", "local"))
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	flags := config.ConfigureServer()
	r := router.Router(flags)
	address := flags.Host + ":" + strconv.Itoa(flags.Port)
	return http.ListenAndServe(address, r)
}
