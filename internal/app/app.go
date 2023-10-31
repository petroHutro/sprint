package app

import (
	"context"
	"fmt"
	"net/http"
	"sprint/internal/authorization"
	"sprint/internal/compression"
	"sprint/internal/config"
	"sprint/internal/handlers"
	"sprint/internal/logger"
	"sprint/internal/router"
	"sprint/internal/storage"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type App struct {
	conf         *config.Flags
	db           *storage.StorageBase
	delQueryChan chan storage.QueryDelete
	router       *chi.Mux
}

func (a *App) FlushDelete() {
	ticker := time.NewTicker(10 * time.Second)

	var requests []storage.QueryDelete

	for {
		select {
		case request := <-a.delQueryChan:
			requests = append(requests, request)
		case <-ticker.C:
			if len(requests) == 0 {
				continue
			}
			go func(requests []storage.QueryDelete) {
				var id []string
				var data []string

				for _, request := range requests {
					id = append(id, request.ID)
					data = append(data, request.Data)
				}

				err := a.db.DeleteURL(context.Background(), a.conf.FilePath, id, data)
				if err != nil {
					logger.Error("cannot delete:%v", err)
				}
			}(requests)
			requests = nil
		}
	}
}

type QD = storage.QueryDelete

func newApp() (*App, error) {
	conf := config.LoadServerConfigure()
	if err := logger.InitLogger(conf.Logger); err != nil {
		return nil, fmt.Errorf("cannot init logger: %w", err)
	}

	storage, err := storage.InitStorage(&conf.Storage)
	if err != nil {
		return nil, fmt.Errorf("cannot init storage: %w", err)
	}

	router := router.CreateRouter()

	logger.Info("Running server: address:%s port:%d", conf.Host, conf.Port)

	return &App{conf: conf, db: storage, router: router, delQueryChan: make(chan QD)}, nil
}

func (a *App) createMiddlewareHandlers() {
	a.router.Use(logger.LoggingMiddleware)
	a.router.Use(authorization.AuthorizationMiddleware)
	a.router.Use(compression.GzipMiddleware)
}

func (a *App) createHandlers() {
	a.router.Route("/", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerPost(w, r, a.conf.BaseURL, a.conf.FileStoragePath, a.db)
		})
	})

	a.router.Route("/ping", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerPing(w, r, a.db)
		})
	})

	a.router.Route("/{id:[a-zA-Z0-9]+}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerGet(w, r, a.db)
		})
	})

	a.router.Route("/api", func(r chi.Router) {
		r.Post("/shorten", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerPostAPI(w, r, a.conf.BaseURL, a.conf.FileStoragePath, a.db)
		})
		r.Post("/shorten/batch", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerPostBatch(w, r, a.conf.BaseURL, a.conf.FileStoragePath, a.db)
		})
		r.Get("/user/urls", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerGetUrls(w, r, a.conf.BaseURL, a.db)
		})
		r.Delete("/user/urls", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerDelete(w, r, a.delQueryChan)
		})
	})
}

func Run() error {
	app, err := newApp()
	if err != nil {
		logger.Panic(err.Error())
	}
	defer logger.Shutdown()

	app.createMiddlewareHandlers()
	app.createHandlers()

	go app.FlushDelete()

	address := app.conf.Host + ":" + strconv.Itoa(app.conf.Port)

	return http.ListenAndServe(address, app.router)
}
