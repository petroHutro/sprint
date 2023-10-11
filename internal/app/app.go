package app

import (
	"net/http"
	"sprint/internal/config"
	"sprint/internal/logger"
	"sprint/internal/router"
	"sprint/internal/storage"
	"strconv"
)

// rename 6:30
func Run() error {
	flags := config.ConfigureServer()
	if err := logger.NewLogger(flags.Logger); err != nil {
		logger.Log.Panic(err.Error())
	}
	defer logger.Log.CloseFileLoger()
	if err := storage.LoadURL(string(flags.FileStoragePath)); err != nil {
		logger.Log.Panic(err.Error())
	}
	r := router.Router(flags)
	address := flags.Host + ":" + strconv.Itoa(flags.Port)
	logger.Log.Info("Running server: address:%s port:%d", flags.Host, flags.Port)
	return http.ListenAndServe(address, r)
}
