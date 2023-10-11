package main

import (
	"log"
	"sprint/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Panic(err.Error())
	}
}
