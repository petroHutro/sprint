package main

import (
	"net/http"
	"sprint/internal/app/config"
	"sprint/internal/app/router"
	"strconv"
)

func main() {
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
