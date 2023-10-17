package app

import (
	"net/http"
	"sprint/internal/config"
	"sprint/internal/logger"
	"sprint/internal/router"
	"sprint/internal/storage"
	"strconv"
)

func Run() error {
	conf := config.LoadServerConfigure()
	if err := logger.InitLogger(conf.Logger); err != nil {
		logger.Panic(err.Error())
	}
	defer logger.Shutdown()

	storage, err := storage.InitStorage(&conf.Storage)
	if err != nil {
		logger.Panic(err.Error())
	}

	r := router.Create(conf, storage)
	address := conf.Host + ":" + strconv.Itoa(conf.Port)
	logger.Info("Running server: address:%s port:%d", conf.Host, conf.Port)

	return http.ListenAndServe(address, r)
}
